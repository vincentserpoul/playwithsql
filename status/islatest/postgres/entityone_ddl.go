package postgres

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
        CREATE TABLE IF NOT EXISTS entityone (
            entityone_id BIGSERIAL NOT NULL,
            time_created DATE NOT NULL DEFAULT CURRENT_DATE,
            PRIMARY KEY (entityone_id)
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx,
		`
        CREATE TABLE IF NOT EXISTS entityone_status (
            entityone_id BIGSERIAL NOT NULL,
            action_id BIGINT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATE NOT NULL DEFAULT CURRENT_DATE,
            is_latest INT NULL DEFAULT 1,
            UNIQUE (is_latest, entityone_id),
            CONSTRAINT es_fk_e
            FOREIGN KEY (entityone_id)
            REFERENCES entityone (entityone_id)
        )
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
