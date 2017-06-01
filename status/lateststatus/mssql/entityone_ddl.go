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
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
            CREATE TABLE entityone_status (
                entityone_status_id BIGINT IDENTITY(1, 1) PRIMARY KEY NOT NULL,
                entityone_id BIGINT NOT NULL,
                action_id INT NOT NULL DEFAULT 1,
                status_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT GETDATE(),
                INDEX es_idx_sid (status_id),
                INDEX es_fk_ei_e_idx (entityone_id),
                CONSTRAINT es_fk_ei_e
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id)
            )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
            CREATE TABLE entityone_lateststatus (
                entityone_id BIGINT NOT NULL,
                entityone_status_id BIGINT NOT NULL,
                CONSTRAINT el_fk_es_esi_idx UNIQUE (entityone_status_id),
                CONSTRAINT el_fk_e_ei_idx UNIQUE (entityone_id),
                PRIMARY KEY(entityone_id, entityone_status_id),
                CONSTRAINT el_fk_e_ei
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id),
                CONSTRAINT el_fk_es_esi
                    FOREIGN KEY (entityone_status_id)
                    REFERENCES entityone_status (entityone_status_id)
            )
    `)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {

	_, errExec = exec.ExecContext(ctx, "DROP TABLE IF EXISTS entityone_lateststatus")
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(ctx, "DROP TABLE IF EXISTS entityone_status")
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(ctx, "DROP TABLE IF EXISTS entityone")
	return errExec
}
