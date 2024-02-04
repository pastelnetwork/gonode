package debug

import (
	"sync"
	"time"
)

// RateLimiter is a simple rate limiter that limits the number of requests per hour
type RateLimiter struct {
	mutex       sync.Mutex
	requests    map[string][]time.Time
	maxRequests int
	interval    time.Duration
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxRequests int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		requests:    make(map[string][]time.Time),
		maxRequests: maxRequests,
		interval:    interval,
	}
}

// CheckRateLimit checks if the rate limit has been reached for the given pid
func (rl *RateLimiter) CheckRateLimit(pid string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()

	// Clean up old requests that are out of the interval
	if reqs, found := rl.requests[pid]; found {
		var updatedReqs []time.Time
		for _, t := range reqs {
			if now.Sub(t) < rl.interval {
				updatedReqs = append(updatedReqs, t)
			}
		}
		rl.requests[pid] = updatedReqs
	}

	// Check if the rate limit has been reached
	if len(rl.requests[pid]) >= rl.maxRequests {
		return false
	}

	// Add the new request
	rl.requests[pid] = append(rl.requests[pid], now)
	return true
}
