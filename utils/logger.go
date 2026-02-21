package utils

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

// InitLogger 初始化 Zap 日志
func InitLogger() {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

// GetLogger 返回 Logger 实例
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}
