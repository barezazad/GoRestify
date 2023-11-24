package param

import (
	"strconv"
	"strings"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_types"

	"github.com/gin-gonic/gin"
)

// Get is a function for filling param.Model
func Get(c *gin.Context, table string) (param Param) {

	generateOrder(c, &param, table)
	generateSelectedColumns(c, &param)
	generateLimit(c, &param)
	generateOffset(c, &param)

	param.Filter = strings.TrimSpace(c.Query("filter"))

	userID, ok := c.Get("USER_ID")
	if ok {
		param.UserID = userID.(uint)
	}

	phone, ok := c.Get("PHONE")
	if ok {
		param.Phone = phone.(string)
	}

	email, ok := c.Get("EMAIL")
	if ok {
		param.Email = email.(string)
	}

	domain, ok := c.Get("X-DOMAIN")
	if ok {
		param.Domain = domain.(pkg_types.Enum)
	} else {
		defaultAPP, exist := c.Get("DEFAULT-Domain")
		if exist {
			param.Domain = pkg_types.Enum(defaultAPP.(string))
		}
	}

	param.Lang = dictionary.GetLang(c)

	param.context = c

	return param
}

func generateOrder(c *gin.Context, param *Param, table string) {
	orderBy := table + ".id"
	direction := "desc"

	if c.Query("order_by") != "" {
		orderBy = c.Query("order_by")
	}

	if c.Query("direction") != "" {
		direction = c.Query("direction")
	}

	param.Order = orderBy + " " + direction
}

func generateSelectedColumns(c *gin.Context, param *Param) {
	param.Select = "*"
	if c.Query("select") != "" {
		param.Select = c.Query("select")
	}
}

func generateLimit(c *gin.Context, param *Param) {
	var err error
	param.Limit = 10
	if c.Query("page_size") != "" {
		param.Limit, err = strconv.Atoi(c.Query("page_size"))
		if err != nil {
			pkg_log.CheckError(err, "Limit is not a number")
			param.Limit = 10
		}
	}
}

func generateOffset(c *gin.Context, param *Param) {
	var page int
	var err error
	page = 1

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			pkg_log.CheckError(err, "Offset is not a positive number")
		}
		if page <= 0 {
			page = 1
		}
	}

	param.Offset = param.Limit * (page - 1)
}

// GetParamUint get uint value in key from context
func (p *Param) GetParamUint(key string) (value uint) {

	tmpValue, ok := p.context.Get(key)
	if ok {
		value = tmpValue.(uint)
	}

	return
}

// GetParamString get String value in key from context
func (p *Param) GetParamString(key string) (value string) {

	tmpValue, ok := p.context.Get(key)
	if ok {
		value = tmpValue.(string)
	}

	return
}
