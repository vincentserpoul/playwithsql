package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vincentserpoul/playwithsql/status/islatest"
)

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(exec sqlx.Ext) (id int64, err error) {

	res, err := exec.Exec(`INSERT INTO entityone DEFAULT VALUES`)
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
	exec *sqlx.Tx,
	entityID int64,
	actionID int,
	statusID int,
) error {
	return islatest.SaveStatus(exec, entityID, actionID, statusID)
}

// SelectEntity returns sqlx.Rows
func (link *Link) SelectEntity(
	q *sqlx.DB,
	entityIDs []int64,
	isStatusIDs []int,
	notStatusIDs []int,
	neverStatusIDs []int,
	hasStatusIDs []int,
	limit int,
) (*sqlx.Rows, error) {
	return islatest.SelectEntity(
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
