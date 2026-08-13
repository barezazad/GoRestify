package pkg_log

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
)

var (
	sensitiveHeaderRE = regexp.MustCompile(`(?im)^(Authorization|Proxy-Authorization|Cookie|Set-Cookie|X-Authorization)\s*:\s*.+$`)
	dsnPasswordRE     = regexp.MustCompile(`(?i)(://[^:\s/]+:)([^@/\s]+)(@)`)
	jsonPasswordRE    = regexp.MustCompile(`(?i)("(?:password|passwd|pwd|secret|token|refresh_token|access_token)"\s*:\s*")([^"]*)(")`)
)

// LoggingTransport logs outbound HTTP traffic with sensitive headers redacted.
type LoggingTransport struct {
	Transport http.RoundTripper
}

// RoundTrip .
func (d *LoggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	requestDump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	requestDumpStr := redactSensitive(string(requestDump))

	if isFile := strings.Contains(strings.ToLower(strings.Join(req.Header["Content-Type"], ",")), "multipart/form-data"); isFile {
		if len(requestDumpStr) > 1000 {
			requestDumpStr = requestDumpStr[:1000] + "......"
		}
	}
	fmt.Println(requestDumpStr)

	trans := d.Transport
	if trans == nil {
		trans = http.DefaultTransport
	}

	resp, err := trans.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		return nil, err
	}

	responseDumpStr := redactSensitive(string(responseDump))
	if len(responseDumpStr) > 1000 {
		responseDumpStr = responseDumpStr[:1000] + "......"
	}
	fmt.Println(responseDumpStr)

	return resp, err
}

func redactSensitive(dump string) string {
	dump = sensitiveHeaderRE.ReplaceAllString(dump, "$1: [REDACTED]")
	dump = dsnPasswordRE.ReplaceAllString(dump, "${1}[REDACTED]${3}")
	dump = jsonPasswordRE.ReplaceAllString(dump, `${1}[REDACTED]${3}`)
	return dump
}
