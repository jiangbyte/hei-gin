-- +goose Up
CREATE TABLE IF NOT EXISTS sys_account_oauth_binding
(
    id          VARCHAR(64)  NOT NULL,
    account_id  VARCHAR(64)  NOT NULL,
    provider    VARCHAR(32)  NOT NULL,
    open_id     VARCHAR(128) NOT NULL,
    union_id    VARCHAR(128),
    nickname    VARCHAR(128),
    avatar      TEXT,
    raw_profile JSONB        NOT NULL DEFAULT '{}'::jsonb,
    bound_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    created_by  VARCHAR(64),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_by  VARCHAR(64),
    PRIMARY KEY (id),
    CONSTRAINT uq_oauth_provider_open_id UNIQUE (provider, open_id),
    CONSTRAINT uq_oauth_account_provider UNIQUE (account_id, provider)
);

CREATE INDEX IF NOT EXISTS idx_oauth_binding_account ON sys_account_oauth_binding (account_id);
CREATE INDEX IF NOT EXISTS idx_oauth_binding_union ON sys_account_oauth_binding (union_id)
    WHERE union_id IS NOT NULL AND union_id <> '';

-- +goose Down
DROP TABLE IF EXISTS sys_account_oauth_binding CASCADE;
