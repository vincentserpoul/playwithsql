package islatest

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
)

// SelectEntity retrieves a slice of entityones
func SelectEntity(
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
                es.entityone_id as status_entityone_id, es.action_id, es.status_id, es.time_created as status_time_created
            FROM entityone e
            INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
                AND es.is_latest = 1
            WHERE 0 = 0
        `

	namedParams, queryFilter := GetFilterSelectEntityOneNamedQuery(entityIDs, isStatusIDs)

	query += queryFilter

	if limit > 0 {
		limitStr := ` LIMIT ` + strconv.Itoa(limit)
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
	return q.Queryx(query, injectedNamedParams...)

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

// SaveStatus will save the status in database for the selected entity
func SaveStatus(
	exec *sqlx.Tx,
	entityID int64,
	actionID int,
	statusID int,
) (err error) {
	typeEntity := "entityone"

	err = resetAllPreviousStatuses(exec, typeEntity, entityID)
	if err != nil {
		return fmt.Errorf("islatest %s SaveStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}

	err = insertNewStatus(exec, typeEntity, entityID, actionID, statusID)
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
	exec *sqlx.Tx,
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

	_, err := exec.NamedExec(queryUpd, map[string]interface{}{"entityID": entityID})
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
	exec *sqlx.Tx,
	typeEntity string,
	entityID int64,
	actionID int,
	statusID int,
) error {
	queryIns := fmt.Sprintf(
		"INSERT INTO %s_status(%s_id, action_id, status_id) "+
			" VALUES (:entityID, :actionID, :statusID)",
		typeEntity,
		typeEntity,
	)

	res, err := exec.NamedExec(queryIns,
		map[string]interface{}{
			"entityID": entityID,
			"actionID": actionID,
			"statusID": statusID,
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
