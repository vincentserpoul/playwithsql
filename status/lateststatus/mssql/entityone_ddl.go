package mssql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {

	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE entityone_status (
            entityone_status_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
            action_id INT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT GETDATE()
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE entityone (
            entityone_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
            time_created DATETIME NOT NULL DEFAULT GETDATE(),
			entityone_status_id BIGINT NOT NULL,
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx_si ON entityone_status(status_id)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE UNIQUE INDEX es_idx_esi ON entityone(entityone_status_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx, "DROP TABLE IF EXISTS entityone")
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(ctx, "DROP TABLE IF EXISTS entityone_status")
	return errExec
}
