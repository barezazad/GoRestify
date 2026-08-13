package pkg_err

import (
	"fmt"
)

// withCode is used for carrying the code of error
type withCode struct {
	Err  error
	Code string
}

func (p *withCode) Error() string {
	return fmt.Sprint(p.Err)
}

// AddCode add custom code to the error, it is useful for tracing an error
// if error pass from different method or function you can gave it multiple error
// and in case error happened easily we can trace which function's passed that error
// similar to panic
func AddCode(err error, code string) error {
	return &withCode{
		Err:  fmt.Errorf("#%v, %w", code, err),
		Code: code,
	}
}

// withMessage keeps the message of the error, each error can have one message
type withMessage struct {
	Err    error
	Msg    string
	Params []interface{}
}

func (p *withMessage) Error() string {
	return fmt.Sprint(p.Err)
}

// addMessage add custom message to error, params can be used inside the translator function
func addMessage(err error, msg string, params ...interface{}) error {
	return &withMessage{
		Err:    err,
		Msg:    msg,
		Params: params,
	}
}

// withType is add type and title to the error
type withType struct {
	Err   error
	Type  string
	Title string
}

func (p *withType) Error() string { return fmt.Sprint(p.Err) }

// addType used for adding type to the error
func addType(err error, errType string, title string) error {
	return &withType{
		Err:   err,
		Type:  errType,
		Title: title,
	}
}

// withPath attach path to the error
type withPath struct {
	Err  error
	Path string
}

func (p *withPath) Error() string { return fmt.Sprint(p.Err) }

// AddPath is used for adding path to the error, useful for REST API
func AddPath(err error, path string) error {
	return &withPath{
		Err:  err,
		Path: path,
	}
}

// withStatus attach status to the error
type withStatus struct {
	Err    error
	Status int
}

func (p *withStatus) Error() string { return fmt.Sprint(p.Err) }

// addStatus can be used for adding HTTP status code like 404 or etc
func addStatus(err error, status int) error {
	return &withStatus{
		Err:    err,
		Status: status,
	}
}

// withInvalidParam holds invalid parameters
type withInvalidParam struct {
	Err    error
	Field  string
	Reason string
	Params []interface{}
}

func (p *withInvalidParam) Error() string { return fmt.Sprint(p.Err) }

// AddInvalidParam is used for specify a field that has an error
func AddInvalidParam(err error, field, reason string, params ...interface{}) error {
	var gErr error
	if err == nil {
		gErr = fmt.Errorf(fmt.Sprintf(reason, params...))
	} else {
		gErr = err
	}

	return &withInvalidParam{
		Err:    gErr,
		Field:  field,
		Reason: reason,
		Params: params,
	}
}

// withCustom is used for holding the uniqError for filling the type and title based on local
// customization
type withCustom struct {
	Err    error
	Custom customError
}

func (p *withCustom) Error() string { return fmt.Sprint(p.Err) }

// SetCustom is used for adding a custom error to the error and it reduce size of the code
func SetCustom(err error, custom customError) error {
	return &withCustom{
		Err:    err,
		Custom: custom,
	}
}
