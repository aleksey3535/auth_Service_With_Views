package repository

import (
	"auth/internal/pkg"
	"fmt"

	"github.com/jmoiron/sqlx"
)

const UsersTable = "users"

type Repository struct {
	db *sqlx.DB
}

func New(db * sqlx.DB) *Repository {
	return &Repository{db : db}
}


func(repo *Repository) CreateUser(login,password string) error {
	query := fmt.Sprintf("INSERT INTO %s (login, password_hash) values ($1, $2)", UsersTable)
	if _, err := repo.db.Exec(query, login, password); err != nil {
		return err
	}
	return nil
}

func(repo *Repository) GetUserIfExist(login string) bool {
	var logins []string
	query := fmt.Sprintf("SELECT login FROM %s WHERE login = $1", UsersTable)
	repo.db.Select(&logins, query, login)
	return len(logins) > 0
	
}

func(repo *Repository) GetUserPass(login string) (string, error) {
	var password string
	query := fmt.Sprintf("SELECT password_hash FROM %s WHERE login = $1", UsersTable)
	row := repo.db.QueryRow(query, login)
	if err := row.Scan(&password); err != nil {
		return "", err
	}
	return password, nil 
}

func(repo *Repository) GetAllUsers() ([]pkg.User, error) {
	var users []pkg.User
	query := fmt.Sprintf("SELECT * FROM %s ", UsersTable)
	err := repo.db.Select(&users, query)
	return users, err
}

func(repo *Repository) DeleteUser(login string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE login=$1", UsersTable)
	if _, err := repo.db.Exec(query, login); err != nil {
		return err
	}
	return nil
}