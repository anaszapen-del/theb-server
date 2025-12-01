package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count      int
	lastSeen   time.Time
	resetTimer *time.Timer
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

// RateLimit implements rate limiting middleware
func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]

		if !exists {
			// New visitor
			v = &visitor{
				count:    1,
				lastSeen: time.Now(),
			}

			// Set timer to clean up after 1 minute
			v.resetTimer = time.AfterFunc(time.Minute, func() {
				mu.Lock()
				delete(visitors, ip)
				mu.Unlock()
			})

			visitors[ip] = v
			mu.Unlock()
			c.Next()
			return
		}

		// Check if we need to reset the counter
		if time.Since(v.lastSeen) > time.Minute {
			v.count = 1
			v.lastSeen = time.Now()

			// Reset the cleanup timer
			if v.resetTimer != nil {
				v.resetTimer.Stop()
			}
			v.resetTimer = time.AfterFunc(time.Minute, func() {
				mu.Lock()
				delete(visitors, ip)
				mu.Unlock()
			})

			mu.Unlock()
			c.Next()
			return
		}

		// Increment request count
		v.count++
		v.lastSeen = time.Now()

		if v.count > requestsPerMinute {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		mu.Unlock()
		c.Next()
	}
}
