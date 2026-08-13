package acc

import "GoRestify/pkg/pkg_types"

// list of resources for acc domain
const (
	Domain string = "accounting"

	TransactionWrite pkg_types.Resource = "transaction:write"
	TransactionRead  pkg_types.Resource = "transaction:read"

	SlotWrite pkg_types.Resource = "slot:write"
	SlotRead  pkg_types.Resource = "slot:read"

	CurrencyWrite pkg_types.Resource = "currency:write"
	CurrencyRead  pkg_types.Resource = "currency:read"

	AccountCreditWrite pkg_types.Resource = "account-credit:write"
	AccountCreditRead  pkg_types.Resource = "account-credit:read"
)
