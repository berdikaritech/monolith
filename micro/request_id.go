package micro

import (
	"context"

	"github.com/koverto/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/server"
)

const requestIDMetadataKey = "request_id"

type ridContextKey struct{}

func ContextWithRequestID(ctx context.Context, rid *uuid.UUID) context.Context {
	return context.WithValue(ctx, ridContextKey{}, rid)
}

func RequestIDFromContext(ctx context.Context) (*uuid.UUID, bool) {
	id, ok := ctx.Value(ridContextKey{}).(*uuid.UUID)
	return id, ok
}

func requestIDHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if rid, ok := metadata.Get(ctx, requestIDMetadataKey); ok {
			rid, err := uuid.Parse(rid)
			if err != nil {
				return err
			}

			ctx = ContextWithRequestID(ctx, rid)
		}

		return fn(ctx, req, rsp)
	}
}

type _requestIDClientWrapper struct {
	client.Client
}

func requestIDClientWrapper(c client.Client) client.Client {
	return &_requestIDClientWrapper{c}
}

func (w *_requestIDClientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if rid, ok := RequestIDFromContext(ctx); ok {
		ctx = metadata.Set(ctx, requestIDMetadataKey, rid.Uuid.String())
	}

	return w.Client.Call(ctx, req, rsp, opts...)
}
