package acc

import "GoRestify/pkg/pkg_types"

// types for acc domain
const (
	CreateTransaction pkg_types.Event = "transaction-create"
	UpdateTransaction pkg_types.Event = "transaction-update"
	DeleteTransaction pkg_types.Event = "transaction-delete"
	ListTransaction   pkg_types.Event = "transaction-list"
	ViewTransaction   pkg_types.Event = "transaction-view"

	CreateSlot pkg_types.Event = "slot-create"
	UpdateSlot pkg_types.Event = "slot-update"
	DeleteSlot pkg_types.Event = "slot-delete"
	ListSlot   pkg_types.Event = "slot-list"
	ViewSlot   pkg_types.Event = "slot-view"

	CreateCurrency pkg_types.Event = "currency-create"
	UpdateCurrency pkg_types.Event = "currency-update"
	DeleteCurrency pkg_types.Event = "currency-delete"
	ListCurrency   pkg_types.Event = "currency-list"
	ViewCurrency   pkg_types.Event = "currency-view"

	CreateAccountCredit pkg_types.Event = "account-credit-create"
	UpdateAccountCredit pkg_types.Event = "account-credit-update"
	DeleteAccountCredit pkg_types.Event = "account-credit-delete"
	ListAccountCredit   pkg_types.Event = "account-credit-list"
	ViewAccountCredit   pkg_types.Event = "account-credit-view"
)
