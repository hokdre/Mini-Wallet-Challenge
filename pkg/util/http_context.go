package util

import (
	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/labstack/echo/v4"
)

const (
	KeyAccountID = "ACCOUNT_ID"
)

func GetAccountID(ctx echo.Context) (uuid.UUID, error) {
	accountIDCtx := ctx.Get(KeyAccountID)
	if accountIDCtx == nil {
		return uuid.Nil, model.ErrLoginInfoUknown
	}

	accountID, ok := accountIDCtx.(uuid.UUID)
	if !ok {
		return uuid.Nil, model.ErrLoginInfoUknown
	}

	return accountID, nil
}

func SetAccountID(ctx echo.Context, accountID uuid.UUID) echo.Context {
	ctx.Set(KeyAccountID, accountID)
	return ctx
}
