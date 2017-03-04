package sqlite

import (
	"context"
	"database/sql"
)

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(
	ctx context.Context,
	db *sql.DB,
	dbName string,
) (errExec error) {
	return nil
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(
	ctx context.Context,
	db *sql.DB,
	dbName string,
) (errExec error) {
	return nil
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(
	ctx context.Context,
	db *sql.DB,
) (errExec error) {
	_, errExec = db.ExecContext(
		ctx,
		`
        CREATE TABLE IF NOT EXISTS entityone (
            entityone_id INTEGER PRIMARY KEY ASC,
            time_created DATETIME DEFAULT (datetime('now','localtime'))
        )
    	`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(
		ctx,
		`
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
    	`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(
		ctx,
		`CREATE INDEX es_idx1 ON entityone_status(status_id, is_latest)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(
		ctx,
		`CREATE INDEX es_idx2 ON entityone_status(entityone_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(
	ctx context.Context,
	db *sql.DB,
) (errExec error) {

	_, errExec = db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS entityone_status`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS entityone`,
	)
	return errExec
}
