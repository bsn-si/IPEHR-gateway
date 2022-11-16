package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/onrik/logrus/filename"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var DefaultLogger *ServiceLogger

func init() {
	DefaultLogger = newServiceLogger()
}

type Formatter string

const (
	FormatterText Formatter = "text"
	FormatterJSON Formatter = "json"
)

type Fields map[string]interface{}

type Level uint32

func (l Level) String() string {
	return logrus.Level(l).String()
}

const (
	PanicLevel = Level(logrus.PanicLevel)
	FatalLevel = Level(logrus.FatalLevel)
	ErrorLevel = Level(logrus.ErrorLevel)
	WarnLevel  = Level(logrus.WarnLevel)
	InfoLevel  = Level(logrus.InfoLevel)
	DebugLevel = Level(logrus.DebugLevel)

	envLogLevel = "LOG_LEVEL"
)

type Logger interface {
	WithField(key string, value interface{}) *ServiceLogger
	WithFields(fields Fields) *ServiceLogger
	WithError(err error) *ServiceLogger
	WithContext(ctx context.Context) *ServiceLogger

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Log(logLevel Level, args ...interface{})

	Debugln(args ...interface{})
	Infoln(args ...interface{})
	Println(args ...interface{})
	Warnln(args ...interface{})
	Errorln(args ...interface{})
	Fatalln(args ...interface{})
	Panicln(args ...interface{})
	Logln(logLevel Level, args ...interface{})
}

type ServiceLogger struct{ entry *logrus.Entry }

func newServiceLogger() *ServiceLogger {
	logger := logrus.New()
	fnHook := filename.NewHook()
	fnHook.Field = "file"
	fnHook.Skip = 8

	fnHook.SkipPrefixes = append(fnHook.SkipPrefixes, "logging/", "logrus/", "logrus@", "gin@v1.7.7/")
	logger.AddHook(fnHook)
	ret := &ServiceLogger{entry: logger.WithFields(nil)}
	ret.SetFormatter(FormatterText)
	level, err := ParseLevel(os.Getenv(envLogLevel)) // TODO replace it
	if err == nil {
		ret.SetLevel(level)
	}
	return ret
}

func ParseLevel(lvl string) (Level, error) {
	level, err := logrus.ParseLevel(lvl)
	return Level(level), err
}

func (l *ServiceLogger) SetLevel(level Level) {
	l.entry.Logger.SetLevel(logrus.Level(level))
}

func (l *ServiceLogger) SetFormatter(ftype Formatter) {
	switch ftype {
	case FormatterText:
		l.entry.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	case FormatterJSON:
		l.entry.Logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyMsg:  "message",
				logrus.FieldKeyTime: "@timestamp",
			},
		})
	}
}

func (l *ServiceLogger) GetLevel() Level {
	return Level(l.entry.Logger.GetLevel())
}

func (l *ServiceLogger) SetOutput(writer io.Writer) {
	l.entry.Logger.SetOutput(writer)
}

func (l *ServiceLogger) WithField(key string, value interface{}) *ServiceLogger {
	return &ServiceLogger{entry: l.entry.WithField(key, value)}
}

func (l *ServiceLogger) WithFields(fields Fields) *ServiceLogger {
	return &ServiceLogger{entry: l.entry.WithFields(logrus.Fields(fields))}
}

func (l *ServiceLogger) WithError(err error) *ServiceLogger {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	log := &ServiceLogger{entry: l.entry.WithError(err)}
	if stackErr, ok := err.(stackTracer); ok {
		log.entry = log.entry.WithField("stacktrace", fmt.Sprintf("%+v", stackErr.StackTrace()))
	}

	return log
}

func (l *ServiceLogger) WithContext(ctx context.Context) *ServiceLogger {
	if ctx == nil {
		return l
	}
	val := ctx.Value(fieldsKey{})
	if val == nil {
		return l
	}
	fields, _ := ctx.Value(fieldsKey{}).(Fields)
	return &ServiceLogger{entry: l.entry.WithFields(logrus.Fields(fields))}
}

