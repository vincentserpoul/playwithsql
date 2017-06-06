package postgres

import (
	"context"
	"fmt"

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
	entityoneID, errInso := insertOne(ctx, tx, actionID, statusID)
	if errInso != nil {
		return 0, fmt.Errorf("Create error: %v", errInso)
	}

	errInsn := history.InsertNewStatus(ctx, tx, entityoneID, actionID, statusID)
	if errInsn != nil {
		return 0, fmt.Errorf("Create error: %v", errInsn)
	}

	return entityoneID, nil
}

// InsertOne will insert a Entityone into db
func insertOne(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {
	stmt, err := tx.PrepareNamed(`
		INSERT INTO entityone(action_id, status_id)
		VALUES(:actionID, :statusID)
		RETURNING entityone_id
	`)
	if err != nil {
		return 0, fmt.Errorf("entityone insertOne(%d, %d): %v", actionID, statusID, err)
	}

	var id int64
	err = stmt.GetContext(
		ctx,
		&id,
		map[string]interface{}{
			"actionID": actionID,
			"statusID": statusID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("entityone insertOne(%d, %d): %v", actionID, statusID, err)
	}

	return id, nil
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
