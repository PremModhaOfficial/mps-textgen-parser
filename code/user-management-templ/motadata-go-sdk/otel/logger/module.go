package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ModuleLevels manages per-module log levels with thread-safe dynamic updates.
// This enables fine-grained control over logging verbosity for different
// parts of an application (e.g., enable debug for "auth" while keeping "http" at info).
//
// Key features:
//   - Per-module log level configuration
//   - Default level for unconfigured modules
//   - Thread-safe dynamic level changes at runtime
//   - Double-check locking for efficient level lookup
//
// Usage:
//
//	// Set specific module levels
//	moduleLevels.SetLevel("database", DebugLevel)
//	moduleLevels.SetLevel("http", WarnLevel)
//
//	// Query module level
//	level := moduleLevels.Level("database")  // Returns DebugLevel
type ModuleLevels struct {
	// mu protects concurrent access to the levels map
	mu sync.RWMutex

	// levels maps module names to their atomic log levels
	// Using zap.AtomicLevel enables lock-free level checks at log time
	levels map[string]zap.AtomicLevel

	// defaultLevel is used for modules without explicit configuration
	// This is an atomic level to allow runtime changes
	defaultLevel zap.AtomicLevel
}

// ModuleLogger wraps a Logger with module-specific level filtering.
// This allows different parts of the application to have independent log levels
// while sharing the underlying logger infrastructure.
//
// Usage:
//
//	// Create a module logger
//	dbLogger := logger.Module("database")
//
//	// Set the module's level independently
//	dbLogger.SetLevel(DebugLevel)
//
//	// Log at the module level
//	dbLogger.Debug(ctx, "query executed", String("table", "users"))
type ModuleLogger struct {
	// Embedded Logger provides all standard logging methods
	*Logger

	// module is the name of this module (e.g., "database", "http", "auth")
	module string

	// moduleLevel is the atomic level for this specific module
	// Changes to this level take effect immediately without locks
	moduleLevel zap.AtomicLevel
}

/* ---------------------------------------- Global State -------------------------------------------------- */

