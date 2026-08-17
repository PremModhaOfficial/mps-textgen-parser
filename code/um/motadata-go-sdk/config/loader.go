package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Environment variable prefixes
const (
	envServicePrefix = "SERVICE_"
	envLoggerPrefix  = "LOG_"
	envMetricsPrefix = "METRICS_"
	envTracerPrefix  = "TRACER_"
	envModulePrefix  = "LOG_MODULE_"
)

// Service environment variables
const (
	EnvServiceName        = envServicePrefix + "NAME"
	EnvServiceVersion     = envServicePrefix + "VERSION"
	EnvServiceEnvironment = envServicePrefix + "ENVIRONMENT"
)

// Logger environment variables
const (
	EnvLogLevel          = envLoggerPrefix + "LEVEL"
	EnvLogConsoleEnabled = envLoggerPrefix + "CONSOLE_ENABLED"
	EnvLogConsoleFormat  = envLoggerPrefix + "CONSOLE_FORMAT"
	EnvLogOTELEnabled    = envLoggerPrefix + "OTEL_ENABLED"
	EnvLogOTELEndpoint   = envLoggerPrefix + "OTEL_ENDPOINT"
	EnvLogOTELInsecure   = envLoggerPrefix + "OTEL_INSECURE"
	EnvLogOTELProtocol   = envLoggerPrefix + "OTEL_PROTOCOL"
	EnvLogOTELDebug      = envLoggerPrefix + "OTEL_DEBUG"
	EnvLogFileEnabled    = envLoggerPrefix + "FILE_ENABLED"
	EnvLogFilePath       = envLoggerPrefix + "FILE_PATH"
	EnvLogFileMaxSizeMB  = envLoggerPrefix + "FILE_MAX_SIZE_MB"
	EnvLogFileMaxBackups = envLoggerPrefix + "FILE_MAX_BACKUPS"
	EnvLogFileMaxAgeDays = envLoggerPrefix + "FILE_MAX_AGE_DAYS"
	EnvLogFileCompress   = envLoggerPrefix + "FILE_COMPRESS"
	EnvLogAddCaller      = envLoggerPrefix + "ADD_CALLER"
	EnvLogCallerSkip     = envLoggerPrefix + "CALLER_SKIP"
)

// Metrics environment variables
const (
	EnvMetricsEnabled           = envMetricsPrefix + "ENABLED"
	EnvMetricsOTELEnabled       = envMetricsPrefix + "OTEL_ENABLED"
	EnvMetricsOTELEndpoint      = envMetricsPrefix + "OTEL_ENDPOINT"
	EnvMetricsOTELInsecure      = envMetricsPrefix + "OTEL_INSECURE"
	EnvMetricsOTELProtocol      = envMetricsPrefix + "OTEL_PROTOCOL"
	EnvMetricsPrometheusEnabled = envMetricsPrefix + "PROMETHEUS_ENABLED"
	EnvMetricsPrometheusPort    = envMetricsPrefix + "PROMETHEUS_PORT"
	EnvMetricsPrometheusPath    = envMetricsPrefix + "PROMETHEUS_PATH"
	EnvMetricsExportInterval    = envMetricsPrefix + "EXPORT_INTERVAL"
)

// Tracer environment variables
const (
	EnvTracerEnabled        = envTracerPrefix + "ENABLED"
	EnvTracerOTELEndpoint   = envTracerPrefix + "OTEL_ENDPOINT"
	EnvTracerOTELInsecure   = envTracerPrefix + "OTEL_INSECURE"
	EnvTracerOTELProtocol   = envTracerPrefix + "OTEL_PROTOCOL"
	EnvTracerOTELDebug      = envTracerPrefix + "OTEL_DEBUG"
	EnvTracerSamplingRatio  = envTracerPrefix + "SAMPLING_RATIO"
	EnvTracerBatchTimeout   = envTracerPrefix + "BATCH_TIMEOUT"
	EnvTracerMaxExportBatch = envTracerPrefix + "MAX_EXPORT_BATCH"
	EnvTracerMaxQueueSize   = envTracerPrefix + "MAX_QUEUE_SIZE"
)

