package models

type ctxKey struct{}

var (
	requestIDKey = &ctxKey{}
	userIDKey    = &ctxKey{}
	roleKey      = &ctxKey{}
)

func GetRequestIDKey() *ctxKey { return requestIDKey }
func GetUserIDKey() *ctxKey    { return userIDKey }
func GetRoleKey() *ctxKey      { return roleKey }
