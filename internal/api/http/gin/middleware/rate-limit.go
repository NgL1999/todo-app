package middleware

import (
	"errors"
	"net/http"
	"todo-app/pkg/client"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

// RateLimiterMiddleware checks if the rate limit is exceeded
func RateLimiter(rateLimiter *limiter.Limiter) func(c *gin.Context) {
	return func(c *gin.Context) {
		ipClient := c.ClientIP()
		limiterCtx, err := rateLimiter.Get(c, ipClient)
		if err != nil {
			c.JSON(http.StatusInternalServerError, client.ErrInternal(errors.New("rate limiter failed")))
			return
		}
		if limiterCtx.Reached {
			c.JSON(http.StatusBadRequest, client.ErrInvalidRequest(errors.New("too many requests")))
			c.Abort()
			return
		}
		c.Next()
	}
}
