package models

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type Form struct {
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
type FormModel struct {
	DB *sql.DB //connection pool
}

func (m *FormModel) Insert(form *Form) (int64, error) {
	var id int64
	//Sql statement that will be ran to insert data, will return the Form_id which is serial and auto increments
	statement := `
	INSERT INTO form (form_status, archive_status,affiant_full_name,other_names,name_change_status,previous_name,reason_for_change,social_security_num, social_security_date,social_security_country,passport_number,passport_date, passport_country,NHI_number, NHI_date, NHI_country, dob,place_of_birth,nationality,acquired_nationality,spouse_name,affiants_address,residencial_phone_number,residenceial_fax_num,residencial_email, user_id)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
	RETURNING form_id
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)
	err := m.DB.QueryRowContext(ctx, statement, form.Form_status, form.Archive_status, form.Affiant_full_name, form.Other_names, form.Name_change_status, form.Previous_name, form.Reason_for_Change, form.Social_security_num, form.Social_security_date, form.Social_security_country, form.Passport_number, form.Passport_date, form.Passport_country, form.NHI_number, form.NHI_date, form.NHI_country, form.Dob, form.Place_of_birth, form.Nationality, form.Acquired_nationality, form.Spouse_name, form.Affiants_address, form.Residencial_phone_number, form.Residenceial_fax_num, form.Residencial_email, form.User_id).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err

	}

	return id, err

}

func (m *FormModel) Read(form_id int64) ([]*Form, error) {
	//create SQL statement
	statement := `
		SELECT *
		FROM form 
		WHERE form_id = $1
	`
	rows, err := m.DB.Query(statement, form_id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//cleanup before we leave our read method
	defer rows.Close()
	// pointer to every piece of data in the form
	formes := []*Form{}

	rows.Next()

	form := &Form{}
	err = rows.Scan(&form.User_id, &form.Form_id, &form.Form_status, &form.Archive_status, &form.Affiant_full_name, &form.Other_names, &form.Previous_name, &form.Name_change_status, &form.Reason_for_Change, &form.Social_security_num, &form.Social_security_date, &form.Social_security_country, &form.Passport_number, &form.Passport_date, &form.Passport_country, &form.NHI_number, &form.NHI_date, &form.NHI_country, &form.Dob, &form.Place_of_birth, &form.Nationality, &form.Acquired_nationality, &form.Spouse_name, &form.Affiants_address, &form.Residencial_phone_number, &form.Residenceial_fax_num, &form.Residencial_email, &form.CreateOn)

	formes = append(formes, form)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//check to see if there were error generated
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return formes, nil
}

func (m *FormModel) Update(form *Form) (int64, error) {
	var id int64

	//Sql statement that will be ran to insert data, will return the Form_id which is serial and auto increments
	statement := `
	UPDATE form
	SET form_status = $1, archive_status = $2, affiant_full_name = $3,other_names = $4,name_change_status = $5,previous_name = $6,reason_for_change = $7,social_security_num = $8, social_security_date = $9,social_security_country = $10,passport_number = $11,passport_date = $12, passport_country = $13,NHI_number = $14, NHI_date = $15, NHI_country = $16, dob = $17,place_of_birth = $18,nationality = $19,acquired_nationality = $20,spouse_name = $21,affiants_address = $22,residencial_phone_number = $23,residenceial_fax_num = $24,residencial_email = $25
	WHERE form_id = $26
	
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)
	_, err := m.DB.ExecContext(ctx, statement, form.Form_status, form.Archive_status, form.Affiant_full_name, form.Other_names, form.Name_change_status, form.Previous_name, form.Reason_for_Change, form.Social_security_num, form.Social_security_date, form.Social_security_country, form.Passport_number, form.Passport_date, form.Passport_country, form.NHI_number, form.NHI_date, form.NHI_country, form.Dob, form.Place_of_birth, form.Nationality, form.Acquired_nationality, form.Spouse_name, form.Affiants_address, form.Residencial_phone_number, form.Residenceial_fax_num, form.Residencial_email, form.Form_id)

	if err != nil {
		fmt.Println(err)
		return 0, err

	}

	return id, err

}

