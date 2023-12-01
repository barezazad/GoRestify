package user_status

import "GoRestify/pkg/pkg_types"

// request status enum
const (
	Active   pkg_types.Enum = "active"
	InActive pkg_types.Enum = "in_active"
)

// List order status list
var List = []pkg_types.Enum{
	Active,
	InActive,
}
