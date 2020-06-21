CREATE TABLE comments (
    id                          UUID            NOT NULL,
    admin_id                    UUID            NOT NULL,
    expert_id                   UUID            NOT NULL,
    body                        TEXT            NOT NULL,
    created_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

ALTER TABLE comments
    ADD CONSTRAINT comments_id_pk PRIMARY KEY(id);

ALTER TABLE comments
    ADD CONSTRAINT comments_admin_id_fk
        FOREIGN KEY(admin_id) REFERENCES admins(id) ON DELETE RESTRICT;

ALTER TABLE comments
    ADD CONSTRAINT comments_expert_id_fk
        FOREIGN KEY(expert_id) REFERENCES experts(id) ON DELETE RESTRICT;

CREATE OR REPLACE VIEW v$comments     AS
SELECT  c.id                          AS id,
        c.expert_id                   AS expert_id,
        c.admin_id                    AS admin_id,
        a.username                    AS admin_username,
        c.body                        AS body,
        c.created_at                  AS created_at,
        c.updated_at                  AS updated_at
FROM comments c
    INNER JOIN admins a ON a.id=c.admin_id;