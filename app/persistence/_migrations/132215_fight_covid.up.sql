CREATE TABLE admins (
    id                  UUID          NOT NULL,
    username            VARCHAR(128)  NOT NULL,
    password            VARCHAR(256)  NOT NULL,
    status              VARCHAR(16)   NOT NULL DEFAULT 'blocked',
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

ALTER TABLE admins
    ADD CONSTRAINT admins_id_pk PRIMARY KEY(id);

ALTER TABLE admins ADD CONSTRAINT admins_status_chk
    CHECK(status IN ('active', 'blocked'));

CREATE TABLE experts (
    id                          UUID            NOT NULL,
    username                    VARCHAR(256)    NOT NULL,
    gender                      VARCHAR(32)     NOT NULL,
    phone                       VARCHAR(32)     NOT NULL,
    email                       VARCHAR(128)    NOT NULL,
    password                    VARCHAR(512)    NOT NULL,
    specializations             VARCHAR(256)[]  NOT NULL,
    education                   VARCHAR(256)    NOT NULL,
    document_urls               VARCHAR(256)[]  NOT NULL,
    status                      VARCHAR(16)     NOT NULL DEFAULT 'on_review',
    created_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

ALTER TABLE experts
    ADD CONSTRAINT experts_id_pk PRIMARY KEY(id);

CREATE UNIQUE INDEX experts_phone_unq ON experts(phone);
CREATE UNIQUE INDEX experts_email_unq ON experts(email);

ALTER TABLE experts ADD CONSTRAINT experts_status_chk
    CHECK(status IN ('active', 'blocked', 'on_review'));

CREATE TABLE requisitions (
    id                          UUID            NOT NULL,
    expert_id                   VARCHAR(256)    NOT NULL,
    username                    VARCHAR(256)    NOT NULL,
	gender                      VARCHAR(32)     NOT NULL,
	phone                       VARCHAR(32)     NOT NULL,
	diagnosis                   VARCHAR(256)    NOT NULL,
	diagnosis_description       TEXT            NOT NULL,
	expert_gender               VARCHAR(32)     NOT NULL,
	feedback_type               VARCHAR(32)     NOT NULL,
	feedback_contact            VARCHAR(128)    NOT NULL,
	feedback_time               VARCHAR(128)    NOT NULL,
	feedback_week_day           VARCHAR(128)    NOT NULL,
	is_adult                    BOOLEAN         NOT NULL,
    status                      VARCHAR(16)     NOT NULL DEFAULT 'created',
    created_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

ALTER TABLE requisitions
    ADD CONSTRAINT requisitions_id_pk PRIMARY KEY(id);

ALTER TABLE requisitions ADD CONSTRAINT requisitions_status_chk
   CHECK(status IN ('created', 'processing', 'completed'));
