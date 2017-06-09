
CREATE TABLE IF NOT EXISTS entityone (
    entityone_id BIGSERIAL NOT NULL,
    time_created DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY (entityone_id)
);

CREATE TABLE IF NOT EXISTS entityone_status (
    entityone_id BIGSERIAL NOT NULL,
    action_id BIGINT NOT NULL DEFAULT 1,
    status_id INT NOT NULL DEFAULT 1,
    time_created DATE NOT NULL DEFAULT CURRENT_DATE,
    is_latest INT NULL DEFAULT 1,
    UNIQUE (is_latest, entityone_id),
    CONSTRAINT es_fk_e
    FOREIGN KEY (entityone_id)
    REFERENCES entityone (entityone_id)
);

CREATE INDEX es_idx1 ON entityone_status(status_id, is_latest);

CREATE INDEX es_idx2 ON entityone_status(entityone_id);
