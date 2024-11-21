package wallet

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/hokdre/mini-ewallet/internal"
	"github.com/hokdre/mini-ewallet/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestWalletRepository(t *testing.T) {
	t.Run("CreateTx", TestCreateTx)
	t.Run("Get", TestGet)
	t.Run("Update", TestUpdate)
	t.Run("Increment", TestIncerement)
	t.Run("Decrement", TestIncerement)
}

func TestCreateTx(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newWallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     model.WalletStatus.Disabled,
			EnabledAt:  nil,
			DisabledAt: nil,
			CreatedAt:  timestamp,
			UpdatedAt:  timestamp,
		}

		mock.
			ExpectPrepare(qCreate).
			ExpectExec().
			WithArgs(
				newWallet.ID,
				newWallet.OwnedBy,
				newWallet.Balance,
				newWallet.Status,
				newWallet.EnabledAt,
				newWallet.DisabledAt,
				newWallet.CreatedAt,
				newWallet.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		errCreate := repo.CreateTx(context.Background(), tx, newWallet)
		assert.NoError(t, errCreate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		errExpected := errors.New("Failed create prepate statement")
		newWallet := model.Wallet{}
		mock.
			ExpectPrepare(qCreate).WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		errCreate := repo.CreateTx(context.Background(), tx, newWallet)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execute", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		errExpected := errors.New("Failed to excecute statement")
		newWallet := model.Wallet{}
		mock.
			ExpectPrepare(qCreate).
			ExpectExec().WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		errCreate := repo.CreateTx(context.Background(), tx, newWallet)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGet(t *testing.T) {
	t.Run("success with filter", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timeStamp := time.Now()
		wallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     "success",
			CreatedAt:  timeStamp,
			UpdatedAt:  timeStamp,
			EnabledAt:  nil,
			DisabledAt: nil,
		}

		expectedRow := sqlmock.NewRows([]string{
			"id",
			"owned_by",
			"balance",
			"status",
			"enabled_at",
			"disabled_at",
			"created_at",
			"updated_at",
		}).AddRow(
			wallet.ID,
			wallet.OwnedBy,
			wallet.Balance,
			wallet.Status,
			wallet.EnabledAt,
			wallet.DisabledAt,
			wallet.CreatedAt,
			wallet.UpdatedAt,
		)

		filter := internal.WalletFilter{
			IDs:       []string{wallet.ID.String()},
			OwnedBies: []string{wallet.OwnedBy.String()},
		}
		mock.ExpectQuery(qGet).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.OwnedBies),
			1,
			0,
		).WillReturnRows(expectedRow)

		repo := &walletRepository{db: db}
		accResult, err := repo.GetOne(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, wallet, accResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success with no filter", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timeStamp := time.Now()
		wallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     "success",
			CreatedAt:  timeStamp,
			UpdatedAt:  timeStamp,
			EnabledAt:  nil,
			DisabledAt: nil,
		}

		expectedRow := sqlmock.NewRows([]string{
			"id",
			"owned_by",
			"balance",
			"status",
			"enabled_at",
			"disabled_at",
			"created_at",
			"updated_at",
		}).AddRow(
			wallet.ID,
			wallet.OwnedBy,
			wallet.Balance,
			wallet.Status,
			wallet.EnabledAt,
			wallet.DisabledAt,
			wallet.CreatedAt,
			wallet.UpdatedAt,
		)

		filter := internal.WalletFilter{
			IDs:       []string{},
			OwnedBies: []string{},
		}
		mock.ExpectQuery(qGet).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.OwnedBies),
			1,
			0,
		).WillReturnRows(expectedRow)

		repo := &walletRepository{db: db}
		accResult, err := repo.GetOne(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, wallet, accResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed, not found", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		filter := internal.WalletFilter{
			IDs:       []string{},
			OwnedBies: []string{},
		}
		mock.ExpectQuery(qGet).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.OwnedBies),
			1,
			0,
		).WillReturnError(sql.ErrNoRows)

		repo := &walletRepository{db: db}
		accResult, err := repo.GetOne(context.Background(), filter)
		assert.Error(t, err)
		assert.Equal(t, model.Wallet{}, accResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timestamp := time.Now()
		newWallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     model.WalletStatus.Disabled,
			EnabledAt:  nil,
			DisabledAt: nil,
			CreatedAt:  timestamp,
			UpdatedAt:  timestamp,
		}

		mock.
			ExpectPrepare(qUpdate).
			ExpectExec().
			WithArgs(
				newWallet.Status,
				newWallet.EnabledAt,
				newWallet.DisabledAt,
				newWallet.UpdatedAt,
				newWallet.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		repo := &walletRepository{db: db}
		errCreate := repo.Update(context.Background(), newWallet)
		assert.NoError(t, errCreate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		errExpected := errors.New("Failed create prepate statement")
		newWallet := model.Wallet{}
		mock.
			ExpectPrepare(qUpdate).WillReturnError(errExpected)

		repo := &walletRepository{db: db}
		errCreate := repo.Update(context.Background(), newWallet)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execure", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		errExpected := errors.New("Failed create prepate statement")
		newWallet := model.Wallet{}
		mock.
			ExpectPrepare(qUpdate).
			ExpectExec().
			WithArgs(
				newWallet.Status,
				newWallet.EnabledAt,
				newWallet.DisabledAt,
				newWallet.UpdatedAt,
				newWallet.ID,
			).WillReturnError(errors.New("unexpected error"))

		repo := &walletRepository{db: db}
		errCreate := repo.Update(context.Background(), newWallet)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestIncerement(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newWallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     model.WalletStatus.Disabled,
			EnabledAt:  nil,
			DisabledAt: nil,
			CreatedAt:  timestamp,
			UpdatedAt:  timestamp,
		}
		amount := 10000

		mock.
			ExpectPrepare(qIncrementWallet).
			ExpectExec().
			WithArgs(
				amount,
				newWallet.UpdatedAt,
				newWallet.ID,
			).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		afftected, errCreate := repo.Increment(context.Background(), tx, newWallet, int64(amount))
		assert.NoError(t, errCreate)
		assert.Equal(t, afftected, int64(1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		newWallet := model.Wallet{}
		amount := 10000

		var errExpected = errors.New("error")
		mock.
			ExpectPrepare(qIncrementWallet).
			WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		afftected, errCreate := repo.Increment(context.Background(), tx, newWallet, int64(amount))
		assert.Error(t, errCreate, errExpected)
		assert.Equal(t, afftected, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execute", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newWallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     model.WalletStatus.Disabled,
			EnabledAt:  nil,
			DisabledAt: nil,
			CreatedAt:  timestamp,
			UpdatedAt:  timestamp,
		}
		amount := 10000

		var errExpected = errors.New("error")
		mock.
			ExpectPrepare(qIncrementWallet).
			ExpectExec().
			WithArgs(
				amount,
				newWallet.UpdatedAt,
				newWallet.ID,
			).
			WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		afftected, errCreate := repo.Increment(context.Background(), tx, newWallet, int64(amount))
		assert.Error(t, errCreate)
		assert.Equal(t, afftected, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}

func TestDecrement(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newWallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     model.WalletStatus.Disabled,
			EnabledAt:  nil,
			DisabledAt: nil,
			CreatedAt:  timestamp,
			UpdatedAt:  timestamp,
		}
		amount := 10000

		mock.
			ExpectPrepare(qDecrementWallet).
			ExpectExec().
			WithArgs(
				amount,
				newWallet.UpdatedAt,
				newWallet.ID,
			).
			WillReturnResult(sqlmock.NewResult(int64(1), 1))

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		afftected, errCreate := repo.Decrement(context.Background(), tx, newWallet, int64(amount))
		assert.NoError(t, errCreate)
		assert.Equal(t, afftected, int64(1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		newWallet := model.Wallet{}
		amount := 10000

		var errExpected = errors.New("error")
		mock.
			ExpectPrepare(qDecrementWallet).
			WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		afftected, errCreate := repo.Decrement(context.Background(), tx, newWallet, int64(amount))
		assert.Error(t, errCreate, errExpected)
		assert.Equal(t, afftected, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execute", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newWallet := model.Wallet{
			ID:         uuid.New(),
			OwnedBy:    uuid.New(),
			Balance:    0,
			Status:     model.WalletStatus.Disabled,
			EnabledAt:  nil,
			DisabledAt: nil,
			CreatedAt:  timestamp,
			UpdatedAt:  timestamp,
		}
		amount := 10000

		var errExpected = errors.New("error")
		mock.
			ExpectPrepare(qDecrementWallet).
			ExpectExec().
			WithArgs(
				amount,
				newWallet.UpdatedAt,
				newWallet.ID,
			).
			WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &walletRepository{db: db}
		afftected, errCreate := repo.Decrement(context.Background(), tx, newWallet, int64(amount))
		assert.Error(t, errCreate)
		assert.Equal(t, afftected, int64(0))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

}
