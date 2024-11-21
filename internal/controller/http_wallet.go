package controller

import (
	"net/http"

	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/hokdre/mini-ewallet/pkg/util"
	"github.com/labstack/echo/v4"
)

type WalletHttpController struct {
	walletService internal.WalletService
}

func NewWalletController(
	walletService internal.WalletService,
) *WalletHttpController {
	return &WalletHttpController{
		walletService: walletService,
	}
}

func (w *WalletHttpController) Init(ctx echo.Context) error {
	payload := new(struct {
		CustomerXID string `form:"customer_xid"`
	})
	if err := ctx.Bind(payload); err != nil {
		return util.SendFailed(
			ctx,
			http.StatusBadRequest,
			map[string]interface{}{
				"body": err.Error(),
			},
		)
	}

	token, err := w.walletService.Init(ctx.Request().Context(), payload.CustomerXID)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	return util.SendSuccess(
		ctx,
		http.StatusOK,
		map[string]interface{}{
			"token": token,
		},
	)
}

func (w *WalletHttpController) Enable(ctx echo.Context) error {
	accountID, err := util.GetAccountID(ctx)
	if err != nil {
		return util.SendError(ctx, http.StatusUnauthorized, err)
	}

	wallet, err := w.walletService.Enable(ctx.Request().Context(), accountID)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	return util.SendSuccess(ctx, http.StatusOK, map[string]interface{}{
		"wallet": map[string]interface{}{
			"id":         wallet.ID,
			"owned_by":   wallet.OwnedBy,
			"status":     wallet.Status,
			"enabled_at": wallet.EnabledAt,
			"balance":    wallet.Balance,
		},
	})
}

func (w *WalletHttpController) Disable(ctx echo.Context) error {
	accountID, err := util.GetAccountID(ctx)
	if err != nil {
		return util.SendError(ctx, http.StatusUnauthorized, err)
	}

	wallet, err := w.walletService.Disable(ctx.Request().Context(), accountID)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	return util.SendSuccess(ctx, http.StatusOK, map[string]interface{}{
		"wallet": map[string]interface{}{
			"id":          wallet.ID,
			"owned_by":    wallet.OwnedBy,
			"status":      wallet.Status,
			"disabled_at": wallet.DisabledAt,
			"balance":     wallet.Balance,
		},
	})
}

func (w *WalletHttpController) Get(ctx echo.Context) error {
	accountID, err := util.GetAccountID(ctx)
	if err != nil {
		return util.SendError(ctx, http.StatusUnauthorized, err)
	}

	wallet, err := w.walletService.Get(ctx.Request().Context(), accountID)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	return util.SendSuccess(ctx, http.StatusOK, map[string]interface{}{
		"wallet": map[string]interface{}{
			"id":         wallet.ID,
			"owned_by":   wallet.OwnedBy,
			"status":     wallet.Status,
			"enabled_at": wallet.EnabledAt,
			"balance":    wallet.Balance,
		},
	})
}

func (w *WalletHttpController) GetTransactions(ctx echo.Context) error {
	accountID, err := util.GetAccountID(ctx)
	if err != nil {
		return util.SendError(ctx, http.StatusUnauthorized, err)
	}

	transactions, err := w.walletService.GetTransactions(ctx.Request().Context(), accountID)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	data := []interface{}{}
	for _, t := range transactions {
		data = append(data, map[string]interface{}{
			"id":            t.ID,
			"status":        t.Status,
			"transacted_at": t.TransactedAt,
			"type":          t.Type,
			"amount":        t.Amount,
			"reference_id":  t.ReferenceID,
		})
	}

	return util.SendSuccess(ctx, http.StatusOK, map[string]interface{}{
		"transactions": data,
	})
}

func (w *WalletHttpController) Deposit(ctx echo.Context) error {
	accountID, err := util.GetAccountID(ctx)
	if err != nil {
		return util.SendError(ctx, http.StatusUnauthorized, err)
	}

	payload := new(struct {
		ReferenceID string `json:"reference_id" form:"reference_id"`
		Amount      int64  `json:"amount" form:"amount"`
	})
	err = ctx.Bind(payload)
	if err != nil {
		return util.SendFailed(
			ctx,
			http.StatusBadRequest,
			map[string]interface{}{
				"body": err.Error(),
			},
		)
	}

	transaction := model.Transaction{
		Amount:      payload.Amount,
		ReferenceID: payload.ReferenceID,
	}
	transaction, err = w.walletService.Deposit(ctx.Request().Context(), accountID, transaction)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	return util.SendSuccess(ctx, http.StatusCreated, map[string]interface{}{
		"deposit": map[string]interface{}{
			"id":           transaction.ID,
			"deposited_by": accountID,
			"status":       transaction.Status,
			"deposited_at": transaction.TransactedAt,
			"amount":       transaction.Amount,
			"reference_id": transaction.ReferenceID,
		},
	})
}

func (w *WalletHttpController) Withdrawal(ctx echo.Context) error {
	accountID, err := util.GetAccountID(ctx)
	if err != nil {
		return util.SendError(ctx, http.StatusUnauthorized, err)
	}

	payload := new(struct {
		ReferenceID string `json:"reference_id" form:"reference_id"`
		Amount      int64  `json:"amount" form:"amount"`
	})
	err = ctx.Bind(payload)
	if err != nil {
		return util.SendFailed(
			ctx,
			http.StatusBadRequest,
			map[string]interface{}{
				"body": err.Error(),
			},
		)
	}

	transaction := model.Transaction{
		Amount:      payload.Amount,
		ReferenceID: payload.ReferenceID,
	}
	transaction, err = w.walletService.Withdrawal(ctx.Request().Context(), accountID, transaction)
	if err != nil {
		return util.SendFailedOrError(ctx, err)
	}

	return util.SendSuccess(ctx, http.StatusCreated, map[string]interface{}{
		"withdrawal": map[string]interface{}{
			"id":            transaction.ID,
			"withdrawal_by": accountID,
			"status":        transaction.Status,
			"withdrawal_at": transaction.TransactedAt,
			"amount":        transaction.Amount,
			"reference_id":  transaction.ReferenceID,
		},
	})
}
