package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	assertions := assert.New(t)

	defaultConfiguration := config.DefaultConfig()

	// Service defaults
	assertions.Equal("app", defaultConfiguration.Service.Name, "Service.Name should default to 'app'")
	assertions.Equal("0.0.0", defaultConfiguration.Service.Version, "Service.Version should default to '0.0.0'")
	assertions.Equal("development", defaultConfiguration.Service.Environment, "Service.Environment should default to 'development'")

	// Logger defaults (flat structure)
	assertions.Equal("info", defaultConfiguration.Logger.Level, "Logger.Level should default to 'info'")
	assertions.True(defaultConfiguration.Logger.ConsoleEnabled, "Logger.ConsoleEnabled should be true by default")
	assertions.Equal("json", defaultConfiguration.Logger.ConsoleFormat, "Logger.ConsoleFormat should default to 'json'")
	assertions.False(defaultConfiguration.Logger.OTELEnabled, "Logger.OTELEnabled should be false by default")
	assertions.False(defaultConfiguration.Logger.FileEnabled, "Logger.FileEnabled should be false by default")
	assertions.True(defaultConfiguration.Logger.AddCaller, "Logger.AddCaller should be true by default")

	// Metrics defaults (flat structure)
	assertions.True(defaultConfiguration.Metrics.Enabled, "Metrics.Enabled should be true by default")
	assertions.False(defaultConfiguration.Metrics.OTELEnabled, "Metrics.OTELEnabled should be false by default")
	assertions.False(defaultConfiguration.Metrics.PrometheusEnabled, "Metrics.PrometheusEnabled should be false by default")
	assertions.Equal(15*time.Second, defaultConfiguration.Metrics.ExportInterval, "Metrics.ExportInterval should default to 15s")
}

func TestLoadFromEnv(t *testing.T) {
	assertions := assert.New(t)

	// Set environment variables
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("SERVICE_VERSION", "1.2.3")
	os.Setenv("SERVICE_ENVIRONMENT", "production")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_CONSOLE_FORMAT", "console")
	os.Setenv("LOG_OTEL_ENABLED", "true")
	os.Setenv("LOG_OTEL_ENDPOINT", "otel.example.com:4317")
	os.Setenv("METRICS_ENABLED", "true")
	os.Setenv("METRICS_OTEL_ENABLED", "true")
	os.Setenv("METRICS_PROMETHEUS_ENABLED", "true")
	os.Setenv("METRICS_PROMETHEUS_PORT", "8080")
	os.Setenv("METRICS_EXPORT_INTERVAL", "30s")

	defer func() {
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("SERVICE_VERSION")
		os.Unsetenv("SERVICE_ENVIRONMENT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_CONSOLE_FORMAT")
		os.Unsetenv("LOG_OTEL_ENABLED")
		os.Unsetenv("LOG_OTEL_ENDPOINT")
		os.Unsetenv("METRICS_ENABLED")
		os.Unsetenv("METRICS_OTEL_ENABLED")
		os.Unsetenv("METRICS_PROMETHEUS_ENABLED")
		os.Unsetenv("METRICS_PROMETHEUS_PORT")
		os.Unsetenv("METRICS_EXPORT_INTERVAL")
	}()

	loadedConfiguration := config.LoadFromEnv()

	// Check service config
	assertions.Equal("test-service", loadedConfiguration.Service.Name, "Service.Name should be 'test-service'")
	assertions.Equal("1.2.3", loadedConfiguration.Service.Version, "Service.Version should be '1.2.3'")
	assertions.Equal("production", loadedConfiguration.Service.Environment, "Service.Environment should be 'production'")

	// Check logger config (flat structure)
	assertions.Equal("debug", loadedConfiguration.Logger.Level, "Logger.Level should be 'debug'")
	assertions.Equal("console", loadedConfiguration.Logger.ConsoleFormat, "Logger.ConsoleFormat should be 'console'")
	assertions.True(loadedConfiguration.Logger.OTELEnabled, "Logger.OTELEnabled should be true")
	assertions.Equal("otel.example.com:4317", loadedConfiguration.Logger.OTELEndpoint, "Logger.OTELEndpoint should be 'otel.example.com:4317'")

	// Check metrics config (flat structure)
	assertions.True(loadedConfiguration.Metrics.Enabled, "Metrics.Enabled should be true")
	assertions.True(loadedConfiguration.Metrics.OTELEnabled, "Metrics.OTELEnabled should be true")
	assertions.True(loadedConfiguration.Metrics.PrometheusEnabled, "Metrics.PrometheusEnabled should be true")
	assertions.Equal(8080, loadedConfiguration.Metrics.PrometheusPort, "Metrics.PrometheusPort should be 8080")
	assertions.Equal(30*time.Second, loadedConfiguration.Metrics.ExportInterval, "Metrics.ExportInterval should be 30s")
}

