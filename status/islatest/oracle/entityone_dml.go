package oracle

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(ctx context.Context, exec sqlx.ExtContext) (id int64, err error) {

	res, err := exec.ExecContext(ctx,
		`
		INSERT INTO entityone(entityone_id) VALUES (default)
		RETURNING entityone_id /*LastInsertId*/ INTO :id
	`, nil)
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

	query := `
            SELECT
                e.entityone_id as "entityone_id", e.time_created as "time_created",
                es.entityone_id as "status_entityone_id", es.action_id as "action_id",
				es.status_id as "status_id", es.time_created as "status_time_created"
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
                AND es.is_latest = 1
            WHERE 0 = 0
        `

	namedParams, queryFilter := islatest.GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

	query += queryFilter

	if limit > 0 {
		limitStr := ` AND ROWNUM <= ` + strconv.Itoa(limit)
		query += limitStr
	}

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
