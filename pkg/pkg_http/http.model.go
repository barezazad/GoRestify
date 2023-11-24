package ghttp

import "GoRestify/pkg/dictionary"

// Request model
type Request struct {
	Method         string // required
	EndPoint       string // required
	Language       dictionary.Lang
	Headers        []Header
	FormData       FormData
	Payload        interface{} // payload of request
	ParsedResponse interface{} // // parsed response data based on provided model
}

// Header model
type Header struct {
	Key   string
	Value string
}

// FormData model
type FormData struct {
	FileKey  string
	FilePath string
	Payload  []PayloadFormData
}

// PayloadFormData .
type PayloadFormData struct {
	Key   string
	Value string
}

// Default headers
var (
	Authorization  = "Authorization"
	XAuthorization = "X-Authorization"
	SecretKey      = "SECRET-KEY"
)
