package sqlite

import (
	"context"
	"fmt"

	"database/sql"

	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(
	ctx context.Context,
	db *sql.DB,
) (id int64, err error) {

	res, err := db.ExecContext(
		ctx,
		`INSERT INTO entityone DEFAULT VALUES`,
	)
	if err != nil {
		return id, fmt.Errorf("entityone Insert(): %v", err)
	}

	id, err = res.LastInsertId()
	if err != nil {
		return id, fmt.Errorf("entityone Insert(): %v", err)
	}

	return id, nil
}

// SaveStatus will save the status in database for the selected entity
func (link *Link) SaveStatus(
	ctx context.Context,
	tx *sql.Tx,
	entityID int64,
	actionID int,
	statusID int,
) error {
	return islatest.SaveStatus(ctx, tx, entityID, actionID, statusID)
}

// SelectEntity returns sql.Rows
func (link *Link) SelectEntity(
	ctx context.Context,
	db *sql.DB,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sql.Rows, error) {
	return islatest.SelectEntity(
		ctx,
		db,
		entityIDs,
		isStatusIDs,
		notStatusIDs,
		neverStatusIDs,
		hasStatusIDs,
		limit,
	)
}
