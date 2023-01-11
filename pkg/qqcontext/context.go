package qqcontext

import (
	"context"
)

const UserIdKey string = "userId"
const DefaultUserIdValue string = ""

func WithUserIdValue(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, UserIdKey, value)
}

func GetUserIdValue(ctx context.Context) string {
	value, ok := ctx.Value(UserIdKey).(string)
	if !ok {
		return ""
	}
	return value
}
