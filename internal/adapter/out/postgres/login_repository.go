package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/core/port/out"
)

type LoginRepository struct {
	connection *sql.DB
}

func NewLoginRepository(connection *sql.DB) out.LoginRepository {
	return &LoginRepository{
		connection: connection,
	}
}

func (lr *LoginRepository) GetLogins() ([]domain.Login, error) {
	query := "SELECT id, email, password FROM login"
	rows, err := lr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []domain.Login{}, err
	}
	defer rows.Close()

	var loginList []domain.Login
	for rows.Next() {
		var loginObj domain.Login
		err = rows.Scan(
			&loginObj.ID,
			&loginObj.Email,
			&loginObj.Password,
		)
		if err != nil {
			fmt.Println(err)
			return []domain.Login{}, err
		}
		loginList = append(loginList, loginObj)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return []domain.Login{}, err
	}

	return loginList, nil
}

func (lr *LoginRepository) GetLoginByID(id int) (domain.Login, error) {
	query := "SELECT id, email, password FROM login WHERE id=$1"
	row := lr.connection.QueryRow(query, id)

	var loginObj domain.Login
	err := row.Scan(
		&loginObj.ID,
		&loginObj.Email,
		&loginObj.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Login{}, domain.ErrLoginNotFound
		}
		fmt.Println(err)
		return domain.Login{}, err
	}

	return loginObj, nil
}

func (lr *LoginRepository) GetLoginByEmail(email string) (domain.Login, error) {
	query := "SELECT id, email, password FROM login WHERE email=$1"
	row := lr.connection.QueryRow(query, email)

	var loginObj domain.Login
	err := row.Scan(
		&loginObj.ID,
		&loginObj.Email,
		&loginObj.Password,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Login{}, domain.ErrLoginNotFound
		}
		fmt.Println(err)
		return domain.Login{}, err
	}

	return loginObj, nil
}

func (lr *LoginRepository) CreateLogin(login domain.Login) (int, error) {
	var id int
	query, err := lr.connection.Prepare("INSERT INTO login (email, password) VALUES ($1, $2) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer query.Close()

	err = query.QueryRow(
		login.Email,
		login.Password,
	).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return id, nil
}

func (lr *LoginRepository) UpdateLogin(login domain.Login) (domain.Login, error) {
	query := "UPDATE login SET email=$1, password=$2 WHERE id=$3"

	stmt, err := lr.connection.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return domain.Login{}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		login.Email,
		login.Password,
		login.ID,
	)
	if err != nil {
		fmt.Println(err)
		return domain.Login{}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.Login{}, domain.ErrLoginNotFound
	}

	return login, nil
}

func (lr *LoginRepository) DeleteLogin(id int) error {
	query := "DELETE FROM login WHERE id=$1"

	stmt, err := lr.connection.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrLoginNotFound
	}

	return nil
}
