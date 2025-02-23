package repository

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var client *redis.Client

func InitRedis() {

	client = NewRedisclient()

}
func NewRedisclient() *redis.Client {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	fmt.Println("✅ Connected to Redis")
	return client
}
func Getclient() *redis.Client {
	fmt.Print("client:", client)
	return client
}
func GetPath(path string) (string, error) {

	val, err := client.Get(path).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}
