package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
	err := r.db.QueryRow(`SELECT login FROM trello.user WHERE id = $1`, id).Scan(&login)
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
	err := r.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM trello."user" WHERE login = $1)`, login).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("execute query: %w", err)
	}

	return exists, nil
}

func (r *PGRepository) AddNewUser(user *model.User) error {
	_, err := r.db.Exec(`INSERT INTO trello.user(id, login, password, tg_name, tg_username, register_time) VALUES ($1,$2,$3,$4,$5,now())`,
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

func (r *PGRepository) AddUserToTaskBar(userID int64) error {
	_, err := r.db.Exec(`INSERT INTO trello.task (user_id) VALUES ($1)`, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PGRepository) UpdateTaskComplexity(userID int64, complexity int) error {
	_, err := r.db.Exec(`UPDATE trello.task SET complexity = $1 WHERE user_id = $2`, complexity, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGRepository) UpdateTaskDeadline(userID int64, deadline time.Time) error {
	_, err := r.db.Exec(`UPDATE trello.task SET deadline = $1 WHERE user_id = $2`, deadline, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *PGRepository) CheckTeam(id int64) (int, error) {
	var teamId int
	err := r.db.QueryRow(`SELECT team_id FROM trello.user_team WHERE user_id = $1;`, id).Scan(&teamId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return teamId, nil
		}
		return 0, err
	}

	return teamId, nil
}

func (r *PGRepository) YourTeam(teamID int) (*model.Team, error) {
	team := &model.Team{}
	err := r.db.QueryRow(`SELECT name FROM trello.team WHERE id = $1`, teamID).Scan(&team.Name)
	if err != nil {
		return nil, err
	}

	row, err := r.db.Query(`SELECT user_id, login FROM trello.user_team LEFT JOIN trello."user" u on u.id = user_team.user_id WHERE team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}

	users, err := Rows(row)
	if err != nil {
		return nil, err
	}

	team.Users = users

	return team, nil
}

func (r *PGRepository) DeleteUserFromTeam(userID int64) error {
	_, err := r.db.Exec(
		`DELETE FROM trello.user_team WHERE user_id = $1`,
		userID,
	)
	if err != nil {
		return err
	}

	return nil
}

func Rows(rows *sql.Rows) ([]*model.User, error) {
	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.Login)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *PGRepository) CreateTeam(id int64, teamName string) error {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		_ = tx.Rollback()
	}(tx)

	var teamId int
	err = tx.QueryRowContext(ctx, `INSERT INTO trello.team (name) VALUES ($1) RETURNING id`, teamName).Scan(&teamId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO trello.user_team (team_id, user_id) VALUES ($1, $2)`, teamId, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PGRepository) AddUserToTeam(teamID int, userID int64) (string, error) {
	_, err := r.db.Exec(`INSERT INTO trello.user_team(team_id, user_id) VALUES ($1, $2)`, teamID, userID)
	if err != nil {
		return "", fmt.Errorf("execute: %w", err)
	}

	var teamName string
	err = r.db.QueryRow(`SELECT name FROM trello.team WHERE id = $1`, teamID).Scan(&teamName)
	if err != nil {
		return "", err
	}

	return teamName, nil
}
