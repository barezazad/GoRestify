package acc_model

// CurrencyTable is used inside the repo layer
const (
	CurrencyTable = "acc_currencies"
)

// Currency model
type Currency struct {
	ID     uint   `json:"id,omitempty"`
	Name   string `gorm:"type:varchar(150)" json:"name,omitempty" bind:"required"`
	Symbol string `gorm:"type:varchar(150)" json:"symbol,omitempty" bind:"required"`
}
