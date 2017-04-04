package oracle

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(exec sqlx.Execer) (errExec error) {
	_, errExec = exec.Exec(`
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

	_, errExec = exec.Exec(`
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

	_, errExec = exec.Exec(`
        CREATE TABLE entityone_status (
            entityone_id NUMBER(10,0) NOT NULL,
            action_id NUMBER(3, 0) NOT NULL,
            status_id NUMBER(3, 0) NOT NULL ,
            time_created DATE DEFAULT SYSDATE NOT NULL,
            is_latest NUMBER(1, 0) DEFAULT 1,
            CONSTRAINT es_fk_e FOREIGN KEY (entityone_id)
            	REFERENCES entityone(entityone_id)
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone_status %v", errExec)
	}

	_, errExec = exec.Exec(`
		CREATE UNIQUE INDEX es_ux_il ON entityone_status(
		(
			case when is_latest is not null
				then entityone_id
			else null
			end
		))
	`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create unique index on entityone_status %v", errExec)
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
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'ENTITYONE_STATUS';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP TABLE ENTITYONE_STATUS';
			END IF;
		END;	
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'ENTITYONE';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP TABLE ENTITYONE';
			END IF;
		END;	
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.Exec(`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_sequences WHERE sequence_name = 'ENTITYONE_SEQ';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP SEQUENCE entityone_seq';
			END IF;
		END;	
	`)

	return errExec
}
