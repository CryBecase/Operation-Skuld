package xlogger

import (
	"errors"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	stdout *zap.SugaredLogger
	stderr *zap.SugaredLogger
}

func (z *ZapLogger) Debug(msg string, kvs ...interface{}) {
	z.stdout.Debugw(msg, kvs...)
}

func (z *ZapLogger) Info(msg string, kvs ...interface{}) {
	z.stdout.Infow(msg, kvs...)
}

func (z *ZapLogger) Warn(msg string, kvs ...interface{}) {
	z.stdout.Warnw(msg, kvs...)
}

func (z *ZapLogger) Error(msg string, kvs ...interface{}) {
	z.stderr.Errorw(msg, kvs...)
}

func (z *ZapLogger) Fatal(msg string, kvs ...interface{}) {
	z.stderr.Fatalw(msg, kvs...)
}

func NewZapLogger(kv ...string) (Logger, error) {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	stdoutCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)
	stderrCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr)),
		zap.InfoLevel,
	)

	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Development(),
		zap.AddStacktrace(zap.ErrorLevel),
	}
	if len(kv)%2 == 1 {
		return nil, errors.New("there is no one-to-one correspondence between key and value")
	}
	fields := make([]zap.Field, 0)
	for i := 0; i < len(kv); i += 2 {
		fields = append(fields, zap.String(kv[i], kv[i+1]))
	}
	if len(fields) > 0 {
		options = append(options, zap.Fields(fields...))
	}

	return &ZapLogger{
		stdout: zap.New(stdoutCore, options...).Sugar(),
		stderr: zap.New(stderrCore, options...).Sugar(),
	}, nil
}
