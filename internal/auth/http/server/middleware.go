package server

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"ride-hail/pkg/logger"
	"time"
)

func (a *API) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = newRequestID()
		}

		w.Header().Set("X-Request-ID", reqID)
		ctx := logger.WithRequestID(r.Context(), reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return a.cfg.JWT.Secret, nil
		})

		if err != nil || !t.Valid {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			userID, okId := claims["user_id"].(string)
			role, okRole := claims["role"].(string)

			if !okId || !okRole {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ctx := logger.WithUserID(r.Context(), userID)
			ctx = logger.WithRole(ctx, role)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
	})
}

func newRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString(fmt.Appendf(nil, "%d", time.Now().UnixNano()))
	}
	return hex.EncodeToString(b)
}
