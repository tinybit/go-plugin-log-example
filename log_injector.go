package main

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/rs/zerolog"
)

const (
	IgnoredLogLine = "received EOF, stopping recv loop: err=\"rpc error: code = Unavailable desc = error reading from server: EOF\""
)

type LogInjector struct {
	hclog.Logger
	lg   *zerolog.Logger
	opts hclog.LoggerOptions
}

func NewLogInjector(logger *zerolog.Logger) *LogInjector {
	hcLogLevel := LogHCLevelFromZerologLevel(logger.GetLevel())

	opts := hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hcLogLevel,
	}

	baseLogger := hclog.New(&opts)

	shedulerLogger := logger.With().Str("app", MainProcessLogLabel).Logger()

	wrapper := &LogInjector{
		baseLogger,
		&shedulerLogger,
		opts,
	}

	return wrapper
}

func (l *LogInjector) Log(level hclog.Level, msg string, args ...interface{}) {
	levelLocal := level
	if levelLocal == hclog.NoLevel {
		levelLocal = hclog.DefaultLevel
	}

	logStr := l.renderLogLineToString(level, msg, args...)

	if logStr == IgnoredLogLine {
		return
	}

	switch levelLocal {
	case hclog.NoLevel:
		l.lg.Info().Msg(logStr)

	case hclog.Trace:
		l.lg.Trace().Msg(logStr)

	case hclog.Debug:
		l.lg.Debug().Msg(logStr)

	case hclog.Info:
		l.lg.Info().Msg(logStr)

	case hclog.Warn:
		l.lg.Warn().Msg(logStr)

	case hclog.Error:
		l.lg.Error().Msg(logStr)
	}
}

func (l *LogInjector) Trace(msg string, args ...interface{}) {
	l.Log(hclog.Trace, msg, args...)
}

func (l *LogInjector) Debug(msg string, args ...interface{}) {
	l.Log(hclog.Debug, msg, args...)
}

func (l *LogInjector) Info(msg string, args ...interface{}) {
	l.Log(hclog.Info, msg, args...)
}

func (l *LogInjector) Warn(msg string, args ...interface{}) {
	l.Log(hclog.Warn, msg, args...)
}

func (l *LogInjector) Error(msg string, args ...interface{}) {
	l.Log(hclog.Error, msg, args...)
}

func (l *LogInjector) Named(name string) hclog.Logger {
	namedWrapper := &LogInjector{
		l.Logger,
		l.lg,
		l.opts,
	}

	return namedWrapper
}

func (l *LogInjector) renderLogLineToString(level hclog.Level, msg string, args ...interface{}) string {
	// explanation:
	// some log strings come in this form:
	// -2023-11-15T10:54:22.783+0800 [DEBUG] plugin: 2023-11-15T10:54:22.783+0800 [DEBUG] plugin: plugin address: network=unix address=/var/folders/xz/_vsfj__d63qc1x2n3lg6646m0000gn/T/plugin2160303011
	//
	// we need to remove all occurences of prefixes in this form:
	// - 2023-11-15T10:54:22.783+0800 [DEBUG]
	// - plugin: 2023-11-15T10:54:22.783+0800 [DEBUG]
	// {anything} [LOG_LEVEL]
	//
	// leaving only log message itself:
	// - plugin: plugin address: network=unix address=/var/folders/xz/_vsfj__d63qc1x2n3lg6646m0000gn/T/plugin2160303011
	//
	// from this line we need to drop prefix "plugin:", resulting in a clean log line:
	// - plugin address: network=unix address=/var/folders/xz/_vsfj__d63qc1x2n3lg6646m0000gn/T/plugin2160303011

	// render log entry to string
	writer := &StringWriter{}
	optsLoc := l.opts
	optsLoc.Output = writer
	hcLogger := hclog.New(&optsLoc)
	hcLogger.Log(level, msg, args...)
	logStr := writer.String()

	levelStr := strings.ToUpper(level.String())

	for {
		foundLevel := false

		if len(logStr) < len(levelStr) {
			break
		}

		// drop prefix (time + log level)
		pos := strings.Index(logStr, levelStr)

		if pos != -1 && pos+1 < len(logStr) {
			logStr = strings.TrimSpace(logStr[pos+len(levelStr)+1:])
			foundLevel = true
		}

		prefixToDrop := "plugin:"
		prefix := logStr[:len(prefixToDrop)]

		if len(logStr) >= len(prefixToDrop) && strings.ToLower(prefix) == prefixToDrop {
			logStr = logStr[len(prefixToDrop):]
			logStr = strings.TrimSpace(logStr)
		}

		if !foundLevel {
			break
		}
	}

	return logStr
}

type StringWriter struct {
	builder strings.Builder
}

func (w *StringWriter) Write(p []byte) (n int, err error) {
	return w.builder.Write(p)
}

func (w *StringWriter) WriteString(s string) (n int, err error) {
	return w.builder.WriteString(s)
}

func (w *StringWriter) WriteByte(b byte) (err error) {
	return w.builder.WriteByte(b)
}

func (w *StringWriter) String() string {
	return w.builder.String()
}

// CustomWriter defines a custom writer type.
type StderrToLogWriter struct {
	lg *zerolog.Logger
}

func NewStderrToLogWriter(logger *zerolog.Logger) *StderrToLogWriter {
	pluginLogger := logger.With().Str("app", "plugin").Logger()
	return &StderrToLogWriter{&pluginLogger}
}

func (s *StderrToLogWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))
	s.lg.Debug().Msg(msg)

	return len(p), nil
}

func LogHCLevelFromZerologLevel(lev zerolog.Level) hclog.Level {
	switch lev {
	case zerolog.PanicLevel:
		fallthrough
	case zerolog.FatalLevel:
		fallthrough
	case zerolog.Disabled:
		fallthrough
	case zerolog.ErrorLevel:
		return hclog.Error

	case zerolog.NoLevel:
		fallthrough
	case zerolog.TraceLevel:
		return hclog.Trace

	case zerolog.DebugLevel:
		return hclog.Debug
	case zerolog.InfoLevel:
		return hclog.Info
	case zerolog.WarnLevel:
		return hclog.Warn

	default:
		return hclog.Trace
	}
}
