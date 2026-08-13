package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"strings"
	"sync/atomic"
	"time"

	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// APILogger is used to save requests and response by using logapi
func APILogger() gin.HandlerFunc {
	var reqID uint64

	level := "debug"

	logger := pkg_log.New(pkg_consts.LogFormat,
		pkg_consts.LogAPIOutput,
		level,
		false,
		true)

	return func(c *gin.Context) {
		start := time.Now()
		buf, _ := io.ReadAll(c.Request.Body)
		reqDataReader := io.NopCloser(bytes.NewBuffer(buf))

		// We have to create a new Buffer, because reqDataReader will be read.
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		id := atomic.AddUint64(&reqID, 1)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		logRequest(logger, c, id, reqDataReader)

		c.Next()

		latency := int(math.Ceil(float64(time.Since(start).Nanoseconds()) / 1000000.0))

		logResponse(logger, c, latency, blw)
	}
}

// Logging Response
func logRequest(logger *logrus.Logger, c *gin.Context, reqID uint64, reqDataReader io.Reader) {

	request := getBody(reqDataReader)

	// prevent to save passwords / credentials
	path := strings.ToLower(c.Request.URL.Path)
	if strings.Contains(path, "login") || strings.Contains(path, "password") || strings.Contains(path, "refresh-token") {
		request = nil
	}

	logger.WithFields(logrus.Fields{
		"reqID":      reqID,
		"ip":         c.Request.Header.Get("X-User-IP"),
		"method":     c.Request.Method,
		"uri":        c.Request.RequestURI,
		"path":       c.Request.URL.Path,
		"request":    request,
		"params":     c.Request.URL.Query(),
		"referer":    c.Request.Referer(),
		"user_agent": c.Request.UserAgent(),
	}).Info("request")
	c.Set("resID", reqID)
}

// Logging Response
func logResponse(logger *logrus.Logger, c *gin.Context, latency int, blw *bodyLogWriter) {

	resID, ok := c.Get("resID")
	if !ok {
		pkg_log.Debug("there is no resIndex for element", getBody(blw.body))
	}

	logger.WithFields(logrus.Fields{
		"resID":       resID,
		"status":      c.Writer.Status(),
		"latency":     latency, // time to process
		"data_length": c.Writer.Size(),
		"response":    getBody(blw.body),
	}).Info("response")

}

// Read body
func getBody(reader io.Reader) interface{} {

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(reader); err != nil {
		pkg_log.Debug(err)
	}

	var obj interface{}

	if err := json.NewDecoder(buf).Decode(&obj); err != nil {
		if err.Error() != "EOF" {

		}
	}

	return obj
}
