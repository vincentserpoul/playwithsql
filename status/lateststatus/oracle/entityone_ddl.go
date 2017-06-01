package oracle

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Link is used to insert and update in mysql
type Link struct{}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {
	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE entityone (
            entityone_id NUMBER(10,0) NOT NULL,
			time_created DATE DEFAULT SYSDATE NOT NULL,
			CONSTRAINT e_pk PRIMARY KEY (entityone_id)
				USING INDEX (CREATE UNIQUE INDEX e_idx_pk ON entityone(entityone_id))
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE SEQUENCE entityone_seq START WITH 1`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create sequence entityone_seq %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
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
		return fmt.Errorf("MigrateUp: create trigger entityone_trig %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
        CREATE TABLE entityone_status (
            entityone_status_id NUMBER(10,0) NOT NULL,
			entityone_id NUMBER(10,0) NOT NULL,
            action_id NUMBER(3, 0) NOT NULL,
            status_id NUMBER(3, 0) NOT NULL ,
            time_created DATE DEFAULT SYSDATE NOT NULL,
			CONSTRAINT es_pk PRIMARY KEY (entityone_status_id)
				USING INDEX (CREATE UNIQUE INDEX es_idx_pk ON entityone_status(entityone_status_id)),
            CONSTRAINT es_fk_e
				FOREIGN KEY (entityone_id)
				REFERENCES entityone(entityone_id)
        )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create table entityone_status %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE SEQUENCE entityone_status_seq START WITH 1`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create entityone_status_seq %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
		CREATE OR REPLACE TRIGGER entityone_status_trig
		BEFORE INSERT ON entityone_status FOR EACH ROW
		BEGIN
			SELECT entityone_status_seq.NEXTVAL
			INTO   :new.entityone_status_id
			FROM   dual;
		END;
	`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create entityone_status_trig %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx_ei ON entityone_status(entityone_id)`,
	)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create es_idx_ei %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`CREATE INDEX es_idx_si ON entityone_status(status_id)`,
	)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create es_idx_si %v", errExec)
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
            CREATE TABLE entityone_lateststatus (
                entityone_id NUMBER(10,0) NOT NULL,
                entityone_status_id NUMBER(10,0) NOT NULL,
                CONSTRAINT el_pk PRIMARY KEY (entityone_id, entityone_status_id)
					USING INDEX (CREATE UNIQUE INDEX el_pk_ee ON entityone_lateststatus(entityone_id, entityone_status_id)),
				CONSTRAINT el_fk_es_esi_idx UNIQUE (entityone_status_id)
					USING INDEX (CREATE UNIQUE INDEX el_idx_esi ON entityone_lateststatus(entityone_status_id)),
                CONSTRAINT el_fk_e_ei_idx UNIQUE (entityone_id)
					USING INDEX (CREATE UNIQUE INDEX el_idx_ei ON entityone_lateststatus(entityone_id)),
                CONSTRAINT el_fk_e_ei
                    FOREIGN KEY (entityone_id)
                    REFERENCES entityone (entityone_id),
                CONSTRAINT el_fk_es_esi
                    FOREIGN KEY (entityone_status_id)
                    REFERENCES entityone_status (entityone_status_id)
            )
    `)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create entityone_lateststatus %v", errExec)
	}

	return nil
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, exec sqlx.ExecerContext) (errExec error) {

	_, errExec = exec.ExecContext(
		ctx,
		`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_tables WHERE table_name = 'ENTITYONE_LATESTSTATUS';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP TABLE ENTITYONE_LATESTSTATUS';
			END IF;
		END;
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
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

	_, errExec = exec.ExecContext(
		ctx,
		`
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

	_, errExec = exec.ExecContext(
		ctx,
		`
		DECLARE cnt NUMBER;
		BEGIN
			SELECT COUNT(*) INTO cnt FROM user_sequences WHERE sequence_name = 'ENTITYONE_STATUS_SEQ';
			IF cnt <> 0 THEN
				EXECUTE IMMEDIATE 'DROP SEQUENCE entityone_status_seq';
			END IF;
		END;
	`)
	if errExec != nil {
		return errExec
	}

	_, errExec = exec.ExecContext(
		ctx,
		`
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
