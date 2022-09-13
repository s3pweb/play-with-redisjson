# 20-read-with-secondary-index-redigo

Sample Golang code writing/searching JSON to/from redis.

Use [redisearch-go](https://github.com/RediSearch/redisearch-go) library which use [Redigo](https://github.com/gomodule/redigo) client.

## Important

Prefer  [GoRedis](https://github.com/go-redis/redis) client which handle better failover scenarii.

## Usage

    20-read-with-secondary-index

    20-read-with-secondary-index -address redis.example.com:6379 -tls -user test -password test
