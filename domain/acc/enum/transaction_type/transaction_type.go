package transaction_type

import "GoRestify/pkg/pkg_types"

// transaction type enum
const (
	Buy pkg_types.Enum = "buy"
)

// List transaction type list
var List = []pkg_types.Enum{
	Buy,
}
