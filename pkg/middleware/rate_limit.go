package middleware

import (
	"fmt"
	"math"
	"strings"
	"time"

	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_err"
	"GoRestify/pkg/rate_limiter"
	"GoRestify/pkg/response"

	"github.com/gin-gonic/gin"
)

var loginLimiter = rate_limiter.NewProgressive(
	pkg_consts.LoginRateLimit,
	time.Duration(pkg_consts.LoginLockoutStage1)*time.Second,
	time.Duration(pkg_consts.LoginLockoutStage2)*time.Second,
	time.Duration(pkg_consts.LoginLockoutStage3)*time.Second,
)

func loginKey(c *gin.Context) string {
	ip := c.ClientIP()
	if ip == "" {
		return "login:unknown"
	}
	return "login:" + ip
}

// LoginRateLimit blocks login while the IP is in a progressive lockout.
func LoginRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := loginKey(c)
		if locked, retryAfter := loginLimiter.Locked(key); locked {
			abortRateLimit(c, retryAfter)
			return
		}
		c.Next()
	}
}

// RecordLoginFailure counts a failed login toward progressive lockout.
func RecordLoginFailure(c *gin.Context) {
	loginLimiter.RecordFailure(loginKey(c))
}

// ResetLoginLimit clears progressive lockout after a successful login.
func ResetLoginLimit(c *gin.Context) {
	loginLimiter.Reset(loginKey(c))
}

// UnblockLoginIP clears lockout state for an IP (admin action).
func UnblockLoginIP(ip string) {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return
	}
	loginLimiter.Reset("login:" + ip)
}

// ListLoginBlocks returns current login rate-limit entries for admin.
func ListLoginBlocks() []rate_limiter.BlockInfo {
	items := loginLimiter.Snapshot()
	for i := range items {
		items[i].Key = strings.TrimPrefix(items[i].Key, "login:")
	}
	return items
}

// RateLimit is a simple fixed-window limiter (non-progressive).
func RateLimit(limiter *rate_limiter.Limiter, keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if keyFn != nil {
			key = keyFn(c)
		}
		if key == "" {
			key = "unknown"
		}

		if !limiter.Allow(key) {
			abortRateLimit(c, time.Duration(pkg_consts.LoginLockoutStage1)*time.Second)
			return
		}

		c.Next()
	}
}

func abortRateLimit(c *gin.Context, retryAfter time.Duration) {
	seconds := int(math.Ceil(retryAfter.Seconds()))
	if seconds < 1 {
		seconds = 1
	}

	err := pkg_err.New(pkg_err.YouVeExceededTheRateLimitForRequestsPleaseTryAgainLater, "E1194301").
		Custom(pkg_err.RateLimitErr).
		Message(pkg_err.YouVeExceededTheRateLimitForRequestsPleaseTryAgainLater).
		Build()

	c.Header("Retry-After", fmt.Sprintf("%d", seconds))
	response.New(c).Error(err).Abort().JSON()
}
