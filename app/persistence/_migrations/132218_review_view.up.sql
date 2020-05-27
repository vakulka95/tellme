CREATE OR REPLACE VIEW v$reviews    AS
SELECT  r.id                          AS id,
        r.expert_id                   AS expert_id,
        e.username                    AS expert_username,
        r.requisition_id              AS requisition_id,
        r.platform_review             AS platform_review,
        r.consultation_count          AS consultation_count,
        r.consultation_review         AS consultation_review,
        r.expert_point                AS expert_point,
        r.expert_review               AS expert_review,
        r.token                       AS token,
        r.status                      AS status,
        r.created_at                  AS created_at,
        r.updated_at                  AS updated_at
FROM reviews r
    INNER JOIN experts e ON e.id=r.expert_id;
