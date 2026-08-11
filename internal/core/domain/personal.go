package domain

import "time"

type Personal struct {
	ID           int       `json:"id_personal"`
	Name         string    `json:"name"`
	Rg           string    `json:"rg"`
	Document     string    `json:"document"`
	Address      string    `json:"address"`
	City         string    `json:"city"`
	Neighborhood string    `json:"neighborhood"`
	State        string    `json:"state"`
	Cep          string    `json:"cep"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	Website      string    `json:"website"`
	Linkedin     string    `json:"linkedin"`
	Github       string    `json:"github"`
	BirthDate    time.Time `json:"birthdate"`
	LoginId      int       `json:"login_id"`
}
