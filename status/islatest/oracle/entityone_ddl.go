package oracle

import (
	"context"
	"database/sql"
	"fmt"
)

// Link is used to insert and update in mysql
type Link struct{}

// InitDB create db if not exists
func (link *Link) InitDB(ctx context.Context, db *sql.DB, dbName string) (errExec error) {
	_, errExec = db.ExecContext(ctx, `CREATE USER `+dbName+` IDENTIFIED BY `+dbName)
	if errExec != nil {
		return errExec
	}

	_, errExec = db.ExecContext(ctx, `ALTER SESSION SET CURRENT_SCHEMA = `+dbName)

	return errExec
}

// DestroyDB destroy db if exists
func (link *Link) DestroyDB(ctx context.Context, db *sql.DB, dbName string) (errExec error) {
	_, errExec = db.ExecContext(ctx, `DROP USER `+dbName)
	return errExec
}

// MigrateUp creates the needed tables
func (link *Link) MigrateUp(ctx context.Context, db *sql.DB) (errExec error) {
	_, errExec = db.ExecContext(ctx, `
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

	_, errExec = db.ExecContext(ctx, `CREATE SEQUENCE entityone_seq START WITH 1`)
	if errExec != nil {
		return fmt.Errorf("MigrateUp: create sequence %v", errExec)
	}

	_, errExec = db.ExecContext(ctx, `
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

	_, errExec = db.ExecContext(ctx, `
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

	_, errExec = db.ExecContext(ctx, `
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

	_, errExec = db.ExecContext(ctx,
		`CREATE INDEX es_idx_ei ON entityone_status(entityone_id)`,
	)
	return errExec
}

// MigrateDown destroys the needed tables
func (link *Link) MigrateDown(ctx context.Context, db *sql.DB) (errExec error) {
	_, errExec = db.ExecContext(ctx, `
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

	_, errExec = db.ExecContext(ctx, `
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

	_, errExec = db.ExecContext(ctx, `
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
