package api

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal/controller"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/hokdre/mini-ewallet/pkg/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func setupRoutes(
	e *echo.Echo,
	walletHandler *controller.WalletHttpController,
	encryption util.Encryption,
) {
	e.Logger.SetLevel(log.DEBUG)
	protected := e.Group("/api/v1/wallet")
	protected.Use(AuthorizationMiddleware(encryption))
	protected.GET("", walletHandler.Get)
	protected.POST("", walletHandler.Enable)
	protected.PATCH("", walletHandler.Disable)
	protected.GET("/transactions", walletHandler.GetTransactions)
	protected.POST("/deposits", walletHandler.Deposit)
	protected.POST("/withdrawals", walletHandler.Withdrawal)

	e.POST("/api/v1/init", walletHandler.Init)
}

func AuthorizationMiddleware(encryption util.Encryption) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// Get the Authorization header
			authHeader := ctx.Request().Header.Get("Authorization")
			if authHeader == "" {
				return util.SendError(
					ctx,
					http.StatusUnauthorized,
					model.ErrLoginInfoUknown,
				)
			}

			headers := strings.Split(authHeader, " ")
			if len(headers) != 2 {
				return util.SendError(
					ctx,
					http.StatusUnauthorized,
					model.ErrLoginInfoUknown,
				)
			}

			if strings.ToLower(headers[0]) != "token" {
				return util.SendError(
					ctx,
					http.StatusUnauthorized,
					model.ErrLoginInfoUknown,
				)
			}
			token, err := encryption.Decrypt(headers[1])
			if err != nil {
				return util.SendError(
					ctx,
					http.StatusUnauthorized,
					model.ErrLoginInfoUknown,
				)
			}

			accountID, err := uuid.Parse(token)
			if err != nil {
				return util.SendError(
					ctx,
					http.StatusUnauthorized,
					model.ErrLoginInfoUknown,
				)
			}

			util.SetAccountID(ctx, accountID)
			return next(ctx)
		}
	}
}
