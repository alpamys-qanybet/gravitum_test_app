package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type Logger struct {
	logger zerolog.Logger
}

func New(level zerolog.Level) *Logger {
	zerolog.SetGlobalLevel(level)
	return &Logger{
		logger: zerolog.New(os.Stdout).With().Timestamp().Logger(),
	}
}

func GetLevelByString(level string) zerolog.Level {
	level = strings.ToLower(level)

	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "":
		return zerolog.NoLevel
	}
	return zerolog.Disabled
}

func (c *Logger) Fatalf(format string, args ...interface{}) {
	c.logger.Fatal().Msgf(format, args...)
}

func (c *Logger) Errorf(format string, args ...interface{}) {
	c.logger.Error().Msgf(format, args...)
}

func (c *Logger) Warnf(format string, args ...interface{}) {
	c.logger.Warn().Msgf(format, args...)
}

func (c *Logger) Infof(format string, args ...interface{}) {
	c.logger.Info().Msgf(format, args...)
}

func (c *Logger) Debugf(format string, args ...interface{}) {
	c.logger.Debug().Msgf(format, args...)
}

func (c *Logger) Fatal(message string) {
	c.logger.Fatal().Msg(message)
}

func (c *Logger) Fatalln(message string) {
	c.logger.Fatal().Msg(message + "\n")
}

func (c *Logger) Error(message string) {
	c.logger.Error().Msg(message)
}

func (c *Logger) Warn(message string) {
	c.logger.Warn().Msg(message)
}

func (c *Logger) Info(message string) {
	c.logger.Info().Msg(message)
}

func (c *Logger) Debug(message string) {
	c.logger.Debug().Msg(message)
}
