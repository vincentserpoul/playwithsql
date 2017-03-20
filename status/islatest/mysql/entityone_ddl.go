package mysql

import "github.com/jmoiron/sqlx"

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`CREATE DATABASE IF NOT EXISTS ` + dbName)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`USE ` + dbName)
	return errExec
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`DROP DATABASE IF EXISTS ` + dbName)
	return errExec
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
            UNIQUE INDEX es_ux (entityone_id ASC, is_latest ASC),
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
