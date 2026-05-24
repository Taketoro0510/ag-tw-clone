package middleware

import (
	"net"
	"net/http"
	"sync"

	"github.com/koitake1/cloudcode-sns/backend/internal/handler/dto"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func newRateLimiter(r rate.Limit, b int) *rateLimiter {
	return &rateLimiter{
		limiters: make(map[string]*rate.Limiter),
		r:        r,
		b:        b,
	}
}

func (rl *rateLimiter) getLimiter(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	limiter, exists := rl.limiters[key]
	if !exists {
		limiter = rate.NewLimiter(rl.r, rl.b)
		rl.limiters[key] = limiter
	}
	return limiter
}

func RateLimit(reqsPerMin int) echo.MiddlewareFunc {
	r := rate.Limit(float64(reqsPerMin) / 60.0)
	rl := newRateLimiter(r, reqsPerMin)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key, ok := c.Get("userID").(string)
			if !ok || key == "" {
				ip, _, _ := net.SplitHostPort(c.Request().RemoteAddr)
				key = ip
			}
			limiter := rl.getLimiter(key)
			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, dto.ErrorResponse{
					Error: dto.ErrorDetail{Code: "RATE_LIMITED", Message: "too many requests"},
				})
			}
			return next(c)
		}
	}
}
