package cockroachdb

import "github.com/jmoiron/sqlx"

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
