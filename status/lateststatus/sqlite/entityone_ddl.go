package sqlite

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
                entityone_id INTEGER PRIMARY KEY ASC,
                time_created DATETIME NOT NULL DEFAULT (datetime('now','localtime'))
            );

            CREATE TABLE IF NOT EXISTS entityone_status (
                entityone_status_id INTEGER PRIMARY KEY ASC,
                entityone_id INTEGER NOT NULL,
                action_id INTEGER NOT NULL DEFAULT 1,
                status_id INTEGER NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT (datetime('now','localtime')),
                CONSTRAINT es_fk_ei_e
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id)
            );

            CREATE TABLE IF NOT EXISTS entityone_lateststatus (
                entityone_id INTEGER NOT NULL,
                entityone_status_id INTEGER NOT NULL,
                UNIQUE(entityone_status_id),
                UNIQUE(entityone_id),
                PRIMARY KEY(entityone_id, entityone_status_id),
                CONSTRAINT el_fk_e_ei
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id),
                CONSTRAINT el_fk_es_esi
                    FOREIGN KEY (entityone_status_id)
                    REFERENCES entityone_status (entityone_status_id)
            );

			CREATE INDEX es_idx_sid ON entityone_status(status_id);

			CREATE INDEX es_fk_ei_e_idx ON entityone_status(entityone_id);
	`)

	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx,
		`
			DROP TABLE IF EXISTS entityone_lateststatus;
			DROP TABLE IF EXISTS entityone_status;
			DROP TABLE IF EXISTS entityone;
	`)
	return errExec
}
