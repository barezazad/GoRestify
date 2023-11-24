package domain_app

import "GoRestify/pkg/pkg_types"

// the type of default status
const (
	Admin pkg_types.Enum = "admin"
	User  pkg_types.Enum = "user"
)

// List to default status
var List = []pkg_types.Enum{
	Admin,
	User,
}
