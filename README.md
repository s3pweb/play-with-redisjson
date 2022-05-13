# play-with-redisjson

Some scripts to explore secondary index with redisjson and redissearch

# Lunch an redis container

```
docker run -p 6379:6379 redislabs/redismod
```
or use redis stack which comes with redis database and RedisInsight editor exposed on [local port 8001](http://localhost:8001)
```
docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
```
