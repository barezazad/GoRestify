package action

import (
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/validator"
)

// Action enums
const (
	Login validator.Action = "login"
)

// MustBeInTypes it uses to check enums for validation
var MustBeInTypes = map[string]interface{}{
	"domain_app": domain_app.List,
}
