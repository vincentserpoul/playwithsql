package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/lateststatus"
)

// Create will insert a new entity in the DB
func (link *Link) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {
	entityoneStatusID, errS := link.insertNewStatus(ctx, tx, actionID, statusID)
	if errS != nil {
		return 0, fmt.Errorf("entityone Create(): %v", errS)
	}

	entityID, errE := link.insertOne(ctx, tx, entityoneStatusID)
	if errE != nil {
		return 0, fmt.Errorf("entityone Create(): %v", errE)
	}

	return entityID, nil
}

// insertOne will insert a Entityone into db
func (link *Link) insertOne(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneStatusID int64,
) (int64, error) {

	stmt, err := tx.PrepareNamed(`
		INSERT INTO entityone (entityone_id, time_created, entityone_status_id)
		VALUES(DEFAULT, DEFAULT, :entityoneStatusID)
		RETURNING entityone_id
	`)
	if err != nil {
		return 0, fmt.Errorf("entityone insertone(): %v", err)
	}

	var id int64
	err = stmt.GetContext(
		ctx,
		&id,
		map[string]interface{}{"entityoneStatusID": entityoneStatusID},
	)
	if err != nil {
		return 0, fmt.Errorf("entityone Insert(): %v", err)
	}

	return id, nil
}

// SaveStatus will save the status in database for the selected entity
func (link *Link) SaveStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityID int64,
	actionID int,
	statusID int,
) error {
	entityStatusID, err := link.insertNewStatus(ctx, tx, actionID, statusID)
	if err != nil {
		return fmt.Errorf("entityone SaveStatus(%d, %d, %d): %v", entityID, actionID, statusID, err)
	}

	return lateststatus.UpdateLatestStatus(ctx, tx, "entityone", entityID, entityStatusID)
}

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

// insertNewStatus will insert a new status into db
func (link *Link) insertNewStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {
	stmt, err := tx.PrepareNamed(`
		INSERT INTO entityone_status(entityone_status_id, time_created, action_id, status_id)
		VALUES(DEFAULT, DEFAULT, :actionID, :statusID)
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
			"actionID": actionID,
			"statusID": statusID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("entityone insertNewStatus(%d, %d): %v", actionID, statusID, err)
	}

	return id, nil
}
