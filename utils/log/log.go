package log

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go-template/utils/activity"
)

const (
	MAX_LOG_ENTRY_SIZE = 8 * 1024
)

var logger *zap.Logger

type Trail struct {
	Label   string
	Payload interface{}
}

// ContextLogger wraps *zap.Logger with helper methods for easier migration from Logrus
type ContextLogger struct {
	*zap.Logger
}

// Error logs an error message with optional error value
func (cl *ContextLogger) Error(msg string, err ...error) {
	if len(err) > 0 && err[0] != nil {
		cl.Logger.Error(msg, zap.Error(err[0]))
	} else {
		cl.Logger.Error(msg)
	}
}

// Info logs an info message
func (cl *ContextLogger) Info(msg string) {
	cl.Logger.Info(msg)
}

// Debug logs a debug message
func (cl *ContextLogger) Debug(msg string) {
	cl.Logger.Debug(msg)
}

// Warn logs a warning message
func (cl *ContextLogger) Warn(msg string) {
	cl.Logger.Warn(msg)
}

// Initialize the logger
func init() {
	var err error
	config := zap.NewProductionConfig()

	// Configure based on environment
	if os.Getenv("APP_MODE") != "release" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	logger, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

// WithContext returns a logger with context fields wrapped in ContextLogger
func WithContext(ctx context.Context) *ContextLogger {
	fields := activity.GetFields(ctx)

	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	return &ContextLogger{Logger: logger.With(zapFields...)}
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return logger
}

// LogOrmer logs ORM debug information
func LogOrmer(obj interface{}, prefix string) {
	logger.Debug(prefix, zap.Any("data", obj))
}

// LogTrail logs trail information
func LogTrail(obj interface{}, prefix string) {
	logger.Info(prefix, zap.Any("data", obj))
}

// LogTrails logs multiple trails
func LogTrails(trails []Trail) {
	for _, trail := range trails {
		go LogTrail(trail.Payload, trail.Label)
	}
}

// Sync flushes any buffered log entries
func Sync() {
	_ = logger.Sync()
}
