package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type DBConfig struct {
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

func DefaultDBConfig() DBConfig {
	return DBConfig{
		Host:     "localhost",
		Username: "user",
		Password: "password",
		DBName:   "mydb",
	}
}

type Config struct {
	Port     int64    `json:"port"`
	Database DBConfig `json:"database"`
}

func (c *DBConfig) ConnectionInfo() string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.DBName)
}

func DefaultConfig() *Config {
	return &Config{
		Port:     3002,
		Database: DefaultDBConfig(),
	}
}

func LoadConfig() *Config {
	f, err := os.Open(".auction_config")
	if err != nil {
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	err = json.NewDecoder(f).Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded .auction_config")
	return &c
}
