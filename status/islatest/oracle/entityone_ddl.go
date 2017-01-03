package oracle

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

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
            entityone_id NUMBER(10,0) NOT NULL,
			time_created DATE DEFAULT SYSDATE NOT NULL,
			CONSTRAINT e_pk PRIMARY KEY (entityone_id)
				USING INDEX (CREATE INDEX e_pk_ei ON entityone(entityone_id))
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone %v", errExec)
	}

	_, errExec = exec.Exec(`CREATE SEQUENCE entityone_seq START WITH 1`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create sequence %v", errExec)
	}

	_, errExec = exec.Exec(
		`
		CREATE OR REPLACE TRIGGER entityone_trig
		BEFORE INSERT ON entityone FOR EACH ROW
		BEGIN
			SELECT entityone_seq.NEXTVAL
			INTO   :new.entityone_id
			FROM   dual;
		END;
	`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create trigger %v", errExec)
	}

	_, errExec = exec.Exec(
		`
        CREATE TABLE entityone_status (
            entityone_id NUMBER(10,0) NOT NULL,
            action_id NUMBER(3, 0) NOT NULL,
            status_id NUMBER(3, 0) NOT NULL ,
            time_created DATE DEFAULT SYSDATE NOT NULL,
            is_latest NUMBER(1, 0) DEFAULT 1,
     		CONSTRAINT es_ux_ilei UNIQUE (is_latest, entityone_id)
				USING INDEX (CREATE UNIQUE INDEX es_ux_ilei ON entityone_status(is_latest, entityone_id)),
            CONSTRAINT es_fk_e FOREIGN KEY (entityone_id)
            	REFERENCES entityone(entityone_id)
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone_status %v", errExec)
	}

	_, errExec = exec.Exec(
		`CREATE INDEX es_idx_ei ON entityone_status(entityone_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'entityone_status';
			IF cnt <> 0 THEN
			EXECUTE IMMEDIATE 'DROP TABLE entityone_status';
			END IF;
		END;	
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'entityone';
			IF cnt <> 0 THEN
			EXECUTE IMMEDIATE 'DROP TABLE entityone';
			END IF;
		END;	
	`)
	return errExec
}