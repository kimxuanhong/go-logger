package logger

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
)

type contextKey string

const requestIDKey contextKey = "requestId"

func InjectRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

func GetRequestID(ctx context.Context) string {
	v := ctx.Value(requestIDKey)
	if v == nil {
		return "unknown"
	}
	return fmt.Sprint(v)
}

func WithContext(ctx context.Context) *logrus.Entry {
	requestID := GetRequestID(ctx)
	return Log.WithFields(logrus.Fields{
		"requestId": requestID,
	})
}
