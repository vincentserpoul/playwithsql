package lateststatus

import (
	"context"
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// SelectEntityone retrieves a slice of entityones
func SelectEntityone(
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
                e.entityone_id, e.time_created,
               	es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
			INNER JOIN entityone_lateststatus el ON el.entityone_id = e.entityone_id
            INNER JOIN entityone_status es ON es.entityone_status_id = el.entityone_status_id
			WHERE 0=0
        `

	namedParams, queryFilter := GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

	query += queryFilter

	if limit > 0 {
		limitStr := ` LIMIT ` + strconv.Itoa(limit)
		query += limitStr
	}

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

// GetFilterSelectEntityOneNamedQuery returns query filter and params for the query
func GetFilterSelectEntityOneNamedQuery(
	entityIDs []int64,
	isStatusIDs []int,
) (map[string]interface{}, string) {

	var queryFilter string
	namedParams := make(map[string]interface{})

	if len(entityIDs) > 0 {
		queryFilter += ` AND e.entityone_id IN (:entityID) `
		namedParams["entityID"] = entityIDs
	}

	if len(isStatusIDs) > 0 {
		queryFilter += `  AND es.status_id IN (:statusID) `
		namedParams["statusID"] = isStatusIDs
	}

	return namedParams, queryFilter
}

// InsertNewStatus will insert a new status into db
func InsertNewStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) (int64, error) {

	res, err := tx.NamedExecContext(
		ctx,
		`
			INSERT INTO entityone_status(entityone_id, action_id, status_id)
			VALUES (:entityoneID, :actionID, :statusID)
		`,
		map[string]interface{}{
			"entityoneID": entityoneID,
			"actionID":    actionID,
			"statusID":    statusID,
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

// UpdateLatestStatus will insert a new status into db
func UpdateLatestStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	entityoneStatusID int64,
) error {

	_, err := tx.NamedExecContext(
		ctx,
		`
			UPDATE entityone_lateststatus
			SET entityone_status_id = :entityoneStatusID
			WHERE entityone_id = :entityoneID
		`,
		map[string]interface{}{
			"entityoneID":       entityoneID,
			"entityoneStatusID": entityoneStatusID,
		},
	)
	if err != nil {
		return fmt.Errorf("entityone updateLatestStatus(%d, %d): %v", entityoneID, entityoneStatusID, err)
	}

	return nil
}

// InsertLatestStatus will insert a new status into db
func InsertLatestStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	entityoneStatusID int64,
) error {

	_, err := tx.NamedExecContext(
		ctx,
		`
			INSERT INTO entityone_lateststatus(entityone_id, entityone_status_id)
			VALUES (:entityoneID, :entityoneStatusID)
		`,
		map[string]interface{}{
			"entityoneID":       entityoneID,
			"entityoneStatusID": entityoneStatusID,
		},
	)
	if err != nil {
		return fmt.Errorf("entityone insertLatestStatus(%d, %d): %v", entityoneID, entityoneStatusID, err)
	}

	return nil
}
