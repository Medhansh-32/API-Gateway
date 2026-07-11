package config

import (
	"sync"

	"github.com/medhansh-32/api-gateway/internal/models"
)

type ConfigManager struct {
	gateWayConfig *models.GatewayConfig
	mu sync.RWMutex
}


func (m *ConfigManager) Get() *models.GatewayConfig {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.gateWayConfig
}

func (m *ConfigManager) Update(gateWayConfig *models.GatewayConfig) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.gateWayConfig = gateWayConfig
}