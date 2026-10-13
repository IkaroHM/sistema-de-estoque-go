package config

import (
	"os"
	"time"
)

type DbConfig struct{
	Db_URL string
	MaxConns int32
	MinConns int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	QueryTimeout time.Duration
}

func carregarConfigsDb() *DbConfig{
	return &DbConfig{
		Db_URL: os.Getenv("DATABASE_URL"),
		MaxConns: 10,
		MinConns: 5,
		MaxConnLifetime: 15 * time.Second,
		MaxConnIdleTime: 10 * time.Second,
		QueryTimeout: 5 * time.Second,
		}
	}