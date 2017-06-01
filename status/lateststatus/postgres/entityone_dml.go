package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/lateststatus"
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
	return lateststatus.SelectEntityone(
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

	entityoneID, errIO := link.insertOne(ctx, tx)
	if errIO != nil {
		return 0, fmt.Errorf("entityone Create(): %v", errIO)
	}

	entityoneStatusID, errS := link.insertNewStatus(ctx, tx, entityoneID, actionID, statusID)
	if errS != nil {
		return 0, fmt.Errorf("entityone Create(): %v", errS)
	}

	errE := link.insertLatestStatus(ctx, tx, entityoneID, entityoneStatusID)
	if errE != nil {
		return 0, fmt.Errorf("entityone Create(): %v", errE)
	}

	return entityoneID, nil
}

// SaveStatus will save the status in database for the selected entity
func (link *Link) SaveStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) error {
	entityStatusID, err := link.insertNewStatus(ctx, tx, entityoneID, actionID, statusID)
	if err != nil {
		return fmt.Errorf("entityone SaveStatus(%d, %d, %d): %v", entityoneID, actionID, statusID, err)
	}

	return link.updateLatestStatus(ctx, tx, entityoneID, entityStatusID)
}

// InsertOne will insert a Entityone into db
func (link *Link) insertOne(ctx context.Context, exec *sqlx.Tx) (id int64, err error) {
	err = exec.QueryRowxContext(
		ctx,
		`
        INSERT INTO entityone(entityone_id, time_created)
        VALUES(DEFAULT, DEFAULT)
        RETURNING entityone_id
    `).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("entityone Insert(): %v", err)
	}

	return id, nil
}

// insertNewStatus will insert a new status into db
func (link *Link) insertNewStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) (int64, error) {
	stmt, err := tx.PrepareNamed(`
		INSERT INTO entityone_status(entityone_status_id, time_created, entityone_id, action_id, status_id)
		VALUES(DEFAULT, DEFAULT, :entityoneID, :actionID, :statusID)
		RETURNING entityone_status_id
	`)
	if err != nil {
		return 0, fmt.Errorf("entityone insertNewStatus(%d, %d): %v", actionID, statusID, err)
	}

	var id int64
	err = stmt.GetContext(
		ctx,
		&id,
		map[string]interface{}{
			"entityoneID": entityoneID,
			"actionID":    actionID,
			"statusID":    statusID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("entityone insertNewStatus(%d, %d): %v", actionID, statusID, err)
	}

	return id, nil
}

// insertLatestStatus will insert a new status into db
func (link *Link) insertLatestStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	entityoneStatusID int64,
) error {
	return lateststatus.InsertLatestStatus(ctx, tx, entityoneID, entityoneStatusID)
}

// updateLatestStatus will insert a new status into db
func (link *Link) updateLatestStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	entityoneStatusID int64,
) error {
	return lateststatus.UpdateLatestStatus(ctx, tx, entityoneID, entityoneStatusID)
}
