package model

import (
	"time"

	"github.com/google/uuid"
)

var WalletStatus = struct {
	Enabled  string
	Disabled string
}{
	Enabled:  "enabled",
	Disabled: "disabled",
}

type Wallet struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	OwnedBy    uuid.UUID  `json:"user_id" db:"user_id" validate:"required"`
	Balance    int64      `json:"balance" db:"balance" validate:"gte=0"`
	Status     string     `json:"status" db:"status" validate:"required,enumWalletStatus"`
	EnabledAt  *time.Time `json:"enabled_at" db:"enabled_at"`
	DisabledAt *time.Time `json:"disabled_at" db:"disabled_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at" validate:"required"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at" validate:"required"`
}
