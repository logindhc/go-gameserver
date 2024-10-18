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
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	FilePath   string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
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

	filePath := LoggerConfig.FilePath

	debug := LoggerConfig.Level == "debug"

	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf(filePath, time.Now().Format("2006-01-02")), // 日志文件路径
		MaxSize:    LoggerConfig.MaxSize,                                   // 每个日志文件保存的大小 单位:M
		MaxAge:     LoggerConfig.MaxAge,                                    // 文件最多保存多少天
		MaxBackups: LoggerConfig.MaxBackups,                                // 日志文件最多保存多少个备份
		Compress:   false,                                                  // 是否压缩
	}
	defer hook.Close()
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "file",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	if debug {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else {
		atomicLevel.SetLevel(zap.InfoLevel)
	}
	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	// 如果是开发环境，同时在控制台上也输出
	if debug {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(writes...),
		atomicLevel,
	)

	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()

	//// 设置初始化字段
	//field := zap.Fields(zap.String("appName", name))
	//global.Logger = zap.New(core, caller, development,field)
	// 构造日志
	Logger = zap.New(core, caller, development)

	defer Logger.Sync()

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
