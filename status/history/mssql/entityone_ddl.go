package mssql

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx,
		`
        CREATE TABLE entityone (
            entityone_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
            status_id INT NOT NULL DEFAULT 1,
            action_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT GETDATE(),
            time_updated DATETIME NOT NULL DEFAULT GETDATE(),
            INDEX e_idx_sid (status_id)
        );

        CREATE TABLE entityone_history (
            entityone_id BIGINT NOT NULL,
            action_id INT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT GETDATE(),
            INDEX es_idx_eid (entityone_id),
            INDEX es_idx_sid (status_id),
            INDEX es_idx_aid (action_id),
            CONSTRAINT es_fk_e_eid
                FOREIGN KEY (entityone_id)
                REFERENCES entityone (entityone_id)
        );
    `)

	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx,
		`
        DROP TABLE IF EXISTS entityone_history;
        DROP TABLE IF EXISTS entityone;
    `)

	return errExec
}
