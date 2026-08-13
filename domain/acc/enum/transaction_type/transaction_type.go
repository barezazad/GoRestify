package transaction_type

import "GoRestify/pkg/pkg_types"

// transaction type enum
const (
	Buy        pkg_types.Enum = "buy"
	Sell       pkg_types.Enum = "sell"
	Transfer   pkg_types.Enum = "transfer"
	Deposit    pkg_types.Enum = "deposit"
	Withdrawal pkg_types.Enum = "withdrawal"
	Fee        pkg_types.Enum = "fee"
)

// List transaction type list
var List = []pkg_types.Enum{
	Buy,
	Sell,
	Transfer,
	Deposit,
	Withdrawal,
	Fee,
}
