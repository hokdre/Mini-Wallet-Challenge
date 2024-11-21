package model

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID                 uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	ExternalCustomerID string    `json:"external_customer_id" db:"external_customer_id" validate:"required"`
	CreatedAt          time.Time `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at" validate:"required"`
}
