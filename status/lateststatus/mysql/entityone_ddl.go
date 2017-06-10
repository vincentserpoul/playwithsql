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
                time_created DATETIME NOT NULL DEFAULT NOW(),
                PRIMARY KEY (entityone_id)
            )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
            CREATE TABLE IF NOT EXISTS entityone_status (
                entityone_status_id BIGINT NOT NULL AUTO_INCREMENT,
                entityone_id BIGINT NOT NULL,
                action_id INT NOT NULL DEFAULT 1,
                status_id INT NOT NULL DEFAULT 1,
                time_created DATETIME NOT NULL DEFAULT NOW(),
                PRIMARY KEY (entityone_status_id),
                INDEX es_idx_sid (status_id),
                INDEX es_fk_ei_e_idx (entityone_id),
                CONSTRAINT es_fk_ei_e
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id)
            );

            CREATE TABLE IF NOT EXISTS entityone_lateststatus (
                entityone_id BIGINT NOT NULL,
                entityone_status_id BIGINT NOT NULL,
                UNIQUE INDEX el_fk_es_esi_idx (entityone_status_id),
                UNIQUE INDEX el_fk_e_ei_idx (entityone_id),
                PRIMARY KEY (entityone_id, entityone_status_id),
                CONSTRAINT el_fk_e_ei
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id),
                CONSTRAINT el_fk_es_esi
                    FOREIGN KEY (entityone_status_id)
                    REFERENCES entityone_status (entityone_status_id)
            );
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