func Printf(format string, args ...interface{}) { DefaultLogger.Printf(format, args...) }
func Println(args ...interface{})               { DefaultLogger.Println(args...) }
func Debug(args ...interface{})                 { DefaultLogger.Debug(args...) }
func Info(args ...interface{})                  { DefaultLogger.Info(args...) }
func Print(args ...interface{})                 { DefaultLogger.Print(args...) }
func Warn(args ...interface{})                  { DefaultLogger.Warn(args...) }
func Error(args ...interface{})                 { func() { DefaultLogger.Error(args...) }() }
func Fatal(args ...interface{})                 { DefaultLogger.Fatal(args...) }
func Panic(args ...interface{})                 { DefaultLogger.Panic(args...) }
func Log(logLevel Level, args ...interface{}) {
	DefaultLogger.Log(Level(logLevel), args...)
}
func WithError(err error) *ServiceLogger             { return DefaultLogger.WithError(err) }
func WithContext(ctx context.Context) *ServiceLogger { return DefaultLogger.WithContext(ctx) }
func WithFields(fields Fields) *ServiceLogger        { return DefaultLogger.WithFields(fields) }

func (l *ServiceLogger) Printf(format string, args ...interface{}) { l.entry.Printf(format, args...) }
func (l *ServiceLogger) Debug(args ...interface{})                 { l.entry.Debug(args...) }
func (l *ServiceLogger) Info(args ...interface{})                  { l.entry.Info(args...) }
func (l *ServiceLogger) Print(args ...interface{})                 { l.entry.Print(args...) }
func (l *ServiceLogger) Warn(args ...interface{})                  { l.entry.Warn(args...) }
func (l *ServiceLogger) Error(args ...interface{})                 { l.entry.Error(args...) }
func (l *ServiceLogger) Fatal(args ...interface{})                 { l.entry.Fatal(args...) }
func (l *ServiceLogger) Panic(args ...interface{})                 { l.entry.Panic(args...) }
func (l *ServiceLogger) Log(logLevel Level, args ...interface{}) {
	l.entry.Log(logrus.Level(logLevel), args...)
}

func (l *ServiceLogger) Debugln(args ...interface{}) { l.entry.Debugln(args...) }
func (l *ServiceLogger) Infoln(args ...interface{})  { l.entry.Infoln(args...) }
func (l *ServiceLogger) Println(args ...interface{}) { l.entry.Println(args...) }
func (l *ServiceLogger) Warnln(args ...interface{})  { l.entry.Warnln(args...) }
func (l *ServiceLogger) Errorln(args ...interface{}) { l.entry.Errorln(args...) }
func (l *ServiceLogger) Fatalln(args ...interface{}) { l.entry.Fatalln(args...) }
func (l *ServiceLogger) Panicln(args ...interface{}) { l.entry.Panicln(args...) }
func (l *ServiceLogger) Logln(logLevel Level, args ...interface{}) {
	l.entry.Logln(logrus.Level(logLevel), args...)
}

// Assert that ServiceLogger implements the Logger interface.
var _ Logger = (*ServiceLogger)(nil)

// -----

type fieldsKey struct{}

func (l *ServiceLogger) ContextWithFields(c context.Context, fields Fields) context.Context {
	return ContextWithFields(c, fields)
}

// ContextWithFields adds logger fields to fields in context
func ContextWithFields(parent context.Context, fields Fields) context.Context {
	var newFields Fields
	val := parent.Value(fieldsKey{})
	if val == nil {
		newFields = fields
	} else {
		newFields = make(Fields)
		oldFields, _ := val.(Fields)
		for k, v := range oldFields {
			newFields[k] = v
		}
		for k, v := range fields {
			newFields[k] = v
		}
	}

	return context.WithValue(parent, fieldsKey{}, newFields)
}

// ContextWithField is like ContextWithFields but adds only one field
func ContextWithField(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithFields(ctx, Fields{key: value})
}

// FieldsFromContext returns logging fields from the context
func FieldsFromContext(ctx context.Context) Fields {
	if ctx == nil {
		return nil
	}
	fields, _ := ctx.Value(fieldsKey{}).(Fields)
	return fields
}
