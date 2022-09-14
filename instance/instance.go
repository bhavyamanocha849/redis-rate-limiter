package instance

import (
	"context"
	"fmt"

	"sync"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type instance struct {
	mongodb *mongo.Client
	rdb     *redis.Client
}

const uri = "mongodb+srv://redis-rate-limiter:mongodb@cluster0.8ww6q.mongodb.net/?retryWrites=true&w=majority"

var singleton = &instance{}

var once sync.Once

func Init() {
	once.Do(func() {
		fmt.Println("Check it out")
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err != nil {
			panic(err)
		}
		defer func() {
			if err = client.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()

		if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		fmt.Println("Successfully connected and pinged.")
		singleton.mongodb = client

		rdb := redis.NewClient(&redis.Options{
			Addr:     "localhost:6397",
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		fmt.Println("Succesfully connected to rdb", rdb)
		singleton.rdb = rdb
	})
}

func Destroy() {
}
func MongoDBClient() *mongo.Client {
	return singleton.mongodb
}

func RedisClient() *redis.Client {
	return singleton.rdb
}
