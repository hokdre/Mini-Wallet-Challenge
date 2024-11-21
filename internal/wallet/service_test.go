package wallet

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	mock "github.com/hokdre/mini-ewallet/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestWalletService(t *testing.T) {
	t.Run("Init", TestInit)
	t.Run("Enable", TestEnable)
	t.Run("Disable", TestDisable)
	t.Run("Get", TestWalletService_Get)
	t.Run("GetTransaction", TestGetTransaction)
	t.Run("Deposit", TestDeposit)
}

func TestInit(t *testing.T) {
	t.Run("failed validation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(errors.New("err")).Times(1)

		w := &walletService{
			cfg: Config{
				Validator: validator,
			},
		}
		token, err := w.Init(context.Background(), "")
		assert.Error(t, err)
		assert.Equal(t, "", token)
	})

	t.Run("failed check account", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		var errExpected = errors.New("err")

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				Validator:   validator,
				AccountRepo: accountRepo,
			},
		}
		token, err := w.Init(context.Background(), externalId)
		assert.Error(t, err, errExpected)
		assert.Equal(t, "", token)
	})

	t.Run("failed create account", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		var errExpected = errors.New("err")

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{}, nil).Times(1)
		accountRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				Validator:    validator,
				AccountRepo:  accountRepo,
				TxRepository: txRepo,
			},
		}
		token, err := w.Init(context.Background(), externalId)
		assert.Error(t, err, errExpected)
		assert.Equal(t, "", token)
	})

	t.Run("failed create wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		var errExpected = errors.New("err")

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{}, nil).Times(1)
		accountRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil).Times(1)

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(errExpected).Times(1)
		w := &walletService{
			cfg: Config{
				Validator:        validator,
				AccountRepo:      accountRepo,
				TxRepository:     txRepo,
				WalletRepository: walletRepo,
			},
		}
		token, err := w.Init(context.Background(), externalId)
		assert.Error(t, err, errExpected)
		assert.Equal(t, "", token)
	})

	t.Run("failed create wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		var errExpected = errors.New("err")

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{}, nil).Times(1)
		accountRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil).Times(1)

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(errExpected).Times(1)
		w := &walletService{
			cfg: Config{
				Validator:        validator,
				AccountRepo:      accountRepo,
				TxRepository:     txRepo,
				WalletRepository: walletRepo,
			},
		}
		token, err := w.Init(context.Background(), externalId)
		assert.Error(t, err, errExpected)
		assert.Equal(t, "", token)
	})

	t.Run("failed create token", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		var errExpected = errors.New("err")

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{}, nil).Times(1)
		accountRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil).Times(1)

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil).Times(1)

		encryption := mock.NewMockEncryption(ctrl)
		encryption.EXPECT().Encrypt(gomock.Any()).Return("", errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				Validator:        validator,
				AccountRepo:      accountRepo,
				TxRepository:     txRepo,
				WalletRepository: walletRepo,
				Encryption:       encryption,
			},
		}
		token, err := w.Init(context.Background(), externalId)
		assert.Error(t, err, errExpected)
		assert.Equal(t, "", token)
	})

	t.Run("Success New Register", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		token := "token"

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{}, nil).Times(1)
		accountRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil).Times(1)

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().CreateTx(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(nil).Times(1)

		encryption := mock.NewMockEncryption(ctrl)
		encryption.EXPECT().Encrypt(gomock.Any()).Return(token, nil).Times(1)

		w := &walletService{
			cfg: Config{
				Validator:        validator,
				AccountRepo:      accountRepo,
				TxRepository:     txRepo,
				WalletRepository: walletRepo,
				Encryption:       encryption,
			},
		}
		res, err := w.Init(context.Background(), externalId)
		assert.NoError(t, err)
		assert.Equal(t, token, res)
	})

	t.Run("Success old register", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		externalId := uuid.New().String()
		token := "token"

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)

		accountRepo := mock.NewMockAccountRepository(ctrl)
		accountRepo.EXPECT().Get(gomock.Any(), internal.AccountFilter{
			ExternalIDs: []string{externalId},
		}).Return(model.Account{
			ID: uuid.New(),
		}, nil).Times(1)

		walletRepo := mock.NewMockWalletRepository(ctrl)

		encryption := mock.NewMockEncryption(ctrl)
		encryption.EXPECT().Encrypt(gomock.Any()).Return(token, nil).Times(1)

		w := &walletService{
			cfg: Config{
				Validator:        validator,
				AccountRepo:      accountRepo,
				TxRepository:     txRepo,
				WalletRepository: walletRepo,
				Encryption:       encryption,
			},
		}
		res, err := w.Init(context.Background(), externalId)
		assert.NoError(t, err)
		assert.Equal(t, token, res)
	})
}

