# 10-read

Sample Golang code writing/reading JSON to/from redis.

Use [go-rejson](https://github.com/nitishm/go-rejson) library wich can use following clients:

- [GoRedis](https://github.com/go-redis/redis)
- [Redigo](https://github.com/gomodule/redigo)

## Usage

    10-read -redigo

    10-read -goredis

    10-read -address redis.example.com:6379 -tls -user test -password test -redigo

    10-read -address redis.example.com:6379 -tls -user test -password test -goredis
