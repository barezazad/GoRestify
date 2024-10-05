package action

import (
	"GoRestify/domain/acc/enum/transaction_type"
	"GoRestify/domain/base/enum/account_status"
	"GoRestify/domain/base/enum/account_type"
	"GoRestify/domain/base/enum/document_type"
	"GoRestify/internal/core/enum/domain_app"

	"GoRestify/pkg/validator"
)

// Action enums
const (
	Login validator.Action = "login"
)

// MustBeInTypes it uses to check enums for validation
var MustBeInTypes = map[string]interface{}{
	"domain_app":       domain_app.List,
	"account_status":   account_status.List,
	"account_type":     account_type.List,
	"transaction_type": transaction_type.List,
	"document_type":    document_type.List,
}
