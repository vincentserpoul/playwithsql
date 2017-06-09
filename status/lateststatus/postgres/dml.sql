--
-- Create
--
INSERT INTO entityone(entityone_id, time_created)
    VALUES(DEFAULT, DEFAULT)
    RETURNING entityone_id;

INSERT INTO entityone_status(entityone_status_id, time_created, entityone_id, action_id, status_id)
    VALUES(DEFAULT, DEFAULT, 1, 1, 1)
    RETURNING entityone_status_id;

INSERT INTO entityone_lateststatus(entityone_id, entityone_status_id)
    VALUES (1, 1);

--
-- SaveStatus
--
INSERT INTO entityone_status(entityone_status_id, time_created, entityone_id, action_id, status_id)
    VALUES(DEFAULT, DEFAULT, 1, 1, 1)
    RETURNING entityone_status_id;

UPDATE entityone_lateststatus
SET entityone_status_id = 2
WHERE entityone_id = 1;

--
-- SelectEntityone
--
SELECT
    e.entityone_id, e.time_created,
    es.action_id, es.status_id, es.time_created as status_time_created
FROM entityone e
INNER JOIN entityone_lateststatus el ON el.entityone_id = e.entityone_id
INNER JOIN entityone_status es ON es.entityone_status_id = el.entityone_status_id
WHERE 0=0
-- selectEntityoneByStatus filter
AND es.status_id IN (1, 2, 3)
-- selectEntityoneByPK filter
AND e.entityone_id IN (1, 2, 3)
LIMIT 3;
