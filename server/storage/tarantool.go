package storage

import (
	"context"
	"log"
	"time"

	"github.com/tarantool/go-tarantool/v2"
)

func NewTarantoolConnection(ctx context.Context, host, port, user, password string) (*tarantool.Connection, error) {
	dialer := tarantool.NetDialer{
		Address:  host + ":" + port,
		User:     user,
		Password: password,
	}
	opts := tarantool.Opts{
		Timeout: 5 * time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		log.Printf("Failed to connect to Tarantool: %v", err)

		return nil, err
	}

	log.Println("Connected to Tarantool")

	return conn, nil
}
