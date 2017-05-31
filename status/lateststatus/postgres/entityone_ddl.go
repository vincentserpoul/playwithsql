package postgres

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
            CREATE TABLE IF NOT EXISTS entityone_status (
				entityone_status_id BIGSERIAL NOT NULL,
				action_id BIGINT NOT NULL DEFAULT 1,
				status_id INT NOT NULL DEFAULT 1,
				time_created DATE NOT NULL DEFAULT CURRENT_DATE,
				PRIMARY KEY (entityone_status_id)
			)
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
            CREATE TABLE IF NOT EXISTS entityone (
				entityone_id BIGSERIAL NOT NULL,
				time_created DATE NOT NULL DEFAULT CURRENT_DATE,
				entityone_status_id BIGINT NOT NULL,
				PRIMARY KEY (entityone_id),
				UNIQUE (entityone_status_id),
				CONSTRAINT e_fk_esi
					FOREIGN KEY (entityone_status_id)
					REFERENCES entityone_status (entityone_status_id)
			)
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx_esi ON entityone(entityone_status_id)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx_si ON entityone_status(status_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(
		ctx,
		"DROP TABLE IF EXISTS entityone",
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		"DROP TABLE IF EXISTS entityone_status",
	)
	return errExec
}
