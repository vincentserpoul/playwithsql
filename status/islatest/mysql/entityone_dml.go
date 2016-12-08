package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(exec sqlx.Ext) (id int64, err error) {

	res, err := exec.Exec(`INSERT INTO entityone VALUES()`)
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
	exec sqlx.Execer,
	entityID int64,
	actionID int,
	statusID int,
) error {
	typeEntity := "entityone"

	queryUpd := fmt.Sprintf(
		"UPDATE %s_status "+
			"SET is_latest = null "+
			"WHERE %s_id= ? AND is_latest = 1",
		typeEntity,
		typeEntity,
	)

	_, err := exec.Exec(queryUpd, entityID)
	if err != nil {
		return fmt.Errorf("infrastructure %s SaveStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}

	queryIns := fmt.Sprintf(
		"INSERT INTO %s_status(%s_id, action_id, status_id) "+
			" VALUES(?, ?, ?)",
		typeEntity,
		typeEntity,
	)

	_, err = exec.Exec(queryIns,
		entityID,
		actionID,
		statusID,
	)

	if err != nil {
		return fmt.Errorf("infrastructure %s SaveStatus(%v, %d, %d): err %v",
			typeEntity,
			entityID,
			actionID,
			statusID,
			err,
		)
	}

	return nil
}

// SelectEntity returns sqlx.Rows
func (link *Link) SelectEntity(
	q sqlx.Queryer,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sqlx.Rows, error) {
	return islatest.SelectEntity(link.IsParamQuestionMark())(
		q,
		entityIDs,
		isStatusIDs,
		notStatusIDs,
		neverStatusIDs,
		hasStatusIDs,
		limit,
	)
}

// IsParamQuestionMark tells if params in SQL are ? or $1
func (link *Link) IsParamQuestionMark() bool {
	return true
}
