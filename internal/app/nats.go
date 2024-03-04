package app

import (
	"errors"
	"github.com/nats-io/nats.go"
)

// ConnectToNats connect to nats
func ConnectToNats(url string) (*nats.Conn, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	if !conn.IsConnected() {
		return nil, errors.New("failed connect to nats")
	}

	return conn, nil
}
