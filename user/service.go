package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) GetUsers(ctx context.Context) ([]User, error) {
	query, args, err := squirrel.Select("id", "email").From("users").ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	query := "SELECT id, email FROM users WHERE id = $1"

	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, errors.New("user not found")
		}
		return User{}, err
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user *User) error {
	user.ID = uuid.New()

	query := "INSERT INTO users (id, email) VALUES ($1, $2)"

	_, err := s.db.ExecContext(ctx, query, user.ID, user.Email)
	return err
}

func (s *Service) UpdateUser(ctx context.Context, id uuid.UUID, email string) error {
	query := "UPDATE users SET email = $1 WHERE id = $2"

	_, err := s.db.ExecContext(ctx, query, email, id)
	return err
}

func (s *Service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := s.db.ExecContext(ctx, query, id)
	return err
}
