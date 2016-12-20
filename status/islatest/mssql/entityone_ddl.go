package mssql

import "github.com/jmoiron/sqlx"

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`CREATE DATABASE ` + dbName)
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
        CREATE TABLE entityone (
            entityone_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
            time_created DATETIME NOT NULL DEFAULT GETDATE()
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`
        CREATE TABLE entityone_status (
            entityone_id BIGINT NOT NULL,
            action_id INT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT GETDATE(),
            is_latest INT NULL DEFAULT 1,
			CONSTRAINT es_ux_ilei 
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
	_, errExec = exec.Exec("DROP TABLE IF EXISTS entityone_status")
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec("DROP TABLE IF EXISTS entityone")
	return errExec
}
