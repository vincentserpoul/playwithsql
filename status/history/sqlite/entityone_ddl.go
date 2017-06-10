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
                entityone_id INTEGER PRIMARY KEY,
                status_id INT NOT NULL DEFAULT 1,
                action_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT (datetime('now','localtime')),
                time_updated DATETIME NOT NULL DEFAULT (datetime('now','localtime'))
            );

			CREATE INDEX e_idx_sid ON entityone(status_id);

            CREATE TABLE IF NOT EXISTS entityone_history (
                entityone_id INTEGER NOT NULL,
                action_id INT NOT NULL DEFAULT 1,
                status_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT (datetime('now','localtime')),
                CONSTRAINT es_fk_e_eid
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id)
            );


			CREATE INDEX es_idx_eid ON entityone_history(entityone_id);

			CREATE INDEX es_idx_sid ON entityone_history(status_id);
			
			CREATE INDEX es_idx_aid ON entityone_history(action_id);
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
