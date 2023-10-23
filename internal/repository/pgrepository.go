package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"tgtrello/internal/model"
)

type PGRepository struct {
	db *sql.DB
}

func NewPgRepository(db *sql.DB) *PGRepository {
	return &PGRepository{db: db}
}

func (r *PGRepository) CheckUserRegister(id int64) (string, error) {
	var login string
	err := r.db.QueryRow(`SELECT login FROM users.user WHERE id = $1`, id).Scan(&login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return login, nil
		}
		return "", fmt.Errorf("execute: %w", err)
	}

	return login, nil
}

func (r *PGRepository) CheckLogin(login string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users."user" WHERE login = $1)`, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("execute query: %w", err)
	}

	return exists, nil
}

func (r *PGRepository) AddNewUser(user *model.User) error {
	_, err := r.db.Exec(`INSERT INTO users.user(id, login, password, tg_name, tg_username, register_time) VALUES ($1,$2,$3,$4,$5,now())`,
		user.ID,
		user.Login,
		user.Password,
		user.TgName,
		user.TgUsername)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	return nil
}
