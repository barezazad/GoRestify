package pkg_err

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_config"
)

// Parse convert chained error to the Final format for send in JSON format
func Parse(err error, lang dictionary.Lang) (int, error) {
	var final final
	var status int

	for err != nil {

		switch e := err.(type) {

		case interface{ Unwrap() error }:
			err = errors.Unwrap(err)

		case *withMessage:
			if final.Message == "" {
				final.Message = dictionary.Translate(lang, e.Msg, e.Params...)
			}
			err = e.Err

		case *withCode:
			if final.Code == "" {
				final.Code = e.Code
			}
			err = e.Err

		case *withStatus:
			if e.Status == 0 {
				e.Status = http.StatusBadRequest
			}

			final.Status = e.Status
			status = e.Status
			err = e.Err

		case *withType:
			final.Type = e.Type
			final.Title = dictionary.Translate(lang, e.Title)
			err = e.Err

		case *withPath:
			final.Path += appendText(final.Path, e.Path)
			err = e.Err

		case *withInvalidParam:
			field := Field{
				Field:        e.Field,
				Reason:       dictionary.Translate(lang, e.Reason, e.Params...),
				ReasonParams: e.Params,
			}
			final.InvalidParams = append(final.InvalidParams, field)
			err = e.Err

		case *withCustom:
			err = e.Err

		case error:
			// disable original_error in prod
			if pkg_config.Config.IsDebug {
				final.OriginalError += e.Error()
			}
			err = errors.Unwrap(err)

		default:
			log.Println("There shouldn't be a default error", err)
			return http.StatusInternalServerError, &final
		}
	}
	return status, &final
}

// GetCustom extract custom error from error's interface
func GetCustom(err error) (customError customError) {
	for err != nil {

		switch e := err.(type) {

		case interface{ Unwrap() error }:
			err = errors.Unwrap(err)

		case *withCustom:
			return e.Custom

		case error:
			if errCast, ok := e.(*withMessage); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*withCode); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*withType); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*withPath); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*withStatus); ok {
				err = errCast.Err
				continue
			}
			if errCast, ok := e.(*withInvalidParam); ok {
				err = errCast.Err
				continue
			}
			return
		default:
			log.Println("There shouldn't be a default for getting custom", err)
			return
		}
	}

	return
}

// ApplyCustom add custom errors to the error's interface
func ApplyCustom(err error) error {

	customError := GetCustom(err)
	theme := UniqErrorMap[customError]

	err = addType(err, ""+theme.Type, theme.Title)
	err = addStatus(err, theme.Status)
	return err

}

func appendText(str string, txt string) (result string) {
	if str == "" {
		result = txt
	} else {
		result = str + ", " + txt
	}
	return
}

// ParseHTTPCallErr to prase json error from http to gerr .
func ParseHTTPCallErr(body string, StatusCode int) (err error) {

	var errRes finalResponseInHTTP
	if jsonErr := json.Unmarshal([]byte(body), &errRes); jsonErr != nil {
		err = New(fmt.Sprintf("Failed to unmarshal error response: %v", jsonErr), "E1181675").
			Custom(InternalServerErr).Message(SomethingWentWrong).Build()
		return
	}

	Log(errors.New("Error"), "E1148015", "http call", errRes.Error)

	textError := errRes.Error.OriginalError
	if textError == "" {
		textError = errRes.Error.Message
	}

	custom := New(textError, "E1618024").
		Custom(getCustomErr(StatusCode)).
		Message(errRes.Error.Message)

	if len(errRes.Error.InvalidParams) > 0 {
		var field Field
		field.Field = errRes.Error.InvalidParams[0].Field
		field.Reason = errRes.Error.InvalidParams[0].Reason
		field.ReasonParams = errRes.Error.MessageParams
		custom.InvalidParam(field.Field, field.Reason, field.ReasonParams...)
	}

	err = custom.Build()

	return
}
