package param

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_log"
	"GoRestify/pkg/pkg_types"

	"github.com/gin-gonic/gin"
)

// orderByPattern allows only a column or table.column identifier (blocks SQL injection in ORDER BY).
var orderByPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*(\.[a-zA-Z_][a-zA-Z0-9_]*)?$`)

// Get is a function for filling param.Model
func Get(c *gin.Context, table string) (param Param) {

	generateOrder(c, &param, table)
	generateSelectedColumns(c, &param)
	generateLimit(c, &param)
	generateOffset(c, &param)

	param.Filter = readFilter(c)

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

	if v := strings.TrimSpace(c.Query("order_by")); v != "" {
		if orderByPattern.MatchString(v) {
			orderBy = v
		} else {
			pkg_log.Info("invalid order_by, using default:", v)
		}
	}

	if v := strings.TrimSpace(c.Query("direction")); v != "" {
		dir := strings.ToLower(v)
		if dir == "asc" || dir == "desc" {
			direction = dir
		} else {
			pkg_log.Info("invalid direction, using default:", v)
		}
	}

	param.Order = orderBy + " " + direction
}

// readFilter reads the filter query value from RawQuery.
// Go's url.ParseQuery skips any key=value that contains a raw ';', which would
// silently drop injection payloads like ...1; DROP TABLE... — so we parse it ourselves.
func readFilter(c *gin.Context) string {
	raw := c.Request.URL.RawQuery
	if raw == "" {
		return ""
	}

	for _, part := range strings.Split(raw, "&") {
		if part == "" {
			continue
		}
		key, value, _ := strings.Cut(part, "=")
		key, err := url.QueryUnescape(key)
		if err != nil || key != "filter" {
			continue
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			continue
		}
		return strings.TrimSpace(value)
	}

	return strings.TrimSpace(c.Query("filter"))
}

func generateSelectedColumns(c *gin.Context, param *Param) {
	param.Select = "*"
	if c.Query("select") != "" {
		param.Select = c.Query("select")
	}
}

const (
	// DefaultPageSize is used when page_size is omitted or invalid.
	DefaultPageSize = 10
	// MaxPageSize caps page_size to avoid oversized list queries (DoS).
	MaxPageSize = 100
)

func generateLimit(c *gin.Context, param *Param) {
	var err error
	param.Limit = DefaultPageSize
	if c.Query("page_size") != "" {
		param.Limit, err = strconv.Atoi(c.Query("page_size"))
		if err != nil || param.Limit <= 0 {
			pkg_log.CheckError(err, "Limit is not a number")
			param.Limit = DefaultPageSize
		} else if param.Limit > MaxPageSize {
			param.Limit = MaxPageSize
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
