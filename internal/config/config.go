package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/medhansh-32/api-gateway/internal/models"
)

type Config struct {
	ServerPort    int    `env:"SERVER_PORT" env-required:"true"`
	RoutingConfig string `env:"ROUTING_CONFIG"`

	DBHost     string `env:"DB_HOST" env-required:"true"`
	DBPort     int    `env:"DB_PORT" env-required:"true"`
	DBUser     string `env:"DB_USER" env-required:"true"`
	DBPassword string `env:"DB_PASSWORD" env-required:"true"`
	DBName     string `env:"DB_NAME" env-required:"true"`
	DBSSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
	JWTSecret  string `env:"JWT_SECRET" env-required:"true"`
}

func LoadApplicationConfig() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadGateWayConfig(path string) (*models.GatewayConfig, error){
	var gatewayConfig models.GatewayConfig
	 err := cleanenv.ReadConfig(path,&gatewayConfig);
	if err!=nil{
		return nil,err
	}

	return &gatewayConfig,nil
}