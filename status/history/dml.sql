
-- CREATE

INSERT INTO entityone(action_id, status_id)
    VALUES (:actionID, :statusID);

INSERT INTO entityone_history(entityone_id, action_id, status_id)
    VALUES (:entityoneID, :actionID, :statusID);


-- UPDATE

UPDATE entityone
SET
    action_id = :actionID,
    status_id = :statusID,
    time_updated=:timeUpdated
WHERE entityone_id = :entityoneID;

INSERT INTO entityone_history(entityone_id, action_id, status_id)
    VALUES (:entityoneID, :actionID, :statusID);


-- SELECT

SELECT
    e.entityone_id, e.time_created,
    e.action_id, e.status_id, e.time_updated as status_time_created
FROM entityone e
WHERE 0=0
-- filter on PK
AND e.entityone_id IN (:entityoneIDs)
-- filter on status
AND e.status_id IN (:statusIDs)
LIMIT 3