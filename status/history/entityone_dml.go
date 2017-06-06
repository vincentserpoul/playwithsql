package history

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

// Create will create an entityone and return an id
func Create(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {
	entityoneID, errInso := insertOne(ctx, tx, actionID, statusID)
	if errInso != nil {
		return 0, fmt.Errorf("Create error: %v", errInso)
	}

	errInsn := InsertNewStatus(ctx, tx, entityoneID, actionID, statusID)
	if errInsn != nil {
		return 0, fmt.Errorf("Create error: %v", errInsn)
	}

	return entityoneID, nil
}

func insertOne(
	ctx context.Context,
	tx *sqlx.Tx,
	actionID int,
	statusID int,
) (int64, error) {

	res, err := tx.NamedExecContext(
		ctx,
		`
			INSERT INTO entityone(action_id, status_id)
			VALUES (:actionID, :statusID)
		`,
		map[string]interface{}{
			"actionID": actionID,
			"statusID": statusID,
		},
	)
	if err != nil {
		return 0, fmt.Errorf("insertOne(%d, %d): %v", actionID, statusID, err)
	}

	id, errL := res.LastInsertId()
	if errL != nil {
		return 0, fmt.Errorf("insertOne(%d, %d): %v", actionID, statusID, errL)
	}

	return id, nil
}

// InsertNewStatus insert a new status in entityone_history
func InsertNewStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) error {

	_, err := tx.NamedExecContext(
		ctx,
		`
			INSERT INTO entityone_history(entityone_id, action_id, status_id)
			VALUES (:entityoneID, :actionID, :statusID)
		`,
		map[string]interface{}{
			"entityoneID": entityoneID,
			"actionID":    actionID,
			"statusID":    statusID,
		},
	)
	if err != nil {
		return fmt.Errorf("insertNewStatus(%d, %d, %d): %v", entityoneID, actionID, statusID, err)
	}

	return nil
}

// updateLatestStatus will update the status in the DB
func updateLatestStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) error {

	_, err := tx.NamedExecContext(
		ctx,
		`
			UPDATE entityone
			SET action_id = :actionID, status_id = :statusID, time_updated=:timeUpdated
			WHERE entityone_id = :entityoneID
		`,
		map[string]interface{}{
			"entityoneID": entityoneID,
			"actionID":    actionID,
			"statusID":    statusID,
			"timeUpdated": time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("entityone updateLatestStatus(%d, %d, %d): %v", entityoneID, actionID, statusID, err)
	}

	return nil
}

// SaveStatus will update the status of an entityone
func SaveStatus(
	ctx context.Context,
	tx *sqlx.Tx,
	entityoneID int64,
	actionID int,
	statusID int,
) error {
	errU := updateLatestStatus(ctx, tx, entityoneID, actionID, statusID)
	if errU != nil {
		return fmt.Errorf("entityone SaveStatus(%d, %d, %d): %v", entityoneID, actionID, statusID, errU)
	}

	errI := InsertNewStatus(ctx, tx, entityoneID, actionID, statusID)
	if errI != nil {
		return fmt.Errorf("entityone SaveStatus(%d, %d, %d): %v", entityoneID, actionID, statusID, errI)
	}

	return nil
}

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
               	e.action_id, e.status_id, e.time_updated as status_time_created
            FROM entityone e
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
	entityoneIDs []int64,
	isStatusIDs []int,
) (map[string]interface{}, string) {

	var queryFilter string
	namedParams := make(map[string]interface{})

	if len(entityoneIDs) > 0 {
		queryFilter += ` AND e.entityone_id IN (:entityoneID) `
		namedParams["entityoneID"] = entityoneIDs
	}

	if len(isStatusIDs) > 0 {
		queryFilter += ` AND e.status_id IN (:statusID) `
		namedParams["statusID"] = isStatusIDs
	}

	return namedParams, queryFilter
}
