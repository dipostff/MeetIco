package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimitMiddleware struct {
	clients sync.Map
}

type RateLimitClient struct {
	timestamps []time.Time
	mu         sync.Mutex
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	rlm := &RateLimitMiddleware{
		clients: sync.Map{},
	}
	go rlm.cleanup()
	return rlm
}

func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			ip = xff
		}

		now := time.Now()
		client, _ := m.clients.LoadOrStore(ip, &RateLimitClient{
			timestamps: []time.Time{},
		})
		rlClient := client.(*RateLimitClient)

		rlClient.mu.Lock()
		defer rlClient.mu.Unlock()

		var validTimestamps []time.Time
		for _, ts := range rlClient.timestamps {
			if now.Sub(ts) < time.Minute {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		rlClient.timestamps = validTimestamps

		if len(rlClient.timestamps) >= 10 {
			http.Error(w, `{"error": "too many requests"}`, http.StatusTooManyRequests)
			return
		}

		rlClient.timestamps = append(rlClient.timestamps, now)
		next.ServeHTTP(w, r)
	})
}

func (m *RateLimitMiddleware) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.clients.Range(func(key, value interface{}) bool {
			client := value.(*RateLimitClient)
			client.mu.Lock()
			now := time.Now()
			var validTimestamps []time.Time
			for _, ts := range client.timestamps {
				if now.Sub(ts) < time.Minute {
					validTimestamps = append(validTimestamps, ts)
				}
			}
			client.timestamps = validTimestamps
			if len(client.timestamps) == 0 {
				m.clients.Delete(key)
			}
			client.mu.Unlock()
			return true
		})
	}
}
