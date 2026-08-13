package acc_model

// CurrencyTable is used inside the repo layer
const (
	CurrencyTable = "acc_currencies"
)

// Currency model
type Currency struct {
	ID     uint   `json:"id,omitempty"`
	Name   string `gorm:"type:varchar(150);uniqueIndex:idx_currency_name" json:"name,omitempty" bind:"required"`
	Symbol string `gorm:"type:varchar(150);uniqueIndex:idx_currency_symbol" json:"symbol,omitempty" bind:"required"`
}
