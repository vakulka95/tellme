CREATE OR REPLACE VIEW v$experts_rating       AS
SELECT  e.id                AS id,
        e.username          AS username,
        e.gender            AS gender,
        e.phone             AS phone,
        e.email             AS email,
        e.status            AS status,
        e.created_at        AS created_at,
        e.updated_at        AS updated_at,
        (SELECT count(r.id) FROM requisitions r WHERE r.expert_id=e.id::VARCHAR AND r.status='processing') AS processing_count,
        (SELECT count(r.id) FROM requisitions r WHERE r.expert_id=e.id::VARCHAR AND r.status='completed') AS completed_count,
        (SELECT count(s.id) FROM sessions s WHERE s.expert_id=e.id) AS session_count,
        (SELECT count(r.id) FROM reviews r WHERE r.expert_id=e.id) AS review_count,
        coalesce((SELECT avg(s.expert_point) FROM reviews s WHERE s.expert_id=e.id AND s.status='completed'), 0.0) AS average_rating
FROM experts e;