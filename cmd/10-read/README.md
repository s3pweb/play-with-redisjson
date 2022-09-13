# 10-read

Sample Golang code writing/reading JSON to/from redis.

Use [go-rejson](https://github.com/nitishm/go-rejson) library wich use following client:

- [GoRedis](https://github.com/go-redis/redis)

## Usage

    10-read

    10-read -address redis.example.com:6379 -tls -user test -password test

    10-read -address redis.example.com:26379 -tls -user test -password test -sentinel

    10-read -address redis.example.com:26379 -tls -user test -password test -sentinel -loop