// Global module levels registry.
// This singleton manages all module-level configurations across the application.
var (
	// moduleLevelsMu protects initialization of moduleLevels
	moduleLevelsMu sync.RWMutex

	// moduleLevels is the global registry of per-module log levels
	moduleLevels *ModuleLevels
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// newModuleLevels creates a new module levels registry with the specified default level.
// The default level is used for any module that doesn't have an explicit level set.
//
// Parameters:
//   - defaultLevel: The default level for unconfigured modules
//
// Returns:
//   - *ModuleLevels: A new module levels registry
func newModuleLevels(defaultLevel Level) *ModuleLevels {

	return &ModuleLevels{
		// Initialize empty levels map - modules are added on demand
		levels: make(map[string]zap.AtomicLevel),

		// Create atomic default level for lock-free runtime changes
		defaultLevel: zap.NewAtomicLevelAt(defaultLevel.zapLevel()),
	}
}

// initModuleLevels initializes the global module levels registry.
// This is called during logger initialization to set up the module system.
//
// Parameters:
//   - defaultLevel: The default level for modules without explicit configuration
func initModuleLevels(defaultLevel Level) {

	moduleLevelsMu.Lock()
	defer moduleLevelsMu.Unlock()

	// Create new registry, replacing any existing one
	// This allows re-initialization on logger reconfiguration
	moduleLevels = newModuleLevels(defaultLevel)
}

// getModuleLevels returns the global module levels registry.
// If the registry hasn't been initialized, it creates one with InfoLevel as default.
//
// This function uses double-check locking for efficient access:
//   - Fast path: read lock to check if initialized
//   - Slow path: write lock to initialize if needed
//
// Returns:
//   - *ModuleLevels: The global module levels registry
func getModuleLevels() *ModuleLevels {

	// Fast path: check with read lock
	moduleLevelsMu.RLock()
	levels := moduleLevels
	moduleLevelsMu.RUnlock()

	if levels != nil {

		return levels
	}

	// Slow path: initialize with write lock
	moduleLevelsMu.Lock()

	// Double-check after acquiring write lock
	// Another goroutine may have initialized while we waited
	if moduleLevels == nil {
		moduleLevels = newModuleLevels(InfoLevel)
	}

	levels = moduleLevels
	moduleLevelsMu.Unlock()

	return levels
}

/* ---------------------------------------- ModuleLevels Methods -------------------------------------------------- */

// SetLevel sets the log level for a specific module.
// If the module already has a level, it's updated atomically.
// If the module is new, a new atomic level is created.
//
// Parameters:
//   - name: The module name (e.g., "database", "http")
//   - level: The desired log level
//
// Example:
//
//	manager.SetLevel("database", DebugLevel)  // Enable debug for database module
func (manager *ModuleLevels) SetLevel(name string, level Level) {

	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Check if module already exists
	if existing, exists := manager.levels[name]; exists {
		// Update existing level atomically (no allocation)
		existing.SetLevel(level.zapLevel())

	} else {
		// Create new atomic level for this module
		manager.levels[name] = zap.NewAtomicLevelAt(level.zapLevel())
	}
}

// Level returns the log level for a specific module.
// If the module doesn't have an explicit level, the default level is returned.
//
// Parameters:
//   - name: The module name to query
//
// Returns:
//   - Level: The module's current level (or default if not configured)
func (manager *ModuleLevels) Level(name string) Level {

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	if level, exists := manager.levels[name]; exists {

		return Level(level.Level())
	}

	// Return default level for unconfigured modules
	return Level(manager.defaultLevel.Level())
}

// getAtomicLevel returns the atomic level for a module using double-check locking.
// If the module doesn't have a level, one is created with the default level.
// This method is optimized for the hot path of module loggers.
//
// Parameters:
//   - name: The module name
//
// Returns:
//   - zap.AtomicLevel: The module's atomic level (never nil)
func (manager *ModuleLevels) getAtomicLevel(name string) zap.AtomicLevel {

	// Fast path: check with read lock
	manager.mu.RLock()

	if level, exists := manager.levels[name]; exists {

		manager.mu.RUnlock()

		return level
	}

	manager.mu.RUnlock()

	// Slow path: create level with write lock
	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Double-check after acquiring write lock
	if level, exists := manager.levels[name]; exists {

		return level
	}

	// Create new level with default value
	newLevel := zap.NewAtomicLevelAt(manager.defaultLevel.Level())
	manager.levels[name] = newLevel

	return newLevel
}

// SetDefaultLevel sets the default level for modules without explicit configuration.
// This affects all unconfigured modules immediately.
//
// Parameters:
//   - level: The new default level
func (manager *ModuleLevels) SetDefaultLevel(level Level) {

	// Atomic update - takes effect immediately
	manager.defaultLevel.SetLevel(level.zapLevel())
}

// DefaultLevel returns the default module level.
//
// Returns:
//   - Level: The current default level
func (manager *ModuleLevels) DefaultLevel() Level {

	return Level(manager.defaultLevel.Level())
}

// ListModules returns all configured module names and their levels.
// This is useful for displaying current module configuration.
//
// Returns:
//   - map[string]Level: Map of module names to their levels
func (manager *ModuleLevels) ListModules() map[string]Level {

	manager.mu.RLock()
	defer manager.mu.RUnlock()

	// Copy map to avoid holding lock during iteration by caller
	result := make(map[string]Level, len(manager.levels))

	for name, level := range manager.levels {

		result[name] = Level(level.Level())
	}

	return result
}

// Reset clears all module-specific levels.
// After reset, all modules will use the default level.
func (manager *ModuleLevels) Reset() {

	manager.mu.Lock()
	defer manager.mu.Unlock()

	// Create new empty map (old one will be garbage collected)
	manager.levels = make(map[string]zap.AtomicLevel)
}

/* ---------------------------------------- Logger Module Methods -------------------------------------------------- */

// Module creates a new logger for a specific module with its own log level.
// The module logger uses the parent's configuration but can have an independent level.
//
// Parameters:
//   - name: The module name (e.g., "database", "http", "auth")
//
// Returns:
//   - *ModuleLogger: A module-specific logger
//
// Example:
//
//	dbLogger := logger.Module("database")
//	dbLogger.SetLevel(DebugLevel)  // Independent level
//	dbLogger.Info(ctx, "connected")  // Output includes "logger":"<parent>.database"
func (logger *Logger) Module(name string) *ModuleLogger {

	// Get or create atomic level for this module
	manager := getModuleLevels()
	level := manager.getAtomicLevel(name)

	// Create named sub-logger (adds module name to logger path)
	named := logger.Named(name)

	return &ModuleLogger{
		Logger:      named,
		module:      name,
		moduleLevel: level,
	}
}

/* ---------------------------------------- ModuleLogger Methods -------------------------------------------------- */

// Debug logs at debug level if enabled for this module.
// This method checks the module-specific level before logging.
//
// Parameters:
//   - ctx: Context for field extraction
//   - message: The log message
//   - fields: Additional structured fields
func (moduleLogger *ModuleLogger) Debug(ctx context.Context, message string, fields ...Field) {

	// Early bailout if level is disabled
	// This check uses the module's atomic level for efficiency
	if moduleLogger.moduleLevel.Enabled(zapcore.DebugLevel) {

		moduleLogger.Logger.Debug(ctx, message, fields...)
	}
}

// Info logs at info level if enabled for this module.
func (moduleLogger *ModuleLogger) Info(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.InfoLevel) {

		moduleLogger.Logger.Info(ctx, message, fields...)
	}
}

