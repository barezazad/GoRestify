package pkg_err

import (
	"fmt"
	"net/http"
)

// finalResponseInHTTP .
type finalResponseInHTTP struct {
	Error final `json:"error"`
}

// final is hold different parts of the error based on rfc 7807
type final struct {
	Code          string        `json:"code,omitempty"`
	Type          string        `json:"type,omitempty"`
	Title         string        `json:"title,omitempty"`
	Message       string        `json:"message,omitempty"`
	MessageParams []interface{} `json:"message_params,omitempty"`
	Path          string        `json:"path,omitempty"`
	InvalidParams []Field       `json:"invalid_params,omitempty"`
	Status        int           `json:"-"`
	OriginalError string        `json:"original_error,omitempty"`
}

// Field is used as an array inside the FieldError
type Field struct {
	Field        string        `json:"field,omitempty"`
	Reason       string        `json:"reason,omitempty"`
	ReasonParams []interface{} `json:"reason_params,omitempty"`
}

func (p *final) Error() string {
	errStr := fmt.Sprintf("#%v, err:%v", p.Code, p.OriginalError)
	return errStr
}

// customError is used for defining errors related to this application, this is a bridge between the
type customError int

// errorTheme hold the error's type and title
type errorTheme struct {
	Type   string
	Title  string
	Status int
}

// customErrorMap is used for defining the error for each domain
type customErrorMap map[customError]errorTheme

// CustomError types
const (
	Nil customError = iota
	UnknownErr
	UnauthorizedErr
	NotFoundErr
	RouteNotFoundErr
	ValidationFailedErr
	ForeignErr
	DuplicateErr
	InternalServerErr
	BindingErr
	ForbiddenErr
	PreDataInsertedErr   //428
	TimeoutErr           //408
	TemporaryRedirectErr //307
	BadRequestErr
)

// getCustomErr get CustomError by http code
func getCustomErr(code int) customError {

	switch code {

	case 401:
		return UnauthorizedErr
	case 404:
		return NotFoundErr
	case 442:
		return ValidationFailedErr
	case 409:
		return DuplicateErr
	case 500:
		return InternalServerErr
	case 422:
		return BindingErr
	case 403:
		return ForbiddenErr
	case 408:
		return TimeoutErr
	case 307:
		return TemporaryRedirectErr
	case 400:
		return BadRequestErr
	}
	return 0
}

// UniqErrorMap is used for categorized errors and connect error with error page also primary fill
// the status code and domain and title
var UniqErrorMap customErrorMap

func init() {
	UniqErrorMap = make(map[customError]errorTheme)

	UniqErrorMap[UnauthorizedErr] = errorTheme{
		Type:   "#UNAUTHORIZED",
		Title:  Unauthorized,
		Status: http.StatusUnauthorized,
	}

	UniqErrorMap[ValidationFailedErr] = errorTheme{
		Type:   "#VALIDATION_FAILED",
		Title:  ValidationFailed,
		Status: http.StatusUnprocessableEntity,
	}

	UniqErrorMap[NotFoundErr] = errorTheme{
		Type:   "#NOT_FOUND",
		Title:  RecordNotFound,
		Status: http.StatusNotFound,
	}

	UniqErrorMap[RouteNotFoundErr] = errorTheme{
		Type:   "#ROUTE_NOT_FOUND",
		Title:  RouteNotFound,
		Status: http.StatusNotFound,
	}

	UniqErrorMap[ForeignErr] = errorTheme{
		Type:   "#FOREIGN_KEY",
		Title:  ForeignKeyError,
		Status: http.StatusConflict,
	}

	UniqErrorMap[InternalServerErr] = errorTheme{
		Type:   "#INTERNAL_SERVER_ERROR",
		Title:  InternalServerError,
		Status: http.StatusInternalServerError,
	}

	UniqErrorMap[DuplicateErr] = errorTheme{
		Type:   "#DUPLICATE_ERROR",
		Title:  DuplicateHappened,
		Status: http.StatusConflict,
	}

	UniqErrorMap[BindingErr] = errorTheme{
		Type:   "#NOT_BIND",
		Title:  BindFailed,
		Status: http.StatusUnprocessableEntity,
	}

	UniqErrorMap[ForbiddenErr] = errorTheme{
		Type:   "#FORBIDDEN",
		Title:  Forbidden,
		Status: http.StatusForbidden,
	}

	UniqErrorMap[TimeoutErr] = errorTheme{
		Type:   "#TIMEOUT",
		Title:  TimeoutHappened,
		Status: http.StatusRequestTimeout,
	}

	UniqErrorMap[TemporaryRedirectErr] = errorTheme{
		Type:   "#TEMPORARY_REDIRECT",
		Title:  TemporaryRedirect,
		Status: http.StatusTemporaryRedirect,
	}
	UniqErrorMap[BadRequestErr] = errorTheme{
		Type:   "#BadRequest",
		Title:  BadRequest,
		Status: http.StatusBadRequest,
	}
}
