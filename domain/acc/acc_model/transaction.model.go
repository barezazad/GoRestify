package acc_model

import (
	"GoRestify/pkg/decimal"
	"GoRestify/pkg/pkg_types"
	"time"
)

// TransactionTable is used inside the repo layer for specify the table name
const (
	TransactionTable = "base_transactions"
)

// Transaction model
type Transaction struct {
	ID          uint             `json:"id,omitempty"`
	SenderID    uint             `gorm:"index:sender_id_idx;not null" json:"sender_id,omitempty" bind:"required"`
	ReceiverID  uint             `gorm:"index:receiver_id_idx;not null" json:"receiver_id,omitempty" `
	CurrencyID  uint             `gorm:"index;not null" json:"currency_id,omitempty" bind:"required"`
	Type        pkg_types.Enum   `gorm:"type:varchar(50);not null;" json:"type,omitempty" bind:"required,one_of=transaction_type"`
	Hash        string           `gorm:"index:hash_idx;type:varchar(70);" json:"hash,omitempty" bind:""`
	Amount      *decimal.Decimal `gorm:"not null;type:decimal(22,6)" json:"amount" bind:"required"`
	Fee         *decimal.Decimal `gorm:"not null;type:decimal(22,6);default:0" json:"fee,omitempty"`
	NetAmount   *decimal.Decimal `gorm:"not null;type:decimal(22,6)" json:"net_amount"`
	Note        string           `gorm:"type:varchar(225);" json:"note,omitempty"`
	Description string           `gorm:"type:varchar(225);" json:"description,omitempty"`
	CreatedBy   uint             `gorm:"index:created_by_idx;not null" json:"created_by,omitempty" bind:"required"`
	CreatedAt   *time.Time       `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	Slots       []Slot           `gorm:"-" json:"slots,omitempty" bind:"-"`
	Currency    string           `gorm:"-:migration;->" json:"currency" table:"acc_currencies.name as currency"`
}
