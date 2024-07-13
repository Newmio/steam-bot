package main

import (
	"bot/internal/app"

	_ "github.com/mattn/go-sqlite3"
)

//redis-cli --cluster add-node <новый-узел:порт> <существующий-узел:порт>
//redis-cli --cluster reshard <существующий-узел:порт>

// chromedriver --port=9515

/*
docker rm -f $(docker ps -a -q)
docker rmi -f $(docker images -a -q)
docker network rm $(docker network ls -q)
docker volume rm $(docker volume ls -q)
*/

/*
redis-cli --cluster create \
    localhost:7001 \
    localhost:7002 \
    localhost:7003 \
    localhost:7004 \
    localhost:7005 \
    localhost:7006 \
    --cluster-replicas 1 --cluster-yes -a yourpassword
*/

/*
redis-cli --cluster create \
    89.28.236.131:6379 \
    89.28.236.131:6380 \
    89.28.236.131:6381 \
    89.28.236.131:6382 \
    89.28.236.131:6383 \
    89.28.236.131:6384 \
    --cluster-replicas 1 -a yourpassword
*/

func main() {
	app.Init()
}
