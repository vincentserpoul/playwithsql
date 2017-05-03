package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(ctx context.Context, exec sqlx.ExtContext) (id int64, err error) {
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
func (link *Link) SelectEntity(
	ctx context.Context,
	q *sqlx.DB,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sqlx.Rows, error) {
	return islatest.SelectEntity(
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
