package postgres

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetPersonalByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPersonalRepository(db)

	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	personalMock := domain.Personal{
		ID:           1,
		Name:         "Test User",
		Rg:           "123456789",
		Document:     "12345678901",
		Address:      "Test Address",
		City:         "Test City",
		Neighborhood: "Test Neighborhood",
		State:        "TS",
		Cep:          "12345-678",
		Phone:        "11999999999",
		Email:        "test@email.com",
		Website:      "test.com",
		Linkedin:     "test-linkedin",
		Github:       "test-github",
		BirthDate:    birthDate,
		LoginId:      1,
	}

	expectedSQL := regexp.QuoteMeta("SELECT id, name, rg, document, address, neighborhood, state, city, cep, phone, email, website, linkedin, github, birthdate, login_id FROM personal WHERE id=$1")

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "rg", "document", "address", "neighborhood", "state", "city", "cep", "phone", "email", "website", "linkedin", "github", "birthdate", "login_id"}).
			AddRow(personalMock.ID, personalMock.Name, personalMock.Rg, personalMock.Document, personalMock.Address, personalMock.Neighborhood, personalMock.State, personalMock.City, personalMock.Cep, personalMock.Phone, personalMock.Email, personalMock.Website, personalMock.Linkedin, personalMock.Github, personalMock.BirthDate, personalMock.LoginId)

		mock.ExpectQuery(expectedSQL).WithArgs(personalMock.ID).WillReturnRows(rows)

		personal, err := repo.GetPersonalByID(personalMock.ID)

		assert.NoError(t, err)
		assert.Equal(t, personalMock, personal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(expectedSQL).WithArgs(personalMock.ID).WillReturnError(sql.ErrNoRows)

		_, err := repo.GetPersonalByID(personalMock.ID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPersonalNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		dbError := errors.New("database error")
		mock.ExpectQuery(expectedSQL).WithArgs(personalMock.ID).WillReturnError(dbError)

		_, err := repo.GetPersonalByID(personalMock.ID)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetPersonals(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPersonalRepository(db)

	birthDate1 := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	birthDate2 := time.Date(1995, 3, 20, 0, 0, 0, 0, time.UTC)

	personalMock1 := domain.Personal{
		ID:           1,
		Name:         "User One",
		Rg:           "111111111",
		Document:     "11111111111",
		Address:      "Address 1",
		City:         "City 1",
		Neighborhood: "Neighborhood 1",
		State:        "S1",
		Cep:          "11111-111",
		Phone:        "11111111111",
		Email:        "user1@email.com",
		Website:      "user1.com",
		Linkedin:     "user1-linkedin",
		Github:       "user1-github",
		BirthDate:    birthDate1,
		LoginId:      1,
	}

	personalMock2 := domain.Personal{
		ID:           2,
		Name:         "User Two",
		Rg:           "222222222",
		Document:     "22222222222",
		Address:      "Address 2",
		City:         "City 2",
		Neighborhood: "Neighborhood 2",
		State:        "S2",
		Cep:          "22222-222",
		Phone:        "22222222222",
		Email:        "user2@email.com",
		Website:      "user2.com",
		Linkedin:     "user2-linkedin",
		Github:       "user2-github",
		BirthDate:    birthDate2,
		LoginId:      2,
	}

	expectedSQL := regexp.QuoteMeta("SELECT id, name, rg, document, address, neighborhood, state, city, cep, phone, email, website, linkedin, github, birthdate, login_id FROM personal")

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "rg", "document", "address", "neighborhood", "state", "city", "cep", "phone", "email", "website", "linkedin", "github", "birthdate", "login_id"}).
			AddRow(personalMock1.ID, personalMock1.Name, personalMock1.Rg, personalMock1.Document, personalMock1.Address, personalMock1.Neighborhood, personalMock1.State, personalMock1.City, personalMock1.Cep, personalMock1.Phone, personalMock1.Email, personalMock1.Website, personalMock1.Linkedin, personalMock1.Github, personalMock1.BirthDate, personalMock1.LoginId).
			AddRow(personalMock2.ID, personalMock2.Name, personalMock2.Rg, personalMock2.Document, personalMock2.Address, personalMock2.Neighborhood, personalMock2.State, personalMock2.City, personalMock2.Cep, personalMock2.Phone, personalMock2.Email, personalMock2.Website, personalMock2.Linkedin, personalMock2.Github, personalMock2.BirthDate, personalMock2.LoginId)

		mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

		personals, err := repo.GetPersonals()

		assert.NoError(t, err)
		assert.Len(t, personals, 2)
		assert.Equal(t, personalMock1, personals[0])
		assert.Equal(t, personalMock2, personals[1])
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Query Error", func(t *testing.T) {
		dbError := errors.New("database error")
		mock.ExpectQuery(expectedSQL).WillReturnError(dbError)

		_, err := repo.GetPersonals()

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty List", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "rg", "document", "address", "neighborhood", "state", "city", "cep", "phone", "email", "website", "linkedin", "github", "birthdate", "login_id"})

		mock.ExpectQuery(expectedSQL).WillReturnRows(rows)

		personals, err := repo.GetPersonals()

		assert.NoError(t, err)
		assert.Len(t, personals, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCreatePersonal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPersonalRepository(db)

	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	personalToCreate := domain.Personal{
		Name:         "New User",
		Rg:           "999999999",
		Document:     "99999999999",
		Address:      "New Address",
		City:         "New City",
		Neighborhood: "New Neighborhood",
		State:        "NC",
		Cep:          "99999-999",
		Phone:        "99999999999",
		Email:        "newuser@email.com",
		Website:      "newuser.com",
		Linkedin:     "newuser-linkedin",
		Github:       "newuser-github",
		BirthDate:    birthDate,
		LoginId:      5,
	}

	expectedSQL := regexp.QuoteMeta("INSERT INTO personal " +
		"(name, rg, document, address, city, neighborhood, state, cep, phone, email, website, linkedin, github, birthdate, login_id) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectQuery().
			WithArgs(
				personalToCreate.Name, personalToCreate.Rg, personalToCreate.Document,
				personalToCreate.Address, personalToCreate.City,
				personalToCreate.Neighborhood, personalToCreate.State,
				personalToCreate.Cep, personalToCreate.Phone, personalToCreate.Email,
				personalToCreate.Website, personalToCreate.Linkedin,
				personalToCreate.Github, personalToCreate.BirthDate, personalToCreate.LoginId,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(10))

		id, err := repo.CreatePersonal(personalToCreate)

		assert.NoError(t, err)
		assert.Equal(t, 10, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		dbError := errors.New("database error")

		mock.ExpectPrepare(expectedSQL).ExpectQuery().WillReturnError(dbError)

		_, err := repo.CreatePersonal(personalToCreate)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUpdatePersonal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPersonalRepository(db)

	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	personalToUpdate := domain.Personal{
		ID:           1,
		Name:         "Updated User",
		Rg:           "888888888",
		Document:     "88888888888",
		Address:      "New Address",
		City:         "New City",
		Neighborhood: "New Neighborhood",
		State:        "NS",
		Cep:          "00000-000",
		Phone:        "999999999",
		Email:        "new@email.com",
		Website:      "new.com",
		Linkedin:     "new-linkedin",
		Github:       "new-github",
		BirthDate:    birthDate,
		LoginId:      2,
	}

	expectedSQL := regexp.QuoteMeta("UPDATE personal SET " +
		"name=$1, rg=$2, document=$3, address=$4, city=$5, neighborhood=$6, state=$7, cep=$8, phone=$9, email=$10, website=$11, linkedin=$12, github=$13, birthdate=$14, login_id=$15 " +
		"WHERE id=$16")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(
				personalToUpdate.Name, personalToUpdate.Rg, personalToUpdate.Document,
				personalToUpdate.Address, personalToUpdate.City,
				personalToUpdate.Neighborhood, personalToUpdate.State,
				personalToUpdate.Cep, personalToUpdate.Phone, personalToUpdate.Email,
				personalToUpdate.Website, personalToUpdate.Linkedin,
				personalToUpdate.Github, personalToUpdate.BirthDate, personalToUpdate.LoginId,
				personalToUpdate.ID,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		updatedPersonal, err := repo.UpdatePersonal(personalToUpdate)

		assert.NoError(t, err)
		assert.Equal(t, personalToUpdate, updatedPersonal)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(
				personalToUpdate.Name, personalToUpdate.Rg, personalToUpdate.Document,
				personalToUpdate.Address, personalToUpdate.City,
				personalToUpdate.Neighborhood, personalToUpdate.State,
				personalToUpdate.Cep, personalToUpdate.Phone, personalToUpdate.Email,
				personalToUpdate.Website, personalToUpdate.Linkedin,
				personalToUpdate.Github, personalToUpdate.BirthDate, personalToUpdate.LoginId,
				personalToUpdate.ID,
			).
			WillReturnResult(sqlmock.NewResult(0, 0))

		_, err := repo.UpdatePersonal(personalToUpdate)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPersonalNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		dbError := errors.New("database error")

		mock.ExpectPrepare(expectedSQL).ExpectExec().WillReturnError(dbError)

		_, err := repo.UpdatePersonal(personalToUpdate)

		assert.Error(t, err)
		assert.Equal(t, dbError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDeletePersonal(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewPersonalRepository(db)
	personalID := 1

	expectedSQL := regexp.QuoteMeta("DELETE FROM personal WHERE id=$1")

	t.Run("Success", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(personalID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.DeletePersonal(personalID)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).
			ExpectExec().
			WithArgs(personalID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeletePersonal(personalID)

		assert.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrPersonalNotFound)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Exec Error", func(t *testing.T) {
		mock.ExpectPrepare(expectedSQL).ExpectExec().WillReturnError(errors.New("db error"))
		err := repo.DeletePersonal(personalID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
