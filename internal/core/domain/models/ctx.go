package models

type contextKey string

const (
	txKey     contextKey = "tx"
	requestID contextKey = "request_id"
	userID    contextKey = "user_id"
	roleKey   contextKey = "role"
)

func GetTxKey() contextKey {
	return txKey
}

// GetRequestIDKey возвращает ключ для request_id в контексте
func GetRequestIDKey() contextKey {
	return requestID
}

// GetUserIDKey возвращает ключ для user_id в контексте
func GetUserIDKey() contextKey {
	return userID
}

// GetRoleKey возвращает ключ для role в контексте
func GetRoleKey() contextKey {
	return roleKey
}
