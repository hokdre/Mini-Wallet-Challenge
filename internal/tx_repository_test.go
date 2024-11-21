package internal

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestTxRepository(t *testing.T) {
	t.Run("Process", TestProcess)
}

func TestProcess(t *testing.T) {
	t.Run("Failed begin", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		var errExpected = errors.New("err")
		mock.ExpectBegin().WillReturnError(errExpected)
		r := txRepository{db: db}
		err = r.Process(context.Background(), nil)
		assert.Error(t, err, errExpected)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success commit", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		f := func(context.Context, *sql.Tx) error {
			return nil
		}

		mock.ExpectCommit()

		r := txRepository{db: db}
		err = r.Process(context.Background(), f)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Success rollback", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()
		mock.ExpectBegin()

		var errExpected = errors.New("err")
		f := func(context.Context, *sql.Tx) error {
			return errExpected
		}

		mock.ExpectRollback()

		r := txRepository{db: db}
		err = r.Process(context.Background(), f)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
