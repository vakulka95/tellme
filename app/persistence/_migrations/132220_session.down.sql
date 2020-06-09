DROP VIEW v$requisitions;

CREATE OR REPLACE VIEW v$experts       AS
SELECT  e.id                AS id,
        e.username          AS username,
        e.gender            AS gender,
        e.phone             AS phone,
        e.email             AS email,
        e.password          AS password,
        e.specializations   AS specializations,
        e.education         AS education,
        e.document_urls     AS document_urls,
        e.status            AS status,
        e.created_at        AS created_at,
        e.updated_at        AS updated_at,
        (SELECT count(r.id) FROM requisitions r WHERE r.expert_id=e.id::VARCHAR AND r.status='processing') AS processing_count,
        (SELECT count(r.id) FROM requisitions r WHERE r.expert_id=e.id::VARCHAR AND r.status='completed') AS completed_count
FROM experts e;

DROP VIEW v$sessions;
DROP TABLE sessions;