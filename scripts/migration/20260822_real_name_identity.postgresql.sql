-- Real-name identity tables + drop profile.name (PostgreSQL)

CREATE TABLE IF NOT EXISTS profile_identity (
    account_id          VARCHAR(64)  NOT NULL,
    status              VARCHAR(32)  NOT NULL DEFAULT 'UNVERIFIED',
    document_type       VARCHAR(32),
    real_name_cipher    TEXT,
    document_no_cipher  TEXT,
    document_no_hash    VARCHAR(128),
    verify_channel      VARCHAR(32),
    provider            VARCHAR(32),
    provider_order_no   VARCHAR(128),
    verified_at         TIMESTAMPTZ,
    source_case_id      VARCHAR(64),
    revoked_at          TIMESTAMPTZ,
    revoked_by          VARCHAR(64),
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by          VARCHAR(64),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by          VARCHAR(64),
    CONSTRAINT pk_profile_identity PRIMARY KEY (account_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS uk_profile_identity_document_hash ON profile_identity (document_no_hash) WHERE document_no_hash IS NOT NULL;

CREATE TABLE IF NOT EXISTS real_name_case (
    case_id                     VARCHAR(64)  NOT NULL,
    business_type               VARCHAR(64)  NOT NULL,
    verify_channel              VARCHAR(32)  NOT NULL,
    status                      VARCHAR(32)  NOT NULL,
    account_id                  VARCHAR(64),
    target_account_hint_cipher  TEXT,
    applicant_contact_cipher    TEXT,
    document_type               VARCHAR(32),
    real_name_cipher            TEXT,
    document_no_cipher          TEXT,
    document_no_hash            VARCHAR(128),
    attachment_ids              TEXT,
    payload_cipher              TEXT,
    handler_dept_id             VARCHAR(64),
    provider                    VARCHAR(32),
    provider_order_no           VARCHAR(128),
    submitter_id                VARCHAR(64),
    reviewer_id                 VARCHAR(64),
    reviewed_at                 TIMESTAMPTZ,
    reject_reason               VARCHAR(512),
    expire_at                   TIMESTAMPTZ,
    created_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    created_by                  VARCHAR(64),
    updated_at                  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_by                  VARCHAR(64),
    CONSTRAINT pk_real_name_case PRIMARY KEY (case_id)
);
CREATE INDEX IF NOT EXISTS idx_real_name_case_account ON real_name_case (account_id);
CREATE INDEX IF NOT EXISTS idx_real_name_case_status ON real_name_case (business_type, status);

CREATE TABLE IF NOT EXISTS real_name_case_record (
    record_id       VARCHAR(64)  NOT NULL,
    case_id         VARCHAR(64)  NOT NULL,
    account_id      VARCHAR(64),
    business_type   VARCHAR(64)  NOT NULL,
    action          VARCHAR(32)  NOT NULL,
    status_before   VARCHAR(32),
    status_after    VARCHAR(32),
    verify_channel  VARCHAR(32),
    provider        VARCHAR(32),
    operator_id     VARCHAR(64),
    dept_id         VARCHAR(64),
    remark          VARCHAR(512),
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_real_name_case_record PRIMARY KEY (record_id)
);
CREATE INDEX IF NOT EXISTS idx_real_name_case_record_case ON real_name_case_record (case_id);

ALTER TABLE profile_user_admin DROP COLUMN IF EXISTS name;
ALTER TABLE profile_user_portal DROP COLUMN IF EXISTS name;

-- Dict seeds
INSERT INTO sys_dict VALUES ('101500', 'REAL_NAME_BUSINESS_TYPE', '实名业务类型', 'REAL_NAME_BUSINESS_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101501', 'REAL_NAME_BUSINESS_ACCOUNT_VERIFY', '账号实名认证', 'ACCOUNT_VERIFY', '#18a058', 'SYS', '101500', 'ENABLED', 1, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101502', 'REAL_NAME_BUSINESS_ACCOUNT_RECOVERY', '实名找回账号', 'ACCOUNT_RECOVERY', '#909399', 'SYS', '101500', 'ENABLED', 2, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_dict VALUES ('101510', 'IDENTITY_VERIFY_STATUS', '实名认证状态', 'IDENTITY_VERIFY_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101511', 'IDENTITY_VERIFY_STATUS_UNVERIFIED', '未认证', 'UNVERIFIED', '#909399', 'SYS', '101510', 'ENABLED', 1, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101512', 'IDENTITY_VERIFY_STATUS_PENDING', '审核中', 'PENDING', '#f0a020', 'SYS', '101510', 'ENABLED', 2, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101513', 'IDENTITY_VERIFY_STATUS_VERIFIED', '已认证', 'VERIFIED', '#18a058', 'SYS', '101510', 'ENABLED', 3, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101514', 'IDENTITY_VERIFY_STATUS_REJECTED', '已驳回', 'REJECTED', '#d03050', 'SYS', '101510', 'ENABLED', 4, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_dict VALUES ('101520', 'IDENTITY_DOCUMENT_TYPE', '证件类型', 'IDENTITY_DOCUMENT_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101521', 'IDENTITY_DOCUMENT_ID_CARD', '居民身份证', 'ID_CARD', '#2080f0', 'SYS', '101520', 'ENABLED', 1, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101522', 'IDENTITY_DOCUMENT_PASSPORT', '护照', 'PASSPORT', '#2080f0', 'SYS', '101520', 'ENABLED', 2, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101523', 'IDENTITY_DOCUMENT_EID_CARD', '电子身份证', 'EID_CARD', '#722ed1', 'SYS', '101520', 'ENABLED', 3, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_dict VALUES ('101530', 'IDENTITY_VERIFY_CHANNEL', '认证通道', 'IDENTITY_VERIFY_CHANNEL', '#2080f0', 'SYS', NULL, 'ENABLED', 0, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101531', 'IDENTITY_VERIFY_CHANNEL_THIRD_PARTY', '第三方实人', 'THIRD_PARTY', '#2080f0', 'SYS', '101530', 'ENABLED', 1, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_dict VALUES ('101532', 'IDENTITY_VERIFY_CHANNEL_MANUAL', '人工审核', 'MANUAL', '#f0a020', 'SYS', '101530', 'ENABLED', 2, NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;

-- Admin menu + permission (dynamic route mode loads from sys_resource)
INSERT INTO sys_resource VALUES ('202230', '200003', 'sys-real-name', '实名认证审核', 'MENU', '210001', '/sys/real-name', '/sys/real-name/index.vue', NULL, 'icon-park-outline:id-card', NULL, NULL, 8, true, false, false, 'ENABLED', '实名认证待审队列与审核', NULL, '{}', NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_resource VALUES ('202231', '202230', 'sys-real-name-review', '审核实名认证', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, false, false, false, 'ENABLED', NULL, NULL, '{}', NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
INSERT INTO sys_iam_relation VALUES ('8106481288132251001', 'RESOURCE', '202231', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:realname:review:verify', 'CASCADE', 'ALL', '[]', false, 0, 'ENABLED', '审核实名认证', NULL, NULL, '{}', NOW(), NULL, NOW(), NULL) ON CONFLICT (id) DO NOTHING;
