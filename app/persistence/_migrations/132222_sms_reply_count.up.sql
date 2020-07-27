ALTER TABLE requisitions ADD COLUMN sms_reply_count INTEGER NOT NULL DEFAULT 0;

CREATE OR REPLACE VIEW v$requisitions       AS
SELECT  r.id                                AS id,
        r.expert_id                         AS expert_id,
        r.username                          AS username,
        r.gender                            AS gender,
        r.phone                             AS phone,
        r.diagnosis                         AS diagnosis,
        r.diagnosis_description             AS diagnosis_description,
        r.expert_gender                     AS expert_gender,
        r.feedback_type                     AS feedback_type,
        r.feedback_contact                  AS feedback_contact,
        r.feedback_time                     AS feedback_time,
        r.feedback_week_day                 AS feedback_week_day,
        (SELECT count(s.id) FROM sessions s WHERE s.requisition_id=r.id) AS session_count,
        r.is_adult                          AS is_adult,
        r.status                            AS status,
        r.created_at                        AS created_at,
        r.updated_at                        AS updated_at,
        r.sms_reply_count                   AS sms_reply_count
FROM requisitions r;