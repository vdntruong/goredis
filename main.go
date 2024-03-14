package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func ExampleClient() *redis.Client {
	opts := &redis.Options{
		Addr:     "localhost:6379",
		Password: "mypassword",
		DB:       0, // use default DB
	}
	return redis.NewClient(opts)
}

func ExampleClientURL() *redis.Client {
	url := "redis://<username>:<password>@localhost:6379/0?protocol=3"
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	return redis.NewClient(opts)
}

func main() {
	rdb := ExampleClient()

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
