package pkg_log

import (
	"GoRestify/pkg/pkg_consts"
	"os"

	"github.com/sirupsen/logrus"
)

// Init initiate the global logger
func Init(output string, isDebug bool) {

	level := "trace"
	if isDebug {
		level = "debug"
	}

	serverLogParam := logParam{
		format:     pkg_consts.LogFormat,
		output:     output,
		level:      level,
		JSONIndent: false,
	}

	logger = initLog(serverLogParam)
}

// New return a pointer to initiated logger
func New(format, output, level string, indent, file bool) *logrus.Logger {
	serverLogParam := logParam{
		format:     format,
		output:     output,
		level:      level,
		JSONIndent: indent,
	}

	return initLog(serverLogParam)
}

func initLog(p logParam) *logrus.Logger {
	log := logrus.New()

	setFormat(log, p)
	setOutput(log, p)
	setLevel(log, p)

	return log
}

func setFormat(log *logrus.Logger, p logParam) {
	switch p.format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: p.JSONIndent,
		})
	}
}

func setOutput(log *logrus.Logger, p logParam) {
	switch p.output {
	case "stdout":
		log.SetOutput(os.Stdout)
	default:
		file, err := os.OpenFile(p.output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			log.Out = file
		} else {
			log.Fatalf("failed to write logs to file %q, [core/logs.go]", p.output)
		}
	}
}

func setLevel(log *logrus.Logger, p logParam) {

	switch p.level {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	}
}
