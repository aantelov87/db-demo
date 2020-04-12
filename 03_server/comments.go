package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Comments struct {
	db *sql.DB
}

func NewComments(params string) (*Comments, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}
	repo := &Comments{db}
	return repo, repo.init()
}

func (repo *Comments) Close() error {
	return repo.db.Close()
}

func (repo *Comments) init() error {
	_, err := repo.db.Exec(`
		CREATE TABLE IF NOT EXISTS Comments (
			"User"    TEXT,
			"Comment" TEXT
		)
	`)
	return err
}

func (repo *Comments) Add(user, comment string) error {
	_, err := repo.db.Exec(`INSERT INTO Comments ("User", "Comment") VALUES ($1, $2)`, user, comment)
	return err
}

func (repo *Comments) List() ([]Comment, error) {
	rows, err := repo.db.Query(`SELECT "User", "Comment" FROM Comments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.User, &comment.Text)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
