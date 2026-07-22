package service

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/medhansh-32/api-gateway/internal/config"
)

type HealthService struct {
	cfg       *config.ConfigManager
	healthMap map[string]bool
	mu        sync.RWMutex
}

func NewHealthService(cfg *config.ConfigManager) *HealthService {
	return &HealthService{cfg: cfg, healthMap: make(map[string]bool)}
}

func (h *HealthService) CheckTargetsHealth() {

	h.mu.Lock()
	defer h.mu.Unlock()

	services := h.cfg.Get().Services

	for _, service := range services {

		log.Println(service)
		healthConfig := service.HealthCheck

		if !healthConfig.Enabled {
			log.Println("Health-check Skiped for Service : ",service.ID)
			continue
		}

		for _, target := range service.Targets {

			log.Println("Checking : ", target)

			client := http.Client{}

			url := "https://" + target.URL + healthConfig.Path

			request, err := http.NewRequest("GET", url, nil)

			ctx, cancel := context.WithTimeout(request.Context(), 20*time.Second)
			defer cancel()

			request = request.WithContext(ctx)

			if err != nil {
				log.Println("Unable to Make request to : ", url)
				h.healthMap[target.URL] = false
				continue
			}

			res, err := client.Do(request)

			if err != nil {
				log.Println(err.Error())
				h.healthMap[target.URL] = false
				continue
			}

			if res.StatusCode != http.StatusOK {
				log.Println(target, "Unhealthy Status : ", res.StatusCode)
				h.healthMap[target.URL] = false
				continue
			}
			
			log.Println(url, " Responded with :", res.Status)

			h.healthMap[target.URL] = true
		}

	}
}

func (h *HealthService) GetServiceHealth(url string) bool {

	h.mu.RLock()
	defer h.mu.RUnlock()

	health, present := h.healthMap[url]

	if present {
		return health
	}
	log.Println("Service is Unhealthy : ", url)
	return false

}
