package islatest

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

// SelectEntity retrieves a slice of entityones
func SelectEntity(
	ctx context.Context,
	db *sql.DB,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sql.Rows, error) {

	query := `
            SELECT
                e.entityone_id, e.time_created,
                es.entityone_id as status_entityone_id, es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
            WHERE es.is_latest = 1
        `

	namedArgs, queryFilter := GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

	query += queryFilter

	if limit > 0 {
		limitStr := ` LIMIT ` + strconv.Itoa(limit)
		query += limitStr
	}

	return db.QueryContext(ctx, query, namedArgs)

}

// GetFilterSelectEntityOneNamedQuery returns query filter and params for the query
func GetFilterSelectEntityOneNamedQuery(
	entityIDs []int64,
	statusIDs []int,
) ([]sql.NamedArg, string) {

	var queryFilter string
	var namedArgs []sql.NamedArg

	if len(entityIDs) > 0 {
		queryFilter += ` AND e.entityone_id IN (@entityIDs) `
		namedArgs = append(namedArgs, sql.Named("entityIDs", entityIDs))
	}

	if len(statusIDs) > 0 {
		queryFilter += `  AND es.status_id IN (@statusIDs) `
		namedArgs = append(namedArgs, sql.Named("statusIDs", statusIDs))
	}

	return namedArgs, queryFilter
}

// SaveStatus will save the status in database for the selected entity
func SaveStatus(
	ctx context.Context,
	tx *sql.Tx,
	entityID int64,
	actionID int,
	statusID int,
) (err error) {
	typeEntity := "entityone"

	err = resetAllPreviousStatuses(ctx, tx, typeEntity, entityID)
	if err != nil {
		return fmt.Errorf("islatest %s SaveStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}

	err = insertNewStatus(ctx, tx, typeEntity, entityID, actionID, statusID)
	if err != nil {
		return fmt.Errorf("islatest %s SaveStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}

	return nil
}

func resetAllPreviousStatuses(
	ctx context.Context,
	tx *sql.Tx,
	typeEntity string,
	entityID int64,
) error {

	queryUpd := fmt.Sprintf(
		"UPDATE %s_status "+
			"SET is_latest = null "+
			"WHERE %s_id=:entityID AND is_latest = 1",
		typeEntity,
		typeEntity,
	)

	_, err := tx.ExecContext(ctx, queryUpd, sql.Named("entityID", entityID))
	if err != nil {
		return fmt.Errorf("islatest resetAllPreviousStatuses(%s, %d): err %v",
			typeEntity,
			entityID,
			err,
		)
	}

	return nil
}

func insertNewStatus(
	ctx context.Context,
	tx *sql.Tx,
	typeEntity string,
	entityID int64,
	actionID int,
	statusID int,
) error {
	queryIns := fmt.Sprintf(
		"INSERT INTO %s_status(%s_id, action_id, status_id) "+
			" VALUES (@entityID, @actionID, @statusID)",
		typeEntity,
		typeEntity,
	)

	res, err := tx.ExecContext(ctx, queryIns,
		[]sql.NamedArg{
			sql.Named("entityID", entityID),
			sql.Named("actionID", actionID),
			sql.Named("statusID", statusID),
		},
	)
	if err != nil {
		return fmt.Errorf("islatest %s insertNewStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}

	rowsAffected, errAff := res.RowsAffected()
	if errAff != nil {
		return fmt.Errorf("islatest %s insertNewStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("islatest %s insertNewStatus(%v, %d, %d): err no rows inserted",
			typeEntity,
			entityID,
			actionID,
			statusID,
		)
	}

	return nil
}
