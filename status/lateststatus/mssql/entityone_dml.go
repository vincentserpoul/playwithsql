package mssql

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/islatest"
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

	res, err := tx.NamedExecContext(
		ctx,
		`INSERT INTO entityone (entityone_status_id) VALUES (:entityoneStatusID)`,
		map[string]interface{}{
			"entityoneStatusID": entityoneStatusID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("entityone Insert(): %v", err)
	}

	id, errL := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("entityone Insert(): %v", errL)
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

// SelectEntityone retrieves a slice of entityones
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

	query := ` SELECT `

	if limit > 0 {
		limitStr := ` TOP ` + strconv.Itoa(limit)
		query += limitStr
	}

	query += `
                e.entityone_id, e.time_created,
                es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_status_id = e.entityone_status_id
            WHERE 0 = 0
        `

	namedParams, queryFilter := islatest.GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

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

// insertNewStatus will insert a new status into db
func (link *Link) insertNewStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {

	res, err := tx.NamedExecContext(
		ctx,
		`INSERT INTO entityone_status(action_id, status_id) VALUES (:actionID, :statusID)`,
		map[string]interface{}{
			"actionID": actionID,
			"statusID": statusID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("entityone insertNewStatus(%d, %d): %v", actionID, statusID, err)
	}

	id, errL := res.LastInsertId()
	if errL != nil {
		return 0, fmt.Errorf("entityone insertNewStatus(%d, %d): %v", actionID, statusID, errL)
	}

	return id, nil
}
