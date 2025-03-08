package entities

import (
	"time"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "PENDING"
	PaymentStatusCompleted PaymentStatus = "COMPLETED"
	PaymentStatusFailed    PaymentStatus = "FAILED"
)

type Payment struct {
	ID            string        `json:"id" gorm:"primaryKey;size:36"`
	OrderID       uint          `json:"order_id" gorm:"not null"`
	Amount        float64       `json:"amount" gorm:"type:decimal(10,2);not null"`
	Method        string        `json:"method" gorm:"size:50;not null"`
	Status        PaymentStatus `json:"status" gorm:"size:20;not null"`
	TransactionID string        `json:"transaction_id,omitempty" gorm:"size:100"`
	ErrorMessage  string        `json:"error_message,omitempty" gorm:"size:255"`
	CreatedAt     time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
}
