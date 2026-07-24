package config

import (
	"sync"

	"github.com/medhansh-32/api-gateway/internal/models"
)

type ConfigManager struct {
	gateWayConfig *models.GatewayConfig
    applicationConfig *Config
	mu sync.RWMutex
}


func (m *ConfigManager) GetGateWayConfig() *models.GatewayConfig {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.gateWayConfig
}

func (m *ConfigManager) UpdateGateWayConfig(gateWayConfig *models.GatewayConfig) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.gateWayConfig = gateWayConfig
}

func (m *ConfigManager) GetApplicationConfig() (*Config) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.applicationConfig
}

func (m *ConfigManager) UpdateApplicationConfig(applicationConfig *Config) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.applicationConfig = applicationConfig
}