// yamlConfig mirrors Config for YAML unmarshaling with proper tags.
type yamlConfig struct {
	Service struct {
		Name        string `yaml:"name"`
		Version     string `yaml:"version"`
		Environment string `yaml:"environment"`
	} `yaml:"service"`

	Logger struct {
		Level   string `yaml:"level"`
		Console struct {
			Enabled *bool  `yaml:"enabled"`
			Format  string `yaml:"format"`
		} `yaml:"console"`
		OTEL struct {
			Enabled  *bool  `yaml:"enabled"`
			Endpoint string `yaml:"endpoint"`
			Insecure *bool  `yaml:"insecure"`
		} `yaml:"otel"`
		File struct {
			Enabled    *bool  `yaml:"enabled"`
			Path       string `yaml:"path"`
			MaxSizeMB  int    `yaml:"max_size_mb"`
			MaxBackups int    `yaml:"max_backups"`
			MaxAgeDays int    `yaml:"max_age_days"`
			Compress   *bool  `yaml:"compress"`
		} `yaml:"file"`
		Caller struct {
			Enabled *bool `yaml:"enabled"`
			Skip    int   `yaml:"skip"`
		} `yaml:"caller"`
		Modules map[string]string `yaml:"modules"`
	} `yaml:"logger"`

	Metrics struct {
		Enabled *bool `yaml:"enabled"`
		OTEL    struct {
			Enabled  *bool  `yaml:"enabled"`
			Endpoint string `yaml:"endpoint"`
			Insecure *bool  `yaml:"insecure"`
			Protocol string `yaml:"protocol"`
		} `yaml:"otel"`
		Prometheus struct {
			Enabled *bool  `yaml:"enabled"`
			Port    int    `yaml:"port"`
			Path    string `yaml:"path"`
		} `yaml:"prometheus"`
		ExportInterval string            `yaml:"export_interval"`
		DefaultLabels  map[string]string `yaml:"default_labels"`
	} `yaml:"metrics"`
}

// Load loads configuration with the following precedence (highest to lowest):
// 1. Environment variables
// 2. YAML config file
// 3. Default values
func Load(configPath string) (Config, error) {
	cfg := DefaultConfig()

	if configPath != "" {
		if err := loadFromYAML(configPath, &cfg); err != nil {
			return cfg, fmt.Errorf("load yaml config: %w", err)
		}
	}

	loadFromEnv(&cfg)
	return cfg, nil
}

// LoadWithEnv loads configuration for a specific environment.
// It looks for config files in the following order:
// 1. {configDir}/config.yaml (base config)
// 2. {configDir}/config.{env}.yaml (environment-specific overrides)
func LoadWithEnv(configDir, env string) (Config, error) {
	cfg := DefaultConfig()

	basePath := filepath.Join(configDir, "config.yaml")
	if fileExists(basePath) {
		if err := loadFromYAML(basePath, &cfg); err != nil {
			return cfg, fmt.Errorf("load base config: %w", err)
		}
	}

	if env != "" {
		envPath := filepath.Join(configDir, fmt.Sprintf("config.%s.yaml", env))
		if fileExists(envPath) {
			if err := loadFromYAML(envPath, &cfg); err != nil {
				return cfg, fmt.Errorf("load env config: %w", err)
			}
		}
	}

	loadFromEnv(&cfg)
	return cfg, nil
}

// LoadFromEnv loads configuration purely from environment variables.
func LoadFromEnv() Config {
	cfg := DefaultConfig()
	loadFromEnv(&cfg)
	return cfg
}

// LoadLoggerConfig loads only logger configuration from a YAML file.
func LoadLoggerConfig(configPath string) (LoggerConfig, error) {
	cfg, err := Load(configPath)
	if err != nil {
		return LoggerConfig{}, err
	}
	return cfg.GetLoggerConfig(), nil
}

// LoadLoggerConfigFromEnv loads logger configuration from environment variables.
func LoadLoggerConfigFromEnv() LoggerConfig {
	cfg := LoadFromEnv()
	return cfg.GetLoggerConfig()
}

// LoadMetricsConfig loads only metrics configuration from a YAML file.
func LoadMetricsConfig(configPath string) (MetricsConfig, error) {
	cfg, err := Load(configPath)
	if err != nil {
		return MetricsConfig{}, err
	}
	return cfg.GetMetricsConfig(), nil
}

