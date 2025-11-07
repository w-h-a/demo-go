package middleware

import (
	"context"

	"github.com/w-h-a/demo-go/api/user"
)

type UserKey struct{}

func GetUserFromCtx(ctx context.Context) (user.User, bool) {
	user, ok := ctx.Value(UserKey{}).(user.User)
	return user, ok
}
