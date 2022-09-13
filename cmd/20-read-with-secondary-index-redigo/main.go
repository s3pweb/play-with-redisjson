package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/RediSearch/redisearch-go/redisearch"
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

	flag.Parse()

	var opts []redis.DialOption

	if *useTLS {
		opts = append(opts, redis.DialUseTLS(true))
	}

	if *user != "" {
		opts = append(opts, redis.DialUsername(*user))
	}

	if *password != "" {
		opts = append(opts, redis.DialPassword(*password))
	}

	conn, err := redis.Dial("tcp", *addr, opts...)

	if err != nil {
		fmt.Println("Unable to connect to redis.")
		flag.PrintDefaults()
		log.Fatalln(err)
	}

	// conn.Close()
	_, err = conn.Do("PING")

	if err != nil {
		fmt.Println("Unable to ping redis.")
		flag.PrintDefaults()
		log.Fatalln(err)
	}

	pool := &redis.Pool{Dial: func() (redis.Conn, error) {
		return redis.Dial("tcp", *addr, opts...)
	}}

	c := redisearch.NewClientFromPool(pool, "idx:input")

	// Create a schema
	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTagFieldOptions("$.resourceId", redisearch.TagFieldOptions{Separator: ';', As: "resourceId"})).
		AddField(redisearch.NewTagFieldOptions("$.entityId", redisearch.TagFieldOptions{Separator: ';', As: "entityId"})).
		AddField(redisearch.NewNumericFieldOptions("$.speed.value", redisearch.NumericFieldOptions{As: "speed"})).
		AddField(redisearch.NewNumericFieldOptions("$.course.value", redisearch.NumericFieldOptions{As: "course"}))

	// Drop an existing index. If the index does not exist an error is returned
	err = c.Drop()

	if err != nil {
		log.Println(err)
	}

	// Create the index with the given schema and definition
	id := &redisearch.IndexDefinition{IndexOn: "JSON", Score: -1.0}
	if err := c.CreateIndexWithIndexDefinition(sc, id.AddPrefix("input")); err != nil {
		log.Fatalln(err)
	}

	rh := rejson.NewReJSONHandler()

	rh.SetRedigoClient(pool.Get())

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
		log.Fatalln(err)
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
		log.Fatalln(err)
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
		log.Fatalln(err)
	}

	fmt.Println("WRITE ->", ret, err)

	docs, total, err := c.Search(redisearch.NewQuery("@resourceId:{5e98612a9a72a30010ec03f5}"))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("READ ONE ->", docs, total, err)

	docs, total, err = c.Search(redisearch.NewQuery("@resourceId:{5e98612a9a72a30010ec03ba | 5e98612a9a72a30010ec03f5}"))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("READ MULTIPLE ->", docs, total, err)

	docs, total, err = c.Search(redisearch.NewQuery("@resourceId:{5e98612a9a72a30010ec03f4} @speed:[0 70]"))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("READ BETWEEN VALUES 1 ->", docs, total, err)

	docs, total, err = c.Search(redisearch.NewQuery("@speed:[0 70] @course:[0 60]"))

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("READ BETWEEN VALUES 2 ->", docs, total, err)
}