func TestEnable(t *testing.T) {
	t.Run("Failed get wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()
		var errExpected = errors.New("err")

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Enable(context.Background(), accountID)
		assert.Error(t, err)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Failed status already enable", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()

		wallet := model.Wallet{
			Status: model.WalletStatus.Enabled,
		}
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Enable(context.Background(), accountID)
		assert.Error(t, err, model.ErrWalletAlreadyEnabled)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Failed Update wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			Status: model.WalletStatus.Disabled,
		}
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		walletRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errExpected).Times(1)
		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Enable(context.Background(), accountID)
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Success Enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()

		wallet := model.Wallet{
			Status: model.WalletStatus.Disabled,
		}
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		walletRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Enable(context.Background(), accountID)
		assert.NoError(t, err)
		assert.Equal(t, res.Status, model.WalletStatus.Enabled)
	})
}

func TestDisable(t *testing.T) {
	t.Run("Failed get wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()
		var errExpected = errors.New("err")

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Disable(context.Background(), accountID)
		assert.Error(t, err)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Failed status already disabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()

		wallet := model.Wallet{
			Status: model.WalletStatus.Disabled,
		}
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Disable(context.Background(), accountID)
		assert.Error(t, err, model.ErrWalletAlreadyDisabled)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Failed Update wallet", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			Status: model.WalletStatus.Enabled,
		}
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		walletRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(errExpected).Times(1)
		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Disable(context.Background(), accountID)
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Success Enabled", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		accountID := uuid.New()

		wallet := model.Wallet{
			Status: model.WalletStatus.Enabled,
		}
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		walletRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil).Times(1)
		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Disable(context.Background(), accountID)
		assert.NoError(t, err)
		assert.Equal(t, res.Status, model.WalletStatus.Disabled)
	})
}

func TestWalletService_Get(t *testing.T) {
	t.Run("failed get wallet", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Get(context.Background(), accountID)
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("failed get wallet status disabled", func(t *testing.T) {
		accountID := uuid.New()

		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{
			Status: model.WalletStatus.Disabled,
		}, nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Get(context.Background(), accountID)
		assert.Error(t, err, model.ErrWalletDisabled)
		assert.Equal(t, model.Wallet{}, res)
	})

	t.Run("Success", func(t *testing.T) {
		accountID := uuid.New()
		wallet := model.Wallet{
			ID:      uuid.New(),
			OwnedBy: accountID,
			Status:  model.WalletStatus.Enabled,
		}

		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Get(context.Background(), accountID)
		assert.NoError(t, err)
		assert.Equal(t, wallet, res)
	})
}

func TestGetTransaction(t *testing.T) {
	t.Run("failed get wallet", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.GetTransactions(context.Background(), accountID)
		assert.Error(t, err, errExpected)
		assert.Nil(t, res)
	})

	t.Run("failed wallet disabled", func(t *testing.T) {
		accountID := uuid.New()

		wallet := model.Wallet{
			ID:      uuid.New(),
			OwnedBy: accountID,
			Status:  model.WalletStatus.Disabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.GetTransactions(context.Background(), accountID)
		assert.Error(t, err, model.ErrWalletDisabled)
		assert.Nil(t, res)
	})

	t.Run("failed get transactions", func(t *testing.T) {
		accountID := uuid.New()
		var errExpect = errors.New("err")

		wallet := model.Wallet{
			ID:      uuid.New(),
			OwnedBy: accountID,
			Status:  model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().List(gomock.Any(), internal.TransactionFilter{
			WalletIDs: []string{wallet.ID.String()},
		}).Return(nil, errExpect).Times(1)
		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				TransactionRepository: transactionRepo,
			},
		}
		res, err := w.GetTransactions(context.Background(), accountID)
		assert.Error(t, err, errExpect)
		assert.Nil(t, res)
	})

	t.Run("Success", func(t *testing.T) {
		accountID := uuid.New()
		transactions := []model.Transaction{}

		wallet := model.Wallet{
			ID:      uuid.New(),
			OwnedBy: accountID,
			Status:  model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)

		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().List(gomock.Any(), internal.TransactionFilter{
			WalletIDs: []string{wallet.ID.String()},
		}).Return(transactions, nil).Times(1)
		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				TransactionRepository: transactionRepo,
			},
		}
		res, err := w.GetTransactions(context.Background(), accountID)
		assert.NoError(t, err)
		assert.Equal(t, transactions, res)
	})

}

