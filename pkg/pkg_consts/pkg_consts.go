package pkg_consts

import "net/http"

// constance for share all projects
const (

	// date and time layout
	DateLayout                = "2006-01-02"
	DateTimeLayout            = "2006-01-02 15:04:05"
	DateTimeLayoutJSON        = "2006-01-02T15:04:05Z"
	DateTimeLayoutZone        = "2006-01-02 15:04:05 -0700"
	DateTimeLayoutTransaction = "060102150405.000000"

	// pagination
	DefaultLimit         = 100
	CheckToUpdateSetting = 60

	// log property and file path
	LogFormat       = "stdout"
	LogServerOutput = "logs/server.log"
	LogAPIOutput    = "logs/api.log"

	TermsFile = "assets/terms/terms.toml"

	// server address
	ServerAddress = "0.0.0.0"
)

// vars
var (
	HTTPClient *http.Client
)
