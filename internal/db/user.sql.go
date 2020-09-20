// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package db

import (
	"context"
)

const admins = `-- name: Admins :many
SELECT id, created_at, name, admin, proposer, email, password, bio FROM users
WHERE admin = true
`

func (q *Queries) Admins(ctx context.Context) ([]User, error) {
	rows, err := q.query(ctx, q.adminsStmt, admins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Admin,
			&i.Proposer,
			&i.Email,
			&i.Password,
			&i.Bio,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const countUsers = `-- name: CountUsers :one
SELECT COUNT(*) FROM users 
WHERE lower(name) = lower($1) OR lower(email) = lower($2)
`

type CountUsersParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (q *Queries) CountUsers(ctx context.Context, arg CountUsersParams) (int64, error) {
	row := q.queryRow(ctx, q.countUsersStmt, countUsers, arg.Username, arg.Email)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (
	email, name, password
) VALUES (
	$1, $2, $3
) 
RETURNING id
`

type CreateUserParams struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int64, error) {
	row := q.queryRow(ctx, q.createUserStmt, createUser, arg.Email, arg.Name, arg.Password)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const proposers = `-- name: Proposers :many
SELECT id, created_at, name, admin, proposer, email, password, bio FROM users 
WHERE proposer = true
`

func (q *Queries) Proposers(ctx context.Context) ([]User, error) {
	rows, err := q.query(ctx, q.proposersStmt, proposers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Admin,
			&i.Proposer,
			&i.Email,
			&i.Password,
			&i.Bio,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setAdmin = `-- name: SetAdmin :exec
UPDATE users SET admin = $2
WHERE id = $1
`

type SetAdminParams struct {
	ID    int64 `json:"id"`
	Admin bool  `json:"admin"`
}

func (q *Queries) SetAdmin(ctx context.Context, arg SetAdminParams) error {
	_, err := q.exec(ctx, q.setAdminStmt, setAdmin, arg.ID, arg.Admin)
	return err
}

const setBio = `-- name: SetBio :exec
UPDATE users SET bio = $2
WHERE id = $1
`

type SetBioParams struct {
	ID  int64  `json:"id"`
	Bio string `json:"bio"`
}

func (q *Queries) SetBio(ctx context.Context, arg SetBioParams) error {
	_, err := q.exec(ctx, q.setBioStmt, setBio, arg.ID, arg.Bio)
	return err
}

const setEmail = `-- name: SetEmail :exec
UPDATE users SET email = $2
WHERE id = $1
RETURNING id, created_at, name, admin, proposer, email, password, bio
`

type SetEmailParams struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

func (q *Queries) SetEmail(ctx context.Context, arg SetEmailParams) error {
	_, err := q.exec(ctx, q.setEmailStmt, setEmail, arg.ID, arg.Email)
	return err
}

const setProposer = `-- name: SetProposer :exec
UPDATE users SET proposer = $2
WHERE id = $1
`

type SetProposerParams struct {
	ID       int64 `json:"id"`
	Proposer bool  `json:"proposer"`
}

func (q *Queries) SetProposer(ctx context.Context, arg SetProposerParams) error {
	_, err := q.exec(ctx, q.setProposerStmt, setProposer, arg.ID, arg.Proposer)
	return err
}

const user = `-- name: User :one
SELECT id, created_at, name, admin, proposer, email, password, bio FROM users 
WHERE id = $1
`

func (q *Queries) User(ctx context.Context, id int64) (User, error) {
	row := q.queryRow(ctx, q.userStmt, user, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Admin,
		&i.Proposer,
		&i.Email,
		&i.Password,
		&i.Bio,
	)
	return i, err
}

const userByEmail = `-- name: UserByEmail :one
SELECT id, created_at, name, admin, proposer, email, password, bio FROM users 
WHERE lower(email) = lower($1)
`

func (q *Queries) UserByEmail(ctx context.Context, email string) (User, error) {
	row := q.queryRow(ctx, q.userByEmailStmt, userByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Admin,
		&i.Proposer,
		&i.Email,
		&i.Password,
		&i.Bio,
	)
	return i, err
}

const userByName = `-- name: UserByName :one
SELECT id, created_at, name, admin, proposer, email, password, bio FROM users 
WHERE lower(name) = lower($1)
`

func (q *Queries) UserByName(ctx context.Context, username string) (User, error) {
	row := q.queryRow(ctx, q.userByNameStmt, userByName, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Admin,
		&i.Proposer,
		&i.Email,
		&i.Password,
		&i.Bio,
	)
	return i, err
}

const users = `-- name: Users :many
SELECT id, created_at, name, admin, proposer, email, password, bio FROM users
`

func (q *Queries) Users(ctx context.Context) ([]User, error) {
	rows, err := q.query(ctx, q.usersStmt, users)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.Name,
			&i.Admin,
			&i.Proposer,
			&i.Email,
			&i.Password,
			&i.Bio,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
