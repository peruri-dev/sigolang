package service

import (
	"context"

	"sigolang/internal/model"
	"sigolang/lib/cache"

	"github.com/go-resty/resty/v2"
	"github.com/uptrace/bun"
)

type Services struct {
	DB    *bun.DB
	Cache *cache.Cache
	Resty *resty.Client
}

// AllServices is an example services interface
type AllServices interface {
	UserList(context.Context) ([]model.User, error)
	UserCreate(ctx context.Context, name string, emails []string) (*model.User, error)
	SendNotification(ctx context.Context) error
}
