package seeds

import (
	"context"

	"sigolang/internal/model"

	"github.com/uptrace/bun"
)

func ResetSchema(ctx context.Context, db *bun.DB) error {
	if err := db.ResetModel(ctx, (*model.User)(nil)); err != nil {
		return err
	}

	users := []model.User{
		{
			Name:   "admin",
			Emails: []string{"admin1@admin", "admin2@admin"},
		},
		{
			Name:   "root",
			Emails: []string{"root1@root", "root2@root"},
		},
	}
	if _, err := db.NewInsert().Model(&users).Exec(ctx); err != nil {
		return err
	}

	return nil
}
