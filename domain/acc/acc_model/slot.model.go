package acc_model

import (
	"GoRestify/pkg/decimal"
	"time"
)

// SlotTable is used inside the repo layer
const (
	SlotTable = "acc_slots"
)

// Slot model — immutable ledger line created only via DoTransaction.
type Slot struct {
	ID            uint             `json:"id,omitempty"`
	AccountID     uint             `gorm:"index;not null" json:"account_id,omitempty" bind:"required"`
	CurrencyID    uint             `gorm:"index;not null" json:"currency_id,omitempty" bind:"required"`
	TransactionID uint             `gorm:"index;not null" json:"transaction_id,omitempty"`
	Debit         *decimal.Decimal `gorm:"not null;type:decimal(22,6);default:0" json:"debit"`
	Credit        *decimal.Decimal `gorm:"not null;type:decimal(22,6);default:0" json:"credit"`
	Balance       *decimal.Decimal `gorm:"not null;type:decimal(22,6)" json:"balance"`
	Notes         string           `gorm:"type:varchar(200)" json:"notes,omitempty" bind:"max=200"`
	CreatedAt     *time.Time       `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	// Joined from transaction header (table name is historically base_transactions).
	Type     string `gorm:"-:migration;->" table:"base_transactions.type as type" json:"type,omitempty"`
	Hash     string `gorm:"-:migration;->" table:"base_transactions.hash as hash" json:"hash,omitempty"`
	Currency string `gorm:"-:migration;->" json:"currency" table:"acc_currencies.name as currency"`
}
