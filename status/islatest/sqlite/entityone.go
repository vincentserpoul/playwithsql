package sqlite

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(exec sqlx.Execer, dbName string) (errExec error) {
	return nil
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(exec sqlx.Execer, dbName string) (errExec error) {
	return nil
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(`
        CREATE TABLE IF NOT EXISTS entityone (
            entityone_id INTEGER PRIMARY KEY ASC,
            time_created DATETIME DEFAULT (datetime('now','localtime'))
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`
        CREATE TABLE IF NOT EXISTS entityone_status (
            entityone_id INT NOT NULL,
            action_id INT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME DEFAULT (datetime('now','localtime')),
            is_latest INT(1) NULL DEFAULT 1,
            UNIQUE (is_latest, entityone_id),
            CONSTRAINT es_fk_e
            FOREIGN KEY (entityone_id)
            REFERENCES entityone (entityone_id)
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`CREATE INDEX es_idx1 ON entityone_status(status_id, is_latest)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`CREATE INDEX es_idx2 ON entityone_status(entityone_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(`DROP TABLE IF EXISTS entityone_status`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`DROP TABLE IF EXISTS entityone`)
	return errExec
}

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
		"INSERT INTO %s_status(%s_id, action_id, status_id) VALUES(?, ?, ?)",
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

// IsParamQuestionMark tells if params in SQL are ? or $1
func (link *Link) IsParamQuestionMark() bool {
	return true
}
