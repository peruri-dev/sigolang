package redisc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sigolang/lib/cache"
	"sigolang/lib/httperror"

	"github.com/redis/go-redis/v9"
)

type CacheStack struct {
	c   *Cache
	key string
}

func (c *Cache) GetStack(key string) cache.ICacheStack {
	return &CacheStack{
		c:   c,
		key: key,
	}
}

func (c *CacheStack) Push(ctx context.Context, items ...string) error {
	rdb := GetRDB(c.c)
	if rdb == nil {
		return httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	interfaceItems := make([]interface{}, len(items))

	// 2. Assign each string to the interface slice
	for i, v := range items {
		interfaceItems[i] = v
	}

	rdb.LPush(ctx, c.key, interfaceItems...)

	return nil
}

func (c *CacheStack) Pop(ctx context.Context) (*string, error) {
	rdb := GetRDB(c.c)
	if rdb == nil {
		return nil, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	r, err := rdb.LPop(ctx, c.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis pop key %s error: %w", c.key, err)
	}

	return &r, err
}

func (c *CacheStack) RPop(ctx context.Context) (*string, error) {
	rdb := GetRDB(c.c)
	if rdb == nil {
		return nil, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	r, err := rdb.RPop(ctx, c.key).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("redis pop key %s error: %w", c.key, err)
	}

	return &r, err
}

func (c *CacheStack) Peek(ctx context.Context, index int) (string, error) {
	rdb := GetRDB(c.c)
	if rdb == nil {
		return "", httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	return rdb.LIndex(ctx, c.key, int64(index)).Result()
}
func (c *CacheStack) Len(ctx context.Context) (int, error) {
	rdb := GetRDB(c.c)
	if rdb == nil {
		return 0, httperror.GenericError("Redis client not available", http.StatusInternalServerError)
	}

	i, err := rdb.LLen(ctx, c.key).Result()
	return int(i), err
}
