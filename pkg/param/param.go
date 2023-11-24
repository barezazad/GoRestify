package param

import (
	"GoRestify/pkg/dictionary"
	"GoRestify/pkg/pkg_consts"
	"GoRestify/pkg/pkg_types"
	"GoRestify/pkg/tx"

	"github.com/gin-gonic/gin"
)

// Param for describing request's parameter
type Param struct {
	Pagination
	Search       string
	Filter       string
	PreCondition string
	UserID       uint
	Phone        string
	Email        string
	Lang         dictionary.Lang
	Domain       pkg_types.Enum
	Dynamic      map[string]string
	DynamicID    map[string]uint
	Tx           tx.Tx
	context      *gin.Context
}

// Pagination is a struct, contains the fields which affected the front-end pagination
type Pagination struct {
	Select string
	Order  string
	Limit  int
	Offset int
}

// New return an initiate of the param with default limit
func New() Param {
	var param Param
	param.Limit = pkg_consts.DefaultLimit
	param.Order = "id"

	return param
}
