package main

import (
	"net/http"

	"go-ratelimit/ratelimit"

	"github.com/gin-gonic/gin"
)

var limiter = ratelimit.NewIPRateLimiter(50, 1)

func rateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := limiter.GetLimiter(c.RemoteIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": http.StatusText(http.StatusTooManyRequests),
			})
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()
	r.Use(rateLimitMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
