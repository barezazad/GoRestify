package base_model

import (
	"GoRestify/pkg/pkg_types"
	"io"
	"mime/multipart"
	"time"
)

// DocumentTable is used inside the repo layer
const (
	DocumentTable = "base_documents"
)

// Document model
type Document struct {
	ID                uint                    `json:"id,omitempty"`
	Name              string                  `gorm:"not null" bind:"required" json:"name,omitempty"`
	ReferenceID       uint                    `gorm:"index:reference_id_idx;not null" bind:"required" json:"reference_id,omitempty"`
	Type              pkg_types.Enum          `bind:"one_of=document_type" json:"type,omitempty"`
	CreatedAt         *time.Time              `gorm:"->;type:timestamp;not null;default:current_timestamp;" json:"created_at,omitempty"`
	UpdatedAt         *time.Time              `gorm:"->;type:timestamp;not null;default:current_timestamp on update current_timestamp;" json:"updated_at,omitempty"`
	Attachments       []*multipart.FileHeader `json:"attachments,omitempty" gorm:"-"`
	AttachmentsBase64 []string                `json:"attachments_base64,omitempty" gorm:"-"`
}

// DocumentReader .
type DocumentReader struct {
	File       io.Reader             `json:"file"`
	Attachment *multipart.FileHeader `json:"attachment"`
}
