package config

import "time"

type Config struct {
	Server ServerConfig
	Database DatabaseConfig
	JWT JWTConfig
}

type ServerConfig struct{
	Port string
	Env string
}

type DatabaseConfig struct{
	URI string
	Name string
	Timeout time.Duration
}