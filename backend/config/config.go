package config

import (
	"context"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Ctx         = context.Background()
	RedisClient *redis.Client
	MongoDB     *mongo.Database
)

func ConnectMongoDB() {
	client, err := mongo.Connect(Ctx, options.Client().ApplyURI("mongodb+srv://poon:Aqwer.2002@works.upw8bbp.mongodb.net/"))
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}
	MongoDB = client.Database("demo")
}

func ConnectRedis() {
	// เชื่อมต่อ Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
}
