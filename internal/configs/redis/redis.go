package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func OpenDb() (*redis.ClusterClient, error) {
	v := viper.New()
	v.AddConfigPath("internal/configs/redis")
	v.SetConfigName("config")

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return initDb(v.GetString("password"), v.GetStringSlice("hosts"))
}

func initDb(pass string, hosts []string) (*redis.ClusterClient, error) {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    hosts,
		Password: pass,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
