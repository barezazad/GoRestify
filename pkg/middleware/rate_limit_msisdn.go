package middleware

import (
	"fmt"

	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/rate_limiter"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

// RateLimitMap define map for limiter
var RateLimitMap *rate_limiter.TTLMap

// InitRateLimiterByMSISDN is used for limit requests from a user to a route
func InitRateLimiterByMSISDN(limit, ttl uint) {
	RateLimitMap = rate_limiter.NewTTLMap(limit, ttl)
}

// RateLimiterByMSISDN is used for limit requests from a user to a route
func RateLimiterByMSISDN() gin.HandlerFunc {

	return func(c *gin.Context) {

		var phone string
		phoneInCtx, ok := c.Get("PHONE")
		if ok {
			phone = phoneInCtx.(string)
		}
		key := fmt.Sprintf("%v-%v", phone, c.Request.URL.String())

		rateCount := RateLimitMap.Get(key)
		rateCount++

		if rateCount > RateLimitMap.Limit {
			err := pkg_err.New(pkg_err.YouVeExceededTheRateLimitForRequestsPleaseTryAgainLater, "E1142747").
				Custom(pkg_err.BadRequestErr).
				Message(pkg_err.YouVeExceededTheRateLimitForRequestsPleaseTryAgainLater).Build()
			response.New(c).Error(err).Abort().JSON()
			return
		}

		RateLimitMap.Set(key, rateCount)

		c.Next()
	}
}
