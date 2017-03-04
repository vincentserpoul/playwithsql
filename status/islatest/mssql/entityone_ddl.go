package mssql

import (
	"context"
	"database/sql"
)

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(ctx context.Context, db *sql.DB, dbName string) (errExec error) {
	_, errExec = db.ExecContext(ctx, `CREATE DATABASE `+dbName)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx, `SET DATABASE = `+dbName)
	return errExec
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(ctx context.Context, db *sql.DB, dbName string) (errExec error) {
	_, errExec = db.ExecContext(ctx, `DROP DATABASE IF EXISTS `+dbName)
	return errExec
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(ctx context.Context, db *sql.DB) (errExec error) {
	_, errExec = db.ExecContext(ctx,
		`
        CREATE TABLE entityone (
            entityone_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
            time_created DATETIME NOT NULL DEFAULT GETDATE()
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx,
		`
        CREATE TABLE entityone_status (
            entityone_id BIGINT NOT NULL,
            action_id INT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT GETDATE(),
            is_latest INT NULL DEFAULT 1,
            CONSTRAINT es_fk_e
            FOREIGN KEY (entityone_id)
            REFERENCES entityone (entityone_id)
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx, `
		CREATE UNIQUE INDEX es_ux_ilei 
		ON entityone_status(entityone_id, is_latest)
		WHERE is_latest IS NOT NULL
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx,
		`CREATE INDEX es_idx1 ON entityone_status(status_id, is_latest)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx,
		`CREATE INDEX es_idx2 ON entityone_status(entityone_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, db *sql.DB) (errExec error) {
	_, errExec = db.ExecContext(ctx, "DROP TABLE IF EXISTS entityone_status")
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx, "DROP TABLE IF EXISTS entityone")
	return errExec
}
