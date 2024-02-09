package acc

import "GoRestify/pkg/pkg_types"

// types for acc domain
const (
	CreateTransaction pkg_types.Event = "transaction-create"
	UpdateTransaction pkg_types.Event = "transaction-update"
	DeleteTransaction pkg_types.Event = "transaction-delete"
	ListTransaction   pkg_types.Event = "transaction-list"
	ViewTransaction   pkg_types.Event = "transaction-view"
)
