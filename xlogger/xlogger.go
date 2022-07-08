package xlogger

type Logger interface {
	Debug(msg string, kvs ...interface{})
	Info(msg string, kvs ...interface{})
	Warn(msg string, kvs ...interface{})
	Error(msg string, kvs ...interface{})
	Fatal(msg string, kvs ...interface{})
}
