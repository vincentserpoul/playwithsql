package oracle

import "github.com/jmoiron/sqlx"

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`CREATE USER ` + dbName + ` IDENTIFIED BY ` + dbName)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`ALTER SESSION SET CURRENT_SCHEMA = ` + dbName)

	return errExec
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(exec sqlx.Execer, dbName string) (errExec error) {
	_, errExec = exec.Exec(`DROP USER ` + dbName)
	return errExec
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(
		`
        CREATE TABLE entityone (
            entityone_id NUMBER(10,0) PRIMARY KEY
				USING INDEX (CREATE INDEX e_pk_ei ON entityone(entityone_id))),
            time_created DATE DEFAULT SYSDATE NOT NULL
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`CREATE SEQUENCE entityone_seq START WITH 1`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`
		CREATE OR REPLACE TRIGGER entityone_trig
		BEFORE INSERT ON entityone 
		FOR EACH ROW
		BEGIN
		SELECT entityone_seq.NEXTVAL
		INTO   :new.id
		FROM   dual;
		END;
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`
        CREATE TABLE entityone_status (
            entityone_id NUMBER(10,0) NOT NULL,
            action_id NUMBER(1, 0) NOT NULL DEFAULT 1,
            status_id NUMBER(1, 0) NOT NULL DEFAULT 1,
            time_created DATE DEFAULT SYSDATE NOT NULL,
            is_latest NUMBER(1, 0) NULL DEFAULT 1,
     		CONSTRAINT es_ux_ilei UNIQUE (is_latest, entityone_id)
				USING INDEX (CREATE UNIQUE INDEX es_ux_ilei ON entityone_status(is_latest, entityone_id)),
            CONSTRAINT es_fk_e FOREIGN KEY (entityone_id)
            	REFERENCES entityone(entityone_id)
        )
    `)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`CREATE INDEX es_idx1 ON entityone_status(status_id, is_latest)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`CREATE INDEX es_idx1 ON entityone_status(status_id, is_latest)`,
	)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(
		`CREATE INDEX es_idx2 ON entityone_status(entityone_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec("DROP TABLE entityone_status")
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec("DROP TABLE entityone")
	return errExec
}
