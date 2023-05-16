package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Comments struct {
	Admin_id  int64
	Form_id   int64
	Comment1  string
	Comment2  string
	Comment3  string
	Comment35 string
	Comment4  string
	Comment5  string
	Comment6  string
	Comment7  string
	Comment8  string
	Comment9  string
	Comment10 string
	Edit_made time.Time
}

// setup dependency injection
type CommentsModel struct {
	DB *sql.DB //connection pool
}

func (m *CommentsModel) InsertComment(comments *Comments) error {

	query := ` 
					INSERT INTO comments(admin_id, form_id, comments1, comments2, comments3, comments35, comments4, comments5, comments6, comments7, comments8, comments9, comments10)
					VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, comments.Admin_id, comments.Form_id, comments.Comment1, comments.Comment2, comments.Comment3, comments.Comment35, comments.Comment4, comments.Comment5, comments.Comment6, comments.Comment7, comments.Comment8, comments.Comment9, comments.Comment10)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (m *CommentsModel) GetComments(form_id int64) ([]*Comments, error) {
	//create SQL statement
	statement := `
		SELECT *
		FROM comments
		WHERE form_id = $1
	`

	rows, err := m.DB.Query(statement, form_id)
	if err != nil {

		commentsByte := []*Comments{}
		comments := &Comments{

			Comment1:  "",
			Comment2:  "",
			Comment3:  "",
			Comment35: "",
			Comment4:  "",
			Comment5:  "",
			Comment6:  "",
			Comment7:  "",
			Comment8:  "",
			Comment9:  "",
			Comment10: "",
		}
		commentsByte = append(commentsByte, comments)
		return commentsByte, nil
	} else {

		//cleanup before we leave our read method
		defer rows.Close()
		// pointer to every piece of data in the form
		commentsByte := []*Comments{}

		rows.Next()

		comment := &Comments{}
		err = rows.Scan(&comment.Admin_id, &comment.Form_id, &comment.Comment1, &comment.Comment2, &comment.Comment3, &comment.Comment35, &comment.Comment4, &comment.Comment5, &comment.Comment6, &comment.Comment7, &comment.Comment8, &comment.Comment9, &comment.Comment10, &comment.Edit_made)

		commentsByte = append(commentsByte, comment)
		return commentsByte, nil
	}
}

func (m *CommentsModel) Exists(form_id int64) bool {

	var count int

	query := ` 
	SELECT COUNT(form_id)
	FROM comments
    WHERE form_id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, form_id).Scan(&count)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(count)
	if count == 0 {
		return false
	} else {
		return true
	}

}

func (m *CommentsModel) UpdateComments(comments *Comments) (int64, error) {
	var id int64

	//Sql statement that will be ran to insert data, will return the Form_id which is serial and auto increments
	statement := `
	UPDATE comments
	SET comments1 = $1, comments2 = $2, comments3 = $3, comments35 = $4, comments4 = $5, comments5 = $6, comments6 = $7, comments7 = $8, comments8 =$9, comments9 =$10, comments10 = $11
	WHERE form_id = $12
	
 `
	//sets the timeout for the DB connection
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//sets the timeout for the DB connection, passes the statements and associates the arguements with the place holders in the SQL ($1, $2 etc)
	_, err := m.DB.ExecContext(ctx, statement, comments.Comment1, comments.Comment2, comments.Comment3, comments.Comment35, comments.Comment4, comments.Comment5, comments.Comment6, comments.Comment7, comments.Comment8, comments.Comment9, comments.Comment10, comments.Form_id)
	if err != nil {
		fmt.Println(err)
		return 0, err

	}

	return id, err

}
