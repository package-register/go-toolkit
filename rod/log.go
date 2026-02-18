package rodutil

import (
	"io"
	"os"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	instance *logrus.Logger
	mux      sync.Mutex
}

// NewLogger create logger func
func NewLogger(logFile string, logToFile, logToConsole bool) (*Logger, error) {
	if logFile == "" {
		logFile = "logs/app.log"
	}

	if !strings.HasSuffix(logFile, ".log") {
		logFile += ".log"
	}

	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false,                 // 是否启用颜色
		FullTimestamp:    true,                  // 是否显示完整时间戳
		DisableTimestamp: false,                 // 是否禁用时间戳
		TimestampFormat:  "2006-01-02 15:04:05", // 设置时间戳的格式
		PadLevelText:     false,                 // 不使用额外的缩进，去掉额外的间距
	})

	logger.SetLevel(logrus.InfoLevel)

	var logOutputs []io.Writer

	if logToFile {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logOutputs = append(logOutputs, file)
	}

	if logToConsole {
		logOutputs = append(logOutputs, os.Stdout)
	}

	if len(logOutputs) > 1 {
		logger.SetOutput(io.MultiWriter(logOutputs...))
	} else if len(logOutputs) == 1 {
		logger.SetOutput(logOutputs[0])
	}

	return &Logger{instance: logger}, nil
}

func (l *Logger) Info(msg string, fields ...logrus.Fields) {
	l.log(logrus.InfoLevel, msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...logrus.Fields) {
	l.log(logrus.WarnLevel, msg, fields...)
}

func (l *Logger) Error(msg string, fields ...logrus.Fields) {
	l.log(logrus.ErrorLevel, msg, fields...)
}

func (l *Logger) Debug(msg string, fields ...logrus.Fields) {
	l.log(logrus.DebugLevel, msg, fields...)
}

func (l *Logger) log(level logrus.Level, msg string, fields ...logrus.Fields) {
	l.mux.Lock()
	defer l.mux.Unlock()

	entry := l.instance.WithFields(logrus.Fields{})
	for _, f := range fields {
		for k, v := range f {
			entry = entry.WithField(k, v)
		}
	}

	switch level {
	case logrus.InfoLevel:
		entry.Info(msg)
	case logrus.WarnLevel:
		entry.Warn(msg)
	case logrus.ErrorLevel:
		entry.Error(msg)
	case logrus.DebugLevel:
		entry.Debug(msg)
	}
}
