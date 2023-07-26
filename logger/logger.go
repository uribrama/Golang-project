package logger

import (
	"encoding/json"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(isDebug bool) *Logger {
	logLevel := zerolog.InfoLevel
	if isDebug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &Logger{logger: &logger}
}

func (l *Logger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

func (l *Logger) Info(msg ...interface{}) {
	info := plainData(msg)
	l.logger.Info().Msg(info)
}

func (l *Logger) Debug(msg ...interface{}) {
	info := plainData(msg)
	l.logger.Debug().Msg(info)
}

func (l *Logger) Warn(msg ...interface{}) {
	info := plainData(msg)
	l.logger.Warn().Msg(info)
}

func (l *Logger) Error(err error) {
	l.logger.Error().Msg(err.Error())
}

func (l *Logger) ErrorL(msg ...interface{}) {
	info := plainData(msg)
	l.logger.Error().Msg(info)
}

func (l *Logger) Fatal(err error) {
	l.logger.Fatal().Msg(err.Error())
}

func plainData(data ...interface{}) string {
	var info string
	for _, d := range data {
		if out, err := json.Marshal(d); err == nil {
			info += string(out)
		}
	}
	return info
}
