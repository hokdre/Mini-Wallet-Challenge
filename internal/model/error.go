package model

import (
	"errors"
	"fmt"
)

var (
	ErrBussiness             = errors.New("Bussiness Error")
	ErrWalletAlreadyEnabled  = fmt.Errorf("%w : Already Enabled", ErrBussiness)
	ErrWalletAlreadyDisabled = fmt.Errorf("%w : Already Disabled", ErrBussiness)
	ErrWalletDisabled        = fmt.Errorf("%w : Wallet Disabled", ErrBussiness)

	ErrLoginInfoUknown = errors.New("Login info unknown")
)
