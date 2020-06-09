 CREATE TABLE sessions (
    id                          UUID            NOT NULL,
    expert_id                   UUID            NOT NULL,
    requisition_id              UUID            NOT NULL,
    created_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

ALTER TABLE sessions
    ADD CONSTRAINT sessions_id_pk PRIMARY KEY(id);

ALTER TABLE sessions
    ADD CONSTRAINT sessions_expert_id_fk
        FOREIGN KEY(expert_id) REFERENCES experts(id) ON DELETE RESTRICT;

ALTER TABLE sessions
    ADD CONSTRAINT sessions_requisition_id_fk
        FOREIGN KEY(requisition_id) REFERENCES requisitions(id) ON DELETE RESTRICT;

CREATE OR REPLACE VIEW v$sessions     AS
SELECT  s.id                          AS id,
        e.id                          AS expert_id,
        e.username                    AS expert_username,
        r.id                          AS requisition_id,
        r.username                    AS requisition_username,
        r.created_at                  AS created_at,
        r.updated_at                  AS updated_at
FROM sessions s
    INNER JOIN experts e ON e.id=s.expert_id
    INNER JOIN requisitions r ON r.id=s.requisition_id;

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
        (SELECT count(r.id) FROM requisitions r WHERE r.expert_id=e.id::VARCHAR AND r.status='completed') AS completed_count,
        (SELECT count(r.id) FROM reviews r WHERE r.expert_id=e.id) AS review_count,
        (SELECT count(s.id) FROM sessions s WHERE s.expert_id=e.id) AS session_count
FROM experts e;

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
        r.updated_at                        AS updated_at
FROM requisitions r;
