package global

import (
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 业务逻辑日志记录器
var Logger *zap.Logger

// ServiceLogger 服务日志记录器
var ServiceLogger *zap.Logger

// SetLogger 设置logger
func SetLogger() error {
	var err error

	// 设置日志记录级别
	var zapLevel zapcore.Level
	switch Config.Logger.Level {
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.DebugLevel
	}

	// 配置级别编码器
	var encodeLevel zapcore.LevelEncoder
	if Config.Logger.ColorLevel == true && Config.Logger.Encode == "console" {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encodeLevel = zapcore.CapitalLevelEncoder
	}

	outputs := strings.Split(Config.Logger.Outputs, "|")

	// 配置编码器的参数
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message", // 消息字段名
		LevelKey:   "level",   // 级别字段名
		TimeKey:    "time",    // 时间字段名
		CallerKey:  "file",    // 记录源码文件的字段名
		// 编码时间字符串的格式
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		EncodeLevel:  encodeLevel,                // 日志级别的编码器
		EncodeCaller: zapcore.ShortCallerEncoder, // Caller的编码器
	}

	var disableStacktrace bool
	var disableCaller bool
	if Config.Logger.EnableCaller == false {
		disableStacktrace = true
	}
	if Config.Logger.EnableCaller == false {
		disableCaller = true
	}

	// 设置Logger
	Logger, err = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // 日志记录级别
		Development:       Config.Service.Debug,           // 开发模式
		Encoding:          Config.Logger.Encode,           // 日志格式json/console
		EncoderConfig:     encoderConfig,                  // 编码器配置
		OutputPaths:       outputs,                        // 输出路径
		DisableStacktrace: disableStacktrace,              // 屏蔽堆栈跟踪
		DisableCaller:     disableCaller,                  // 屏蔽调用信息
	}.Build()
	if err != nil {
		return err
	}

	// 设置ServiceLogger
	ServiceLogger, err = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // 日志记录级别
		Development:       Config.Service.Debug,           // 开发模式
		Encoding:          Config.Logger.Encode,           // 日志格式json/console
		EncoderConfig:     encoderConfig,                  // 编码器配置
		OutputPaths:       outputs,                        // 输出路径
		DisableStacktrace: true,                           // 屏蔽堆栈跟踪
		DisableCaller:     true,                           // 屏蔽跟踪
	}.Build()
	if err != nil {
		return err
	}

	return err
}
