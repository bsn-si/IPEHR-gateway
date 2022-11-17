package log

import (
	"context"
	"io"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	LogLevel  string    `json:"log_level" mapstructure:"log_level" default:"debug"`
	Formatter Formatter `json:"formatter,omitempty" mapstructure:"formatter" default:"text"`
}

var DefaultLogger *ServiceLogger

func init() {
	pflag.String("log_level", "debug", "set up logger level")
	pflag.String("formatter", "text", "set up logger messages format")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		Panic(err)
	}

	DefaultLogger = NewLoggerWithConfig(cfg)
}

func NewLoggerWithConfig(cfg Config) *ServiceLogger {
	srvLog := NewLog()
	srvLog.SetFormatter(cfg.Formatter)

	level, err := ParseLevel(cfg.LogLevel)
	if err == nil {
		srvLog.SetLevel(level)
	}

	return srvLog
}

func NewLog() *ServiceLogger {
	log := logrus.New()

	srvLog := &ServiceLogger{
		entry: log.WithFields(nil),
	}

	return srvLog
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func GetLevel() Level {
	return DefaultLogger.GetLevel()
}

func SetOutput(writer io.Writer) {
	DefaultLogger.SetOutput(writer)
}

func SetDefaultFields(fields Fields) {
	DefaultLogger.entry.Data = logrus.Fields(fields)
}

func SetFormatter(ftype Formatter) {
	DefaultLogger.SetFormatter(ftype)
}

func ParseLevel(lvl string) (Level, error) {
	level, err := logrus.ParseLevel(lvl)
	return Level(level), err
}

func Printf(format string, args ...interface{}) { DefaultLogger.Printf(format, args...) }

func Debug(args ...interface{})                 { DefaultLogger.Debug(args...) }
func Info(args ...interface{})                  { DefaultLogger.Info(args...) }
func Print(args ...interface{})                 { DefaultLogger.Print(args...) }
func Warn(args ...interface{})                  { DefaultLogger.Warn(args...) }
func Error(args ...interface{})                 { DefaultLogger.Error(args...) }
func Fatal(args ...interface{})                 { DefaultLogger.Fatal(args...) }
func Fatalf(format string, args ...interface{}) { DefaultLogger.Fatalf(format, args...) }
func Panic(args ...interface{})                 { DefaultLogger.Panic(args...) }

func Debugln(args ...interface{}) { DefaultLogger.Debugln(args...) }
func Infoln(args ...interface{})  { DefaultLogger.Infoln(args...) }
func Println(args ...interface{}) { DefaultLogger.Println(args...) }
func Warnln(args ...interface{})  { DefaultLogger.Warnln(args...) }
func Errorln(args ...interface{}) { DefaultLogger.Errorln(args...) }
func Fatalln(args ...interface{}) { DefaultLogger.Fatalln(args...) }
func Panicln(args ...interface{}) { DefaultLogger.Panicln(args...) }

func WithError(err error) *ServiceLogger             { return DefaultLogger.WithError(err) }
func WithContext(ctx context.Context) *ServiceLogger { return DefaultLogger.WithContext(ctx) }
func WithFields(fields Fields) *ServiceLogger        { return DefaultLogger.WithFields(fields) }

func WithField(key string, value interface{}) *ServiceLogger {
	return DefaultLogger.WithField(key, value)
}
