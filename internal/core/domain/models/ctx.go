package models

type ctxKey struct{}

var (
	requestIDKey = &ctxKey{}
	userIDKey    = &ctxKey{}
	roleKey      = &ctxKey{}
	txKey        = &ctxKey{}
)

func GetRequestIDKey() *ctxKey { return requestIDKey }
func GetUserIDKey() *ctxKey    { return userIDKey }
func GetRoleKey() *ctxKey      { return roleKey }
func GetTxKey() *ctxKey        { return txKey }
