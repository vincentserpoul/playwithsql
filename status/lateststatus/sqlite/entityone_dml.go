package sqlite

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

// insertOne will insert a Entityone into db
func (link *Link) insertOne(ctx context.Context, exec *sqlx.Tx) (id int64, err error) {

	res, err := exec.ExecContext(ctx, `INSERT INTO entityone DEFAULT VALUES`)
	if err != nil {
		return id, fmt.Errorf("entityone Insert(): %v", err)
	}

	id, err = res.LastInsertId()
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
	return lateststatus.InsertNewStatus(ctx, tx, entityoneID, actionID, statusID)
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
