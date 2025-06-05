package service

import (
	"context"

	"sigolang/internal/model"
)

func (svc *Services) UserList(ctx context.Context) ([]model.User, error) {
	users := make([]model.User, 0)
	if err := svc.DB.NewSelect().Model(&users).OrderExpr("id ASC").Scan(ctx); err != nil {
		return nil, err
	}
	return users, nil
}

func (svc *Services) UserCreate(ctx context.Context, name string, emails []string) (*model.User, error) {
	var user = model.User{}
	user.Name = name
	user.Emails = emails

	if _, err := svc.DB.NewInsert().Model(&user).Exec(ctx); err != nil {
		return nil, err
	}
	return &user, nil
}
