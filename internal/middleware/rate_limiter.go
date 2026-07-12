package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/medhansh-32/api-gateway/internal/config"
	"github.com/medhansh-32/api-gateway/internal/models"
)

type RateLimitMiddleware struct {
	cfg *config.ConfigManager
	rlt *RateLimiter
}

func NewRateLimitingMiddleware(cfg *config.ConfigManager) *RateLimitMiddleware {
	return &RateLimitMiddleware{cfg: cfg, rlt: &RateLimiter{BucketHolder: make(map[string]*Bucket)}}
}

type RateLimiter struct {
	BucketHolder map[string]*Bucket
	mu           sync.RWMutex
}

type Bucket struct {
	LastRefillTime time.Time
	Tokens         int
	mu             sync.Mutex
}

func (R RateLimitMiddleware) RateLimitCheck(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rateLimitEnabled := checkRateLimitingEnabledForURL(r.URL, R.cfg)
		clientIP := r.RemoteAddr

		if !rateLimitEnabled {
			log.Println("Rate Limiting Not Enabled for this Route")
			next.ServeHTTP(w, r)
			return
		}

		allowed := R.rlt.rateLimit(clientIP, &(R.cfg.Get().RateLimit))

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Too Many Requests",
			})
			return
		}

		next.ServeHTTP(w, r)
	})

}

func checkRateLimitingEnabledForURL(url *url.URL, cfgManager *config.ConfigManager) bool {
	cfg := cfgManager.Get()
	routes := cfg.Routes

	for _, route := range routes {

		if matchURL(url, route.Path) {
			log.Println("Route Matched : ", route, " Rate Limiting enabled : ", route.RateLimit)

			return route.RateLimit
		}

	}
	log.Print("No Route Found For : ", url)
	return false
}

func (R *RateLimiter) rateLimit(clientIP string, rateLimit *models.RateLimit) bool {

	bucket := R.getOrCreateBucket(clientIP, rateLimit)
	b, _ := json.MarshalIndent(bucket, "", "  ")
	log.Println(string(b))
	error := R.checkBucket(bucket, clientIP, rateLimit)

	if error != nil {
		log.Print(error.Error())
		return false
	}

	return true

}

func (R *RateLimiter) checkBucket(bucket *Bucket, clientIP string, rateLimit *models.RateLimit) error {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	if time.Since(bucket.LastRefillTime) >= rateLimit.Window.Duration {
		bucket.Tokens = rateLimit.Requests
		bucket.LastRefillTime = time.Now()

		log.Println("Bucket refilled for:", clientIP)
	}

	if bucket.Tokens == 0 {
		return errors.New("Too many requests")
	}

	bucket.Tokens--

	log.Println("Bucket token deducted for:", clientIP)

	return nil
}

func (R *RateLimiter) getOrCreateBucket(clientIP string, rateLimit *models.RateLimit) *Bucket {

	log.Println("rate Limit : ", rateLimit)

	R.mu.RLock()
	bucket, exists := R.BucketHolder[clientIP]
	R.mu.RUnlock()

	if exists {
		return bucket
	}

	R.mu.Lock()
	defer R.mu.Unlock()

	if bucket, exists = R.BucketHolder[clientIP]; exists {
		return bucket
	}

	bucket = &Bucket{Tokens: rateLimit.Requests, LastRefillTime: time.Now()}
	R.BucketHolder[clientIP] = bucket
	return bucket

}
