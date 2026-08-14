package postgres

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetLoginByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewLoginRepository(db)

	loginMock := domain.Login{
		ID:       1,
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedSQL := regexp.QuoteMeta("SELECT id, email, password FROM login WHERE id=$1")

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(loginMock.ID, loginMock.Email, loginMock.Password)

		mock.ExpectQuery(expectedSQL).WithArgs(loginMock.ID).WillReturnRows(rows)

		login, err := repo.GetLoginByID(loginMock.ID)

		assert.NoError(t, err)
		assert.Equal(t, loginMock, login)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(expectedSQL).WithArgs(loginMock.ID).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetLoginByID(loginMock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrLoginNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		dbError := errors.New("database error")
		mock.ExpectQuery(expectedSQL).WithArgs(loginMock.ID).WillReturnError(dbError)

		_, err := repo.GetLoginByID(loginMock.ID)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetLoginByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewLoginRepository(db)

	loginMock := domain.Login{
		ID:       1,
		Email:    "test@example.com",
		Password: "password123",
	}

	expectedSQL := regexp.QuoteMeta("SELECT id, email, password FROM login WHERE email=$1")

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(loginMock.ID, loginMock.Email, loginMock.Password)

		mock.ExpectQuery(expectedSQL).WithArgs(loginMock.Email).WillReturnRows(rows)

		login, err := repo.GetLoginByEmail(loginMock.Email)

		assert.NoError(t, err)
		assert.Equal(t, loginMock, login)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(expectedSQL).WithArgs(loginMock.Email).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetLoginByEmail(loginMock.Email)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrLoginNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		dbError := errors.New("database error")
		mock.ExpectQuery(expectedSQL).WithArgs(loginMock.Email).WillReturnError(dbError)

		_, err := repo.GetLoginByEmail(loginMock.Email)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetLogins(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewLoginRepository(db)

	loginMock1 := domain.Login{
		ID:       1,
		Email:    "user1@example.com",
		Password: "pass1",
	}

	loginMock2 := domain.Login{
		ID:       2,
		Email:    "user2@example.com",
		Password: "pass2",
	}

	expectedSQL := regexp.QuoteMeta("SELECT id, email, password FROM login")

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"}).
			AddRow(loginMock1.ID, loginMock1.Email, loginMock1.Password).
			AddRow(loginMock2.ID, loginMock2.Email, loginMock2.Password)

		mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

		logins, err := repo.GetLogins()

		assert.NoError(t, err)
		assert.Len(t, logins, 2)
		assert.Equal(t, loginMock1, logins[0])
		assert.Equal(t, loginMock2, logins[1])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		dbError := errors.New("database error")
		mock.ExpectQuery(expectedSQL).WillReturnError(dbError)

		_, err := repo.GetLogins()

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty List", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "password"})

		mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

		logins, err := repo.GetLogins()

		assert.NoError(t, err)
		assert.Len(t, logins, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreateLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewLoginRepository(db)

	loginToCreate := domain.Login{
		Email:    "newuser@example.com",
		Password: "newpass123",
	}

	expectedSQL := regexp.QuoteMeta("INSERT INTO login (email, password) VALUES ($1, $2) RETURNING id")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectQuery().
			WithArgs(loginToCreate.Email, loginToCreate.Password).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

		id, err := repo.CreateLogin(loginToCreate)

		assert.NoError(t, err)
		assert.Equal(t, 10, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		dbError := errors.New("database error")

		mock.ExpectPrepare(expectedSQL).ExpectQuery().WillReturnError(dbError)

		_, err := repo.CreateLogin(loginToCreate)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdateLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewLoginRepository(db)

	loginToUpdate := domain.Login{
		ID:       1,
		Email:    "updated@example.com",
		Password: "updatedpass",
	}

	expectedSQL := regexp.QuoteMeta("UPDATE login SET email=$1, password=$2 WHERE id=$3")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(loginToUpdate.Email, loginToUpdate.Password, loginToUpdate.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		updatedLogin, err := repo.UpdateLogin(loginToUpdate)

		assert.NoError(t, err)
		assert.Equal(t, loginToUpdate, updatedLogin)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(loginToUpdate.Email, loginToUpdate.Password, loginToUpdate.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := repo.UpdateLogin(loginToUpdate)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrLoginNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		dbError := errors.New("database error")

		mock.ExpectPrepare(expectedSQL).ExpectExec().WillReturnError(dbError)

		_, err := repo.UpdateLogin(loginToUpdate)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeleteLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewLoginRepository(db)
	loginID := 1

	expectedSQL := regexp.QuoteMeta("DELETE FROM login WHERE id=$1")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(loginID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeleteLogin(loginID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(loginID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteLogin(loginID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrLoginNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).ExpectExec().WillReturnError(errors.New("db error"))
		err := repo.DeleteLogin(loginID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
