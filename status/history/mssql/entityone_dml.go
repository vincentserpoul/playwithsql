package mssql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/history"
)

// SelectEntityone retrieves a slice of entityones
func (link *Link) SelectEntityone(
	ctx context.Context,
	q *sqlx.DB,
	entityoneIDs []int64,
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
               	e.action_id, e.status_id, e.time_updated as status_time_created
            FROM entityone e
			WHERE 0=0
        `

	namedParams, queryFilter := history.GetFilterSelectEntityOneNamedQuery(entityoneIDs, isStatusIDs)

	query += queryFilter

	query, injectedNamedParams, err := sqlx.Named(query, namedParams)
	if err != nil {
		return nil, fmt.Errorf("SelectEntityone error: %v", err)
	}

	query, injectedNamedParams, err = sqlx.In(query, injectedNamedParams...)
	if err != nil {
		return nil, fmt.Errorf("SelectEntityone error: %v", err)
	}

	query = q.Rebind(query)
	return q.QueryxContext(ctx, query, injectedNamedParams...)

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
