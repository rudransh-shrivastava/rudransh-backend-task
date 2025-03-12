package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")
		if idToken == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Properly extract Bearer token
		parts := strings.Split(idToken, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization format, expected 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		// Verify the token with Firebase
		token, err := s.authClient.VerifyIDToken(context.Background(), tokenStr)
		if err != nil {
			s.logger.Errorf("Token verification failed: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// You can attach token claims or user info to the context if needed
		ctx := context.WithValue(r.Context(), "userID", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Server) newRateLimitMiddleware() mux.MiddlewareFunc {
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  1,
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)
	middleware := stdlib.NewMiddleware(rateLimiter)

	return middleware.Handler
}
