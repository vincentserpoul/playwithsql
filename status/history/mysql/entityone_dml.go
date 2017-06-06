package mysql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/history"
)

// SelectEntityone returns sqlx.Rows
func (link *Link) SelectEntityone(
	ctx context.Context,
	q *sqlx.DB,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sqlx.Rows, error) {
	return history.SelectEntityone(
		ctx,
		q,
		entityIDs,
		isStatusIDs,
		notStatusIDs,
		neverStatusIDs,
		hasStatusIDs,
		limit,
	)
}

// Create will insert a new entity in the DB
func (link *Link) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {
	return history.Create(ctx, tx, actionID, statusID)
}

// SaveStatus will save the status in database for the selected entity
func (link *Link) SaveStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) error {
	return history.SaveStatus(ctx, tx, entityoneID, actionID, statusID)
}
