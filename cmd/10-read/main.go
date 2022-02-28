package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
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

	useRedigo := flag.Bool("redigo", false, "use Redigo client")
	useGoRedis := flag.Bool("goredis", false, "use GoRedis client")

	flag.Parse()

	if !*useRedigo && !*useGoRedis {
		fmt.Println("Please tell what client to use.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rh := rejson.NewReJSONHandler()

	// Redigo Client
	if *useRedigo {
		conn, err := redis.Dial("tcp", *addr, redis.DialUseTLS(*useTLS), redis.DialUsername(*user), redis.DialPassword(*password))

		if err != nil {
			log.Fatalln(err)
		}

		rh.SetRedigoClient(conn)
	}

	// GoRedis Client
	if *useGoRedis {
		opts := &goredis.Options{Addr: *addr, Username: *user, Password: *password}
		if *useTLS {
			opts.TLSConfig = &tls.Config{}
		}
		cli := goredis.NewClient(opts)

		rh.SetGoRedisClient(cli)
	}

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
		log.Fatalln(err)
	}

	inputJSON, err := redis.Bytes(rh.JSONGet("input:5e986128435ad3fafbeb88dc", "$"))

	if err != nil {
		log.Fatalln(err)
	}

	readInput := []Input{}
	err = json.Unmarshal(inputJSON, &readInput)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("1 -> %+v\n", readInput[0])

	inputDescription, err := redis.Bytes(rh.JSONGet("input:5e986128435ad3fafbeb88dc", ".description"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("2 -> %s\n", string(inputDescription))

	inputCourseTS, err := redis.Bytes(rh.JSONGet("input:5e986128435ad3fafbeb88dc", ".course.ts"))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("3 -> %s\n", string(inputCourseTS))
}
