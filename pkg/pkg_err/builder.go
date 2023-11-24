package pkg_err

import (
	"GoRestify/pkg/pkg_log"
	"errors"
)

// PkgErr is main type for controlling the package
type PkgErr struct {
	err error
}

// New return an initiate of the PkgErr
func New(errStr, code string, data ...interface{}) *PkgErr {
	var pkg_err PkgErr
	err := errors.New(errStr)

	if len(code) > 2 {
		pkg_err.err = AddCode(err, code)
	} else {
		pkg_err.err = err
	}

	pkg_log.LogError(pkg_err.err, pkg_err.err.Error(), data...)

	return &pkg_err
}

// Take initiate the
func Take(err error, code string, data ...interface{}) *PkgErr {
	var pkg_err PkgErr

	if len(code) > 2 {
		pkg_err.err = AddCode(err, code)
	} else {
		pkg_err.err = err
	}

	pkg_log.LogError(pkg_err.err, pkg_err.err.Error(), data...)

	return &pkg_err
}

// Message append a message to the error
func (p *PkgErr) Message(message string, params ...interface{}) *PkgErr {
	p.err = addMessage(p.err, message, params...)
	return p
}

// Custom is used when some value like status code and basic data needs to be appended to the error
func (p *PkgErr) Custom(custom customError) *PkgErr {
	p.err = SetCustom(p.err, custom)
	return p
}

// InvalidParam is used when want to pint to a field which caused the error
func (p *PkgErr) InvalidParam(field, reason string, params ...interface{}) *PkgErr {
	p.err = AddInvalidParam(p.err, field, reason, params...)
	return p
}

// Build return an initiate of the struct
func (p *PkgErr) Build() error {
	return p.err
}
