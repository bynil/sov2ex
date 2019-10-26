package log

import (
	"io"
	"os"
	"time"

	"github.com/bynil/sov2ex/pkg/config"

	"github.com/arthurkiller/rollingwriter"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	logger, _ = zap.NewDevelopment()
}

var logger *zap.Logger
var Level = zap.NewAtomicLevelAt(zap.DebugLevel)

func Sync() {
	logger.Sync()
}

func GetLogger() *zap.Logger {
	return logger
}

// Debug uses fmt.Sprint to construct and log a message.
func Debug(args ...interface{}) {
	logger.Sugar().Debug(args)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	logger.Sugar().Info(args)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	logger.Sugar().Warn(args)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	logger.Sugar().Error(args)
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanic(args ...interface{}) {
	logger.Sugar().DPanic(args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func Panic(args ...interface{}) {
	logger.Sugar().Panic(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	logger.Sugar().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	logger.Sugar().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	logger.Sugar().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func DPanicf(template string, args ...interface{}) {
	logger.Sugar().DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	logger.Sugar().Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	logger.Sugar().Fatalf(template, args...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-Level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Debugw(msg, keysAndValues...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Infow(msg, keysAndValues...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Warnw(msg, keysAndValues...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Errorw(msg, keysAndValues...)
}

// DPanicw logs a message with some additional context. In development, the
// logger then panics. (See DPanicLevel for details.) The variadic key-value
// pairs are treated as they are in With.
func DPanicw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().DPanicw(msg, keysAndValues...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Panicw(msg, keysAndValues...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Fatalw(msg, keysAndValues...)
}

func InitLog() {
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Local().Format("2006-01-02 15:04:05.000"))
	}
	encoderCfg := zapcore.EncoderConfig{
		NameKey:        "Name",
		StacktraceKey:  "Stack",
		MessageKey:     "Msg",
		LevelKey:       "Level",
		TimeKey:        "Time",
		CallerKey:      "Caller",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     timeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var (
		enc    zapcore.Encoder
		writer io.Writer
	)

	if config.C.LogStdout {
		enc = zapcore.NewConsoleEncoder(encoderCfg)
		writer = os.Stdout
	} else {
		enc = zapcore.NewJSONEncoder(encoderCfg)
		conf := rollingwriter.Config{
			LogPath:                config.C.LogDir,
			TimeTagFormat:          "060102150405",
			FileName:               "sov2ex",
			MaxRemain:              365,
			RollingPolicy:          rollingwriter.TimeRolling,
			RollingTimePattern:     "0 0 0 * * *",
			RollingVolumeSize:      "500M",
			WriterMode:             "lock",
			BufferWriterThershould: 8 * 1024 * 1024,
			Compress:               false,
		}
		var err error
		writer, err = rollingwriter.NewWriterFromConfig(&conf)
		if err != nil {
			panic(errors.Errorf("initialize log writer error: %s", err))
		}
	}

	var zapOpts []zap.Option
	if config.C.Debug {
		zapOpts = append(zapOpts, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	zapOpts = append(zapOpts, zap.AddCallerSkip(1))
	sink := zapcore.AddSync(writer)
	logger = zap.New(
		zapcore.NewCore(enc, sink, Level),
		zapOpts...,
	)
	zap.RedirectStdLog(logger)
}
