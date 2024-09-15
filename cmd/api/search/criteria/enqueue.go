package criteria

import (
	"context"

	"ahbcc/internal/log"
	"ahbcc/internal/scrapper"
)

// Enqueue retrieves the criteria by ID from the database and enqueues its information
type Enqueue func(ctx context.Context, criteriaID int, forced bool) error

// MakeEnqueue creates a new Enqueue
func MakeEnqueue(selectCriteriaByID SelectByID, selectExecutionsByStatuses SelectExecutionsByStatuses, enqueueCriteria scrapper.EnqueueCriteria) Enqueue {
	return func(ctx context.Context, criteriaID int, forced bool) error {
		criteriaDAO, err := selectCriteriaByID(ctx, criteriaID)
		if err != nil {
			log.Error(ctx, err.Error())
			return FailedToExecuteSelectCriteriaByID
		}

		if !forced {
			executionsDAO, err := selectExecutionsByStatuses(ctx, []string{"PENDING", "IN PROGRESS"})
			if err != nil {
				log.Error(ctx, err.Error())
				return FailedToExecuteSelectExecutionsByStatuses
			}

			for _, execution := range executionsDAO {
				if execution.SearchCriteriaID == criteriaID {
					return AnExecutionOfThisCriteriaIDIsAlreadyEnqueued
				}
			}
		}

		body := scrapper.CriteriaDTO{
			ID:               criteriaDAO.ID,
			Name:             criteriaDAO.Name,
			AllOfTheseWords:  criteriaDAO.AllOfTheseWords,
			ThisExactPhrase:  criteriaDAO.ThisExactPhrase,
			AnyOfTheseWords:  criteriaDAO.AllOfTheseWords,
			NoneOfTheseWords: criteriaDAO.NoneOfTheseWords,
			TheseHashtags:    criteriaDAO.TheseHashtags,
			Language:         criteriaDAO.Language,
			Since:            criteriaDAO.Since,
			Until:            criteriaDAO.Until,
		}

		err = enqueueCriteria(ctx, body)
		if err != nil {
			log.Error(ctx, err.Error())
			return FailedToExecuteEnqueueCriteria
		}

		return nil
	}
}
