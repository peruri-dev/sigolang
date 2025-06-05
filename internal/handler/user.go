package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/peruri-dev/inalog"
	"github.com/peruri-dev/inatrace"
	"github.com/danielgtaylor/huma/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/rand"
)

type UserResponseBody struct {
	Name   string   `json:"name"`
	Emails []string `json:"emails"`
}

type UserCreateBody struct {
	Name   string   `json:"name" example:"john" required:"true"`
	Emails []string `json:"emails" example:"a@a.com,b@b.com" required:"true"`
}

type CreateUserInput struct {
	Body UserCreateBody `required:"true"`
}

type ListUserOutput struct {
	Body   []UserResponseBody
	Status int
}

type CreateUserOutput struct {
	Body   UserResponseBody
	Status int
}

func (h *Handler) GetUsers(ctx context.Context, input *struct{}) (*ListUserOutput, error) {
	users, err := h.svc.UserList(ctx)

	if err != nil {
		inalog.LogWith(inalog.WithCfg{Ctx: ctx}).Error("Failed to list story", slog.Any("error", err))
		return nil, huma.Error400BadRequest("fail")
	}

	var data []UserResponseBody
	for _, user := range users {
		var r UserResponseBody
		r.Name = user.Name
		r.Emails = user.Emails
		data = append(data, r)
	}

	status := http.StatusOK
	if len(data) == 0 {
		status = http.StatusNoContent
	}

	return &ListUserOutput{
		Body:   data,
		Status: status,
	}, nil
}

func (h *Handler) CreateUser(ctx context.Context, input *CreateUserInput) (*CreateUserOutput, error) {
	user, err := h.svc.UserCreate(ctx, input.Body.Name, input.Body.Emails)
	if err != nil {
		inalog.LogWith(inalog.WithCfg{Ctx: ctx}).Error("Failed to create user", slog.Any("error", err))
		return nil, huma.Error400BadRequest("fail")
	}
	inalog.LogWith(inalog.WithCfg{Ctx: ctx}).Info("New user created", slog.Int("id", int(user.ID)))

	var r UserResponseBody
	r.Name = user.Name
	r.Emails = user.Emails

	_ = h.svc.SendNotification(ctx)

	return &CreateUserOutput{
		Body:   r,
		Status: http.StatusCreated,
	}, nil
}

func (h *Handler) NotifyUser(ctx context.Context, input *struct{}) (*struct {
	Body []byte
}, error) {
	_, span := inatrace.Start(ctx, "sendNotification", trace.WithAttributes(attribute.String("id", "id")))
	defer span.End()

	n := rand.Intn(3) // n will be between 0 and 2
	inalog.LogWith(inalog.WithCfg{Ctx: ctx}).Info("Sleeping", slog.Int("durataion", n))
	time.Sleep(time.Duration(n) * time.Second)

	n = rand.Intn(3) // n will be between 0 and 2
	if n == 0 {
		fmt.Println("div", 1/n)
	}

	return &struct {
		Body []byte
	}{
		Body: []byte("."),
	}, nil
}

func (h *Handler) RegisterUser(api huma.API) {
	huma.Register(api,
		huma.Operation{
			OperationID: "get-users",
			Method:      http.MethodGet,
			Path:        "/api/users",
			Summary:     "Get a bunch of users",
		}, h.GetUsers,
	)

	huma.Register(api,
		huma.Operation{
			OperationID: "create-user",
			Method:      http.MethodPost,
			Path:        "/api/users",
			Summary:     "Create a user",
		}, h.CreateUser,
	)

	huma.Register(api,
		huma.Operation{
			OperationID: "notify-user",
			Method:      http.MethodPut,
			Path:        "/api/users/notify",
			Summary:     "Notify a user",
		}, h.NotifyUser,
	)
}
