package mylogger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FwdToZapWriter struct {
	logger *zap.SugaredLogger
}

func (fw *FwdToZapWriter) Write(p []byte) (n int, err error) {
	fw.logger.Errorw(string(p))
	return len(p), nil
}

func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return cfg.Build()
}
