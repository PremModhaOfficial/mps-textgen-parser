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
	applyYAMLServiceConfig(yc, cfg)
	applyYAMLLoggerConfig(yc, cfg)
	applyYAMLMetricsConfig(yc, cfg)
}

// applyYAMLServiceConfig applies service-level YAML values to the config.
// Service values are propagated to logger and metrics configs for consistency.
func applyYAMLServiceConfig(yc *yamlConfig, cfg *Config) {
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
}

// applyYAMLLoggerConfig applies logger-level YAML values to the config.
func applyYAMLLoggerConfig(yc *yamlConfig, cfg *Config) {
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
}

// applyYAMLMetricsConfig applies metrics-level YAML values to the config.
func applyYAMLMetricsConfig(yc *yamlConfig, cfg *Config) {
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
		if duration, err := time.ParseDuration(yc.Metrics.ExportInterval); err == nil {
			cfg.Metrics.ExportInterval = duration
		}
	}
	if len(yc.Metrics.DefaultLabels) > 0 {
		if cfg.Metrics.DefaultLabels == nil {
			cfg.Metrics.DefaultLabels = make(map[string]string)
		}
		for labelKey, labelValue := range yc.Metrics.DefaultLabels {
			cfg.Metrics.DefaultLabels[labelKey] = labelValue
		}
	}
}

func loadFromEnv(cfg *Config) {
	loadServiceFromEnv(cfg)
	loadLoggerFromEnv(cfg)
	loadModuleLevelsFromEnv(cfg)
	loadMetricsFromEnv(cfg)
	loadTracerFromEnv(cfg)
	applyServiceDefaultsToTracer(cfg)
}

// loadServiceFromEnv loads service-level config from environment variables.
// Service values are propagated to logger and metrics configs for consistency.
func loadServiceFromEnv(cfg *Config) {
	if envValue := os.Getenv(EnvServiceName); envValue != "" {
		cfg.Service.Name = envValue
		cfg.Logger.ServiceName = envValue
		cfg.Metrics.ServiceName = envValue
	}
	if envValue := os.Getenv(EnvServiceVersion); envValue != "" {
		cfg.Service.Version = envValue
		cfg.Logger.ServiceVersion = envValue
		cfg.Metrics.ServiceVersion = envValue
	}
	if envValue := os.Getenv(EnvServiceEnvironment); envValue != "" {
		cfg.Service.Environment = envValue
		cfg.Logger.Environment = envValue
		cfg.Metrics.Environment = envValue
	}
}

// loadLoggerFromEnv loads logger config from environment variables.
func loadLoggerFromEnv(cfg *Config) {
	envStr(EnvLogLevel, &cfg.Logger.Level)
	envBool(EnvLogConsoleEnabled, &cfg.Logger.ConsoleEnabled)
	envStr(EnvLogConsoleFormat, &cfg.Logger.ConsoleFormat)
	envBool(EnvLogOTELEnabled, &cfg.Logger.OTELEnabled)
	envStr(EnvLogOTELEndpoint, &cfg.Logger.OTELEndpoint)
	envBool(EnvLogOTELInsecure, &cfg.Logger.OTELInsecure)
	envStr(EnvLogOTELProtocol, &cfg.Logger.OTELProtocol)
	envBool(EnvLogOTELDebug, &cfg.Logger.OTELDebug)
	envBool(EnvLogFileEnabled, &cfg.Logger.FileEnabled)
	envStr(EnvLogFilePath, &cfg.Logger.FilePath)
	envInt(EnvLogFileMaxSizeMB, &cfg.Logger.FileMaxSizeMB)
	envInt(EnvLogFileMaxBackups, &cfg.Logger.FileMaxBackups)
	envInt(EnvLogFileMaxAgeDays, &cfg.Logger.FileMaxAgeDays)
	envBool(EnvLogFileCompress, &cfg.Logger.FileCompress)
	envBool(EnvLogAddCaller, &cfg.Logger.AddCaller)
	envInt(EnvLogCallerSkip, &cfg.Logger.CallerSkip)
}

// loadMetricsFromEnv loads metrics config from environment variables.
func loadMetricsFromEnv(cfg *Config) {
	envBool(EnvMetricsEnabled, &cfg.Metrics.Enabled)
	envBool(EnvMetricsOTELEnabled, &cfg.Metrics.OTELEnabled)
	envStr(EnvMetricsOTELEndpoint, &cfg.Metrics.OTELEndpoint)
	envBool(EnvMetricsOTELInsecure, &cfg.Metrics.OTELInsecure)
	envStr(EnvMetricsOTELProtocol, &cfg.Metrics.OTELProtocol)
	envBool(EnvMetricsPrometheusEnabled, &cfg.Metrics.PrometheusEnabled)
	envInt(EnvMetricsPrometheusPort, &cfg.Metrics.PrometheusPort)
	envStr(EnvMetricsPrometheusPath, &cfg.Metrics.PrometheusPath)
	envDuration(EnvMetricsExportInterval, &cfg.Metrics.ExportInterval)
}

// loadTracerFromEnv loads tracer config from environment variables.
func loadTracerFromEnv(cfg *Config) {
	envBool(EnvTracerEnabled, &cfg.Tracer.Enabled)
	envStr(EnvTracerOTELEndpoint, &cfg.Tracer.OTELEndpoint)
	envBool(EnvTracerOTELInsecure, &cfg.Tracer.OTELInsecure)
	envStr(EnvTracerOTELProtocol, &cfg.Tracer.OTELProtocol)
	envBool(EnvTracerOTELDebug, &cfg.Tracer.OTELDebug)
	envFloat(EnvTracerSamplingRatio, &cfg.Tracer.SamplingRatio)
	envDuration(EnvTracerBatchTimeout, &cfg.Tracer.BatchTimeout)
	envInt(EnvTracerMaxExportBatch, &cfg.Tracer.MaxExportBatch)
	envInt(EnvTracerMaxQueueSize, &cfg.Tracer.MaxQueueSize)
}

// applyServiceDefaultsToTracer copies service-level values to the tracer config
// if they were not explicitly set for the tracer.
func applyServiceDefaultsToTracer(cfg *Config) {
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

/* ---------------------------------------- Environment Variable Helpers -------------------------------------------------- */

// envStr sets target to the environment variable value if set and non-empty.
func envStr(key string, target *string) {
	if envValue := os.Getenv(key); envValue != "" {
		*target = envValue
	}
}

// envBool sets target to the parsed boolean from the environment variable if set.
func envBool(key string, target *bool) {
	if envValue := os.Getenv(key); envValue != "" {
		normalized := strings.ToLower(strings.TrimSpace(envValue))
		*target = normalized == "true" || normalized == "1" || normalized == "yes" || normalized == "on"
	}
}

// envInt sets target to the parsed integer from the environment variable if valid.
func envInt(key string, target *int) {
	if envValue := os.Getenv(key); envValue != "" {
		if parsedInt, err := strconv.Atoi(envValue); err == nil {
			*target = parsedInt
		}
	}
}

// envFloat sets target to the parsed float from the environment variable if valid.
func envFloat(key string, target *float64) {
	if envValue := os.Getenv(key); envValue != "" {
		if parsedFloat, err := strconv.ParseFloat(envValue, 64); err == nil {
			*target = parsedFloat
		}
	}
}

// envDuration sets target to the parsed duration from the environment variable if valid.
func envDuration(key string, target *time.Duration) {
	if envValue := os.Getenv(key); envValue != "" {
		if duration, err := time.ParseDuration(envValue); err == nil {
			*target = duration
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
