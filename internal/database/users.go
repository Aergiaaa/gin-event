package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

func (um *UserModel) Insert(u *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`

	return um.DB.QueryRowContext(ctx, query, u.Email, u.Password, u.Name).Scan(&u.Id)
}

func (um *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT * FROM users`
	rows, err := um.DB.QueryContext(ctx, query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var u User

		err := rows.Scan(&u.Id, &u.Email, &u.Name, &u.Password)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (um *UserModel) Get(id int) (*User, error) {
	query := `SELECT id, email, name, password FROM users WHERE id = $1`
	return um.getUser(query, id)
}

func (um *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, name, password FROM users WHERE email = $1`
	return um.getUser(query, email)
}

func (um *UserModel) getUser(query string, args ...any) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User
	err := um.DB.QueryRowContext(ctx, query, args...).
		Scan(&u.Id, &u.Email, &u.Name, &u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &u, nil
}