func TestLoadFromYAML(t *testing.T) {
	assertions := assert.New(t)

	temporaryDirectory := t.TempDir()
	configFilePath := filepath.Join(temporaryDirectory, "config.yaml")

	yamlConfigContent := `
service:
  name: yaml-service
  version: 2.0.0
  environment: staging

logger:
  level: warn
  console:
    enabled: true
    format: json
  otel:
    enabled: true
    endpoint: collector:4317
  file:
    enabled: true
    path: /var/log/app.log
    max_size_mb: 50
  modules:
    database: debug
    http: info

metrics:
  enabled: true
  otel:
    enabled: true
    endpoint: collector:4317
  prometheus:
    enabled: true
    port: 9091
  export_interval: 10s
  default_labels:
    region: us-west-2
`
	writeError := os.WriteFile(configFilePath, []byte(yamlConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write config file")

	loadedConfiguration, loadError := config.Load(configFilePath)
	assertions.NoError(loadError, "Load should not fail")

	// Check service config
	assertions.Equal("yaml-service", loadedConfiguration.Service.Name, "Service.Name should be 'yaml-service'")
	assertions.Equal("2.0.0", loadedConfiguration.Service.Version, "Service.Version should be '2.0.0'")
	assertions.Equal("staging", loadedConfiguration.Service.Environment, "Service.Environment should be 'staging'")

	// Check logger config (flat structure)
	assertions.Equal("warn", loadedConfiguration.Logger.Level, "Logger.Level should be 'warn'")
	assertions.True(loadedConfiguration.Logger.OTELEnabled, "Logger.OTELEnabled should be true")
	assertions.True(loadedConfiguration.Logger.FileEnabled, "Logger.FileEnabled should be true")
	assertions.Equal(50, loadedConfiguration.Logger.FileMaxSizeMB, "Logger.FileMaxSizeMB should be 50")

	// Check module levels
	assertions.Equal("debug", loadedConfiguration.Logger.ModuleLevels["database"], "Logger.ModuleLevels[database] should be 'debug'")
	assertions.Equal("info", loadedConfiguration.Logger.ModuleLevels["http"], "Logger.ModuleLevels[http] should be 'info'")

	// Check metrics config (flat structure)
	assertions.True(loadedConfiguration.Metrics.OTELEnabled, "Metrics.OTELEnabled should be true")
	assertions.Equal(9091, loadedConfiguration.Metrics.PrometheusPort, "Metrics.PrometheusPort should be 9091")
	assertions.Equal(10*time.Second, loadedConfiguration.Metrics.ExportInterval, "Metrics.ExportInterval should be 10s")
	assertions.Equal("us-west-2", loadedConfiguration.Metrics.DefaultLabels["region"], "Metrics.DefaultLabels[region] should be 'us-west-2'")
}

func TestLoadWithEnv(t *testing.T) {
	assertions := assert.New(t)

	temporaryDirectory := t.TempDir()

	// Base config
	baseConfigContent := `
service:
  name: base-service
  environment: development

logger:
  level: info
  console:
    enabled: true

metrics:
  enabled: true
`
	writeError := os.WriteFile(filepath.Join(temporaryDirectory, "config.yaml"), []byte(baseConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write base config")

	// Production override
	productionConfigContent := `
service:
  environment: production

logger:
  level: warn
  otel:
    enabled: true
    endpoint: prod-collector:4317

metrics:
  otel:
    enabled: true
`
	writeError = os.WriteFile(filepath.Join(temporaryDirectory, "config.production.yaml"), []byte(productionConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write production config")

	loadedConfiguration, loadError := config.LoadWithEnv(temporaryDirectory, "production")
	assertions.NoError(loadError, "LoadWithEnv should not fail")

	// Should have base values where not overridden
	assertions.Equal("base-service", loadedConfiguration.Service.Name, "Service.Name should be 'base-service' from base config")

	// Should have production overrides
	assertions.Equal("production", loadedConfiguration.Service.Environment, "Service.Environment should be 'production' from override")
	assertions.Equal("warn", loadedConfiguration.Logger.Level, "Logger.Level should be 'warn' from production config")
	assertions.True(loadedConfiguration.Logger.OTELEnabled, "Logger.OTELEnabled should be true from production config")
	assertions.Equal("prod-collector:4317", loadedConfiguration.Logger.OTELEndpoint, "Logger.OTELEndpoint should be 'prod-collector:4317'")
	assertions.True(loadedConfiguration.Metrics.OTELEnabled, "Metrics.OTELEnabled should be true from production config")
}

func TestEnvOverridesYAML(t *testing.T) {
	assertions := assert.New(t)

	temporaryDirectory := t.TempDir()

	yamlConfigContent := `
service:
  name: yaml-service
  environment: staging

logger:
  level: info
`
	writeError := os.WriteFile(filepath.Join(temporaryDirectory, "config.yaml"), []byte(yamlConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write config file")

	// Set env vars that should override YAML
	os.Setenv("SERVICE_NAME", "env-service")
	os.Setenv("LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("SERVICE_NAME")
		os.Unsetenv("LOG_LEVEL")
	}()

	loadedConfiguration, loadError := config.Load(filepath.Join(temporaryDirectory, "config.yaml"))
	assertions.NoError(loadError, "Load should not fail")

	// Env vars should override YAML
	assertions.Equal("env-service", loadedConfiguration.Service.Name, "Service.Name should be 'env-service' (env overrides yaml)")
	assertions.Equal("debug", loadedConfiguration.Logger.Level, "Logger.Level should be 'debug' (env overrides yaml)")

	// Non-overridden values should come from YAML
	assertions.Equal("staging", loadedConfiguration.Service.Environment, "Service.Environment should be 'staging' from YAML")
}

func TestModuleLevelsFromEnv(t *testing.T) {
	assertions := assert.New(t)

	os.Setenv("LOG_MODULE_DATABASE", "debug")
	os.Setenv("LOG_MODULE_HTTP", "warn")
	os.Setenv("LOG_MODULE_CACHE", "error")
	defer func() {
		os.Unsetenv("LOG_MODULE_DATABASE")
		os.Unsetenv("LOG_MODULE_HTTP")
		os.Unsetenv("LOG_MODULE_CACHE")
	}()

	loadedConfiguration := config.LoadFromEnv()

	assertions.Equal("debug", loadedConfiguration.Logger.ModuleLevels["database"], "Logger.ModuleLevels[database] should be 'debug'")
	assertions.Equal("warn", loadedConfiguration.Logger.ModuleLevels["http"], "Logger.ModuleLevels[http] should be 'warn'")
	assertions.Equal("error", loadedConfiguration.Logger.ModuleLevels["cache"], "Logger.ModuleLevels[cache] should be 'error'")
}

func TestEnvVarExpansionInYAML(t *testing.T) {
	assertions := assert.New(t)

	temporaryDirectory := t.TempDir()

	yamlConfigContent := `
service:
  name: ${MY_SERVICE_NAME}
  environment: ${MY_ENV}

logger:
  otel:
    endpoint: ${OTEL_ENDPOINT}
`
	writeError := os.WriteFile(filepath.Join(temporaryDirectory, "config.yaml"), []byte(yamlConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write config file")

	os.Setenv("MY_SERVICE_NAME", "expanded-service")
	os.Setenv("MY_ENV", "test")
	os.Setenv("OTEL_ENDPOINT", "collector.test:4317")
	defer func() {
		os.Unsetenv("MY_SERVICE_NAME")
		os.Unsetenv("MY_ENV")
		os.Unsetenv("OTEL_ENDPOINT")
	}()

	loadedConfiguration, loadError := config.Load(filepath.Join(temporaryDirectory, "config.yaml"))
	assertions.NoError(loadError, "Load should not fail")

	assertions.Equal("expanded-service", loadedConfiguration.Service.Name, "Service.Name should be 'expanded-service' from env expansion")
	assertions.Equal("test", loadedConfiguration.Service.Environment, "Service.Environment should be 'test' from env expansion")
	assertions.Equal("collector.test:4317", loadedConfiguration.Logger.OTELEndpoint, "Logger.OTELEndpoint should be 'collector.test:4317' from env expansion")
}

func TestMustLoad(t *testing.T) {
	assertions := assert.New(t)

	temporaryDirectory := t.TempDir()

	yamlConfigContent := `
service:
  name: must-load-service
`
	configFilePath := filepath.Join(temporaryDirectory, "config.yaml")
	writeError := os.WriteFile(configFilePath, []byte(yamlConfigContent), 0644)
	assertions.NoError(writeError, "Failed to write config file")

	loadedConfiguration := config.MustLoad(configFilePath)
	assertions.Equal("must-load-service", loadedConfiguration.Service.Name, "Service.Name should be 'must-load-service'")
}

func TestMustLoadPanicsOnError(t *testing.T) {
	assertions := assert.New(t)

	assertions.Panics(func() {
		config.MustLoad("/nonexistent/config.yaml")
	}, "MustLoad should panic on invalid path")
}

func TestInvalidYAML(t *testing.T) {
	assertions := assert.New(t)

	temporaryDirectory := t.TempDir()
	configFilePath := filepath.Join(temporaryDirectory, "config.yaml")

	writeError := os.WriteFile(configFilePath, []byte("invalid: yaml: [unclosed"), 0644)
	assertions.NoError(writeError, "Failed to write config file")

	_, loadError := config.Load(configFilePath)
	assertions.Error(loadError, "Load should return error for invalid YAML")
}

func TestBoolParsing(t *testing.T) {
	testCases := []struct {
		environmentVariable string
		environmentValue    string
		expectedResult      bool
	}{
		{"LOG_OTEL_ENABLED", "true", true},
		{"LOG_OTEL_ENABLED", "TRUE", true},
		{"LOG_OTEL_ENABLED", "1", true},
		{"LOG_OTEL_ENABLED", "yes", true},
		{"LOG_OTEL_ENABLED", "on", true},
		{"LOG_OTEL_ENABLED", "false", false},
		{"LOG_OTEL_ENABLED", "0", false},
		{"LOG_OTEL_ENABLED", "no", false},
		{"LOG_OTEL_ENABLED", "off", false},
		{"LOG_OTEL_ENABLED", "invalid", false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.environmentValue, func(t *testing.T) {
			assertions := assert.New(t)

			os.Setenv(testCase.environmentVariable, testCase.environmentValue)
			defer os.Unsetenv(testCase.environmentVariable)

			loadedConfiguration := config.LoadFromEnv()
			assertions.Equal(testCase.expectedResult, loadedConfiguration.Logger.OTELEnabled,
				"Logger.OTELEnabled should be %v for value %q", testCase.expectedResult, testCase.environmentValue)
		})
	}
}

func TestFileConfigDefaults(t *testing.T) {
	assertions := assert.New(t)

	defaultConfiguration := config.DefaultConfig()

	assertions.Equal("logs/app.log", defaultConfiguration.Logger.FilePath, "Logger.FilePath should default to 'logs/app.log'")
	assertions.Equal(100, defaultConfiguration.Logger.FileMaxSizeMB, "Logger.FileMaxSizeMB should default to 100")
	assertions.Equal(5, defaultConfiguration.Logger.FileMaxBackups, "Logger.FileMaxBackups should default to 5")
	assertions.Equal(7, defaultConfiguration.Logger.FileMaxAgeDays, "Logger.FileMaxAgeDays should default to 7")
	assertions.True(defaultConfiguration.Logger.FileCompress, "Logger.FileCompress should be true by default")
}

func TestGetLoggerConfig(t *testing.T) {
	assertions := assert.New(t)

	unifiedConfiguration := config.DefaultConfig()
	unifiedConfiguration.Service.Name = "my-service"
	unifiedConfiguration.Service.Version = "1.0.0"
	unifiedConfiguration.Service.Environment = "production"

	extractedLoggerConfig := unifiedConfiguration.GetLoggerConfig()

	assertions.Equal("my-service", extractedLoggerConfig.ServiceName, "ServiceName should be 'my-service'")
	assertions.Equal("1.0.0", extractedLoggerConfig.ServiceVersion, "ServiceVersion should be '1.0.0'")
	assertions.Equal("production", extractedLoggerConfig.Environment, "Environment should be 'production'")
}

func TestGetMetricsConfig(t *testing.T) {
	assertions := assert.New(t)

	unifiedConfiguration := config.DefaultConfig()
	unifiedConfiguration.Service.Name = "my-service"
	unifiedConfiguration.Service.Version = "1.0.0"
	unifiedConfiguration.Service.Environment = "production"

	extractedMetricsConfig := unifiedConfiguration.GetMetricsConfig()

	assertions.Equal("my-service", extractedMetricsConfig.ServiceName, "ServiceName should be 'my-service'")
	assertions.Equal("1.0.0", extractedMetricsConfig.ServiceVersion, "ServiceVersion should be '1.0.0'")
	assertions.Equal("production", extractedMetricsConfig.Environment, "Environment should be 'production'")
}
