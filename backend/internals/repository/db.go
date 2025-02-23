package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConfig struct {
	URI      string
	Database string
}

var Client *mongo.Client

func InitMongo(uri string, mongodb string) {
	config := MongoDBConfig{
		URI:      uri,
		Database: mongodb,
	}
	_, database, err := InitializeMongoDB(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(database)
	if err != nil {
		log.Print(err)
	}
	fmt.Println("✅ Connected to mongo")
}

func InitializeMongoDB(config MongoDBConfig) (*mongo.Client, *mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(config.URI)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, nil, err
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}
	database := client.Database(config.Database)
	Client = client
	return client, database, nil
}

func GetRateLimiterCollection() *mongo.Collection {
	return Client.Database("rate_limiter").Collection("rate_limiter")
}

func GetDocuments(endpoint string) (*bson.M, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	filter := bson.M{"endpoint": bson.M{"$eq": endpoint}}

	defer cancel()

	collection := GetRateLimiterCollection()
	var result bson.M
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("❌ No document found")
		} else {
			fmt.Println("❌ Error finding document:", err)
		}
		return nil, err
	}

	fmt.Println("✅ Found document:")
	return &result, nil

}
func TrackIPForPath(path, ip string) (int, error) {
	key := fmt.Sprintf("ip_path_count:%s:%s", ip, path)

	count, err := client.Incr(key).Result()
	if err != nil {
		return 0, err
	}

	if count == 1 {
		client.Expire(key, time.Hour)
	}

	log.Printf("📢 IP: %s -> Path: %s -> Requests in last hour: %d", ip, path, count)
	return int(count), nil

}
