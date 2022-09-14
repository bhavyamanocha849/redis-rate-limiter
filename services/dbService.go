package services

import (
	"context"
	"fmt"

	"github.com/bhavyamanocha849/redis-rate-limiter/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoDbService struct {
	mongoDb    *mongo.Client
	collection *mongo.Collection
}

type MongoDbService interface {
	FetchFromDB(key string) error
	AddToDB(ctx context.Context, inputData models.ApiConfig) (string, int64, error)
	UpdateDB(ctx context.Context, keyId string, value int64) error
}

func (svc *mongoDbService) FetchFromDB(key string) error {
	return nil
}

func (svc *mongoDbService) AddToDB(ctx context.Context, insertData models.ApiConfig) (string, int64, error) {
	_, err := svc.collection.InsertOne(ctx, insertData)
	if err != nil {
		fmt.Println("err", err)
		return "", 0, err
	}
	return "", 0, nil
}

func (svc *mongoDbService) UpdateDB(ctx context.Context, keyId string, value int64) error {
	//write through cache :(?
	return nil
}

func NewDbService(mongoDb *mongo.Client) MongoDbService {
	return &mongoDbService{
		mongoDb: mongoDb,
	}
}
