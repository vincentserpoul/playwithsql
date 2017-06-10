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
			CREATE TABLE IF NOT EXISTS entityone (
				entityone_id BIGSERIAL NOT NULL,
				time_created DATE NOT NULL DEFAULT CURRENT_DATE,
				PRIMARY KEY (entityone_id)
			);

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
			);

			CREATE INDEX es_idx_sid_il ON entityone_status(status_id, is_latest);

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
