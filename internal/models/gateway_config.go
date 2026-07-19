package models

import (
	"fmt"
	"time"
)

// Duration wraps time.Duration so it can unmarshal YAML strings like "15s", "200ms", "1m".
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", s, err)
	}
	d.Duration = parsed
	return nil
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return d.Duration.String(), nil
}

// Config is the root of the gateway configuration.
type GatewayConfig struct {
	Gateway       GatewayInfo           `yaml:"gateway"`
	Server        ServerConfig          `yaml:"server"`
	Services      []ServiceConfig       `yaml:"services"`
	Routes        []RouteConfig         `yaml:"routes"`
	Middlewares   map[string]Middleware `yaml:"middlewares"`
	CORS          CORSConfig            `yaml:"cors"`
	Headers       HeadersConfig         `yaml:"headers"`
	RateLimit	  RateLimit	            `yaml:"rateLimit"`
	JWT           JWTConfig             `yaml:"jwt"`
	Observability ObservabilityConfig   `yaml:"observability"`
	Compression   ToggleConfig          `yaml:"compression"`
	RequestID     ToggleConfig          `yaml:"requestID"`
	AccessLog     ToggleConfig          `yaml:"accessLog"`
}

// ---------------------------------------------------------------------------
// Top-level info & server
// ---------------------------------------------------------------------------

type GatewayInfo struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type ServerConfig struct {
	Host         string   `yaml:"host"`
	Port         int      `yaml:"port"`
	ReadTimeout  Duration `yaml:"readTimeout"`
	WriteTimeout Duration `yaml:"writeTimeout"`
	IdleTimeout  Duration `yaml:"idleTimeout"`
}

// ---------------------------------------------------------------------------
// Services
// ---------------------------------------------------------------------------

type ServiceConfig struct {
	ID             string               `yaml:"id"`
	Targets        []TargetConfig       `yaml:"targets"`
	LoadBalancer   LoadBalancerConfig   `yaml:"loadBalancer"`
	Timeout        ServiceTimeout       `yaml:"timeout"`
	Retries        RetryConfig          `yaml:"retries"`
	HealthCheck    HealthCheckConfig    `yaml:"healthCheck"`
	CircuitBreaker CircuitBreakerConfig `yaml:"circuitBreaker"`
}

type TargetConfig struct {
	URL    string `yaml:"url"`
	Weight int    `yaml:"weight"`
}

type LoadBalancerConfig struct {
	Strategy string `yaml:"strategy"` // e.g. "weighted-round-robin", "round-robin"
}

type ServiceTimeout struct {
	Connect  Duration `yaml:"connect"`
	Response Duration `yaml:"response"`
}

type RetryConfig struct {
	Attempts int      `yaml:"attempts"`
	Backoff  Duration `yaml:"backoff"`
}

type HealthCheckConfig struct {
	Enabled  bool     `yaml:"enabled"`
	Path     string   `yaml:"path"`
	Interval Duration `yaml:"interval"`
	Timeout  Duration `yaml:"timeout"`
}

type CircuitBreakerConfig struct {
	Enabled          bool     `yaml:"enabled"`
	FailureThreshold int      `yaml:"failureThreshold"`
	ResetTimeout     Duration `yaml:"resetTimeout"`
}

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

type RouteConfig struct {
	ID           string        `yaml:"id"`
	Service      string        `yaml:"service"`
	Path         string        `yaml:"path"`
	Methods      []string      `yaml:"methods"`
	StripPrefix  bool          `yaml:"stripPrefix"`
	Auth         AuthConfig    `yaml:"auth"`
	Middleware   []string      `yaml:"middleware"`
	RateLimit    bool    `yaml:"rateLimit,omitempty"`
	Cache        *CacheConfig  `yaml:"cache,omitempty"`
}

type AuthConfig struct {
	Enabled bool     `yaml:"enabled"`
	Type    string   `yaml:"type"` // e.g. "jwt"
	Roles   []string `yaml:"roles,omitempty"`
}

type RateLimit struct {
	Requests int      `yaml:"requests"`
	Window   Duration `yaml:"window"`
}

type CacheConfig struct {
	Enabled bool     `yaml:"enabled"`
	TTL     Duration `yaml:"ttl"`
}

// ---------------------------------------------------------------------------
// Middleware (reusable, named)
// ---------------------------------------------------------------------------

type Middleware struct {
	Enabled bool `yaml:"enabled"`
}

// ---------------------------------------------------------------------------
// CORS / Headers / JWT
// ---------------------------------------------------------------------------

type CORSConfig struct {
	Enabled      bool     `yaml:"enabled"`
	AllowOrigins []string `yaml:"allowOrigins"`
	AllowMethods []string `yaml:"allowMethods"`
	AllowHeaders []string `yaml:"allowHeaders"`
}

type HeadersConfig struct {
	Request  HeaderOps `yaml:"request"`
	Response HeaderOps `yaml:"response"`
}

type HeaderOps struct {
	Add    map[string]string `yaml:"add,omitempty"`
	Remove []string          `yaml:"remove,omitempty"`
}

type JWTConfig struct {
	Issuer   string `yaml:"issuer"`
	Audience string `yaml:"audience"`
	Secret   string `yaml:"secret"`
}

// ---------------------------------------------------------------------------
// Observability / toggles
// ---------------------------------------------------------------------------

type ObservabilityConfig struct {
	Metrics MetricsConfig `yaml:"metrics"`
	Tracing ToggleConfig  `yaml:"tracing"`
}

type MetricsConfig struct {
	Enabled  bool   `yaml:"enabled"`
	Endpoint string `yaml:"endpoint"`
}

// ToggleConfig covers simple {enabled: bool} blocks (compression, requestID, accessLog, tracing).
type ToggleConfig struct {
	Enabled bool `yaml:"enabled"`
}