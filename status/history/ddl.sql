
CREATE TABLE IF NOT EXISTS entityone (
    entityone_id BIGSERIAL NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    action_id INT NOT NULL DEFAULT 1,
    time_created DATE NOT NULL DEFAULT CURRENT_DATE,
    time_updated DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (entityone_id)
);

CREATE TABLE IF NOT EXISTS entityone_history (
    entityone_id BIGINT NOT NULL,
    action_id INT NOT NULL DEFAULT 1,
    status_id INT NOT NULL DEFAULT 1,
    time_created DATE NOT NULL DEFAULT CURRENT_DATE,
    CONSTRAINT es_fk_e_eid
        FOREIGN KEY (entityone_id)
        REFERENCES entityone (entityone_id)
);

CREATE INDEX e_idx_sid ON entityone(status_id);

CREATE INDEX es_idx_eid ON entityone_history(entityone_id);

CREATE INDEX es_idx_sid ON entityone_history(status_id);

CREATE INDEX es_idx_aid ON entityone_history(action_id);
