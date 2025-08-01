package apiserver

import (
	"async_api/store"
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func NewLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("http request", "method", r.Method, "path", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

type userCtxKey struct{}

func ContextWithUser(ctx context.Context, user *store.User) context.Context {
	return context.WithValue(ctx, userCtxKey{}, user)
}

func NewAuthMiddleware(JwtManager *JwtManager, userStore *store.UserStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth") {
				next.ServeHTTP(w, r)
				return
			}
			// authorization header
			// Should be: Authorization: Bearer <access_token>
			authHeader := r.Header.Get("Authorization")
			var token string
			if parts := strings.Split(authHeader, "Bearer "); len(parts) == 2 {
				token = parts[1]
			}
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			parsedToken, err := JwtManager.Parse(token)
			if err != nil {
				slog.Error("failed to parse token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !JwtManager.IsAccessToken(parsedToken) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not an access token"))
				return
			}

			userIDStr, err := parsedToken.Claims.GetSubject()
			if err != nil {
				slog.Error("failed to extract subject claim from token", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				slog.Error("token subject is not valid uuid", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			user, err := userStore.ByID(r.Context(), userID)
			if err != nil {
				slog.Error("failed to get user by id", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r.WithContext(ContextWithUser(r.Context(), user)))
		})
	}
}
