package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/joelsouza82/curriculum-go/internal/core/domain"
	"github.com/joelsouza82/curriculum-go/internal/core/port/out"
)

type PersonalRepository struct {
	connection *sql.DB
}

func NewPersonalRepository(connection *sql.DB) out.PersonalRepository {
	return &PersonalRepository{
		connection: connection,
	}
}

func (pr *PersonalRepository) GetPersonalByID(id int) (domain.Personal, error) {
	query := "SELECT id, name, rg, document, address, neighborhood, state, city, cep, phone, email, website, linkedin, github, birthdate, login_id FROM personal WHERE id=$1"
	row := pr.connection.QueryRow(query, id)

	var personalObj domain.Personal
	err := row.Scan(
		&personalObj.ID,
		&personalObj.Name,
		&personalObj.Rg,
		&personalObj.Document,
		&personalObj.Address,
		&personalObj.Neighborhood,
		&personalObj.State,
		&personalObj.City,
		&personalObj.Cep,
		&personalObj.Phone,
		&personalObj.Email,
		&personalObj.Website,
		&personalObj.Linkedin,
		&personalObj.Github,
		&personalObj.BirthDate,
		&personalObj.LoginId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Personal{}, domain.ErrPersonalNotFound
		}
		fmt.Println(err)
		return domain.Personal{}, err
	}

	return personalObj, nil
}

func (pr *PersonalRepository) GetPersonals() ([]domain.Personal, error) {
	query := "SELECT id, name, rg, document, address, neighborhood, state, city, cep, phone, email, website, linkedin, github, birthdate, login_id FROM personal"
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []domain.Personal{}, err
	}

	var personalList []domain.Personal
	var personalObj domain.Personal

	for rows.Next() {
		err = rows.Scan(
			&personalObj.ID,
			&personalObj.Name,
			&personalObj.Rg,
			&personalObj.Document,
			&personalObj.Address,
			&personalObj.Neighborhood,
			&personalObj.State,
			&personalObj.City,
			&personalObj.Cep,
			&personalObj.Phone,
			&personalObj.Email,
			&personalObj.Website,
			&personalObj.Linkedin,
			&personalObj.Github,
			&personalObj.BirthDate,
			&personalObj.LoginId,
		)

		if err != nil {
			fmt.Println(err)
			return []domain.Personal{}, err
		}

		personalList = append(personalList, personalObj)
	}

	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return []domain.Personal{}, err
	}

	rows.Close()

	return personalList, nil
}

func (pr *PersonalRepository) CreatePersonal(personal domain.Personal) (int, error) {
	var id int
	query, err := pr.connection.Prepare("INSERT INTO personal " +
		"(name, rg, document, address, city, neighborhood, state, cep, phone, email, website, linkedin, github, birthdate, login_id) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(
		personal.Name,
		personal.Rg,
		personal.Document,
		personal.Address,
		personal.City,
		personal.Neighborhood,
		personal.State,
		personal.Cep,
		personal.Phone,
		personal.Email,
		personal.Website,
		personal.Linkedin,
		personal.Github,
		personal.BirthDate,
		personal.LoginId,
	).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()
	return id, nil
}

func (pr *PersonalRepository) UpdatePersonal(personal domain.Personal) (domain.Personal, error) {
	query := "UPDATE personal SET " +
		"name=$1, rg=$2, document=$3, address=$4, city=$5, neighborhood=$6, state=$7, cep=$8, phone=$9, email=$10, website=$11, linkedin=$12, github=$13, birthdate=$14, login_id=$15 " +
		"WHERE id=$16"

	stmt, err := pr.connection.Prepare(query)
	if err != nil {
		fmt.Println(err)
		return domain.Personal{}, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		personal.Name,
		personal.Rg,
		personal.Document,
		personal.Address,
		personal.City,
		personal.Neighborhood,
		personal.State,
		personal.Cep,
		personal.Phone,
		personal.Email,
		personal.Website,
		personal.Linkedin,
		personal.Github,
		personal.BirthDate,
		personal.LoginId,
		personal.ID,
	)
	if err != nil {
		fmt.Println(err)
		return domain.Personal{}, err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.Personal{}, domain.ErrPersonalNotFound
	}

	return personal, nil
}

func (pr *PersonalRepository) DeletePersonal(id int) error {
	query := "DELETE FROM personal WHERE id=$1"

	stmt, err := pr.connection.Prepare(query)
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
		return domain.ErrPersonalNotFound
	}

	return nil
}
