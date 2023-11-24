package utils

import (
	"GoRestify/pkg/pkg_log"
	"crypto/tls"
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"
)

// ConfigEmail is used for initiate the ConfigEmail design pattern
type ConfigEmail struct {
	Host     string
	Port     int
	Username string
	Password string
}

// SendEmail send email
func (c *ConfigEmail) SendEmail(toArr, ccArr []string, aliasFrom, subject, body, attachment string) (err error) {

	// example "Account Verification <%v>"
	From := fmt.Sprintf("%v <%v>", aliasFrom, c.Username)

	m := gomail.NewMessage()
	m.SetHeader("Subject", subject)
	m.SetHeader("From", From)
	m.SetHeader("To", toArr...)
	for _, v := range ccArr {
		m.SetAddressHeader("Cc", v, strings.Split(v, "@")[0])
	}
	if attachment != "" {
		m.Attach(attachment)
	}
	m.SetBody("text/html", body)

	d := gomail.NewDialer(c.Host, c.Port, c.Username, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = d.DialAndSend(m)
	if err != nil {
		pkg_log.CheckError(err, fmt.Sprint("Error In sending Email:", err, c))
	}

	return
}
