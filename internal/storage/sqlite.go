package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"messenger/internal/models_bd"
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
)

var DB *sql.DB 

func NewSQLiteStorage(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}
	return Migrate()
}

func Migrate() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password_hash BLOB NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func SaveUser(ctx context.Context, email string, hash []byte) (int64, error) {
	result, err := DB.ExecContext(ctx, 
		`INSERT INTO users (email, password_hash) VALUES (?, ?)`, 
		email, hash)
	if err != nil {
		return 0, fmt.Errorf("save user: %w", err)
	}
	return result.LastInsertId()
}

func GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash FROM users WHERE email = ?`, 
		email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func GetUserByID(ctx context.Context, id int64) (models.User, error) {
	var user models.User
	err := DB.QueryRowContext(ctx,
		`SELECT id, email, password_hash FROM users WHERE id = ?`, 
		id).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

func Close() error {
	return DB.Close()
}