CREATE TABLE reviews (
    id                          UUID            NOT NULL,
    expert_id                   UUID            NOT NULL,
    requisition_id              UUID            NOT NULL,
    platform_review             TEXT            NOT NULL,
    consultation_count          INTEGER         NOT NULL,
    consultation_review         TEXT            NOT NULL,
    expert_point                INTEGER         NOT NULL,
    expert_review               TEXT            NOT NULL,
    token                       VARCHAR(64)     NOT NULL,
    status                      VARCHAR(16)     NOT NULL DEFAULT 'requested',
    created_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

ALTER TABLE reviews
    ADD CONSTRAINT reviews_id_pk PRIMARY KEY(id);

ALTER TABLE reviews
    ADD CONSTRAINT reviews_expert_id_fk
        FOREIGN KEY(expert_id) REFERENCES experts(id) ON DELETE RESTRICT;

ALTER TABLE reviews
    ADD CONSTRAINT reviews_requisition_id_fk
        FOREIGN KEY(requisition_id) REFERENCES requisitions(id) ON DELETE RESTRICT;

    CREATE UNIQUE INDEX reviews_token_unq ON reviews(token);

ALTER TABLE reviews ADD CONSTRAINT reviews_status_chk
    CHECK(status IN ('requested', 'completed'));

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
        (SELECT count(r.id) FROM reviews r WHERE r.expert_id=e.id) AS review_count
FROM experts e;