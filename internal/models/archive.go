package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type archive struct {
	User_id                  int32
	Form_id                  int64
	Form_status              string
	Archive_status           bool
	Affiant_full_name        string
	Other_names              string
	Name_change_status       bool
	Previous_name            string
	Reason_for_Change        string
	Social_security_num      int64
	Social_security_date     time.Time
	Social_security_country  string
	Passport_number          int64
	Passport_date            time.Time
	Passport_country         string
	NHI_number               int64
	NHI_date                 time.Time
	NHI_country              string
	Dob                      time.Time
	Place_of_birth           string
	Nationality              string
	Acquired_nationality     string
	Spouse_name              string
	Affiants_address         string
	Residencial_phone_number int64
	Residenceial_fax_num     int64
	Residencial_email        string
	CreateOn                 time.Time
}

// setup dependency injection
type ArchiveModel struct {
	DB *sql.DB //connection pool
}

func (m *ArchiveModel) ArchiveForm(form_id int64) error {
	//Sql statement that will be ran to change the form status to accepted
	statement1 := `
	INSERT INTO archive
	SELECT * FROM form
	WHERE form_id = $1; 
	`
	statement2 := `
	DELETE FROM form 
	WHERE form_id = $1
	`

	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)

	_, err := m.DB.ExecContext(ctx, statement1, form_id)
	_, err = m.DB.ExecContext(ctx, statement2, form_id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (m *ArchiveModel) UnArchiveForm(form_id int64) error {
	//Sql statement that will be ran to change the form status to accepted
	statement1 := `
	INSERT INTO form
	SELECT * FROM archive
	WHERE form_id = $1; 
	`
	statement2 := `
	DELETE FROM archive
	WHERE form_id = $1
	`

	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)

	_, err := m.DB.ExecContext(ctx, statement1, form_id)
	_, err = m.DB.ExecContext(ctx, statement2, form_id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (m *FormModel) LoadArchiveData() ([]*Form, error) {
	//create SQL statement
	statement := `
		SELECT *
		FROM archive
	`
	rows, err := m.DB.Query(statement)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//cleanup before we leave our read method
	defer rows.Close()
	// pointer to every piece of data in the form
	formes := []*Form{}

	for rows.Next() {
		form := &Form{}
		err = rows.Scan(&form.User_id, &form.Form_id, &form.Form_status, &form.Archive_status, &form.Affiant_full_name, &form.Other_names, &form.Previous_name, &form.Name_change_status, &form.Reason_for_Change, &form.Social_security_num, &form.Social_security_date, &form.Social_security_country, &form.Passport_number, &form.Passport_date, &form.Passport_country, &form.NHI_number, &form.NHI_date, &form.NHI_country, &form.Dob, &form.Place_of_birth, &form.Nationality, &form.Acquired_nationality, &form.Spouse_name, &form.Affiants_address, &form.Residencial_phone_number, &form.Residenceial_fax_num, &form.Residencial_email, &form.CreateOn)
		formes = append(formes, form) //contain first row
	}

	if err != nil {
		return nil, err
	}
	//check to see if there were error generated
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return formes, nil
}

func (m *ArchiveModel) Stats() (string, error) {

	var Total int

	query := ` 
	SELECT COUNT(*)
	FROM form
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query).Scan(&Total)

	if err != nil {
		return "", err
	}

	return strconv.Itoa(Total), nil

}
