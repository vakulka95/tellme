CREATE OR REPLACE VIEW v$requisition_not_reviewed    AS
SELECT  r.id                          AS id,
        v.id                          AS review_id,
        r.expert_id                   AS expert_id,
        r.phone                       AS phone,
        r.status                      AS status,
        r.created_at                  AS created_at,
        r.updated_at                  AS updated_at
FROM requisitions r
    LEFT OUTER JOIN reviews v ON r.id=v.requisition_id;
