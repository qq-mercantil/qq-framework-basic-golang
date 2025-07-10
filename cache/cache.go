package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	trace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

type ICacheClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Ping(ctx context.Context) error
	HSet(ctx context.Context, key string, values map[string]interface{}, expiration time.Duration) error
	HSetField(ctx context.Context, key string, field string, value interface{}, expiration time.Duration) error
	HGet(ctx context.Context, key, field string) (string, error)
}

type CacheClient struct {
	cc *redis.Client
}

type CacheOptions struct {
    Config      ICacheProvider
    ServiceName string
}

//dd:ignore
func NewCacheClient(opts CacheOptions) (*CacheClient, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Config.GetHost(), opts.Config.GetPort()),
		Password: opts.Config.GetPassword(),
		DB:       opts.Config.GetDB(),
	})

	serviceName := opts.ServiceName
    if serviceName == "" {
        serviceName = os.Getenv("APP_NAME") + "-redis"
    }

	trace.WrapClient(client, trace.WithServiceName(serviceName))

	pong, err := client.Ping(context.Background()).Result()

	if err != nil {
		return nil, err
	}

	fmt.Println(pong)

	return &CacheClient{
		cc: client,
	}, nil
}

func (c CacheClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var data []byte
	var err error

	// Check if the value is a string
	switch v := value.(type) {
	case string:
		data = []byte(v)
	default:
		// Marshal the value to JSON
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	// Store the data in the cache
	err = c.cc.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (c CacheClient) Get(ctx context.Context, key string) (string, error) {
	keyValue, err := c.cc.Get(ctx, key).Result()

	if err != nil {
		return "", err
	}

	fmt.Println(keyValue)

	return keyValue, nil

}

func (c CacheClient) Ping(ctx context.Context) error {
	_, err := c.cc.Ping(ctx).Result()

	if err != nil {
		return err
	}

	return nil
}

func (c CacheClient) HSet(ctx context.Context, key string, values map[string]interface{}, expiration time.Duration) error {
	if len(values) == 0 {
		return nil
	}

	finalValues := make([]interface{}, 0, len(values)*2)
	// Converter os valores para strings (serialização)
	for field, value := range values {
		var serializedValue string

		switch v := value.(type) {
		case string:
			// Não precisa serializar strings
			serializedValue = v
		default:
			// Serializar outros tipos como JSON
			data, err := json.Marshal(v)
			if err != nil {
				return err
			}
			serializedValue = string(data)
		}

		finalValues = append(finalValues, field, serializedValue)
	}

	err := c.cc.HSet(ctx, key, finalValues...).Err()
	if err != nil {
		return err
	}

	// Aplicar a expiração, se necessária
	if expiration > 0 {
		err = c.cc.Expire(ctx, key, expiration).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CacheClient) HSetField(ctx context.Context, key string, field string, value interface{}, expiration time.Duration) error {
	var fieldValue string
	switch v := value.(type) {
	case string:
		fieldValue = v
	default:
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		fieldValue = string(data)
	}

	err := c.cc.HSet(ctx, key, field, fieldValue).Err()
	if err != nil {
		return err
	}

	if expiration > 0 {
		err = c.cc.Expire(ctx, key, expiration).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c CacheClient) HGet(ctx context.Context, key, field string) (string, error) {
	value, err := c.cc.HGet(ctx, key, field).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
