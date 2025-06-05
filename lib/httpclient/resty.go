package httpclient

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/dubonzi/otelresty"
	"github.com/go-resty/resty/v2"
	"github.com/peruri-dev/inalog"
	"github.com/tidwall/pretty"
)

func InitRestyClient() *resty.Client {
	client := resty.New()

	client.SetTimeout(30 * time.Second)

	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		reqBody, _ := json.Marshal(req.Body)

		slog.DebugContext(req.Context(), fmt.Sprintf("----> Request: Method=%s, URL=%s", req.Method, req.URL),
			slog.Group("restyRequest",
				slog.String("method", req.Method),
				slog.String("uRL", req.URL),
				inalog.HttpHeaderToSlog(req.Header),
				slog.String("body", string(pretty.Ugly(reqBody))),
			))

		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		slog.DebugContext(resp.Request.Context(), fmt.Sprintf("<---- Response: StatusCode=%d, Status=%s", resp.StatusCode(), resp.Status()),
			slog.Group("restyResponse",
				slog.Int("statusCode", resp.StatusCode()),
				slog.String("status", resp.Status()),
				inalog.HttpHeaderToSlog(resp.Header()),
				slog.String("body", string(pretty.Ugly(resp.Body()))),
			),
		)

		return nil

	})

	opts := []otelresty.Option{otelresty.WithTracerName("sigolang-resty")}
	otelresty.TraceClient(client, opts...)

	return client
}
