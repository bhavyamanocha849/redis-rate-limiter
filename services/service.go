package services

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

type CacheService interface {
	Ping() error
	GetKey(ctx context.Context, key string) (*string, error)
	GetKeys(ctx context.Context, keys []string) ([]interface{}, error)
	SetKey(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	EvictKey(ctx context.Context, key string) error
	AddKeys(ctx context.Context) error
	RemoveKeys(ctx context.Context) error
}

type redisCache struct {
	client *redis.Client
}

func (s *redisCache) Ping() error {
	// errorName := "CACHE_PING"
	_, err := s.client.Ping().Result()
	if err != nil {
		// logger.Log.WithError(err).Error(errorName)
		return err
	}
	return nil
}

func (s *redisCache) SetKey(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// errorName := "Cache_Service_SetKey"
	// add prefix
	key = fmt.Sprintf("%v_%v", "ABC", key)
	val, err := s.client.Set(key, value, ttl).Result()
	if err != nil {
		// logger.Log.WithError(err).Error(errorName)
		return err
	}

	fmt.Println("val", val)

	return nil
}

func (s *redisCache) GetKey(ctx context.Context, key string) (*string, error) {
	// errorName := "CACHE_SERVICE_GetKey"
	key = fmt.Sprintf("%v_%v", "", key)
	data, err := s.client.Get(key).Result()
	if err != nil {
		// logger.Log.WithError(err).Error(errorName)
		return nil, err
	}
	return &data, nil
}

func (s *redisCache) GetKeys(ctx context.Context, keys []string) ([]interface{}, error) {
	// errorName := "CACHE_GetKeys"

	// add prefix to keys
	for i := range keys {
		keys[i] = fmt.Sprintf("%v_%v", "CacheKeyPrefix", keys[i])
	}
	res := s.client.MGet(keys...)

	if res.Err() != nil {
		// logger.Log.WithError(res.Err()).Error(errorName)
		return nil, res.Err()
	}
	result, err := res.Result()
	if err != nil {
		// logger.Log.WithError(err).Error(errorName)
	}
	return result, nil
}
func (s *redisCache) EvictKey(ctx context.Context, key string) error {
	key = fmt.Sprintf("%v_%v", "constants.CacheKeyPrefix", key)
	// logger.Log.Info("Evicting ", key)
	redisReponse := s.client.Del(key)
	if redisReponse.Err() != nil {

		return redisReponse.Err()
	}
	res, err := redisReponse.Result()
	if err != nil {
		// logger.Log.WithError(err)
	}
	fmt.Printf("redisReponse: %+v\n", res)
	return nil
}

func (s *redisCache) AddKeys(ctx context.Context) error {

	return nil
}

func (s *redisCache) RemoveKeys(ctx context.Context) error {

	return nil
}

func NewCacheService(rdb *redis.Client) CacheService {
	return &redisCache{client: rdb}
}
