package config

import (
	"fmt"
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

func carregarConfigsDb() (*DbConfig, error){

	db_url := os.Getenv("DATABASE_URL")
	if db_url == ""{
		return nil, fmt.Errorf("A variavel de ambiente DATABASE_URL está vazia!")
	}

	return &DbConfig{
		Db_URL: db_url,
		MaxConns: 10,
		MinConns: 5,
		MaxConnLifetime: 15 * time.Second,
		MaxConnIdleTime: 10 * time.Second,
		QueryTimeout: 5 * time.Second,
		}, nil
	}