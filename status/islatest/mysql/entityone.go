package mysql

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

	_, errExec = exec.Exec(`USE ` + dbName)
	if errExec != nil {
		return errExec
	}

	return nil
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`DROP DATABASE IF EXISTS ` + dbName)
	if errExec != nil {
		return errExec
	}

	return nil
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(`
        CREATE TABLE IF NOT EXISTS entityone (
            entityone_id BIGINT NOT NULL AUTO_INCREMENT,
            time_created DATETIME NOT NULL DEFAULT NOW(),
            PRIMARY KEY (entityone_id)
        )
        ENGINE = InnoDB
        DEFAULT CHARACTER SET = utf8
        COLLATE = utf8_bin
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`
        CREATE TABLE IF NOT EXISTS entityone_status (
            entityone_id BIGINT NOT NULL,
            action_id BIGINT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT NOW(),
            is_latest TINYINT(1) NULL DEFAULT 1 COMMENT 'can be null',
            UNIQUE INDEX es_ux (is_latest ASC, entityone_id ASC),
            INDEX es_ix (status_id ASC, is_latest ASC),
                CONSTRAINT es_fk_e
                FOREIGN KEY (entityone_id)
                REFERENCES entityone (entityone_id)
                ON DELETE NO ACTION
                ON UPDATE NO ACTION
        )
        ENGINE = InnoDB
        DEFAULT CHARACTER SET = utf8;
    `)
	if errExec != nil {
		return errExec
	}

	return nil
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(`DROP TABLE IF EXISTS entityone_status`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`DROP TABLE IF EXISTS entityone`)
	if errExec != nil {
		return errExec
	}

	return nil
}

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
