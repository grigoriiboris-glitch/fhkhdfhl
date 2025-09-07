package auth

import (
	"sync"
	"time"
)

// RateLimiter implements a simple in-memory rate limiter for authentication attempts
type RateLimiter struct {
	attempts map[string]*attemptRecord
	mutex    sync.RWMutex
	
	// Configuration
	maxAttempts int
	window      time.Duration
	blockTime   time.Duration
}

type attemptRecord struct {
	count     int
	firstTry  time.Time
	blockedAt time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxAttempts int, window, blockTime time.Duration) *RateLimiter {
	rl := &RateLimiter{
		attempts:    make(map[string]*attemptRecord),
		maxAttempts: maxAttempts,
		window:      window,
		blockTime:   blockTime,
	}
	
	// Start cleanup goroutine
	go rl.cleanup()
	
	return rl
}

// IsAllowed checks if the identifier is allowed to make an attempt
func (rl *RateLimiter) IsAllowed(identifier string) bool {
	rl.mutex.RLock()
	record, exists := rl.attempts[identifier]
	rl.mutex.RUnlock()
	
	now := time.Now()
	
	if !exists {
		return true
	}
	
	// Check if still blocked
	if !record.blockedAt.IsZero() && now.Sub(record.blockedAt) < rl.blockTime {
		return false
	}
	
	// Check if window has expired
	if now.Sub(record.firstTry) > rl.window {
		return true
	}
	
	// Check if under limit
	return record.count < rl.maxAttempts
}

// RecordAttempt records a failed attempt
func (rl *RateLimiter) RecordAttempt(identifier string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	record, exists := rl.attempts[identifier]
	
	if !exists {
		rl.attempts[identifier] = &attemptRecord{
			count:    1,
			firstTry: now,
		}
		return
	}
	
	// Reset if window expired
	if now.Sub(record.firstTry) > rl.window {
		record.count = 1
		record.firstTry = now
		record.blockedAt = time.Time{}
		return
	}
	
	record.count++
	
	// Block if exceeded limit
	if record.count >= rl.maxAttempts {
		record.blockedAt = now
	}
}

// Reset clears attempts for an identifier (e.g., on successful login)
func (rl *RateLimiter) Reset(identifier string) {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	delete(rl.attempts, identifier)
}

// cleanup removes expired records periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		
		for identifier, record := range rl.attempts {
			// Remove if window expired and not blocked, or block time expired
			if (now.Sub(record.firstTry) > rl.window && record.blockedAt.IsZero()) ||
				(!record.blockedAt.IsZero() && now.Sub(record.blockedAt) > rl.blockTime) {
				delete(rl.attempts, identifier)
			}
		}
		
		rl.mutex.Unlock()
	}
}