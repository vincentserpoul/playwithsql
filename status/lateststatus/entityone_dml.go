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
            INNER JOIN entityone_status es ON es.entityone_status_id = e.entityone_status_id
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

// UpdateLatestStatus will update the entityone row to the latest updated status
func UpdateLatestStatus(
	ctx context.Context,
	exec *sqlx.Tx,
	typeEntity string,
	entityID int64,
	entityStatusID int64,
) error {
	queryIns := fmt.Sprintf(
		"UPDATE %s SET entityone_status_id = :entityStatusID"+
			" WHERE %s_id = :entityID",
		typeEntity,
		typeEntity,
	)

	res, err := exec.NamedExecContext(
		ctx,
		queryIns,
		map[string]interface{}{
			"entityID":       entityID,
			"entityStatusID": entityStatusID,
		},
	)
	if err != nil {
		return fmt.Errorf("lateststatus %s updateLatestStatus(%d, %d): err %v",
			typeEntity,
			entityID,
			entityStatusID,
			err,
		)
	}

	rowsAffected, errAff := res.RowsAffected()
	if errAff != nil {
		return fmt.Errorf("lateststatus %s updateLatestStatus(%d, %d): err %v",
			typeEntity,
			entityID,
			entityStatusID,
			err,
		)
	}
	if rowsAffected != 1 {
		return fmt.Errorf("lateststatus %s updateLatestStatus(%d, %d): err no rows inserted",
			typeEntity,
			entityID,
			entityStatusID,
		)
	}

	return nil
}
