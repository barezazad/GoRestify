package acc

import "GoRestify/pkg/pkg_types"

// list of resources for acc domain
const (
	Domain string = "accounting"

	TransactionWrite pkg_types.Resource = "transaction:write"
	TransactionRead  pkg_types.Resource = "transaction:read"
)
