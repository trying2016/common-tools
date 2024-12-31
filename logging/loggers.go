package logging

import (
	"bytes"
	"io"
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

// const
const (
	PanicLevel = "panic"
	FatalLevel = "fatal"
	ErrorLevel = "error"
	WarnLevel  = "warn"
	InfoLevel  = "info"
	DebugLevel = "debug"
	TraceLevel = "trace"
)
const (
	//PANIC log level
	PANIC uint32 = iota
	//FATAL has list msg
	FATAL
	//ERROR has list msg
	ERROR
	//WARN only log
	WARN
	//INFO only log
	INFO
	//DEBUG only log
	DEBUG
	//TRACE only log
	TRACE
)
const (
	//MsgFormatSingle use info
	MsgFormatSingle uint32 = iota
	//MsgFormatMulti use show all func call relation
	MsgFormatMulti
)

var (
	logFile *os.File
)

// LogFormat is to log format
type LogFormat = map[string]interface{}

type emptyWriter struct{}

func (ew emptyWriter) Write(p []byte) (int, error) {
	return 0, nil
}

type Logger struct {
	*logrus.Logger
	//CallRelation to show stack list
	CallRelation uint32
}

func NewLogger() *Logger {
	return &Logger{
		Logger: logrus.New(),
	}
}

// SetCallList to set CallList
func (logger *Logger) SetCallRelation(button uint32) {
	logger.CallRelation = button
}

// logger pointer must be initialized, else would panic.
var clog *Logger
var vlog *Logger

func convertLevel(level string) logrus.Level {
	switch level {
	case PanicLevel:
		return logrus.PanicLevel
	case FatalLevel:
		return logrus.FatalLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case WarnLevel:
		return logrus.WarnLevel
	case InfoLevel:
		return logrus.InfoLevel
	case DebugLevel:
		return logrus.DebugLevel
	case TraceLevel:
		return logrus.TraceLevel
	default:
		return logrus.InfoLevel
	}
}

// Init loggers
func Init(path, filename string, level string, age uint32, disableCPrint bool) {
	vlog = NewLogger()
	LoadFunctionHooker(vlog)

	var fileHooker logrus.Hook
	if path != "" {
		fileHooker = NewFileRotateHooker(path, filename, age, nil)
		vlog.Hooks.Add(fileHooker)
	}

	vlog.Out = &emptyWriter{}
	vlog.Formatter = &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04 05",
	}
	vlog.Level = convertLevel(level)

	if !disableCPrint {
		clog = NewLogger()
		LoadFunctionHooker(clog)
		if path != "" {
			clog.Hooks.Add(fileHooker)
		}
		clog.Out = os.Stdout
		clog.Formatter = &logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04 05",
		}
		clog.Level = convertLevel(level)
	} else {
		clog = vlog
	}

	vlog.WithFields(logrus.Fields{
		"path":  path,
		"level": level,
	}).Info("Logger Configuration.")
}

// InitV2 loggers
func InitV2(dir, filename string, level string, age uint32, disableCPrint bool) {
	file, err := os.OpenFile(path.Join(dir, filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	logFile = file
	// Set logrus to write to both the file and the console
	multiWriter := io.MultiWriter(file, os.Stdout)
	logrus.SetOutput(multiWriter)
	// Set logrus to write to the file

	// Set log level
	logrus.SetLevel(convertLevel(level))

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

// GetGID return gid
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// CPrint into stdout + log
func CPrint(level uint32, msg string, formats ...LogFormat) {
	if logFile != nil {
		data := mergeLogFormats(formats...)
		switch level {
		case PANIC:
			{
				logrus.WithFields(data).Panic(msg)
				break
			}
		case FATAL:
			{
				logrus.WithFields(data).Fatal(msg)
				break
			}
		case ERROR:
			{
				logrus.WithFields(data).Error(msg)
				break
			}
		case WARN:
			{
				logrus.WithFields(data).Warn(msg)
				break
			}
		case INFO:
			{
				logrus.WithFields(data).Info(msg)
				break
			}
		case DEBUG:
			{
				logrus.WithFields(data).Debug(msg)
				break
			}
		case TRACE:
			{
				logrus.WithFields(data).Trace(msg)
				break
			}
		default:
			{
				logrus.WithFields(data).Error(msg)
				return
			}
		}
		return
	}

	if clog == nil {
		Init("", "miner.log", "info", 0, false)
	}
	data := mergeLogFormats(formats...)
	switch level {
	case PANIC:
		{
			clog.SetCallRelation(MsgFormatMulti)
			clog.WithFields(data).Panic(msg)
			break
		}
	case FATAL:
		{
			clog.SetCallRelation(MsgFormatMulti)
			clog.WithFields(data).Fatal(msg)
			break
		}
	case ERROR:
		{
			clog.SetCallRelation(MsgFormatMulti)
			clog.WithFields(data).Error(msg)
			break
		}
	case WARN:
		{
			clog.SetCallRelation(MsgFormatSingle)
			clog.WithFields(data).Warn(msg)
			break
		}
	case INFO:
		{
			clog.SetCallRelation(MsgFormatSingle)
			clog.WithFields(data).Info(msg)
			break
		}
	case DEBUG:
		{
			clog.SetCallRelation(MsgFormatSingle)
			clog.WithFields(data).Debug(msg)
			break
		}
	case TRACE:
		{
			clog.SetCallRelation(MsgFormatSingle)
			clog.WithFields(data).Trace(msg)
			break
		}
	default:
		{
			clog.SetCallRelation(MsgFormatMulti)
			clog.WithFields(data).Error(msg)
		}
	}
}

// VPrint into log
func VPrint(level uint32, msg string, formats ...LogFormat) {
	if vlog == nil {
		Init("", "tmp-mass.log", "info", 0, false)
	}
	data := mergeLogFormats(formats...)
	switch level {
	case PANIC:
		{
			vlog.SetCallRelation(MsgFormatMulti)
			vlog.WithFields(data).Panic(msg)
			break
		}
	case FATAL:
		{
			vlog.SetCallRelation(MsgFormatMulti)
			vlog.WithFields(data).Fatal(msg)
			break
		}
	case ERROR:
		{
			vlog.SetCallRelation(MsgFormatMulti)
			vlog.WithFields(data).Error(msg)
			break
		}
	case WARN:
		{
			vlog.SetCallRelation(MsgFormatSingle)
			vlog.WithFields(data).Warn(msg)
			break
		}
	case INFO:
		{
			vlog.SetCallRelation(MsgFormatSingle)
			vlog.WithFields(data).Info(msg)
			break
		}
	case DEBUG:
		{
			vlog.SetCallRelation(MsgFormatSingle)
			vlog.WithFields(data).Debug(msg)
			break
		}
	case TRACE:
		{
			vlog.SetCallRelation(MsgFormatSingle)
			vlog.WithFields(data).Trace(msg)
			break
		}
	default:
		{
			vlog.SetCallRelation(MsgFormatMulti)
			vlog.WithFields(data).Error(msg)
		}
	}
}

// mergeLogFormats merges LogFormats.
// Same key would be covered by later-presented values.
func mergeLogFormats(formats ...LogFormat) LogFormat {
	format := LogFormat{}
	for _, data := range formats {
		if data == nil {
			continue
		}
		for k, v := range data {
			vv := v
			format[k] = vv
		}
	}
	//format["tid"] = GetGID()
	return format
}
