package logger

import (
	"context"
	"fmt"
	"gameserver/conf"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

var (
	Logger *zap.Logger
)

func Init() {
	coreList := make([]zapcore.Core, 0)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = CustomTimeFormatEncoder
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	log := conf.GameConfig.Log

	logName := fmt.Sprintf(log.FilePath, time.Now().Format("2006-01-02"))
	infoLumberjack := lumberjack.Logger{
		Filename:   logName,        // 日志文件路径
		MaxSize:    log.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     log.MaxAge,     // 文件最多保存多少天
		MaxBackups: log.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   false,          // 是否压缩
	}
	coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(&infoLumberjack), zapcore.InfoLevel))

	if log.Level == "debug" {
		coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
	}
	errLogName := fmt.Sprintf(log.ErrFilePath, time.Now().Format("2006-01-02"))
	errLumberjack := lumberjack.Logger{
		Filename:   errLogName,     // 日志文件路径
		MaxSize:    log.MaxSize,    // 每个日志文件保存的大小 单位:M
		MaxAge:     log.MaxAge * 2, // 文件最多保存多少天
		MaxBackups: log.MaxBackups, // 日志文件最多保存多少个备份
		Compress:   false,          // 是否压缩
	}
	coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(&errLumberjack), zapcore.ErrorLevel))

	core := zapcore.NewTee(coreList...)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()

	addCallerSkip := zap.AddCallerSkip(1)

	addStacktrace := zap.AddStacktrace(zap.ErrorLevel)

	// 开启文件及行号
	development := zap.Development()
	// 构造日志
	Logger = zap.New(core, caller, addCallerSkip, development, addStacktrace)

	Logger.Info("logger init success")
}

func getGoroutineID() uint64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := fmt.Sscanf(idField, "%d", new(uint64))
	if err == nil {
		return uint64(id)
	}
	return 0
}

func CustomTimeFormatEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(t.Format(conf.GameConfig.Log.TimeFormat))
}

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
			zap.Uint64("goroutine_id", getGoroutineID()),
		)
	}
}

func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()),
							"connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					_ = c.Error(err.(error))
					c.Abort()
					return
				}
				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}

type GormLoggerAdapter struct {
	Logger *zap.SugaredLogger
}

func (l *GormLoggerAdapter) LogMode(logLevel logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLoggerAdapter) Info(ctx context.Context, msg string, values ...interface{}) {
	l.Logger.Infow(msg, values...)
}

func (l *GormLoggerAdapter) Warn(ctx context.Context, msg string, values ...interface{}) {
	l.Logger.Warnw(msg, values...)
}

func (l *GormLoggerAdapter) Error(ctx context.Context, msg string, values ...interface{}) {
	l.Logger.Errorw(msg, values...)
}

func (l *GormLoggerAdapter) Trace(ctx context.Context, begin time.Time, fc func() (sql string, row int64), err error) {
	sql, rows := fc()
	if err != nil {
		l.Logger.Errorw("Database operation failed", zap.String("sql", sql), zap.Int64("rows_affected", rows), zap.Error(err))
	} else {
		l.Logger.Infow("Database operation succeeded", zap.String("sql", sql), zap.Int64("rows_affected", rows))
	}
}
