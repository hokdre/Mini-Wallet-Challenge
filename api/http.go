package api

import (
	"context"
	"net/http"
	"time"

	"github.com/hokdre/mini-ewallet/internal/controller"
	"github.com/hokdre/mini-ewallet/pkg/util"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

var (
	e *echo.Echo
)

type Config struct {
	PORT          string
	ReadTimeOut   time.Duration
	WriteTimeOut  time.Duration
	WalletHandler *controller.WalletHttpController
	Encryption    util.Encryption
}

func HTTPStart(cfg Config) {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPatch,
		},
	}))

	setupRoutes(
		e,
		cfg.WalletHandler,
		cfg.Encryption,
	)

	server := &http.Server{
		Addr:         cfg.PORT,
		ReadTimeout:  cfg.ReadTimeOut,
		WriteTimeout: cfg.WriteTimeOut,
	}
	go func() {
		if err := e.StartServer(server); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Shutting down the server: ", err)
		}
	}()

}

func HttpDown(ctx context.Context) {
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatalf("Ungrafully shutdown : %s \n", err)
	}

	e.Logger.Printf("Successfully shutdown the server")
}
