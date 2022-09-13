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
	"github.com/gomodule/redigo/redis"
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
	useTLS := flag.Bool("tls", false, "use TLS connection")
	user := flag.String("user", "default", "redis user")
	password := flag.String("password", "", "redis password")

	useSentinel := flag.Bool("sentinel", false, "use sentinel to retrieve master")

	loop := flag.Bool("loop", false, "loop?")

	flag.Parse()

	rh := rejson.NewReJSONHandler()

	var goredisClient *goredis.Client

	if *useSentinel {
		opts := &goredis.FailoverOptions{MasterName: "default", SentinelAddrs: strings.Split(*addr, ","), Username: *user, Password: *password, SentinelUsername: *user, SentinelPassword: *password}
		if *useTLS {
			opts.TLSConfig = &tls.Config{
				InsecureSkipVerify: true, // sentinels returns IP instead of hostname
			}
		}
		goredisClient = goredis.NewFailoverClient(opts)
	} else {
		opts := &goredis.Options{Addr: *addr, Username: *user, Password: *password}
		if *useTLS {
			opts.TLSConfig = &tls.Config{}
		}
		goredisClient = goredis.NewClient(opts)
	}

	rh.SetGoRedisClient(goredisClient)

	do := func() {
		ret, err := rh.JSONSet("input:5e986128435ad3fafbeb88dc", "$", &Input{
			ResourceID: "5e98612a9a72a30010ec03f4", // could be a secondary index
			EntityID:   "5e9572de2b4aae0010433600", // could be a secondary index,

			Description: "BLABLABLA",
			Speed: InputSpeed{
				Ts:    time.Now().Unix(),
				Value: 25,
			},
			Course: InputCourse{
				Ts:    time.Now().Unix(),
				Value: 200,
			},
			Sensor: []int{1001, 1002, 1003},
		})

		fmt.Println("WRITE ->", ret, err)

		if err != nil {
			log.Println(err)
		}

		inputJSON, err := redis.Bytes(rh.JSONGet("input:5e986128435ad3fafbeb88dc", "$"))

		if err != nil {
			log.Println(err)
		}

		readInput := []Input{}
		err = json.Unmarshal(inputJSON, &readInput)
		if err != nil {
			log.Println(err)
		}

		if len(readInput) > 0 {
			fmt.Printf("1 -> %+v\n", readInput[0])
		}

		inputDescription, err := redis.Bytes(rh.JSONGet("input:5e986128435ad3fafbeb88dc", ".description"))
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("2 -> %s\n", string(inputDescription))

		inputCourseTS, err := redis.Bytes(rh.JSONGet("input:5e986128435ad3fafbeb88dc", ".course.ts"))
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("3 -> %s\n", string(inputCourseTS))
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