func (m *FormModel) Delete(form_id int64) error {
	//Sql statement that will be ran to insert data, will return the Form_id which is serial and auto increments
	statement := `
	DELETE FROM form 
	WHERE form_id = $1
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)
	_, err := m.DB.ExecContext(ctx, statement, form_id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err

}

func (m *FormModel) LoadData() ([]*Form, error) {
	//create SQL statement
	statement := `
		SELECT *
		FROM form 

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

func (m *FormModel) IndividualLoadData(ID int) ([]*Form, error) {
	//create SQL statement
	statement := `
		SELECT *
		FROM form 
		WHERE user_id = $1
		
	`
	rows, err := m.DB.Query(statement, ID)

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

func (m *FormModel) AcceptForm(form_id int64) error {
	//Sql statement that will be ran to change the form status to accepted
	statement := `
	UPDATE form
	SET form_status = $1
	WHERE form_id = $2
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)
	status := "accepted"
	_, err := m.DB.ExecContext(ctx, statement, status, form_id)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

func (m *FormModel) DenyForm(form_id int64) error {
	//Sql statement that will be ran to change the form status to accepted
	statement := `
	UPDATE form
	SET form_status = $1
	WHERE form_id = $2
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)
	status := "denied"
	_, err := m.DB.ExecContext(ctx, statement, status, form_id)

	if err != nil {
		fmt.Println(err)

		fmt.Println(err)
		return err
	}

	return err
}

func (m *FormModel) Stats() (string, string, string, string, error) {
	var Accepted int
	var unverified int
	var Denied int
	var Total int

	StatUnverified := "unverified"
	StatAccepted := "accepted"
	StatDenied := "denied"

	query := ` 
	SELECT COUNT(form_status)
	FROM form
    WHERE form_status = $1
	`

	query2 := ` 
	SELECT COUNT(form_status)
	FROM form
    WHERE form_status = $1
	`

	query3 := ` 
	SELECT COUNT(form_status)
	FROM form
    WHERE form_status = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, StatUnverified).Scan(&unverified)

	if err != nil {
		return "", "", "", "", err
	}
	err = m.DB.QueryRowContext(ctx, query2, StatDenied).Scan(&Denied)

	if err != nil {
		return "", "", "", "", err
	}
	err = m.DB.QueryRowContext(ctx, query3, StatAccepted).Scan(&Accepted)

	Total = Accepted + unverified + Denied

	return strconv.Itoa(Accepted), strconv.Itoa(unverified), strconv.Itoa(Denied), strconv.Itoa(Total), nil

}

func (m *FormModel) NormalStats(ID int) (string, string, string, string, error) {
	var Accepted int
	var unverified int
	var Denied int
	var Total int

	StatUnverified := "unverified"
	StatAccepted := "accepted"
	StatDenied := "denied"

	query := ` 
	SELECT COUNT(form_status)
	FROM form
    WHERE form_status = $1
	AND user_id = $2
	`

	query2 := ` 
	SELECT COUNT(form_status)
	FROM form
    WHERE form_status = $1
	AND user_id = $2
	`

	query3 := ` 
	SELECT COUNT(form_status)
	FROM form
    WHERE form_status = $1
	AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, StatUnverified, ID).Scan(&unverified)

	if err != nil {
		return "", "", "", "", err
	}
	err = m.DB.QueryRowContext(ctx, query2, StatDenied, ID).Scan(&Denied)

	if err != nil {
		return "", "", "", "", err
	}
	err = m.DB.QueryRowContext(ctx, query3, StatAccepted, ID).Scan(&Accepted)

	Total = Accepted + unverified + Denied

	return strconv.Itoa(Accepted), strconv.Itoa(unverified), strconv.Itoa(Denied), strconv.Itoa(Total), nil

}

func (m *FormModel) FormOwner(form_id int64) (int, error) {

	var ID int
	//Sql statement that will be ran to change the form status to accepted
	statement := `
	SELECT user_id
	FROM form
	WHERE form_id = $1
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, statement, form_id).Scan(&ID)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return ID, err
}
