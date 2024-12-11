package user_test

import (
	"MicroF1-test-case/user"
	"context"
	"github.com/google/uuid"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := user.NewService(db)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO users (id, email) VALUES ($1, $2)")).
		WithArgs(sqlmock.AnyArg(), "newuser@example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.CreateUser(context.Background(), &user.User{Email: "newuser@example.com"})
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := user.NewService(db)

	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow("00000000-0000-0000-0000-000000000000", "test@example.com")
	mock.ExpectQuery("SELECT id, email FROM users").WillReturnRows(rows)

	users, err := s.GetUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, "test@example.com", users[0].Email)
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := user.NewService(db)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET email = $1 WHERE id = $2")).
		WithArgs("updated@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id := uuid.New()
	err = s.UpdateUser(context.Background(), id, "updated@example.com")
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	s := user.NewService(db)

	id := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM users WHERE id = $1")).
		WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = s.DeleteUser(context.Background(), id)
	assert.NoError(t, err)
}
