package global

import (
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger struct {
	Default  *zap.Logger // 默认logger，不带caller信息
	Caller   *zap.Logger // 带caller信息的logger
	StdError *log.Logger // 实现标准包log的logger，级别为error
}

// 设置logger
func SetLogger() error {
	var err error

	level := strings.ToLower(Config.Logger.Level)
	// 设置日志记录级别
	var zapLevel zapcore.Level
	switch level {
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	case "dpanic":
		zapLevel = zapcore.DPanicLevel
	case "panic":
		zapLevel = zapcore.PanicLevel
	case "fatal":
		zapLevel = zapcore.FatalLevel
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

	// 设置Logger
	Logger.Default, err = zap.Config{
		Level:             zap.NewAtomicLevelAt(zapLevel), // 日志记录级别
		Development:       Config.Service.Debug,           // 开发模式
		Encoding:          Config.Logger.Encode,           // 日志格式json/console
		EncoderConfig:     encoderConfig,                  // 编码器配置
		OutputPaths:       outputs,                        // 输出路径
		DisableStacktrace: true,                           // 屏蔽堆栈跟踪
		DisableCaller:     true,                           // 屏蔽调用信息
	}.Build()
	if err != nil {
		return err
	}

	// 将Logger转为log.Logger，级别为error
	Logger.StdError, err = zap.NewStdLogAt(Logger.Default, zap.ErrorLevel)
	if err != nil {
		return err
	}

	// 带caller的logger
	Logger.Caller = Logger.Default.WithOptions(zap.AddCaller())
	return err
}
