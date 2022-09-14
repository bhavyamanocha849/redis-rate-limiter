package limiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/bhavyamanocha849/redis-rate-limiter/datatypes"
	"github.com/bhavyamanocha849/redis-rate-limiter/services"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type limiter struct {
	redisCacheService services.CacheService
	mongoDbService    services.MongoDbService
	redisClient       *redis.Client
}

const (
	sortedSetMax = "+inf"
	sortedSetMin = "-inf"
)

func (svc *limiter) Run(ctx context.Context, response *datatypes.Request) (*datatypes.Response, error) {

	//find in the cache the value having the key
	//if not in cache make a call to mongoDB
	currTime := time.Now().UTC()
	window := time.Duration.Minutes(5)
	fmt.Println(currTime, window)
	key := "cdcd"
	expiresAt := time.Now().Add(time.Duration(window))
	minimumValid := time.Now().Add(-time.Duration(window))
	limit := uint64(123)
	//remove all the values from the set having score less than t;

	result, err := svc.redisClient.ZCount(key, strconv.FormatInt(minimumValid.UnixMilli(), 10), sortedSetMax).Result()
	if err == nil && uint64(result) >= limit {
		return &datatypes.Response{
			IsValid:   false,
			Count:     0,
			ExpiresAt: time.Duration(expiresAt.Second()),
		}, nil
	}

	item := uuid.New()

	p := svc.redisClient.Pipeline()

	_, err = p.ZRemRangeByScore(key, "0", strconv.FormatInt(minimumValid.UnixMilli(), 10)).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to remove items from key %v", key)
	}

	_, err = p.ZAdd(key, redis.Z{
		Score:  float64(time.Now().UnixMilli()),
		Member: item.String(),
	}).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add item to key %v", key)
	}

	// count how many non-expired requests we have on the sorted set
	count := p.ZCount(key, sortedSetMin, sortedSetMax)

	_, err = p.Exec()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute sorted set pipeline for key: %v", key)
	}

	totalRequests, err := count.Result()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to count items for key %v", key)
	}

	requests := uint64(totalRequests)

	if requests > limit {
		return &datatypes.Response{
			IsValid:   false,
			Count:     requests,
			ExpiresAt: time.Duration(expiresAt.Second()),
		}, nil
	}

	return &datatypes.Response{
		IsValid:   false,
		Count:     requests,
		ExpiresAt: time.Duration(expiresAt.Second()),
	}, nil
}

func NewLimiter(redisCacheService services.CacheService, mongoDbService services.MongoDbService, redisClient *redis.Client) Limiter {
	return &limiter{
		redisCacheService: redisCacheService,
		mongoDbService:    mongoDbService,
		redisClient:       redisClient,
	}
}
