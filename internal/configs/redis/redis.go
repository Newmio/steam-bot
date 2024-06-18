package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DbName   string
}

func OpenDb() (*redis.Client, error) {
	v := viper.New()
	v.AddConfigPath("internal/app/storage/redis")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return initDb(Config{
		Host:     v.GetString("host"),
		Port:     v.GetString("port"),
		Password: v.GetString("password"),
		DbName:   v.GetString("dbName"),
	})
}

func initDb(c Config) (*redis.Client, error) {
	dbName, err := strconv.Atoi(c.DbName)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Host, c.Port),
		Password: c.Password,
		DB:       dbName,
	})

	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
