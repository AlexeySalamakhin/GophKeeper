// Package logger содержит настройку структурированного логирования.
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func Init(level string) error {
	config := zap.NewProductionConfig()

	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		atomicLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
	config.Level = atomicLevel

	// Настройка энкодера для читаемого вывода
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.CallerKey = "caller"

	Logger, err = config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	return nil
}

// InitDevelopment инициализирует логгер для разработки (более читаемый вывод).
func InitDevelopment() error {
	var err error
	Logger, err = zap.NewDevelopment(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	return err
}

func Sync() {
	if Logger != nil {
		Logger.Sync()
	}
}
