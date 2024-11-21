package account

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
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository(t *testing.T) {
	t.Run("CreateTx", TestCreateTx)
	t.Run("Get", TestGet)
}

func TestCreateTx(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		timestamp := time.Now()
		newAcc := model.Account{
			ID:                 uuid.New(),
			ExternalCustomerID: uuid.NewString(),
			CreatedAt:          timestamp,
			UpdatedAt:          timestamp,
		}
		mock.
			ExpectPrepare(qCreate).
			ExpectExec().
			WithArgs(
				newAcc.ID,
				newAcc.ExternalCustomerID,
				newAcc.CreatedAt,
				newAcc.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &accountRepository{}
		errCreate := repo.CreateTx(context.Background(), tx, newAcc)
		assert.NoError(t, errCreate)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Prepare", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		errExpected := errors.New("Failed create prepate statement")
		newAcc := model.Account{}
		mock.
			ExpectPrepare(qCreate).WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &accountRepository{}
		errCreate := repo.CreateTx(context.Background(), tx, newAcc)
		assert.Error(t, errCreate, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed Execute", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		errExpected := errors.New("Failed to excecute statement")
		newAcc := model.Account{}
		mock.
			ExpectPrepare(qCreate).
			ExpectExec().WillReturnError(errExpected)

		tx, err := db.Begin()
		assert.NoError(t, err)

		repo := &accountRepository{}
		errCreate := repo.CreateTx(context.Background(), tx, newAcc)
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
		acc := model.Account{
			ID:                 uuid.New(),
			ExternalCustomerID: uuid.NewString(),
			CreatedAt:          timeStamp,
			UpdatedAt:          timeStamp,
		}

		expectedRow := sqlmock.NewRows([]string{
			"id", "external_id", "created_at", "updated_at",
		}).AddRow(acc.ID, acc.ExternalCustomerID, acc.CreatedAt, acc.UpdatedAt)

		filter := internal.AccountFilter{
			IDs:         []string{acc.ID.String()},
			ExternalIDs: []string{acc.ExternalCustomerID},
		}
		mock.ExpectQuery(qGet).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.ExternalIDs),
			1,
			defaultOffset,
		).WillReturnRows(expectedRow)

		repo := &accountRepository{db: db}
		accResult, err := repo.Get(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, acc, accResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success with no filter", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		timeStamp := time.Now()
		acc := model.Account{
			ID:                 uuid.New(),
			ExternalCustomerID: uuid.NewString(),
			CreatedAt:          timeStamp,
			UpdatedAt:          timeStamp,
		}

		expectedRow := sqlmock.NewRows([]string{
			"id", "external_id", "created_at", "updated_at",
		}).AddRow(acc.ID, acc.ExternalCustomerID, acc.CreatedAt, acc.UpdatedAt)

		filter := internal.AccountFilter{}
		mock.ExpectQuery(qGet).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.ExternalIDs),
			1,
			defaultOffset,
		).WillReturnRows(expectedRow)

		repo := &accountRepository{db: db}
		accResult, err := repo.Get(context.Background(), filter)
		assert.NoError(t, err)
		assert.Equal(t, acc, accResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Failed, not found", func(t *testing.T) {
		db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
		assert.NoError(t, err)
		defer db.Close()

		filter := internal.AccountFilter{}
		mock.ExpectQuery(qGet).WithArgs(
			pq.Array(filter.IDs),
			pq.Array(filter.ExternalIDs),
			1,
			defaultOffset,
		).WillReturnError(sql.ErrNoRows)

		repo := &accountRepository{db: db}
		accResult, err := repo.Get(context.Background(), filter)
		assert.Error(t, err, sql.ErrNoRows)
		assert.Equal(t, model.Account{}, accResult)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
