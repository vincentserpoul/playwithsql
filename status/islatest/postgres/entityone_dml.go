package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// Create will insert a new entity in the DB along with the status
func (link *Link) Create(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {
	id, err := link.insertOne(ctx, tx)
	if err != nil {
		return id, fmt.Errorf("entityone Create(): %v", err)
	}

	err = link.SaveStatus(ctx, tx, id, actionID, statusID)
	if err != nil {
		return id, fmt.Errorf("entityone Create(): %v", err)
	}

	return id, nil
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

// SaveStatus will save the status in database for the selected entity
func (link *Link) SaveStatus(
	ctx context.Context,
	exec *sqlx.Tx,
	entityID int64,
	actionID int,
	statusID int,
) error {
	return islatest.SaveStatus(ctx, exec, entityID, actionID, statusID)
}

// SelectEntity returns sqlx.Rows
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
	return islatest.SelectEntityone(
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
