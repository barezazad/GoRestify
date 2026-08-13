package acc_term

// List of messages and errors for acc domain
const (
	Transaction  = "transaction"
	Transactions = "transactions"

	Slot  = "slot"
	Slots = "slots"

	Currency   = "currency"
	Currencies = "currencies"

	AccountCredit  = "account credit"
	AccountCredits = "account credits"

	YouHaveNotEnoughCredit     = "you do not have enough balance"
	LedgerRecordsAreImmutable  = "ledger records are immutable; use DoTransaction to post money movement"
	TransactionDBTxRequired    = "database transaction is required for posting"
	SlotsRequired              = "at least one slot is required"
	SlotsMustBalance           = "slot debits must equal slot credits"
	SlotsMustMatchAmount       = "sum of slot debits must equal transaction amount"
	SlotCurrencyMismatch       = "slot currency must match transaction currency"
	SlotAccountNotAllowed      = "slot account must be sender or receiver"
	InvalidTransactionAmounts  = "amount must be > 0 and fee must be >= 0 and <= amount"
	InvalidSlotAmounts         = "slot debit/credit cannot be negative; one side must be non-zero"
	CurrencyInUse              = "currency is referenced by credits, transactions, or slots"
	AccountCreditAlreadyExists = "account credit already exists for this account and currency"
)
