package ghttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_err"
)

// Do http call for a server
func Do(request Request) (err error) {

	payload, err := generatePayload(&request)
	if err != nil {
		return
	}

	req, err := http.NewRequest(request.Method, request.EndPoint, payload)
	if err != nil {
		errorMessage := fmt.Sprintf("Init HTTP call to method %v and API %v: %v", request.Method, request.EndPoint, err)
		err = pkg_err.New(errorMessage, "E1122422").Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")

	// set language
	if request.Language == "" {
		request.Language = dictionary.En
	}
	req.Header.Set("X-LANGUAGE", string(request.Language))

	// set headers
	for _, v := range request.Headers {
		if v.Key == Authorization && !strings.Contains(strings.ToLower(v.Value), "bearer ") {
			v.Value = fmt.Sprintf("bearer %v", v.Value)
		}

		req.Header.Set(v.Key, v.Value)
	}

	var res *http.Response
	if res, err = pkg_consts.HTTPClient.Do(req); err != nil {
		errorMessage := fmt.Sprintf("DO HTTP call to method %v and API %v: %v", request.Method, request.EndPoint, err)
		err = pkg_err.New(errorMessage, "E1198982").Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}

	// read Body
	body, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		message := fmt.Sprintf("Failed to read error response body: %v", readErr)
		err = pkg_err.New(message, "E1196236", request.EndPoint).Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
		return
	}
	defer res.Body.Close()

	bodyStr := string(body)

	if res.StatusCode == 200 {
		// parse success response
		if request.ParsedResponse != nil {
			if err = json.Unmarshal([]byte(bodyStr), &request.ParsedResponse); err != nil {
				err = pkg_err.New(fmt.Sprintf("Failed to unmarshal error response: %v", err), "E1130716").
					Custom(pkg_err.InternalServerErr).Message(pkg_err.SomethingWentWrong).Build()
				return
			}
		}
	} else {
		// parse error response
		err = pkg_err.ParseHTTPCallErr(bodyStr, res.StatusCode)

	}

	return err
}

// generatePayload factory payload
func generatePayload(request *Request) (payload io.Reader, err error) {

	payload = nil

	if request.FormData.FilePath != "" { // if payload is DataForm

		// create bufferData multipart
		bufferData := &bytes.Buffer{}
		writer := multipart.NewWriter(bufferData)

		// open file
		file, errFile1 := os.Open(request.FormData.FilePath)
		defer file.Close()

		part1, errFile1 := writer.CreateFormFile(request.FormData.FileKey, filepath.Base(request.FormData.FilePath))
		if _, errFile1 = io.Copy(part1, file); err != nil {
			err = pkg_err.Take(errFile1, "E1168269", "error in set file in payload").
				Message(pkg_err.SomethingWentWrong).Custom(pkg_err.InternalServerErr).Build()
			return
		}

		for _, v := range request.FormData.Payload {
			_ = writer.WriteField(v.Key, v.Value)
		}

		err = writer.Close()
		if err != nil {
			err = pkg_err.Take(err, "E1127211", "error in close write file").
				Message(pkg_err.SomethingWentWrong).Custom(pkg_err.InternalServerErr).Build()
			return
		}

		request.Headers = append(request.Headers,
			Header{
				Key:   "Content-Type",
				Value: writer.FormDataContentType(),
			},
		)

		payload = bufferData

	} else if request.Payload == "" || request.Payload == nil { // if get request
		payload = nil

	} else { // if payload is json
		payloadByte, _ := json.Marshal(request.Payload)
		payload = bytes.NewReader(payloadByte)
	}

	return
}
