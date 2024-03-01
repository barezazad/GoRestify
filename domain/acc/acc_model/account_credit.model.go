package acc_model

import (
	"GoRestify/pkg/decimal"
	"time"
)

// AccountCreditTable is used inside the repo layer
const (
	AccountCreditTable = "acc_account_credits"
)

// AccountCredit model
type AccountCredit struct {
	ID         uint             `json:"id,omitempty"`
	CurrencyID uint             `gorm:"index:currency_id_idx;not null" json:"currency_id,omitempty" bind:"required"`
	AccountID  uint             `gorm:"index:account_id_idx;not null" json:"account_id,omitempty" bind:"required"`
	Balance    *decimal.Decimal `gorm:"not null;type:decimal(22,6)" json:"balance"`
	CreatedAt  *time.Time       `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt  *time.Time       `gorm:"->;type:timestamp;not null;default:current_timestamp on update current_timestamp;" json:"updated_at,omitempty"`
	Currency   string           `gorm:"-:migration;->" json:"currency" table:"acc_currencies.name as currency"`
}
