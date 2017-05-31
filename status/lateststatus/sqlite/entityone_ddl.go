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
            CREATE TABLE IF NOT EXISTS entityone_status (
                entityone_status_id INTEGER PRIMARY KEY ASC,
                action_id INT NOT NULL DEFAULT 1,
                status_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT (datetime('now','localtime'))
            )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(ctx,
		`
            CREATE TABLE IF NOT EXISTS entityone (
                entityone_id INTEGER PRIMARY KEY ASC,
                time_created DATETIME NOT NULL DEFAULT (datetime('now','localtime')),
                entityone_status_id INTEGER NOT NULL,
				UNIQUE(entityone_status_id),
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
		`CREATE INDEX es_idx1 ON entityone_status(status_id)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx2 ON entityone(entityone_status_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS entityone`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`DROP TABLE IF EXISTS entityone_status`,
	)
	return errExec
}
