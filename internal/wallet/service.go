package wallet

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/hokdre/mini-ewallet/pkg/util"
)

type Config struct {
	AccountRepo           internal.AccountRepository
	WalletRepository      internal.WalletRepository
	TransactionRepository internal.TransactionRepository
	Validator             util.Validator
	Encryption            util.Encryption
	TxRepository          internal.TxRepository
}

type walletService struct {
	cfg Config
}

func NewWalletService(cfg Config) *walletService {
	return &walletService{cfg: cfg}
}

func (w *walletService) Init(ctx context.Context, externalID string) (string, error) {
	newAccount := model.Account{
		ID:                 uuid.New(),
		ExternalCustomerID: externalID,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	err := w.cfg.Validator.Validate(newAccount)
	if err != nil {
		return "", err
	}

	existingAcc, errGetAcc := w.cfg.AccountRepo.Get(ctx, internal.AccountFilter{
		ExternalIDs: []string{externalID},
	})
	if errGetAcc != nil && errGetAcc != sql.ErrNoRows {
		return "", errGetAcc
	}

	accountID := existingAcc.ID
	if accountID == uuid.Nil {
		v, _, err := w.createAccountAndWallet(ctx, externalID)
		if err != nil {
			return "", err
		}
		accountID = v
	}

	return w.createToken(accountID)
}

func (w *walletService) createAccountAndWallet(
	ctx context.Context,
	externalID string) (uuid.UUID, uuid.UUID, error) {
	timeStamp := time.Now()
	newAccount := model.Account{
		ID:                 uuid.New(),
		ExternalCustomerID: externalID,
		CreatedAt:          timeStamp,
		UpdatedAt:          timeStamp,
	}

	newWallet := model.Wallet{
		ID:        uuid.New(),
		OwnedBy:   newAccount.ID,
		Status:    model.WalletStatus.Disabled,
		Balance:   0,
		CreatedAt: timeStamp,
		UpdatedAt: timeStamp,
	}

	errCreate := w.cfg.TxRepository.Process(ctx, func(ctx context.Context, tx *sql.Tx) error {
		errAccount := w.cfg.AccountRepo.CreateTx(ctx, tx, newAccount)
		if errAccount != nil {
			return errAccount
		}

		errWallet := w.cfg.WalletRepository.CreateTx(ctx, tx, newWallet)
		if errWallet != nil {
			return errWallet
		}

		return nil
	})
	if errCreate != nil {
		return uuid.Nil, uuid.Nil, errCreate
	}

	return newAccount.ID, newWallet.ID, nil
}

func (w *walletService) createToken(accountID uuid.UUID) (string, error) {
	return w.cfg.Encryption.Encrypt(accountID.String())
}

func (w *walletService) Enable(ctx context.Context, accountID uuid.UUID) (model.Wallet, error) {
	wallet, err := w.cfg.WalletRepository.GetOne(ctx, internal.WalletFilter{
		OwnedBies: []string{accountID.String()},
	})
	if err != nil {
		return model.Wallet{}, err
	}
	if wallet.Status == model.WalletStatus.Enabled {
		return model.Wallet{}, model.ErrWalletAlreadyEnabled
	}

	timestamp := time.Now()
	wallet.Status = model.WalletStatus.Enabled
	wallet.EnabledAt = &timestamp
	wallet.DisabledAt = nil
	wallet.UpdatedAt = timestamp
	err = w.cfg.WalletRepository.Update(ctx, wallet)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletService) Disable(ctx context.Context, accountID uuid.UUID) (model.Wallet, error) {
	wallet, err := w.cfg.WalletRepository.GetOne(ctx, internal.WalletFilter{
		OwnedBies: []string{accountID.String()},
	})
	if err != nil {
		return model.Wallet{}, err
	}
	if wallet.Status == model.WalletStatus.Disabled {
		return model.Wallet{}, model.ErrWalletAlreadyDisabled
	}

	timestamp := time.Now()
	wallet.Status = model.WalletStatus.Disabled
	wallet.EnabledAt = nil
	wallet.DisabledAt = &timestamp
	wallet.UpdatedAt = timestamp
	err = w.cfg.WalletRepository.Update(ctx, wallet)
	if err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil
}

func (w *walletService) Get(ctx context.Context, accountID uuid.UUID) (model.Wallet, error) {
	wallet, err := w.cfg.WalletRepository.GetOne(ctx, internal.WalletFilter{
		OwnedBies: []string{accountID.String()},
	})
	if err != nil {
		return model.Wallet{}, err
	}

	if wallet.Status == model.WalletStatus.Disabled {
		return model.Wallet{}, model.ErrWalletDisabled
	}

	return wallet, nil
}

func (w *walletService) GetTransactions(ctx context.Context, accountID uuid.UUID) ([]model.Transaction, error) {
	wallet, err := w.Get(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if wallet.Status == model.WalletStatus.Disabled {
		return nil, model.ErrWalletDisabled
	}

	transactions, err := w.cfg.TransactionRepository.List(ctx, internal.TransactionFilter{
		WalletIDs: []string{wallet.ID.String()},
	})
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (w *walletService) Deposit(ctx context.Context, accountID uuid.UUID, transaction model.Transaction) (model.Transaction, error) {
	wallet, err := w.Get(ctx, accountID)
	if err != nil {
		return model.Transaction{}, err
	}

	timestamp := time.Now()
	transaction.ID = uuid.New()
	transaction.WalletID = wallet.ID
	transaction.CreatedAt = timestamp
	transaction.UpdatedAt = timestamp
	transaction.TransactedAt = nil
	transaction.Status = model.TransactionStatus.Pending
	transaction.Type = model.TransactionType.Deposit
	err = w.cfg.Validator.Validate(transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	err = w.cfg.TransactionRepository.Create(ctx, transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	err = w.cfg.TxRepository.Process(ctx, func(ctx context.Context, tx *sql.Tx) error {
		_, errIncrement := w.cfg.WalletRepository.Increment(ctx, tx, wallet, transaction.Amount)
		timestamp := time.Now()
		transaction.Status = model.TransactionStatus.Success
		transaction.TransactedAt = &timestamp
		if errIncrement != nil {
			transaction.Status = model.TransactionStatus.Failed
			transaction.TransactedAt = nil
		}

		errTransaction := w.cfg.TransactionRepository.UpdateTx(ctx, tx, transaction)
		if errTransaction != nil {
			return errTransaction
		}

		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}

func (w *walletService) Withdrawal(ctx context.Context, accountID uuid.UUID, transaction model.Transaction) (model.Transaction, error) {
	wallet, err := w.Get(ctx, accountID)
	if err != nil {
		return model.Transaction{}, err
	}

	timestamp := time.Now()
	transaction.ID = uuid.New()
	transaction.WalletID = wallet.ID
	transaction.CreatedAt = timestamp
	transaction.UpdatedAt = timestamp
	transaction.TransactedAt = nil
	transaction.Status = model.TransactionStatus.Pending
	transaction.Type = model.TransactionType.Withdrawal
	err = w.cfg.Validator.Validate(transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	err = w.cfg.TransactionRepository.Create(ctx, transaction)
	if err != nil {
		return model.Transaction{}, err
	}

	err = w.cfg.TxRepository.Process(ctx, func(ctx context.Context, tx *sql.Tx) error {
		affected, errDecrement := w.cfg.WalletRepository.Decrement(ctx, tx, wallet, transaction.Amount)
		timestamp := time.Now()
		transaction.Status = model.TransactionStatus.Success
		transaction.TransactedAt = &timestamp
		if errDecrement != nil || affected == 0 {
			transaction.Status = model.TransactionStatus.Failed
			transaction.TransactedAt = nil
		}

		errTransaction := w.cfg.TransactionRepository.UpdateTx(ctx, tx, transaction)
		if errTransaction != nil {
			return errTransaction
		}

		return nil
	})
	if err != nil {
		return model.Transaction{}, err
	}

	return transaction, nil
}
