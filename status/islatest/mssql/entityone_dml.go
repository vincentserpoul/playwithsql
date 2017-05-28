package mssql

import (
	"context"
	"fmt"
	"strconv"

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

// SelectEntity retrieves a slice of entityones
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

	query := ` SELECT `

	if limit > 0 {
		limitStr := ` TOP ` + strconv.Itoa(limit)
		query += limitStr
	}

	query += `
                e.entityone_id, e.time_created,
                es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
                AND es.is_latest = 1
            WHERE 0 = 0
        `

	namedParams, queryFilter := islatest.GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

	query += queryFilter

	query, injectedNamedParams, err := sqlx.Named(query, namedParams)
	if err != nil {
		return nil, fmt.Errorf("SelectEntity error: %v", err)
	}

	query, injectedNamedParams, err = sqlx.In(query, injectedNamedParams...)
	if err != nil {
		return nil, fmt.Errorf("SelectEntity error: %v", err)
	}

	query = q.Rebind(query)
	return q.QueryxContext(ctx, query, injectedNamedParams...)

}
