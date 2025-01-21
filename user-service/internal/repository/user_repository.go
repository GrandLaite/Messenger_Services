package repository

import (
	"database/sql"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(d *sql.DB) *UserRepository {
	return &UserRepository{db: d}
}

func (r *UserRepository) Create(u User) (int, error) {
	q := `INSERT INTO users (username, password_hash, role) 
	      VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRow(q, u.Username, u.PasswordHash, u.Role).Scan(&id)
	return id, err
}

func (r *UserRepository) GetByUsername(name string) (User, error) {
	var us User
	q := `SELECT id, username, password_hash, role FROM users WHERE username = $1`
	err := r.db.QueryRow(q, name).Scan(&us.ID, &us.Username, &us.PasswordHash, &us.Role)
	return us, err
}
