package pkg_err

import "GoRestify/pkg/pkg_log"

// Log is combining AddCode and LogError in services
func Log(err error, code string, message string, data ...interface{}) error {

	if code != "" {
		err = AddCode(err, code)
	}

	pkg_log.LogError(err, message, data...)

	return err
}

// TickValidate is automatically add validation error custom to the error
func TickValidate(err error, code string, data ...interface{}) error {
	Log(err, code, ValidationFailed, data...)
	return SetCustom(err, ValidationFailedErr)
}
