package request

import (
	"context"
	"fmt"
	"go-test/lib/logger"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type (
	ctxKeyReqId struct{}
	Generator   func(*http.Request) (string, error)
)

const xRequestIDHeaderKey string = "X-Request-Id"

var CtxKeyReqId ctxKeyReqId = ctxKeyReqId{}

func DefaultGenerator(r *http.Request) (string, error) {
	id := r.Header.Get(xRequestIDHeaderKey)
	if id != "" {
		return id, nil
	}

	if id, ok := r.Context().Value(xRequestIDHeaderKey).(string); ok {
		if id != "" {
			return id, nil
		}
	}

	uid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return uid.String(), nil
}

func RequestIDContext(generator Generator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqId, err := generator(r)
			if err != nil {
				logger.GetLogger(r.Context()).Error("RequestId generator err: ", err)
				reqId = fmt.Sprintf("%d", time.Now().UnixNano())
			}

			ctx := context.WithValue(r.Context(), CtxKeyReqId, reqId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(CtxKeyReqId).(string); ok {
		return v
	}
	return ""
}
