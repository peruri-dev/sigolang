package redisc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"sigolang/config"
	"sigolang/lib/cache"
	"sigolang/lib/httperror"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Impl interface{}
}

func RegisterRedisCache() {
	cache.RegisterCache(&cache.CacheFactory{
		Prefixes: []string{"redis://"},
		Create: func(c *config.CacheConfig) (cache cache.ICache, err error) {
			opts, err := redis.ParseURL(c.CacheUri)
			if err != nil {
				return nil, err
			}

			if c.CachePoolSize != 0 {
				opts.PoolSize = c.CachePoolSize
			}

			if c.CachePoolTimeout != 0 {
				opts.PoolTimeout = c.CachePoolTimeout
			}

			ctx := context.Background()

			rdb := redis.NewClient(opts)
			err = rdb.Ping(ctx).Err()
			if err != nil {
				return nil, fmt.Errorf("redis ping error: %w", err)
			}

			// Enable tracing instrumentation.
			if err := redisotel.InstrumentTracing(rdb); err != nil {
				slog.Error("unable to tracing redis", slog.Any("error", err))
			}

			// Enable metrics instrumentation.
			if err := redisotel.InstrumentMetrics(rdb); err != nil {
				slog.Error("unable to metric redis", slog.Any("error", err))
			}

			res, err := rdb.Do(ctx, "INFO").Result()
			if err != nil {
				panic(err)
			}

			version := ""
			info := strings.SplitSeq(res.(string), "\n")
			for line := range info {
				if strings.HasPrefix(line, "redis_version:") {
					version = strings.TrimSpace(strings.Split(line, ":")[1])
				}
			}

			slog.Info(fmt.Sprintf("redis connected URI:%s Version:%s", c.CacheUri, version), slog.String("version", version))

			redisCache := &Cache{}

			redisCache.Impl = rdb

			return redisCache, nil
		},
	})
}

///

func GetRDB(cache *Cache) *redis.Client {
	if rdb, ok := cache.Impl.(*redis.Client); ok {
		return rdb
	}
	return nil
}

func (c *Cache) Write(ctx context.Context, key string, value string, expiration time.Duration) error {
	rdb := GetRDB(c)
	if rdb == nil {
		return httperror.GenericError("Redis client not available", http.StatusInternalServerError)

	}
	return rdb.Set(ctx, key, value, expiration).Err()
}

func (c *Cache) Lock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return false, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}
	return rdb.SetNX(ctx, key, value, expiration).Result()
}

func (c *Cache) UnLock(ctx context.Context, key string) (int64, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return 0, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}
	return rdb.Del(ctx, key).Result()
}

func (c *Cache) Read(ctx context.Context, key string) (string, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return "", httperror.GenericError("Redis client not available", http.StatusInternalServerError)

	}
	r, err := rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("redis get key %s missing: %w", key, err)
		}
		return "", fmt.Errorf("redis get key %s error: %w", key, err)
	}

	return r, nil
}

func (c *Cache) GetHash(ctx context.Context, key string) (map[string]string, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return nil, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	fields, err := rdb.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get hash: %w", err)
	}

	return fields, nil
}

func (c *Cache) GetHashField(ctx context.Context, key string, field string) (string, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return "", httperror.GenericError("Redis client not available", http.StatusInternalServerError)

	}

	data, err := rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get hash: %w", err)
	}

	return data, nil
}

func (c *Cache) SetHash(ctx context.Context, key string, fields map[string]string, expiration time.Duration) error {
	rdb := GetRDB(c)
	if rdb == nil {
		return httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	// TODO: this should be deprecated
	// if err := rdb.HMSet(ctx, key, fields).Err(); err != nil {
	// 	return fmt.Errorf("failed to set hash: %w", err)
	// }

	if err := rdb.HSet(ctx, key, fields).Err(); err != nil {
		return fmt.Errorf("failed to set hash: %w", err)
	}

	if err := rdb.Expire(ctx, key, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set expiration: %w", err)
	}

	return nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	rdb := GetRDB(c)
	if rdb == nil {
		return httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	if err := rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key: %w", err)
	}

	return nil
}

type FieldValue struct {
	Field string
	Value string
	TTL   time.Duration
}

// SetHashFieldsTTL perform HSET then HEXPIRE
// NOTE: require redis version 7.4+
func (c *Cache) SetHashFieldsTTL(ctx context.Context, hashKey string, fieldValues []cache.FieldValue) error {
	rdb := GetRDB(c)
	if rdb == nil {
		return httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	if len(fieldValues) == 0 {
		return nil
	}

	sets := map[string]string{}
	for _, v := range fieldValues {
		sets[v.Field] = v.Value
	}

	if err := rdb.HSet(ctx, hashKey, sets).Err(); err != nil {
		return fmt.Errorf("failed to set hash fields: %w", err)
	}

	pipe := rdb.Pipeline()

	for _, fv := range fieldValues {
		pipe.HExpire(ctx, hashKey, fv.TTL, fv.Field)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	return nil
}

func (c *Cache) Incr(ctx context.Context, key string) (int64, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return 0, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	return rdb.Incr(ctx, key).Result()
}

func (c *Cache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	rdb := GetRDB(c)
	if rdb == nil {
		return httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	return rdb.Expire(ctx, key, expiration).Err()
}

func (c *Cache) TTL(ctx context.Context, key string) (time.Duration, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return 0, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	return rdb.PTTL(ctx, key).Result()
}

func (c *Cache) ReadWithTTL(ctx context.Context, key string) (string, time.Duration, error) {
	rdb := GetRDB(c)
	if rdb == nil {
		return "", 0, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	pipe := rdb.Pipeline()
	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", 0, err
	}

	value, err := getCmd.Result()
	if err != nil {
		return "", 0, err
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return "", 0, err
	}

	return value, ttl, nil
}
