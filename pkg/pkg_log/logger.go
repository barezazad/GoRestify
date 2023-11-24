package pkg_log

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

// Debug print struct with details with logrus ability
func Debug(objs ...interface{}) {
	for _, v := range objs {
		parts := make(map[string]interface{}, 2)
		parts["type"] = fmt.Sprintf("%T", v)
		parts["value"] = v
		dataInJSON, _ := json.Marshal(parts)

		logger.Debug(string(dataInJSON))
	}
}

// CheckError print all errors which happened inside the services, mainly they just have
// an error and a message
func CheckError(err error, message string, data ...interface{}) {
	if err != nil {
		LogError(err, message, data...)
	}
}

// LogError record an error with message and variadic parameters
func LogError(err error, message string, data ...interface{}) {
	if data == nil {
		logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Error(message)
	} else {
		logger.WithFields(logrus.Fields{
			"err":  err.Error(),
			"data": fmt.Sprintf("%+v", data),
		}).Error(message)

	}
}

// CheckInfo record the info
func CheckInfo(err error, message string, data ...interface{}) {
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Info(message)
		if data != nil {
			logger.Debug(data...)
		}
	}
}

// Info is information
func Info(data ...interface{}) {
	logger.Info(data...)
}

// Error is error
func Error(data ...interface{}) {
	logger.Error(data...)
}

// Fatal stop application
func Fatal(data ...interface{}) {
	logger.Fatal(data...)
}
