package cockroachdb

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`CREATE DATABASE IF NOT EXISTS ` + dbName)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`SET DATABASE = ` + dbName)
	return errExec
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`DROP DATABASE IF EXISTS ` + dbName)
	return errExec
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(
		`
        CREATE TABLE IF NOT EXISTS entityone (
            entityone_id BIGSERIAL NOT NULL,
            time_created DATE NOT NULL DEFAULT CURRENT_DATE,
            PRIMARY KEY (entityone_id)
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`
        CREATE TABLE IF NOT EXISTS entityone_status (
            entityone_id BIGSERIAL NOT NULL,
            action_id BIGINT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATE NOT NULL DEFAULT CURRENT_DATE,
            is_latest INT NULL DEFAULT 1,
            UNIQUE (is_latest, entityone_id),
            INDEX (status_id, is_latest),
            CONSTRAINT es_fk_e
            FOREIGN KEY (entityone_id)
            REFERENCES entityone (entityone_id),
            INDEX (entityone_id)
        )
    `)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec("DROP TABLE IF EXISTS entityone_status")
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec("DROP TABLE IF EXISTS entityone")
	return errExec
}

// InsertOne will insert a Entityone into db
func (link *Link) InsertOne(exec sqlx.Ext) (id int64, err error) {
	err = exec.QueryRowx(`
		INSERT INTO entityone(entityone_id, time_created)
		VALUES(DEFAULT, DEFAULT)
		RETURNING entityone_id
	`).Scan(&id)
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
			"WHERE %s_id=$1 AND is_latest = 1",
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
		"INSERT INTO %s_status(%s_id, action_id, status_id) VALUES($1, $2, $3)",
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
	return false
}
