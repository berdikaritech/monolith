package micro

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type request interface {
	Endpoint() string
	Service() string
}

const (
	durationLogKey     = "duration_ms"
	grpcEndpointLogKey = "grpc.endpoint"
	grpcErrorLogKey    = "grpc.error"
	grpcServiceLogKey  = "grpc.service"
	grpcSpanLogKey     = "grpc.span"
	requestIdLogKey    = "request_id"
	startTimeLogKey    = "start_time"
)

const (
	grpcServiceClient = "client"
	grpcServiceServer = "server"
)

func logHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		event, now := logWrapper(ctx, req)
		event.Str(grpcSpanLogKey, grpcServiceServer)

		err := fn(ctx, req, rsp)
		event.AnErr(grpcErrorLogKey, err)
		event.Dur(durationLogKey, time.Since(now)).Send()

		return err
	}
}

type _logClientWrapper struct {
	client.Client
}

func logClientWrapper(c client.Client) client.Client {
	return &_logClientWrapper{c}
}

func (w *_logClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	event, now := logWrapper(ctx, req)
	event.Str(grpcSpanLogKey, grpcServiceClient)

	err := w.Client.Call(ctx, req, rsp, opts...)
	event.AnErr(grpcErrorLogKey, err)
	event.Dur(durationLogKey, time.Since(now)).Send()

	return err
}

func logWrapper(ctx context.Context, req request) (*zerolog.Event, time.Time) {
	now := time.Now()
	event := log.Info().Time(startTimeLogKey, now)
	event.Str(grpcEndpointLogKey, req.Endpoint())
	event.Str(grpcServiceLogKey, req.Service())

	if rid, ok := RequestIDFromContext(ctx); ok {
		event.Str(requestIdLogKey, rid.Uuid.String())
	}

	return event, now
}