func TestDeposit(t *testing.T) {
	t.Run("failed get wallet", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Deposit(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed validate", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
				Validator:        validator,
			},
		}
		res, err := w.Deposit(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed create pending tx", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
			},
		}
		res, err := w.Deposit(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed increment, error update status transaction", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Increment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), errExpected).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Deposit(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed increment, success update status transaction", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Increment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), errExpected).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Deposit(context.Background(), accountID, model.Transaction{})
		assert.Nil(t, err)
		assert.Equal(t, model.TransactionStatus.Failed, res.Status)
	})

	t.Run("Success Deposit", func(t *testing.T) {
		accountID := uuid.New()

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Increment(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Deposit(context.Background(), accountID, model.Transaction{})
		assert.Nil(t, err)
		assert.Equal(t, model.TransactionStatus.Success, res.Status)
	})
}

func TestWithdrawal(t *testing.T) {
	t.Run("failed get wallet", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(model.Wallet{}, errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed validate", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository: walletRepo,
				Validator:        validator,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed create pending tx", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed decrment, error update status transaction", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Decrement(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), errExpected).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(errExpected).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Error(t, err, errExpected)
		assert.Equal(t, model.Transaction{}, res)
	})

	t.Run("failed decrment, success update status transaction", func(t *testing.T) {
		accountID := uuid.New()
		var errExpected = errors.New("err")

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Decrement(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), errExpected).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Nil(t, err)
		assert.Equal(t, model.TransactionStatus.Failed, res.Status)
	})

	t.Run("failed decrment simulate conccurent issue, success update status transaction", func(t *testing.T) {
		accountID := uuid.New()

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Decrement(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(0), nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Nil(t, err)
		assert.Equal(t, model.TransactionStatus.Failed, res.Status)
	})

	t.Run("Success Withdrawal", func(t *testing.T) {
		accountID := uuid.New()

		wallet := model.Wallet{
			ID:     uuid.New(),
			Status: model.WalletStatus.Enabled,
		}
		ctrl := gomock.NewController(t)
		walletRepo := mock.NewMockWalletRepository(ctrl)
		walletRepo.EXPECT().GetOne(gomock.Any(), internal.WalletFilter{
			OwnedBies: []string{accountID.String()},
		}).Return(wallet, nil).Times(1)

		validator := mock.NewMockValidator(ctrl)
		validator.EXPECT().Validate(gomock.Any()).Return(nil).Times(1)

		transactionRepo := mock.NewMockTransactionRepository(ctrl)
		transactionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(1)

		walletRepo.EXPECT().Decrement(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(int64(1), nil).Times(1)

		txRepo := mock.NewMockTxRepository(ctrl)
		txRepo.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
			return fn(ctx, nil)
		}).Times(1)

		transactionRepo.EXPECT().UpdateTx(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil).Times(1)

		w := &walletService{
			cfg: Config{
				WalletRepository:      walletRepo,
				Validator:             validator,
				TransactionRepository: transactionRepo,
				TxRepository:          txRepo,
			},
		}
		res, err := w.Withdrawal(context.Background(), accountID, model.Transaction{})
		assert.Nil(t, err)
		assert.Equal(t, model.TransactionStatus.Success, res.Status)
	})
}
