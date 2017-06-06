package mysql

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
                entityone_id BIGINT NOT NULL AUTO_INCREMENT,
                status_id INT NOT NULL DEFAULT 1,
                action_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT NOW(),
                time_updated DATETIME NOT NULL DEFAULT NOW(),
                PRIMARY KEY (entityone_id),
                INDEX e_idx_sid (status_id ASC)
            )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
            CREATE TABLE IF NOT EXISTS entityone_history (
                entityone_id BIGINT NOT NULL,
                action_id INT NOT NULL DEFAULT 1,
                status_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT NOW(),
                INDEX es_idx_eid (entityone_id),
                INDEX es_idx_sid (status_id ASC),
                INDEX es_idx_aid (action_id ASC),
                CONSTRAINT es_fk_e_eid
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id)
            )
    `)

	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx, `DROP TABLE IF EXISTS entityone_history`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(ctx, `DROP TABLE IF EXISTS entityone`)
	return errExec
}
