package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Muaz717/sso/internal/config"
	"github.com/Muaz717/sso/internal/domain/models"
	"github.com/Muaz717/sso/internal/storage"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg config.DB) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username,
		cfg.DBPassword,
		cfg.Host,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect db: %w", err)
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "postgres.SaveUser"

	query := `INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id`

	row := s.db.QueryRow(ctx, query, email, passHash)

	var userId int64

	err := row.Scan(&userId)
	if err != nil {
		if pgErr, ok := err.(*pgx.PgError); ok {
			return 0, fmt.Errorf("%s: SQL Error: %s, Detail: %s, Where: %s", op, pgErr.Message, pgErr.Detail, pgErr.Where)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.User"

	query := `SELECT id, email, pass_hash FROM users WHERE email=$1`

	row := s.db.QueryRow(ctx, query, email)

	var user models.User

	err := row.Scan(&user.Id, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `SELECT is_admin FROM users WHERE id = $1`

	row := s.db.QueryRow(ctx, query, userId)

	var isAdmin bool
	err := row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, err
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.sqlite.App"

	query := `SELECT id, name, secret FROM apps WHERE id = $1`

	row := s.db.QueryRow(ctx, query, id)

	var app models.App
	err := row.Scan(&app.Id, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
