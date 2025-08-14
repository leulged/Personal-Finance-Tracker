package services

import (
    "sync"
    "time"
)

// RateLimiter handles rate limiting for various operations
type RateLimiter struct {
    requests map[string][]time.Time
    mutex    sync.RWMutex
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
    return &RateLimiter{
        requests: make(map[string][]time.Time),
    }
}

// IsAllowed checks if a request is allowed based on rate limiting rules
func (r *RateLimiter) IsAllowed(key string, maxRequests int, window time.Duration) bool {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    now := time.Now()
    windowStart := now.Add(-window)

    // Get existing requests for this key
    requests, exists := r.requests[key]
    if !exists {
        requests = []time.Time{}
    }

    // Filter out old requests outside the window
    var validRequests []time.Time
    for _, reqTime := range requests {
        if reqTime.After(windowStart) {
            validRequests = append(validRequests, reqTime)
        }
    }

    // Check if we're under the limit
    if len(validRequests) < maxRequests {
        // Add current request
        validRequests = append(validRequests, now)
        r.requests[key] = validRequests
        return true
    }

    // Update requests list (remove old ones)
    r.requests[key] = validRequests
    return false
}

// GetRemainingTime returns how long until the next request is allowed
func (r *RateLimiter) GetRemainingTime(key string, window time.Duration) time.Duration {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    requests, exists := r.requests[key]
    if !exists {
        return 0
    }

    if len(requests) == 0 {
        return 0
    }

    // Find the oldest request
    oldest := requests[0]
    for _, req := range requests {
        if req.Before(oldest) {
            oldest = req
        }
    }

    windowEnd := oldest.Add(window)
    now := time.Now()

    if windowEnd.After(now) {
        return windowEnd.Sub(now)
    }

    return 0
}

// Cleanup removes old entries to prevent memory leaks
func (r *RateLimiter) Cleanup() {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    now := time.Now()
    cutoff := now.Add(-24 * time.Hour) // Keep only last 24 hours

    for key, requests := range r.requests {
        var validRequests []time.Time
        for _, reqTime := range requests {
            if reqTime.After(cutoff) {
                validRequests = append(validRequests, reqTime)
            }
        }

        if len(validRequests) == 0 {
            delete(r.requests, key)
        } else {
            r.requests[key] = validRequests
        }
    }
}

// StartCleanup starts a background cleanup routine
func (r *RateLimiter) StartCleanup() {
    go func() {
        ticker := time.NewTicker(1 * time.Hour)
        defer ticker.Stop()

        for range ticker.C {
            r.Cleanup()
        }
    }()
} 
