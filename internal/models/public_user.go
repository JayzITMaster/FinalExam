package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNoRecord           = errors.New("no matching record found")
	ErrInvalidCredentials = errors.New("invalid credential")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidAuth        = errors.New("Authentication Code is not Valid")
)

type PublicUser struct {
	id          int32
	email       string
	pu_password int32
}

// setup dependency injection
type PublicUserModel struct {
	DB *sql.DB //connection pool
}

func (m *PublicUserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	query := ` 
					INSERT INTO public_user(users_name, email, pu_password_hash)
					VALUES($1, $2, $3)
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, query, name, email, hashedPassword)
	fmt.Println(err)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "public_user_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m *PublicUserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	query := `
		SELECT id, pu_password_hash
		FROM public_user
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&id, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	//password is correct
	return id, nil

}

func (m *PublicUserModel) UserData(ID int) (string, string, error) {
	var userName string
	var email string
	query := `
		SELECT users_name, email
		FROM public_user
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, ID).Scan(&userName, &email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", err
		}
	}
	//we got the admin
	return userName, email, nil

}
