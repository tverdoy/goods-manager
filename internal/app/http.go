package app

import (
	"database/sql"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"goods-manager/internal/cache"
	"goods-manager/internal/good/controller"
	"goods-manager/internal/good/repository"
	"goods-manager/internal/good/usecase"
	repository2 "goods-manager/internal/logger/repository"
	usecase2 "goods-manager/internal/logger/usecase"
	"goods-manager/internal/logger/workers"
	"goods-manager/internal/transactor"
	"log"

	_ "goods-manager/internal/docs"
)

// RunHTTPServe run HTTP server at `address`
//
// @title			Goods manager
// @version		1.0
// @description	Goods manager APIr
//
// @host		localhost:8080
// @BasePath	/
func RunHTTPServe(address string, db *sql.DB, cache cache.Cache, nats *nats.Conn, clickhouse driver.Conn) error {
	r := gin.New()

	// Init middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Init newTransactor
	newTransactor := transactor.NewTransactor(db)

	// Init repository layer
	goodRepo := repository.NewGoodRepository(newTransactor)
	goodRepoCache := repository.NewGoodRepositoryCache(cache, goodRepo)

	loggerRepo := repository2.NewLoggerRepository(clickhouse)

	// Init usecase layer
	loggerUsecase := usecase2.NewLoggerUsecase(nats, loggerRepo)

	goodUsecase := usecase.NewGoodUsecase(goodRepoCache, loggerUsecase, newTransactor)

	// Init controller layer
	goodController := controller.NewGoodController(goodUsecase)

	// Add route
	goodR := r.Group("/good")

	goodR.POST("/create", goodController.Create)
	goodR.GET("/list", goodController.List)
	goodR.PATCH("/update", goodController.Update)
	goodR.DELETE("/remove", goodController.Delete)

	goodR.PATCH("/reprioritiize", goodController.Reprioritize)

	// Init swagger doc
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Println("starting logger worker...")
	if err := workers.NewLoggerWorker(nats, loggerUsecase).Run(); err != nil {
		log.Panicln("failed start logger worker", err)
	}

	log.Println("starting server...")
	return r.Run(address)
}
