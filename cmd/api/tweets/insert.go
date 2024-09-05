package tweets

import (
	"context"
	"fmt"
	"strings"

	"ahbcc/cmd/api/tweets/quotes"
	"ahbcc/internal/database"
	"ahbcc/internal/log"
)

// Insert inserts a new TweetDTO into 'tweets' table
type Insert func(ctx context.Context, tweet []TweetDTO) error

// MakeInsert creates a new Insert
func MakeInsert(db database.Connection, insertQuote quotes.InsertSingle, deleteOrphanQuotes quotes.DeleteOrphans) Insert {
	const (
		query string = `
			INSERT INTO tweets(hash, is_a_reply, text_content, images, quote_id, search_criteria_id) 
			VALUES %s
		    ON CONFLICT (hash, search_criteria_id) DO NOTHING;
		`
		parameters = 6
	)

	return func(ctx context.Context, tweets []TweetDTO) error {
		placeholders := make([]string, 0, len(tweets)*parameters)
		values := make([]any, 0, len(tweets)*parameters)
		quoteIDs := make([]int, 0, len(tweets))
		for i, tweet := range tweets {
			idx := i * parameters
			placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6))
			values = append(values, tweet.Hash, tweet.IsAReply, tweet.TextContent, tweet.Images)

			quoteID, err := insertQuote(ctx, tweet.Quote)
			if err != nil {
				values = append(values, nil)
			} else {
				values = append(values, quoteID)
				quoteIDs = append(quoteIDs, quoteID)
			}

			values = append(values, tweet.SearchCriteriaID)
		}

		queryToExecute := fmt.Sprintf(query, strings.Join(placeholders, ","))

		_, err := db.Exec(ctx, queryToExecute, values...)
		if err != nil {
			log.Error(ctx, err.Error())
			return FailedToInsertTweets
		}

		if len(quoteIDs) > 0 {
			err = deleteOrphanQuotes(ctx, quoteIDs)
		}

		return nil
	}
}
