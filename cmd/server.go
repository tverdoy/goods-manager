package main

import (
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"goods-manager/internal/app"
	"goods-manager/internal/cache/redis"
	"log"
	"os"
	"strconv"
)

const (
	defaultAddress = ":8080"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// RunHTTPServe run HTTP server at `address`
//
// @title			Goods manager
// @version		1.0
// @description	Goods manager APIr
//
// @host		localhost:8080
// @BasePath	/
func main() {
	//prepare database
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	pgInfo := fmt.Sprintf("host = %s port = %s "+
		"user = %s password = %s dbname = %s sslmode = disable", host, port, user, password, dbname)

	db, err := app.ConnectToPostgres(pgInfo)
	if err != nil {
		log.Panicln("failed connect to database:", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Panicln("failed close database connection:", err)
		}
	}()

	// connect to redis
	redisAddr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	redisDBInt, err := strconv.Atoi(redisDB)
	if err != nil {
		log.Panicln("failed convert redis db to int:", err)
	}

	redisClient, err := app.ConnectToRedis(redisAddr, redisPassword, redisDBInt)
	if err != nil {
		log.Panicln("failed connect to redis:", err)
	}

	cache := redis.NewCache(redisClient)

	// connect to clickhouse
	clickhouseAddr := os.Getenv("CLICKHOUSE_HOST") + ":" + os.Getenv("CLICKHOUSE_PORT")
	clickhouseClient, err := app.ConnectToClickHouse(clickhouseAddr)
	if err != nil {
		log.Panicln("failed connect to clickhouse:", err)
	}

	// connect to nats
	natsAddr := os.Getenv("NATS_HOST") + ":" + os.Getenv("NATS_PORT")
	natsClient, err := app.ConnectToNats(natsAddr)
	if err != nil {
		log.Panicln("failed connect to nats:", err)
	}

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}

	if err := app.RunHTTPServe(address, db, cache, natsClient, clickhouseClient); err != nil {
		log.Panic(err)
	}
}
