package logger

import (
	"errors"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	// DebugLevel has verbose message
	DebugLevel = "debug"

	// InfoLevel is default log level
	InfoLevel = "info"

	// WarnLevel is for logging messages about possible issues
	WarnLevel = "warn"

	// ErrorLevel is for logging errors
	ErrorLevel = "error"

	// FatalLevel is for logging fatal messages.
	// The system shutsdown after logging the message.
	FatalLevel = "fatal"
)

// logger instance
const (
	InstanceZapLogger int = iota
)

// default config log timestamp
const (
	TimestampFormat    = "15:04:05.999 02/01/2006 (UTC)"
	DEFAULTLOGFILEPATH = "./github-audit.log"
	DEFAULTLOGLEVEL    = InfoLevel
)

var (
	log                      Logger
	once                     sync.Once
	config                   Configuration
	ConfigFile               string
	errInvalidLoggerInstance = errors.New("Invalid logger instance")
)

// Logger is our contract for the logger
type Logger interface {
	Print(...interface{})

	Printf(string, ...interface{})

	Println(...interface{})

	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})

	Debug(args ...interface{})

	Info(args ...interface{})

	Warn(args ...interface{})

	Error(args ...interface{})

	Fatal(args ...interface{})

	Panic(args ...interface{})

	WithFields(keyValues Fields) Logger
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers,
// for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
}

// NewLogger returns an instance of logger
func NewLogger(config Configuration, loggerInstance int) error {
	switch loggerInstance {
	case InstanceZapLogger:
		logger, err := newZapLogger(config)
		if err != nil {
			return err
		}
		log = logger
		return nil
	default:
		return errInvalidLoggerInstance
	}
}

// Printf to satisfy std logger
func Printf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Print to satisfy std logger
func Print(args ...interface{}) {
	log.Debug(args...)
}

// Println to satisfy std logger
func Println(args ...interface{}) {
	log.Debug(args...)
}

// Debugf ...
func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

// Infof ...
func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

// Warnf ...
func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

// Errorf ...
func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

// Fatalf ...
func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

// Panicf ...
func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

// Debug ...
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info ...
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn ...
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error ...
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal ...
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// Panic ...
func Panic(args ...interface{}) {
	log.Panic(args...)
}

// WithFields ...
func WithFields(keyValues Fields) Logger {
	return log.WithFields(keyValues)
}

// DefaultLoggerConfig is used if default config is not provided
func DefaultLoggerConfig() Configuration {
	return Configuration{
		EnableConsole:     false,
		ConsoleLevel:      DEFAULTLOGLEVEL,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         DEFAULTLOGLEVEL,
		FileJSONFormat:    false,
		FileLocation:      DEFAULTLOGFILEPATH,
	}
}

func LoadLogConfig(filepath string) (Configuration, error) {
	type Config struct {
		LogPath  string `json:"logpath" yaml:"logpath"`
		LogLevel string `json:"loglevel" yaml:"loglevel"`
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return Configuration{}, err
	}
	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Configuration{}, err
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = DEFAULTLOGLEVEL
	}
	if cfg.LogPath == "" {
		cfg.LogPath = DEFAULTLOGFILEPATH
	}

	logconfig := Configuration{
		EnableConsole:     false,
		ConsoleLevel:      cfg.LogLevel,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         cfg.LogLevel,
		FileJSONFormat:    false,
		FileLocation:      cfg.LogPath,
	}

	return logconfig, nil
}

// GetLogger return logger instance
// initialize global logger only onces
func GetLogger() Logger {
	once.Do(
		func() {
			var err error
			config, err = LoadLogConfig(ConfigFile)
			if err != nil {
				config = DefaultLoggerConfig()
			}
			err = NewLogger(config, InstanceZapLogger)
			if err != nil {
				Panic("Could not instantiate log %s", err.Error())
			}
		},
	)
	return log
}
