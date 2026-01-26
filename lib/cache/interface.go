package cache

import (
	"context"
	"time"
)

type ICache interface {
	Write(ctx context.Context, key string, value string, expiration time.Duration) error
	Lock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	UnLock(ctx context.Context, key string) (int64, error)
	Read(ctx context.Context, key string) (string, error)
	GetHash(ctx context.Context, key string) (map[string]string, error)
	GetHashField(ctx context.Context, key string, field string) (string, error)
	SetHash(ctx context.Context, key string, fields map[string]string, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	SetHashFieldsTTL(ctx context.Context, hashKey string, fieldValues []FieldValue) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
}
