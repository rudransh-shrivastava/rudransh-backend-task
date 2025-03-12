package middleware

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func NewRateLimitMiddleware() mux.MiddlewareFunc {
	rate := limiter.Rate{
		Period: 1 * time.Second,
		Limit:  1,
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)
	middleware := stdlib.NewMiddleware(rateLimiter)

	return middleware.Handler
}
