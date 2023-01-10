package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Logger struct {
	config zap.Config
	base   *zap.Logger
	*zap.SugaredLogger
}

func (l *Logger) Init() {
	l.config = zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			MessageKey:     "msg",
			CallerKey:      "",
			StacktraceKey:  "",
			FunctionKey:    zapcore.OmitKey,
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths: []string{
			"stdout",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}
	l.base, _ = l.config.Build()
	l.SugaredLogger = l.base.Sugar()
}

func (l *Logger) SwitchLevel(level string) {
	if strings.EqualFold(l.config.Level.String(), level) {
		l.Infow("switching logging level", "not needed, same")
		return
	}

	l.Infow("switching logging level", "from", l.config.Level.String(), "to", level)

	err := l.config.Level.UnmarshalText([]byte(level))
	if err != nil {
		l.Fatalw("failed to switch level", "from", l.config.Level, "to", level)
	}
}
