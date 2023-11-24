package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"strings"
	"time"

	"GoRestify/pkg/pkg_config"
	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// APILogger is used to save requests and response by using logapi
func APILogger() gin.HandlerFunc {
	var reqID uint

	level := "debug"

	logger := pkg_log.New(pkg_consts.LogFormat,
		pkg_consts.LogAPIOutput,
		level,
		false,
		true)

	return func(c *gin.Context) {

		if pkg_config.Config.IsDebug {
			c.Next()
		}

		start := time.Now()
		buf, _ := io.ReadAll(c.Request.Body)
		reqDataReader := io.NopCloser(bytes.NewBuffer(buf))

		//We have to create a new Buffer, because reqDataReader will be read.
		c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))
		reqID++

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		logRequest(logger, c, reqID, reqDataReader)

		c.Next()

		latency := int(math.Ceil(float64(time.Since(start).Nanoseconds()) / 1000000.0))

		logResponse(logger, c, latency, blw)

	}
}

// Logging Response
func logRequest(logger *logrus.Logger, c *gin.Context, reqID uint, reqDataReader io.Reader) {

	request := getBody(reqDataReader)

	// prevent to save the passwords
	if strings.Contains(c.Request.URL.Path, "login") {
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
