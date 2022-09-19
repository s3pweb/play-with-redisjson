package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	goredis "github.com/go-redis/redis/v8"
	"github.com/s3pweb/play-with-redisjson/utils"

	"github.com/nitishm/go-rejson/v4"
)

type InputSpeed struct {
	Ts    int64 `json:"ts"`
	Value int   `json:"value"`
}

type InputCourse struct {
	Ts    int64 `json:"ts"`
	Value int   `json:"value"`
}

type Input struct {
	ID          string      `json:"id"`
	ResourceID  string      `json:"resourceId"`
	EntityID    string      `json:"entityId"`
	Description string      `json:"description"`
	Speed       InputSpeed  `json:"speed"`
	Course      InputCourse `json:"course"`
	Sensor      []int       `json:"sensor"`
}

func main() {
	addr := flag.String("address", "localhost:6379", "redis server address")
	useSentinel := flag.Bool("sentinel", false, "use sentinel to retrieve master")
	useTLS := flag.Bool("tls", false, "use TLS connection")
	user := flag.String("user", "default", "redis user")
	password := flag.String("password", "", "redis password")
	loop := flag.Bool("loop", false, "loop?")

	flag.Parse()

	var goredisClient *goredis.Client

	if *useSentinel {
		opts := &goredis.FailoverOptions{MasterName: "default", SentinelAddrs: strings.Split(*addr, ","), Username: *user, Password: *password, SentinelUsername: *user, SentinelPassword: *password}
		if *useTLS {
			opts.TLSConfig = &tls.Config{}
		}
		goredisClient = goredis.NewFailoverClient(opts)
	} else {
		opts := &goredis.Options{Addr: *addr, Username: *user, Password: *password}
		if *useTLS {
			opts.TLSConfig = &tls.Config{}
		}
		goredisClient = goredis.NewClient(opts)
	}

	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClient(goredisClient)

	do := func() {
		res, err := goredisClient.Do(context.Background(), "FT.DROP", "idx:input").Result()

		log.Println("DROP INDEX", res, err)

		// c := redisearch.NewClientFromPool(pool, "idx:input")

		res, err = goredisClient.Do(context.Background(),
			utils.ToSliceOfAny(strings.Split("FT.CREATE idx:input ON JSON PREFIX 1 input SCHEMA $.resourceId AS resourceId TAG SEPARATOR ; $.entityId AS entityId TAG SEPARATOR ; $.speed.value AS speed NUMERIC $.course.value AS course NUMERIC", " "))...,
		).Result()
		log.Println("CREATE INDEX", res, err)

		ret, err := rh.JSONSet("input:5e98612a9a72a30010ec0001", "$", &Input{
			ID:         "5e98612a9a72a30010ec0001",
			ResourceID: "5e98612a9a72a30010ec03f5", // could be a secondary index
			EntityID:   "5e9572de2b4aae0010433600", // could be a secondary index,

			Speed: InputSpeed{
				Ts:    time.Now().Unix(),
				Value: 50,
			},
			Course: InputCourse{
				Ts:    time.Now().Unix(),
				Value: 100,
			},
		})

		if err != nil {
			log.Println(err)
		}

		fmt.Println("WRITE ->", ret, err)

		ret, err = rh.JSONSet("input:5e98612a9a72a30010ec0003", "$", &Input{
			ID:         "5e98612a9a72a30010ec0003",
			ResourceID: "5e98612a9a72a30010ec03ba", // could be a secondary index
			EntityID:   "5e9572de2b4aae0010433600", // could be a secondary index,

			Speed: InputSpeed{
				Ts:    time.Now().Unix(),
				Value: 75,
			},
			Course: InputCourse{
				Ts:    time.Now().Unix(),
				Value: 20,
			},
		})

		if err != nil {
			log.Println(err)
		}

		fmt.Println("WRITE ->", ret, err)

		ret, err = rh.JSONSet("input:5e98612a9a72a30010ec0004", "$", &Input{
			ID:         "5e98612a9a72a30010ec0004",
			ResourceID: "5e98612a9a72a30010ec03ba", // could be a secondary index
			EntityID:   "5e9572de2b4aae0010433600", // could be a secondary index,

			Speed: InputSpeed{
				Ts:    time.Now().Unix(),
				Value: 18,
			},
			Course: InputCourse{
				Ts:    time.Now().Unix(),
				Value: 45,
			},
		})

		if err != nil {
			log.Println(err)
		}

		fmt.Println("WRITE ->", ret, err)

		res, err = goredisClient.Do(context.Background(), "FT.SEARCH", "idx:input", "@resourceId:{5e98612a9a72a30010ec03f5}").Result()
		if pretty, err2 := json.MarshalIndent(res, "", "  "); err2 == nil {
			fmt.Println("READ ONE ->", string(pretty), err)
		} else {
			fmt.Println("READ ONE ->", res, err)
		}

		res, err = goredisClient.Do(context.Background(), "FT.SEARCH", "idx:input", "@resourceId:{5e98612a9a72a30010ec03ba | 5e98612a9a72a30010ec03f5}").Result()
		if pretty, err2 := json.MarshalIndent(res, "", "  "); err2 == nil {
			fmt.Println("READ MULTIPLE ->", string(pretty), err)
		} else {
			fmt.Println("READ MULTIPLE ->", res, err)
		}

		// res, err = goredisClient.Do(context.Background(), "FT.SEARCH", "idx:input", "@resourceId:{5e98612a9a72a30010ec03f4} @speed:[0 70]").Result()
		res, err = goredisClient.Do(context.Background(), "FT.SEARCH", "idx:input", "@resourceId:{5e98612a9a72a30010ec03f4}").Result()
		if pretty, err2 := json.MarshalIndent(res, "", "  "); err2 == nil {
			fmt.Println("READ BETWEEN VALUES 1 ->", string(pretty), err)
		} else {
			fmt.Println("READ BETWEEN VALUES 1 ->", res, err)
		}

		res, err = goredisClient.Do(context.Background(), "FT.SEARCH", "idx:input", "@speed:[0 70] @course:[0 60]").Result()
		if pretty, err2 := json.MarshalIndent(res, "", "  "); err2 == nil {
			fmt.Println("READ BETWEEN VALUES 2 ->", string(pretty), err)
		} else {
			fmt.Println("READ BETWEEN VALUES 2 ->", res, err)
		}
	}

	if *loop {
		for {

			fmt.Println(goredisClient.Info(context.Background(), "replication").Result())

			do()
			time.Sleep(500 * time.Millisecond)
		}
	}

	do()
}
