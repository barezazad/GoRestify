package pkg_log

import "github.com/sirupsen/logrus"

var logger *logrus.Logger

// logParam used for parameter between start and initLog
type logParam struct {
	format     string
	output     string
	level      string
	JSONIndent bool
}
