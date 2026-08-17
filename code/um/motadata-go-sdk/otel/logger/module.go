package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ModuleLevels manages per-module log levels with thread-safe dynamic updates
type ModuleLevels struct {
	mu sync.RWMutex

	levels map[string]zap.AtomicLevel

	defaultLevel zap.AtomicLevel
}

// ModuleLogger wraps a Logger with module-specific level filtering
type ModuleLogger struct {
	*Logger

	module string

	moduleLevel zap.AtomicLevel
}

/* ---------------------------------------- Global State -------------------------------------------------- */

// Global module levels registry
var (
	moduleLevelsMu sync.RWMutex

	moduleLevels *ModuleLevels
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// newModuleLevels creates a new module levels registry
func newModuleLevels(defaultLevel Level) *ModuleLevels {

	return &ModuleLevels{
		levels:       make(map[string]zap.AtomicLevel),
		defaultLevel: zap.NewAtomicLevelAt(defaultLevel.zapLevel()),
	}
}

// initModuleLevels initializes the global module levels registry
func initModuleLevels(defaultLevel Level) {
	moduleLevelsMu.Lock()
	defer moduleLevelsMu.Unlock()

	moduleLevels = newModuleLevels(defaultLevel)
}

// getModuleLevels returns the global module levels registry
func getModuleLevels() *ModuleLevels {
	moduleLevelsMu.RLock()
	levels := moduleLevels
	moduleLevelsMu.RUnlock()

	if levels == nil {
		moduleLevelsMu.Lock()

		if moduleLevels == nil {
			moduleLevels = newModuleLevels(InfoLevel)
		}

		levels = moduleLevels
		moduleLevelsMu.Unlock()
	}

	return levels
}

/* ---------------------------------------- ModuleLevels Methods -------------------------------------------------- */

// SetLevel sets the log level for a specific module
func (manager *ModuleLevels) SetLevel(name string, level Level) {

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if existing, exists := manager.levels[name]; exists {

		existing.SetLevel(level.zapLevel())

	} else {

		manager.levels[name] = zap.NewAtomicLevelAt(level.zapLevel())
	}
}

// Level returns the log level for a specific module
func (manager *ModuleLevels) Level(name string) Level {

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	if level, exists := manager.levels[name]; exists {

		return Level(level.Level())
	}

	return Level(manager.defaultLevel.Level())
}

// getAtomicLevel returns the atomic level for a module using double-check locking
func (manager *ModuleLevels) getAtomicLevel(name string) zap.AtomicLevel {

	manager.mu.RLock()

	if level, exists := manager.levels[name]; exists {

		manager.mu.RUnlock()

		return level
	}

	manager.mu.RUnlock()

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if level, exists := manager.levels[name]; exists {

		return level
	}

	newLevel := zap.NewAtomicLevelAt(manager.defaultLevel.Level())
	manager.levels[name] = newLevel

	return newLevel
}

// SetDefaultLevel sets the default level for modules without explicit configuration
func (manager *ModuleLevels) SetDefaultLevel(level Level) {
	manager.defaultLevel.SetLevel(level.zapLevel())
}

// DefaultLevel returns the default module level
func (manager *ModuleLevels) DefaultLevel() Level {
	return Level(manager.defaultLevel.Level())
}

// ListModules returns all configured module names and their levels
func (manager *ModuleLevels) ListModules() map[string]Level {

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	result := make(map[string]Level, len(manager.levels))

	for name, level := range manager.levels {

		result[name] = Level(level.Level())
	}

	return result
}

// Reset clears all module-specific levels
func (manager *ModuleLevels) Reset() {

	manager.mu.Lock()
	defer manager.mu.Unlock()

	manager.levels = make(map[string]zap.AtomicLevel)
}

/* ---------------------------------------- Logger Module Methods -------------------------------------------------- */

// Module creates a new logger for a specific module with its own log level
func (logger *Logger) Module(name string) *ModuleLogger {

	manager := getModuleLevels()
	level := manager.getAtomicLevel(name)

	named := logger.Named(name)

	return &ModuleLogger{
		Logger:      named,
		module:      name,
		moduleLevel: level,
	}
}

/* ---------------------------------------- ModuleLogger Methods -------------------------------------------------- */

// Debug logs at debug level if enabled for this module
func (moduleLogger *ModuleLogger) Debug(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.DebugLevel) {

		moduleLogger.Logger.Debug(ctx, message, fields...)
	}
}

// Info logs at info level if enabled for this module
func (moduleLogger *ModuleLogger) Info(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.InfoLevel) {

		moduleLogger.Logger.Info(ctx, message, fields...)
	}
}

// Warn logs at warn level if enabled for this module
func (moduleLogger *ModuleLogger) Warn(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.WarnLevel) {

		moduleLogger.Logger.Warn(ctx, message, fields...)
	}
}

// Error logs at error level if enabled for this module
func (moduleLogger *ModuleLogger) Error(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.ErrorLevel) {

		moduleLogger.Logger.Error(ctx, message, fields...)
	}
}

// Fatal logs at fatal level if enabled for this module
func (moduleLogger *ModuleLogger) Fatal(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.FatalLevel) {

		moduleLogger.Logger.Fatal(ctx, message, fields...)
	}
}

// With returns a ModuleLogger with additional preset fields
func (moduleLogger *ModuleLogger) With(fields ...Field) *ModuleLogger {
	return &ModuleLogger{
		Logger:      moduleLogger.Logger.With(fields...),
		module:      moduleLogger.module,
		moduleLevel: moduleLogger.moduleLevel,
	}
}

// SetLevel sets the log level for this module dynamically
func (moduleLogger *ModuleLogger) SetLevel(level Level) {
	moduleLogger.moduleLevel.SetLevel(level.zapLevel())
}

// Level returns the current log level for this module
func (moduleLogger *ModuleLogger) Level() Level {
	return Level(moduleLogger.moduleLevel.Level())
}

// ModuleName returns the name of this module
func (moduleLogger *ModuleLogger) ModuleName() string {

	return moduleLogger.module
}

/* ---------------------------------------- Package-level Module Functions -------------------------------------------------- */

// Module returns a module logger from the global logger
func Module(name string) *ModuleLogger {

	return L().Module(name)
}

// SetModuleLevel sets the log level for a specific module globally
func SetModuleLevel(module string, level Level) {

	getModuleLevels().SetLevel(module, level)
}

// ModuleLevel returns the current log level for a specific module
func ModuleLevel(module string) Level {

	return getModuleLevels().Level(module)
}

// SetModuleLevels sets multiple module levels at once
func SetModuleLevels(levels map[string]Level) {

	manager := getModuleLevels()

	for name, level := range levels {

		manager.SetLevel(name, level)
	}
}

// ListModuleLevels returns all configured module levels
func ListModuleLevels() map[string]Level {

	return getModuleLevels().ListModules()
}

// ResetModuleLevels clears all module-specific levels
func ResetModuleLevels() {
	getModuleLevels().Reset()
}
