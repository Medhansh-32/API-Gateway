package schedular

import (
	"github.com/medhansh-32/api-gateway/internal/service"
	"github.com/robfig/cron/v3"
)

type HealthCheckSchedular struct{
	healthService *service.HealthService
}

func NewHealthCheckSchedular(healthService *service.HealthService) (*HealthCheckSchedular){
	return &HealthCheckSchedular{healthService: healthService}
}

func (H HealthCheckSchedular) InitHealthCheckSchedular() {
	c := cron.New()

	c.AddFunc("@every 10s", func() {
		H.healthService.CheckTargetsHealth()
	})

	c.Start()

}