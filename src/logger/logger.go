package logger

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"zues/src/config"
)
// copy to https://blog.csdn.net/qq_23113053/article/details/84101560




func init(){
	log.SetLevel(parseLogLevel(config.Cfg.LogConfig.Level))
	log.SetOutput(&lumberjack.Logger{
		Filename:   config.Cfg.LogConfig.FileName,
		MaxSize:    config.Cfg.LogConfig.MaxSize, // megabytes
		MaxBackups: config.Cfg.LogConfig.MaxBackups,
		MaxAge:     config.Cfg.LogConfig.MaxAge, //days
		Compress:   config.Cfg.LogConfig.Compress, // disabled by default
	})
	log.SetFormatter(&Formatter{})
}

func parseLogLevel(logLevel string) log.Level {
	switch logLevel {
	case "DEBUG":
		return log.DebugLevel
	case "INFO":
		return log.InfoLevel
	case "ERROR":
		return log.ErrorLevel
	case "FATAL":
		return log.FatalLevel
	case "PANIC":
		return log.PanicLevel
	case "TRACE":
		return log.TraceLevel
	case "WARN":
		return log.WarnLevel
	default:
		panic("日志级别错误！")
	}
}


func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatal(format, args)
}

