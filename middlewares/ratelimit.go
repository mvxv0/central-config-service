package middlewares

import (
	"net"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

var (
	visitors = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		// Stavljen je na 1 zahtev u sekundi da bi ga ja lakse testirao
		limiter = rate.NewLimiter(1, 1)
		visitors[ip] = limiter
	}

	return limiter
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded: previse zahteva", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
