package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hokdre/mini-ewallet/api"
	"github.com/hokdre/mini-ewallet/config"
	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/account"
	"github.com/hokdre/mini-ewallet/internal/controller"
	"github.com/hokdre/mini-ewallet/internal/transaction"
	"github.com/hokdre/mini-ewallet/internal/wallet"
	"github.com/hokdre/mini-ewallet/pkg/persistence"
	"github.com/hokdre/mini-ewallet/pkg/util"
)

func main() {
	cfg := config.Init()

	// persistence
	db, err := persistence.OpenPostgreDB(
		persistence.Config{
			Host:        cfg.PostgreHost,
			Username:    cfg.PostgreUsername,
			Password:    cfg.PostgrePassword,
			DB:          cfg.PostgreDB,
			Port:        cfg.PostgrePort,
			SSLMode:     cfg.PostgreSSLMode,
			MaxIdleConn: cfg.PostgreMaxIdleConn,
			MaxOpenConn: cfg.PostgreMaxOpenConn,
		},
	)
	if err != nil {
		log.Fatalf("failed open db : %s", err)
	}

	// repository
	accountRepo := account.NewAccountRepo(db)
	walletRepo := wallet.NewWalletRepository(db)
	transactionRepo := transaction.NewAccountRepo(db)
	txRepo := internal.NewTxRepository(db)

	// util
	validator := util.NewValidator()
	encryption, err := util.NewAesEncryption(cfg.AESSecret)
	if err != nil {
		log.Fatalf("failed construct encryption : %s", err)
	}

	// service
	walletService := wallet.NewWalletService(
		wallet.Config{
			AccountRepo:           accountRepo,
			WalletRepository:      walletRepo,
			TransactionRepository: transactionRepo,
			TxRepository:          txRepo,
			Validator:             validator,
			Encryption:            encryption,
		},
	)

	// http handler
	walletHandler := controller.NewWalletController(walletService)

	// start server
	api.HTTPStart(api.Config{
		PORT:          ":" + cfg.RestPORT,
		ReadTimeOut:   cfg.RestReadTimeOut,
		WriteTimeOut:  cfg.RestWriteTimeOut,
		WalletHandler: walletHandler,
		Encryption:    encryption,
	})

	// shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	signal := <-quit
	log.Printf("receive kill signal : %s \n", signal)
	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.RestShoutDownTimeOut,
	)
	defer cancel()
	api.HttpDown(ctx)
}