// LoadMetricsConfigFromEnv loads metrics configuration from environment variables.
func LoadMetricsConfigFromEnv() MetricsConfig {
	cfg := LoadFromEnv()
	return cfg.GetMetricsConfig()
}

// LoadTracerConfig loads only tracer configuration from a YAML file.
func LoadTracerConfig(configPath string) (TracerConfig, error) {
	cfg, err := Load(configPath)
	if err != nil {
		return TracerConfig{}, err
	}
	return cfg.GetTracerConfig(), nil
}

// LoadTracerConfigFromEnv loads tracer configuration from environment variables.
func LoadTracerConfigFromEnv() TracerConfig {
	cfg := LoadFromEnv()
	return cfg.GetTracerConfig()
}

// MustLoad loads configuration and panics on error.
func MustLoad(configPath string) Config {
	cfg, err := Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// MustLoadWithEnv loads environment-specific configuration and panics on error.
func MustLoadWithEnv(configDir, env string) Config {
	cfg, err := LoadWithEnv(configDir, env)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

func loadFromYAML(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file %s: %w", path, err)
	}

	expanded := os.ExpandEnv(string(data))

	var yc yamlConfig
	if err := yaml.Unmarshal([]byte(expanded), &yc); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	applyYAMLConfig(&yc, cfg)
	return nil
}

func applyYAMLConfig(yc *yamlConfig, cfg *Config) {
	// Service config
	if yc.Service.Name != "" {
		cfg.Service.Name = yc.Service.Name
		cfg.Logger.ServiceName = yc.Service.Name
		cfg.Metrics.ServiceName = yc.Service.Name
	}
	if yc.Service.Version != "" {
		cfg.Service.Version = yc.Service.Version
		cfg.Logger.ServiceVersion = yc.Service.Version
		cfg.Metrics.ServiceVersion = yc.Service.Version
	}
	if yc.Service.Environment != "" {
		cfg.Service.Environment = yc.Service.Environment
		cfg.Logger.Environment = yc.Service.Environment
		cfg.Metrics.Environment = yc.Service.Environment
	}

	// Logger config
	if yc.Logger.Level != "" {
		cfg.Logger.Level = yc.Logger.Level
	}
	if yc.Logger.Console.Enabled != nil {
		cfg.Logger.ConsoleEnabled = *yc.Logger.Console.Enabled
	}
	if yc.Logger.Console.Format != "" {
		cfg.Logger.ConsoleFormat = yc.Logger.Console.Format
	}
	if yc.Logger.OTEL.Enabled != nil {
		cfg.Logger.OTELEnabled = *yc.Logger.OTEL.Enabled
	}
	if yc.Logger.OTEL.Endpoint != "" {
		cfg.Logger.OTELEndpoint = yc.Logger.OTEL.Endpoint
	}
	if yc.Logger.OTEL.Insecure != nil {
		cfg.Logger.OTELInsecure = *yc.Logger.OTEL.Insecure
	}
	if yc.Logger.File.Enabled != nil {
		cfg.Logger.FileEnabled = *yc.Logger.File.Enabled
	}
	if yc.Logger.File.Path != "" {
		cfg.Logger.FilePath = yc.Logger.File.Path
	}
	if yc.Logger.File.MaxSizeMB > 0 {
		cfg.Logger.FileMaxSizeMB = yc.Logger.File.MaxSizeMB
	}
	if yc.Logger.File.MaxBackups > 0 {
		cfg.Logger.FileMaxBackups = yc.Logger.File.MaxBackups
	}
	if yc.Logger.File.MaxAgeDays > 0 {
		cfg.Logger.FileMaxAgeDays = yc.Logger.File.MaxAgeDays
	}
	if yc.Logger.File.Compress != nil {
		cfg.Logger.FileCompress = *yc.Logger.File.Compress
	}
	if yc.Logger.Caller.Enabled != nil {
		cfg.Logger.AddCaller = *yc.Logger.Caller.Enabled
	}
	if yc.Logger.Caller.Skip > 0 {
		cfg.Logger.CallerSkip = yc.Logger.Caller.Skip
	}
	if len(yc.Logger.Modules) > 0 {
		if cfg.Logger.ModuleLevels == nil {
			cfg.Logger.ModuleLevels = make(map[string]string)
		}
		for module, level := range yc.Logger.Modules {
			cfg.Logger.ModuleLevels[module] = level
		}
	}

	// Metrics config
	if yc.Metrics.Enabled != nil {
		cfg.Metrics.Enabled = *yc.Metrics.Enabled
	}
	if yc.Metrics.OTEL.Enabled != nil {
		cfg.Metrics.OTELEnabled = *yc.Metrics.OTEL.Enabled
	}
	if yc.Metrics.OTEL.Endpoint != "" {
		cfg.Metrics.OTELEndpoint = yc.Metrics.OTEL.Endpoint
	}
	if yc.Metrics.OTEL.Insecure != nil {
		cfg.Metrics.OTELInsecure = *yc.Metrics.OTEL.Insecure
	}
	if yc.Metrics.OTEL.Protocol != "" {
		cfg.Metrics.OTELProtocol = yc.Metrics.OTEL.Protocol
	}
	if yc.Metrics.Prometheus.Enabled != nil {
		cfg.Metrics.PrometheusEnabled = *yc.Metrics.Prometheus.Enabled
	}
	if yc.Metrics.Prometheus.Port > 0 {
		cfg.Metrics.PrometheusPort = yc.Metrics.Prometheus.Port
	}
	if yc.Metrics.Prometheus.Path != "" {
		cfg.Metrics.PrometheusPath = yc.Metrics.Prometheus.Path
	}
	if yc.Metrics.ExportInterval != "" {
		if d, err := time.ParseDuration(yc.Metrics.ExportInterval); err == nil {
			cfg.Metrics.ExportInterval = d
		}
	}
	if len(yc.Metrics.DefaultLabels) > 0 {
		if cfg.Metrics.DefaultLabels == nil {
			cfg.Metrics.DefaultLabels = make(map[string]string)
		}
		for k, v := range yc.Metrics.DefaultLabels {
			cfg.Metrics.DefaultLabels[k] = v
		}
	}
}

func loadFromEnv(cfg *Config) {
	// Service config
	if v := os.Getenv(EnvServiceName); v != "" {
		cfg.Service.Name = v
		cfg.Logger.ServiceName = v
		cfg.Metrics.ServiceName = v
	}
	if v := os.Getenv(EnvServiceVersion); v != "" {
		cfg.Service.Version = v
		cfg.Logger.ServiceVersion = v
		cfg.Metrics.ServiceVersion = v
	}
	if v := os.Getenv(EnvServiceEnvironment); v != "" {
		cfg.Service.Environment = v
		cfg.Logger.Environment = v
		cfg.Metrics.Environment = v
	}

	// Logger config
	if v := os.Getenv(EnvLogLevel); v != "" {
		cfg.Logger.Level = v
	}
	if v := os.Getenv(EnvLogConsoleEnabled); v != "" {
		cfg.Logger.ConsoleEnabled = parseBool(v)
	}
	if v := os.Getenv(EnvLogConsoleFormat); v != "" {
		cfg.Logger.ConsoleFormat = v
	}
	if v := os.Getenv(EnvLogOTELEnabled); v != "" {
		cfg.Logger.OTELEnabled = parseBool(v)
	}
	if v := os.Getenv(EnvLogOTELEndpoint); v != "" {
		cfg.Logger.OTELEndpoint = v
	}
	if v := os.Getenv(EnvLogOTELInsecure); v != "" {
		cfg.Logger.OTELInsecure = parseBool(v)
	}
	if v := os.Getenv(EnvLogOTELProtocol); v != "" {
		cfg.Logger.OTELProtocol = v
	}
	if v := os.Getenv(EnvLogOTELDebug); v != "" {
		cfg.Logger.OTELDebug = parseBool(v)
	}
	if v := os.Getenv(EnvLogFileEnabled); v != "" {
		cfg.Logger.FileEnabled = parseBool(v)
	}
	if v := os.Getenv(EnvLogFilePath); v != "" {
		cfg.Logger.FilePath = v
	}
	if v := os.Getenv(EnvLogFileMaxSizeMB); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Logger.FileMaxSizeMB = i
		}
	}
	if v := os.Getenv(EnvLogFileMaxBackups); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Logger.FileMaxBackups = i
		}
	}
	if v := os.Getenv(EnvLogFileMaxAgeDays); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Logger.FileMaxAgeDays = i
		}
	}
	if v := os.Getenv(EnvLogFileCompress); v != "" {
		cfg.Logger.FileCompress = parseBool(v)
	}
	if v := os.Getenv(EnvLogAddCaller); v != "" {
		cfg.Logger.AddCaller = parseBool(v)
	}
	if v := os.Getenv(EnvLogCallerSkip); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Logger.CallerSkip = i
		}
	}

	// Logger module levels
	loadModuleLevelsFromEnv(cfg)

	// Metrics config
	if v := os.Getenv(EnvMetricsEnabled); v != "" {
		cfg.Metrics.Enabled = parseBool(v)
	}
	if v := os.Getenv(EnvMetricsOTELEnabled); v != "" {
		cfg.Metrics.OTELEnabled = parseBool(v)
	}
	if v := os.Getenv(EnvMetricsOTELEndpoint); v != "" {
		cfg.Metrics.OTELEndpoint = v
	}
	if v := os.Getenv(EnvMetricsOTELInsecure); v != "" {
		cfg.Metrics.OTELInsecure = parseBool(v)
	}
	if v := os.Getenv(EnvMetricsOTELProtocol); v != "" {
		cfg.Metrics.OTELProtocol = v
	}
	if v := os.Getenv(EnvMetricsPrometheusEnabled); v != "" {
		cfg.Metrics.PrometheusEnabled = parseBool(v)
	}
	if v := os.Getenv(EnvMetricsPrometheusPort); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			cfg.Metrics.PrometheusPort = port
		}
	}
	if v := os.Getenv(EnvMetricsPrometheusPath); v != "" {
		cfg.Metrics.PrometheusPath = v
	}
	if v := os.Getenv(EnvMetricsExportInterval); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Metrics.ExportInterval = d
		}
	}

	// Tracer config
	if v := os.Getenv(EnvTracerEnabled); v != "" {
		cfg.Tracer.Enabled = parseBool(v)
	}
	if v := os.Getenv(EnvTracerOTELEndpoint); v != "" {
		cfg.Tracer.OTELEndpoint = v
	}
	if v := os.Getenv(EnvTracerOTELInsecure); v != "" {
		cfg.Tracer.OTELInsecure = parseBool(v)
	}
	if v := os.Getenv(EnvTracerOTELProtocol); v != "" {
		cfg.Tracer.OTELProtocol = v
	}
	if v := os.Getenv(EnvTracerOTELDebug); v != "" {
		cfg.Tracer.OTELDebug = parseBool(v)
	}
	if v := os.Getenv(EnvTracerSamplingRatio); v != "" {
		if ratio, err := strconv.ParseFloat(v, 64); err == nil {
			cfg.Tracer.SamplingRatio = ratio
		}
	}
	if v := os.Getenv(EnvTracerBatchTimeout); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			cfg.Tracer.BatchTimeout = d
		}
	}
	if v := os.Getenv(EnvTracerMaxExportBatch); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Tracer.MaxExportBatch = i
		}
	}
	if v := os.Getenv(EnvTracerMaxQueueSize); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cfg.Tracer.MaxQueueSize = i
		}
	}

	// Also apply service config to tracer
	if cfg.Tracer.ServiceName == "" {
		cfg.Tracer.ServiceName = cfg.Service.Name
	}
	if cfg.Tracer.ServiceVersion == "" {
		cfg.Tracer.ServiceVersion = cfg.Service.Version
	}
	if cfg.Tracer.Environment == "" {
		cfg.Tracer.Environment = cfg.Service.Environment
	}
}

func loadModuleLevelsFromEnv(cfg *Config) {
	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, envModulePrefix) {
			continue
		}

		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}

		moduleName := strings.TrimPrefix(parts[0], envModulePrefix)
		moduleName = strings.ToLower(moduleName)
		level := strings.ToLower(strings.TrimSpace(parts[1]))

		if moduleName != "" && level != "" {
			if cfg.Logger.ModuleLevels == nil {
				cfg.Logger.ModuleLevels = make(map[string]string)
			}
			cfg.Logger.ModuleLevels[moduleName] = level
		}
	}
}

func parseBool(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "1" || s == "yes" || s == "on"
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
