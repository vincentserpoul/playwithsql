package mssql

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(ctx context.Context, db *sql.DB) (id int64, err error) {

	res, err := db.ExecContext(ctx, `INSERT INTO entityone DEFAULT VALUES`)
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

// SelectEntity retrieves a slice of entityones
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

	query := ` SELECT `

	if limit > 0 {
		limitStr := ` TOP ` + strconv.Itoa(limit)
		query += limitStr
	}

	query += `
                e.entityone_id, e.time_created,
                es.entityone_id as status_entityone_id, es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
                AND es.is_latest = 1
            WHERE 0 = 0
        `

	namedArgs, queryFilter := islatest.GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

	query += queryFilter

	return db.QueryContext(ctx, query, namedArgs)
}
