package session_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ahbcc/cmd/api/user/session"
	"ahbcc/internal/database"
)

func TestDelete_success(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPostgresConnection.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	deleteSession := session.MakeDelete(mockPostgresConnection)

	got := deleteSession(context.Background(), "token")

	assert.Nil(t, got)
	mockPostgresConnection.AssertExpectations(t)
}

func TestDelete_failsWhenDeleteOperationThrowsError(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPostgresConnection.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("failed to delete user session"))

	deleteSession := session.MakeDelete(mockPostgresConnection)

	want := session.FailedToDeleteUserSession
	got := deleteSession(context.Background(), "token")

	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
}

func TestDeleteExpiredSessions_success(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPostgresConnection.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, nil)

	deleteUserExpiredSessions := session.MakeDeleteUserExpiredSessions(mockPostgresConnection)

	got := deleteUserExpiredSessions(context.Background(), 1234)

	assert.Nil(t, got)
	mockPostgresConnection.AssertExpectations(t)
}

func TestDeleteExpiredSessions_failsWhenDeleteOperationThrowsError(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPostgresConnection.On("Exec", mock.Anything, mock.Anything, mock.Anything).Return(pgconn.CommandTag{}, errors.New("failed to delete user expired sessions"))

	deleteUserExpiredSessions := session.MakeDeleteUserExpiredSessions(mockPostgresConnection)

	want := session.FailedToDeleteExpiredUserSessions
	got := deleteUserExpiredSessions(context.Background(), 1234)

	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
}
