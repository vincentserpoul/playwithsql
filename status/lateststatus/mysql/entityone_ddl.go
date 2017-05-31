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
        CREATE TABLE IF NOT EXISTS entityone_status (
            entityone_status_id BIGINT NOT NULL AUTO_INCREMENT,
            action_id BIGINT NOT NULL DEFAULT 1,
            status_id INT NOT NULL DEFAULT 1,
            time_created DATETIME NOT NULL DEFAULT NOW(),
            INDEX es_idx_sid (status_id ASC),
            PRIMARY KEY (entityone_status_id)
        )
        ENGINE = InnoDB
        DEFAULT CHARACTER SET = utf8;
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE IF NOT EXISTS entityone (
            entityone_id BIGINT NOT NULL AUTO_INCREMENT,
            time_created DATETIME NOT NULL DEFAULT NOW(),
            entityone_status_id BIGINT NOT NULL,
            PRIMARY KEY (entityone_id),
            UNIQUE INDEX e_idx_esi (entityone_status_id ASC),
            CONSTRAINT e_fk_esi
                FOREIGN KEY (entityone_status_id)
                REFERENCES entityone_status (entityone_status_id)
        )
        ENGINE = InnoDB
        DEFAULT CHARACTER SET = utf8
        COLLATE = utf8_bin
    `)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(ctx, `DROP TABLE IF EXISTS entityone`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(ctx, `DROP TABLE IF EXISTS entityone_status`)
	return errExec
}
