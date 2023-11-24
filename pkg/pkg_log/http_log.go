package pkg_log

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
)

// LoggingTransport .
type LoggingTransport struct {
	Transport http.RoundTripper
}

// RoundTrip .
func (d *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	requestDumpStr := string(requestDump)

	if isFile := strings.Contains(strings.ToLower(strings.Join(req.Header["Content-Type"], ",")), "multipart/form-data"); isFile {
		if len(requestDumpStr) > 1000 {
			requestDumpStr = requestDumpStr[:1000] + "......"
		}
	}
	fmt.Println(requestDumpStr)

	// Use underlying transport or default if nil
	trans := d.Transport
	if trans == nil {
		trans = http.DefaultTransport
	}

	// Get response
	resp, err := trans.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Print response
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	responseDumpStr := string(responseDump)
	if len(responseDumpStr) > 1000 {
		responseDumpStr = responseDumpStr[:1000] + "......"
	}
	fmt.Println(responseDumpStr)

	return resp, err
}
