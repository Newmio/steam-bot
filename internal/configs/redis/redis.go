package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type Config struct {
	Host     string
	Port     string
	Password string
	DbName   string
}

func OpenDb() (*redis.ClusterClient, error) {
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

func initDb(c Config) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{fmt.Sprintf("%s:%s", c.Host, c.Port)},
		Password: c.Password,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
