package auth

import (
	"context"
	"net/http"

	"github.com/w-h-a/demo-go/api/user"
	httphandler "github.com/w-h-a/demo-go/internal/handler/http"
	"github.com/w-h-a/demo-go/internal/middleware"
)

type authMiddleware struct {
	handler http.Handler
}

func (m *authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := httphandler.ReqToCtx(r)

	var authenticatedUser user.User
	// var authErr error

	// TODO: fill in authenticatedUser & handle authErr

	ctxWithUser := context.WithValue(ctx, middleware.UserKey{}, authenticatedUser)
	rWithUser := r.WithContext(ctxWithUser)

	m.handler.ServeHTTP(w, rWithUser)
}

// TODO: pass in dependencies (auth service)
func New() func(h http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return &authMiddleware{
			handler: handler,
		}
	}
}
