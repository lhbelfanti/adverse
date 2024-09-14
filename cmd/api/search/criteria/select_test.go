package criteria_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ahbcc/cmd/api/search/criteria"
	"ahbcc/internal/database"
)

func TestSelectByID_success(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPgxRow := new(database.MockPgxRow)
	mockCriteria := criteria.MockCriteriaDAO()
	database.MockScan[criteria.DAO](mockPgxRow, mockCriteria, t)
	mockPostgresConnection.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockPgxRow)

	selectCriteriaByID := criteria.MakeSelectByID(mockPostgresConnection)

	want := mockCriteria
	got, err := selectCriteriaByID(context.Background(), 1)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
	mockPgxRow.AssertExpectations(t)

}

func TestSelectByID_failsWhenSelectOperationFails(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPgxRow := new(database.MockPgxRow)
	mockPgxRow.On("Scan", mock.Anything).Return(pgx.ErrNoRows)
	mockPostgresConnection.On("QueryRow", mock.Anything, mock.Anything, mock.Anything).Return(mockPgxRow)

	selectCriteriaByID := criteria.MakeSelectByID(mockPostgresConnection)

	want := criteria.FailedToRetrieveCriteriaData
	_, got := selectCriteriaByID(context.Background(), 1)

	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
	mockPgxRow.AssertExpectations(t)
}

func TestSelectExecutionsByState_success(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPgxRows := new(database.MockPgxRows)
	mockPostgresConnection.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockPgxRows, nil)
	mockExecutionDAOSlice := criteria.MockExecutionDAOSlice()
	mockCollectRows := database.MockCollectRows[criteria.ExecutionDAO](mockExecutionDAOSlice, nil)

	selectExecutionsByState := criteria.MakeSelectExecutionsByState(mockPostgresConnection, mockCollectRows)

	want := mockExecutionDAOSlice
	got, err := selectExecutionsByState(context.Background(), "DONE")

	assert.Nil(t, err)
	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
	mockPgxRows.AssertExpectations(t)
}

func TestSelectExecutionsByState_failsWhenSelectOperationThrowsError(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPgxRows := new(database.MockPgxRows)
	mockPostgresConnection.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockPgxRows, errors.New("failed to select executions by state"))
	mockExecutionDAOSlice := criteria.MockExecutionDAOSlice()
	mockCollectRows := database.MockCollectRows[criteria.ExecutionDAO](mockExecutionDAOSlice, nil)

	selectExecutionsByState := criteria.MakeSelectExecutionsByState(mockPostgresConnection, mockCollectRows)

	want := criteria.FailedToExecuteSelectSearchCriteriaExecutionByState
	_, got := selectExecutionsByState(context.Background(), "DONE")

	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
	mockPgxRows.AssertExpectations(t)
}

func TestSelectExecutionsByState_failsWhenCollectRowsThrowsError(t *testing.T) {
	mockPostgresConnection := new(database.MockPostgresConnection)
	mockPgxRows := new(database.MockPgxRows)
	mockPostgresConnection.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(mockPgxRows, nil)
	mockCollectRows := database.MockCollectRows[criteria.ExecutionDAO](nil, errors.New("failed to collect rows"))

	selectExecutionsByState := criteria.MakeSelectExecutionsByState(mockPostgresConnection, mockCollectRows)

	want := criteria.FailedToExecuteSelectCollectRowsInSelectExecutionByState
	_, got := selectExecutionsByState(context.Background(), "DONE")

	assert.Equal(t, want, got)
	mockPostgresConnection.AssertExpectations(t)
	mockPgxRows.AssertExpectations(t)
}
