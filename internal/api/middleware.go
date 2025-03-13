package api

import (
	"context"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/rudransh-shrivastava/rudransh-backend-task/internal/schema"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"gorm.io/gorm"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")
		if idToken == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Properly parse Bearer token
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

// RBACMiddleware returns a middleware that only allows users with one of the allowed roles.
func RBACMiddleware(db *gorm.DB, allowedRoles ...schema.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uid, ok := r.Context().Value("userID").(string)
			if !ok || uid == "" {
				http.Error(w, "Unauthorized: missing user ID", http.StatusUnauthorized)
				return
			}

			var user schema.User
			if err := db.Where("uid = ?", uid).First(&user).Error; err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			// Check if the user's role is allowed.
			allowed := slices.Contains(allowedRoles, user.Role)
			if !allowed {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// The corsMiddleware adds the necessary headers to enable CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// mockAuthMiddleware is a mock auth middleware that sets a fixed user ID in the context
// This is used for testing purposes only.
func mockAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "userID", "mock-user-id")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
