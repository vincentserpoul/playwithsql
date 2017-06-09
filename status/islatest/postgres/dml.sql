--
-- Create
--
INSERT INTO entityone(entityone_id, time_created)
VALUES(DEFAULT, DEFAULT)
RETURNING entityone_id;

UPDATE entityone_status
    SET is_latest = NULL
    WHERE entityone_id = 1 AND is_latest = 1;

INSERT INTO entityone_status(entityone_id, action_id, status_id)
    VALUES (1, 1, 1);

--
-- SaveStatus
--
UPDATE entityone_status
    SET is_latest = NULL
    WHERE entityone_id = 1 AND is_latest = 1;

INSERT INTO entityone_status(entityone_id, action_id, status_id)
    VALUES (1, 1, 1);

--
-- SelectEntityone
--
SELECT
    e.entityone_id, e.time_created,
    es.action_id, es.status_id, es.time_created as status_time_created
FROM entityone e
INNER JOIN entityone_status es ON es.entityone_id = e.entityone_id
WHERE es.is_latest = 1
-- selectEntityoneByStatus filter
AND es.status_id IN (1, 2, 3)
-- selectEntityoneByPK filter
AND e.entityone_id IN (1, 2, 3)
LIMIT 3