// Warn logs at warn level if enabled for this module.
func (moduleLogger *ModuleLogger) Warn(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.WarnLevel) {

		moduleLogger.Logger.Warn(ctx, message, fields...)
	}
}

// Error logs at error level if enabled for this module.
func (moduleLogger *ModuleLogger) Error(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.ErrorLevel) {

		moduleLogger.Logger.Error(ctx, message, fields...)
	}
}

// Fatal logs at fatal level if enabled for this module.
// Note: This will terminate the application after logging.
func (moduleLogger *ModuleLogger) Fatal(ctx context.Context, message string, fields ...Field) {

	if moduleLogger.moduleLevel.Enabled(zapcore.FatalLevel) {

		moduleLogger.Logger.Fatal(ctx, message, fields...)
	}
}

// With returns a ModuleLogger with additional preset fields.
// The new logger shares the same module level.
//
// Parameters:
//   - fields: Fields to add to all future log entries
//
// Returns:
//   - *ModuleLogger: A new module logger with preset fields
func (moduleLogger *ModuleLogger) With(fields ...Field) *ModuleLogger {

	return &ModuleLogger{
		Logger:      moduleLogger.Logger.With(fields...),
		module:      moduleLogger.module,
		moduleLevel: moduleLogger.moduleLevel,
	}
}

// SetLevel sets the log level for this module dynamically.
// Changes take effect immediately for this and all other loggers using this module.
//
// Parameters:
//   - level: The new log level
func (moduleLogger *ModuleLogger) SetLevel(level Level) {

	moduleLogger.moduleLevel.SetLevel(level.zapLevel())
}

// Level returns the current log level for this module.
//
// Returns:
//   - Level: The module's current level
func (moduleLogger *ModuleLogger) Level() Level {

	return Level(moduleLogger.moduleLevel.Level())
}

// ModuleName returns the name of this module.
//
// Returns:
//   - string: The module name
func (moduleLogger *ModuleLogger) ModuleName() string {

	return moduleLogger.module
}

/* ---------------------------------------- Package-level Module Functions -------------------------------------------------- */

// Module returns a module logger from the global logger.
// This is a convenience function for creating module loggers.
//
// Parameters:
//   - name: The module name
//
// Returns:
//   - *ModuleLogger: A module-specific logger
//
// Example:
//
//	dbLogger := logger.Module("database")
func Module(name string) *ModuleLogger {

	return L().Module(name)
}

// SetModuleLevel sets the log level for a specific module globally.
// This affects all loggers using the specified module.
//
// Parameters:
//   - module: The module name
//   - level: The desired log level
func SetModuleLevel(module string, level Level) {

	getModuleLevels().SetLevel(module, level)
}

// ModuleLevel returns the current log level for a specific module.
//
// Parameters:
//   - module: The module name
//
// Returns:
//   - Level: The module's current level
func ModuleLevel(module string) Level {

	return getModuleLevels().Level(module)
}

// SetModuleLevels sets multiple module levels at once.
// This is useful for bulk configuration from configuration files.
//
// Parameters:
//   - levels: Map of module names to levels
func SetModuleLevels(levels map[string]Level) {

	manager := getModuleLevels()

	for name, level := range levels {

		manager.SetLevel(name, level)
	}
}

// ListModuleLevels returns all configured module levels.
//
// Returns:
//   - map[string]Level: Map of module names to their levels
func ListModuleLevels() map[string]Level {

	return getModuleLevels().ListModules()
}

// ResetModuleLevels clears all module-specific levels.
// After reset, all modules will use the default level.
func ResetModuleLevels() {

	getModuleLevels().Reset()
}
