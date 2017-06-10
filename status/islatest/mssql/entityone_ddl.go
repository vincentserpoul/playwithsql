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
			CREATE TABLE entityone (
				entityone_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
				time_created DATETIME NOT NULL DEFAULT GETDATE()
			);

			CREATE TABLE entityone_status (
				entityone_id BIGINT NOT NULL,
				action_id INT NOT NULL DEFAULT 1,
				status_id INT NOT NULL DEFAULT 1,
				time_created DATETIME NOT NULL DEFAULT GETDATE(),
				is_latest INT NULL DEFAULT 1,
				CONSTRAINT es_fk_e
				FOREIGN KEY (entityone_id)
				REFERENCES entityone (entityone_id)
			);

			CREATE UNIQUE INDEX es_ux_ilei 
				ON entityone_status(entityone_id, is_latest)
				WHERE is_latest IS NOT NULL;

			CREATE INDEX es_idx_sid_is ON entityone_status(status_id, is_latest);

			CREATE INDEX es_idx_eid ON entityone_status(entityone_id);
	`)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx,
		`
		DROP TABLE IF EXISTS entityone_status;
		DROP TABLE IF EXISTS entityone;
	`)

	return errExec
}
