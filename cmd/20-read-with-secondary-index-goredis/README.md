# 20-read-with-secondary-index-goredis

Sample Golang code writing/searching JSON to/from redis.

## Important

Does not use any library for redisearch (like [redisearch-go](https://github.com/RediSearch/redisearch-go) which use [Redigo](https://github.com/gomodule/redigo) client).
Instead send raw Redis commands.
[GoRedis](https://github.com/go-redis/redis) client handle better failover scenarii.

## Usage

    20-read-with-secondary-index

    20-read-with-secondary-index -address redis.example.com:6379 -tls -user test -password test

    20-read-with-secondary-index -address redis.example.com:26379 -tls -user test -password test -sentinel

    20-read-with-secondary-index -address redis.example.com:26379 -tls -user test -password test -sentinel -loop
