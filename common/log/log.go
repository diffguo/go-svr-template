package log

import (
	"fmt"
	"go-svr-template/common/goroutineid"
	"log"
	"os"
	"strings"
)

// 没有日志切割，切割任务丢给shell脚本

var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	Reset   = string([]byte{27, 91, 48, 109})
)

func ColorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

func ColorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return Reset
	}
}

func ColorForReset() string {
	return Reset
}

const (
	LogLevelError = 1 << iota
	LogLevelWarn
	LogLevelInfo
	LogLevelDebug
)

var GLog *SimpleLog

type SimpleLog struct {
	LogLevel int
	Log      *log.Logger
}

func InitLog(logDir string, logFile string, logStrLevel string) (*SimpleLog, error) {
	logStrLevel = strings.ToLower(logStrLevel)

	if logStrLevel != "debug" && logStrLevel != "info" && logStrLevel != "warn" && logStrLevel != "error" {
		return nil, fmt.Errorf("wrong log level")
	}

	simpleLog := SimpleLog{}

	if logStrLevel == "debug" {
		simpleLog.LogLevel = LogLevelDebug
	}

	if logStrLevel == "info" {
		simpleLog.LogLevel = LogLevelInfo
	}

	if logStrLevel == "warn" {
		simpleLog.LogLevel = LogLevelWarn
	}

	if logStrLevel == "error" {
		simpleLog.LogLevel = LogLevelError
	}

	filename := logDir + "/" + logFile
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	simpleLog.Log = log.New(f, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	if GLog == nil {
		GLog = &simpleLog
	}

	return &simpleLog, nil
}

func (slog *SimpleLog) Debug(v ...interface{}) {
	if slog.LogLevel == LogLevelDebug {
		slog.Log.Output(3, fmt.Sprintf("[DEBUG][%d] %s", goroutineid.GetGoID(), fmt.Sprint(v...)))
	}
}

func (slog *SimpleLog) Debugf(format string, args ...interface{}) {
	if slog.LogLevel == LogLevelDebug {
		msg := fmt.Sprintf(format, args...)
		slog.Log.Output(3, fmt.Sprintf("[DEBUG][%d] %s", goroutineid.GetGoID(), msg))
	}
}

func (slog *SimpleLog) Info(v ...interface{}) {
	if slog.LogLevel >= LogLevelInfo {
		slog.Log.Output(3, fmt.Sprintf("[INFO][%d] %s", goroutineid.GetGoID(), fmt.Sprint(v...)))
	}
}

func (slog *SimpleLog) Infof(format string, args ...interface{}) {
	if slog.LogLevel >= LogLevelInfo {
		msg := fmt.Sprintf(format, args...)
		slog.Log.Output(3, fmt.Sprintf("[INFO][%d] %s", goroutineid.GetGoID(), msg))
	}
}

func (slog *SimpleLog) Warn(v ...interface{}) {
	if slog.LogLevel >= LogLevelWarn {
		slog.Log.Output(3, fmt.Sprintf("[WARN][%d] %s", goroutineid.GetGoID(), fmt.Sprint(v...)))
	}
}

func (slog *SimpleLog) Warnf(format string, args ...interface{}) {
	if slog.LogLevel >= LogLevelWarn {
		msg := fmt.Sprintf(format, args...)
		slog.Log.Output(3, fmt.Sprintf("[WARN][%d] %s", goroutineid.GetGoID(), msg))
	}
}

func (slog *SimpleLog) Error(v ...interface{}) {
	if slog.LogLevel >= LogLevelError {
		slog.Log.Output(3, fmt.Sprintf("[ERRO][%d] %s", goroutineid.GetGoID(), fmt.Sprint(v...)))
	}
}

func (slog *SimpleLog) Errorf(format string, args ...interface{}) {
	if slog.LogLevel >= LogLevelError {
		msg := fmt.Sprintf(format, args...)
		slog.Log.Output(3, fmt.Sprintf("[ERRO][%d] %s", goroutineid.GetGoID(), msg))
	}
}

func Debug(v ...interface{}) {
	GLog.Debug(v...)
}

func Debugf(format string, args ...interface{}) {
	GLog.Debugf(format, args...)
}

func Info(v ...interface{}) {
	GLog.Info(v...)
}

func Infof(format string, args ...interface{}) {
	GLog.Infof(format, args...)
}

func Warn(v ...interface{}) {
	GLog.Warn(v...)
}

func Warnf(format string, args ...interface{}) {
	GLog.Warnf(format, args...)
}

func Error(v ...interface{}) {
	GLog.Error(v...)
}

func Errorf(format string, args ...interface{}) {
	GLog.Errorf(format, args...)
}
