package global

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// 使用默认参数设置logger，用于没有读取到配置文件时替代标准包的log使用
func SetDefaultLogger() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	zerolog.TimeFieldFormat = Now("y-m-d h:i:s")

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: zerolog.TimeFieldFormat,
	}

	log.Logger = log.Output(output)
}

// 当前时间转为字符串
func Now(str string) string {
	str = strings.Replace(str, "y", "2006", -1)
	str = strings.Replace(str, "m", "01", -1)
	str = strings.Replace(str, "d", "02", -1)
	str = strings.Replace(str, "h", "15", -1)
	str = strings.Replace(str, "i", "04", -1)
	str = strings.Replace(str, "s", "05", -1)
	return str
}

// 设置logger
func SetLogger() error {
	// 设置级别
	level := strings.ToLower(Config.Logger.Level)
	switch level {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "empty":
		zerolog.SetGlobalLevel(zerolog.NoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
		return nil
	}

	// 设置时间格式
	if Config.Logger.TimeFormat == "timestamp" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	} else {
		zerolog.TimeFieldFormat = Now(Config.Logger.TimeFormat)
	}

	// 设置日志输出方式
	var (
		output  io.Writer
		logFile *os.File
		err     error
	)
	// 设置日志文件
	if Config.Logger.FilePath != "" {
		// 输出到文件
		if Config.Logger.FileMode == 0 {
			Config.Logger.FileMode = os.FileMode(0600)
		}
		logFile, err = os.OpenFile(Config.Logger.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, Config.Logger.FileMode)
		if nil != err {
			return err
		}
	}
	switch Config.Logger.Encode {
	// console编码
	case "console":
		if logFile != nil {
			output = zerolog.ConsoleWriter{
				Out:        logFile,
				NoColor:    true,
				TimeFormat: zerolog.TimeFieldFormat,
			}
		} else {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				NoColor:    Config.Logger.NoColor,
				TimeFormat: zerolog.TimeFieldFormat,
			}
		}
	// json编码
	case "json":
		if logFile != nil {
			output = logFile
		} else {
			output = os.Stdout
		}
	default:
		return errors.New("logger.encode配置参数值只支持json和console")
	}

	log.Logger = log.Output(output)

	return nil
}
