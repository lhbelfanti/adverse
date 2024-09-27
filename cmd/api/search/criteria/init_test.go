package criteria_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"ahbcc/cmd/api/search/criteria"
)

func TestInit_success(t *testing.T) {
	mockExecutionsDAO := criteria.MockExecutionsDAO()
	mockSelectExecutionsByStatuses := criteria.MockSelectExecutionsByStatuses(mockExecutionsDAO, nil)
	mockResume := criteria.MockResume(nil)

	init := criteria.MakeInit(mockSelectExecutionsByStatuses, mockResume)

	got := init(context.Background())

	assert.Nil(t, got)
}

func TestInit_failsWhenSelectExecutionsByStatusesThrowsError(t *testing.T) {
	mockSelectExecutionsByStatuses := criteria.MockSelectExecutionsByStatuses(nil, errors.New("failed while executing select executions by statuses"))
	mockResume := criteria.MockResume(nil)

	init := criteria.MakeInit(mockSelectExecutionsByStatuses, mockResume)

	want := criteria.FailedToExecuteSelectExecutionsByStatuses
	got := init(context.Background())

	assert.Equal(t, want, got)
}

func TestInit_failsWhenEnqueueThrowsError(t *testing.T) {
	mockExecutionsDAO := criteria.MockExecutionsDAO()
	mockSelectExecutionsByStatuses := criteria.MockSelectExecutionsByStatuses(mockExecutionsDAO, nil)
	mockResume := criteria.MockResume(errors.New("failed while executing resume"))

	init := criteria.MakeInit(mockSelectExecutionsByStatuses, mockResume)

	want := criteria.FailedToExecuteEnqueueCriteria
	got := init(context.Background())

	assert.Equal(t, want, got)
}
