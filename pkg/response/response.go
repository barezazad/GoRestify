package response

import (
	"errors"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/param"
	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_err"

	"github.com/gin-gonic/gin"
)

// Result is a standard output for success and failed requests
type Result struct {
	Message string                 `json:"message,omitempty"`
	Data    interface{}            `json:"data,omitempty"`
	Error   error                  `json:"error,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
	// CustomError pkg_err.CustomError     `json:"custom_error,omitempty"`
}

// Response holding method related to response
type Response struct {
	Result  Result
	status  int
	Engine  *pkg_config.ConfigInstance
	Context *gin.Context
	abort   bool
	params  param.Param
}

// New initiate the Response object
func New(context *gin.Context) *Response {
	return &Response{
		Engine:  pkg_config.Config,
		Context: context,
	}
}

// NewParam initiate the Response object and params
func NewParam(context *gin.Context, table string) (*Response, param.Param) {
	params := param.Get(context, table)
	return &Response{
		Engine:  pkg_config.Config,
		Context: context,
		params:  params,
	}, params
}

// Error is used for add error to the result
func (r *Response) Error(err interface{}, data ...interface{}) *Response {
	if errCast, ok := err.(string); ok {
		r.Result.Error = errors.New(errCast)
	}
	if errCast, ok := err.(error); ok {
		r.Result.Error = errCast
	}

	r.Result.Data = data
	return r
}

// Status is used for add error to the result
func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}

// Message get a message and params then translate it
func (r *Response) Message(message string, params ...interface{}) *Response {
	r.Result.Message = dictionary.Translate(dictionary.GetLang(r.Context), message,
		params...)
	return r
}

// Abort prepare response to abort instead of return in last step (JSON)
func (r *Response) Abort() *Response {
	r.abort = true
	return r
}

// JSON write output as a json
func (r *Response) JSON(data ...interface{}) {
	var parsedError error
	lang := dictionary.GetLang(r.Context)
	if r.Result.Error != nil {
		r.Result.Error = pkg_err.AddPath(r.Result.Error, r.Context.Request.RequestURI)

		r.Result.Error = pkg_err.ApplyCustom(r.Result.Error)

		r.status, parsedError = pkg_err.Parse(r.Result.Error, lang)
	}

	// if data is one element don't put it in array
	var finalData interface{}
	if data != nil {
		finalData = data
		if len(data) == 1 {
			finalData = data[0]
		}
	}

	if r.abort {
		r.Context.AbortWithStatusJSON(r.status, &Result{
			Message: r.Result.Message,
			Error:   parsedError,
			Data:    finalData,
		})
	} else {
		r.Context.JSON(r.status, &Result{
			Message: r.Result.Message,
			Error:   parsedError,
			Data:    finalData,
			// CustomError: r.Result.Error,
		})
	}
}
