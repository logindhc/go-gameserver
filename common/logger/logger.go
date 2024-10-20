package logger

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var (
	name, cfgType, path = "logger", "yml", "."
	Logger              *zap.Logger
	LoggerConfig        *Config
)

type Config struct {
	Level       string `mapstructure:"level" json:"level" yaml:"level"`
	FilePath    string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
	ErrFilePath string `mapstructure:"err_file_path" json:"err_file_path" yaml:"err_file_path"`
	MaxAge      int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	MaxSize     int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups  int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
}

func initConfig() {
	viper.SetConfigName(name)
	viper.SetConfigType(cfgType)
	viper.AddConfigPath(path)
	if err := viper.ReadInConfig(); err != nil {
		_ = fmt.Errorf("read %s config error: %v", name, err)
		return
	}
	loggerConfig := &Config{}
	if err := viper.Unmarshal(loggerConfig); err != nil {
		_ = fmt.Errorf("%s config unbale to decode into struct: %v", name, err)
		return
	}
	LoggerConfig = loggerConfig
}

func init() {

	initConfig()

	coreList := make([]zapcore.Core, 0)
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	infoLumberjack := lumberjack.Logger{
		Filename:   fmt.Sprintf(LoggerConfig.FilePath, time.Now().Format("2006-01-02")), // 日志文件路径
		MaxSize:    LoggerConfig.MaxSize,                                                // 每个日志文件保存的大小 单位:M
		MaxAge:     LoggerConfig.MaxAge,                                                 // 文件最多保存多少天
		MaxBackups: LoggerConfig.MaxBackups,                                             // 日志文件最多保存多少个备份
		Compress:   false,                                                               // 是否压缩
	}
	coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(&infoLumberjack), zapcore.InfoLevel))

	if LoggerConfig.Level == "debug" {
		coreList = append(coreList, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel))
	}

	errLumberjack := lumberjack.Logger{
		Filename:   fmt.Sprintf(LoggerConfig.ErrFilePath, time.Now().Format("2006-01-02")), // 日志文件路径
		MaxSize:    LoggerConfig.MaxSize,                                                   // 每个日志文件保存的大小 单位:M
		MaxAge:     LoggerConfig.MaxAge * 2,                                                // 文件最多保存多少天
		MaxBackups: LoggerConfig.MaxBackups,                                                // 日志文件最多保存多少个备份
		Compress:   false,                                                                  // 是否压缩
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
		//l.Logger.Infow("Database operation succeeded", zap.String("sql", sql), zap.Int64("rows_affected", rows))
	}
}
