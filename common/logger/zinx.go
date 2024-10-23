package logger

import (
	"context"
	"fmt"
)

type ZapLogger struct {
}

func (log *ZapLogger) InfoF(format string, v ...interface{}) {
	Logger.Info(fmt.Sprintf(format, v...))
}

func (log *ZapLogger) ErrorF(format string, v ...interface{}) {
	Logger.Error(fmt.Sprintf(format, v...))
}

func (log *ZapLogger) DebugF(format string, v ...interface{}) {
	Logger.Debug(fmt.Sprintf(format, v...))
}

func (log *ZapLogger) InfoFX(ctx context.Context, format string, v ...interface{}) {
	Logger.Info(fmt.Sprintf(format, v...))
}

func (log *ZapLogger) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	Logger.Error(fmt.Sprintf(format, v...))
}

func (log *ZapLogger) DebugFX(ctx context.Context, format string, v ...interface{}) {
	Logger.Debug(fmt.Sprintf(format, v...))
}
