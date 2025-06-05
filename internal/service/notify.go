package service

import (
	"context"
	"fmt"

	"sigolang/config"
)

func (svc *Services) SendNotification(ctx context.Context) error {
	c := config.Get()
	baseUrl := fmt.Sprintf("http://%s:%d", c.Host, c.Port)

	_, err := svc.Resty.R().SetContext(ctx).
		EnableTrace().
		Put(baseUrl + "/api/users/notify")

	if err != nil {
		return fmt.Errorf("fail to send notification: %w", err)
	}

	return nil
}
