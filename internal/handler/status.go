package handler

import (
	"context"
	"net/http"
	"sigolang/config"

	"github.com/danielgtaylor/huma/v2"
)

type HealthResp struct {
	Body struct {
		Message string `json:"message"`
	}
}

type StatusResp struct {
	Body struct {
		Version string `json:"version"`
	}
}

func (h *Handler) Version(ctx context.Context, input *struct{}) (*StatusResp, error) {
	c := config.Get()
	s := &StatusResp{}
	s.Body.Version = c.AppVersion
	return s, nil
}

func (h *Handler) HealthZ(ctx context.Context, input *struct{}) (*HealthResp, error) {
	s := &HealthResp{}
	s.Body.Message = "OK"
	return s, nil
}

func (h *Handler) RoutesStatus(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "Version",
		Path:        "/version",
		Method:      http.MethodGet,
		Tags:        []string{"Status:Version"},
	}, h.Version)

	huma.Register(api, huma.Operation{
		OperationID: "HealthZ",
		Path:        "/healthz",
		Method:      http.MethodGet,
		Tags:        []string{"Status:HealthZ"},
	}, h.HealthZ)
}
