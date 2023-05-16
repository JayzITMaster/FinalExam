package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type admin struct {
	id          int32
	name        string
	email       string
	au_password int32
	created_at  time.Time
}

// setup dependency injection
type AdminModel struct {
	DB *sql.DB //connection pool
}

func (m *AdminModel) InsertAdmin(name, email, password, auth string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		return err
	}

	if auth != "1234" {
		return ErrInvalidAuth
	}

	query := ` 
					INSERT INTO admin_users(users_name, email, au_password_hash)
					VALUES($1, $2, $3)
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = m.DB.ExecContext(ctx, query, name, email, hashedPassword)
	fmt.Println(err)
	if err != nil {
		switch {
		case err.Error() == `ERROR: duplicate key value violates unique constraint "admin_users_email_key" (SQLSTATE 23505)`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}
func (m *AdminModel) AuthenticateAdmin(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	query := `
		SELECT id, au_password_hash
		FROM admin_users
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

func (m *AdminModel) AdminData(ID int) (string, string, error) {
	var userName string
	var email string
	query := `
		SELECT users_name, email
		FROM admin_users
		WHERE ID = $1
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
