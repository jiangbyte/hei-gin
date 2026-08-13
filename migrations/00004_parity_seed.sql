-- +goose Up
-- Parity seed: audit alert / auth runtime config keys + notice/audit performance indexes.
-- Author: Charlie

-- Audit alert keys (idempotent)
INSERT INTO sys_config (id, config_key, config_value, category, remark, sort_code, value_type, label, scope, scene, is_builtin, ext_json, created_at, updated_at)
VALUES
  ('cfg_parity_audit_enabled', 'AUDIT_ALERT_ENABLED', 'true', 'AUDIT_ALERT', '审计告警总开关', 1, 'BOOL', '审计告警总开关', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_bf', 'AUDIT_ALERT_RULE_BRUTE_FORCE', 'true', 'AUDIT_ALERT', '暴力破解检测', 10, 'BOOL', '暴力破解检测', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_interval', 'AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS', '60', 'AUDIT_ALERT', '分析周期(秒)', 7, 'INT', '分析周期(秒)', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_threshold', 'AUDIT_ALERT_BRUTE_FORCE_THRESHOLD', '10', 'AUDIT_ALERT', '暴力破解频率', 20, 'INT', '暴力破解频率', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_cooldown', 'AUDIT_ALERT_ALERT_COOLDOWN_SECONDS', '1800', 'AUDIT_ALERT', '告警冷却(秒)', 8, 'INT', '告警冷却(秒)', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_push', 'AUDIT_ALERT_NOTIFY_PUSH', 'true', 'AUDIT_ALERT', '推送通知', 3, 'BOOL', '推送通知', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_email', 'AUDIT_ALERT_NOTIFY_EMAIL', 'true', 'AUDIT_ALERT', '邮件通知', 2, 'BOOL', '邮件通知', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_email_to', 'AUDIT_ALERT_NOTIFY_EMAIL_TO', '', 'AUDIT_ALERT', '审计告警收件邮箱', 2, 'STRING', '告警收件邮箱', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_webhook_en', 'AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK', 'false', 'AUDIT_ALERT', '自定义 Webhook 通知', 4, 'BOOL', '自定义 Webhook 通知', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_audit_webhook_url', 'AUDIT_ALERT_WEBHOOK_URL', '', 'AUDIT_ALERT', 'Webhook 地址', 5, 'STRING', 'Webhook 地址', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW())
ON CONFLICT (config_key) DO NOTHING;

-- Auth-related keys commonly required by boot V3 closed-loop (idempotent)
INSERT INTO sys_config (id, config_key, config_value, category, remark, sort_code, value_type, label, scope, scene, is_builtin, ext_json, created_at, updated_at)
VALUES
  ('cfg_parity_auth_portal_reg', 'AUTH_REGISTER_PORTAL_ENABLED', 'true', 'AUTH_REGISTER', 'PORTAL 开放注册', 6, 'BOOL', 'PORTAL 开放注册', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_auth_portal_email', 'AUTH_REGISTER_PORTAL_REQUIRE_EMAIL', 'true', 'AUTH_REGISTER', 'PORTAL 注册要求邮箱', 8, 'BOOL', 'PORTAL 注册要求邮箱', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_auth_portal_phone', 'AUTH_REGISTER_PORTAL_REQUIRE_PHONE', 'false', 'AUTH_REGISTER', 'PORTAL 注册不要求手机', 7, 'BOOL', 'PORTAL 注册要求手机号', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW()),
  ('cfg_parity_auth_default_pwd', 'AUTH_DEFAULT_PASSWORD', '123456', 'AUTH_PASSWORD', '新建账户默认密码', 1, 'STRING', '新建账户默认密码', NULL, NULL, TRUE, '{}'::jsonb, NOW(), NOW())
ON CONFLICT (config_key) DO NOTHING;

-- Performance indexes (boot V6 conceptual)
CREATE INDEX IF NOT EXISTS idx_msg_notice_status_kind_publish
    ON msg_notice (status, kind, publish_at);

CREATE INDEX IF NOT EXISTS idx_msg_notice_status_pinned_publish
    ON msg_notice (status, is_pinned, publish_at DESC);

CREATE INDEX IF NOT EXISTS idx_msg_notice_target_account_types_gin
    ON msg_notice USING GIN (target_account_types);

CREATE INDEX IF NOT EXISTS idx_msg_notice_target_account_ids_gin
    ON msg_notice USING GIN (target_account_ids);

CREATE INDEX IF NOT EXISTS idx_msg_notice_read_account
    ON msg_notice_read (account_type, account_id);

CREATE INDEX IF NOT EXISTS idx_sys_account_identity_account_id
    ON sys_account_identity (account_id);

CREATE UNIQUE INDEX IF NOT EXISTS uq_sys_operation_audit_log_request_id
    ON sys_operation_audit_log (request_id)
    WHERE request_id IS NOT NULL AND btrim(request_id) <> '';

-- +goose Down
DROP INDEX IF EXISTS uq_sys_operation_audit_log_request_id;
DROP INDEX IF EXISTS idx_sys_account_identity_account_id;
DROP INDEX IF EXISTS idx_msg_notice_read_account;
DROP INDEX IF EXISTS idx_msg_notice_target_account_ids_gin;
DROP INDEX IF EXISTS idx_msg_notice_target_account_types_gin;
DROP INDEX IF EXISTS idx_msg_notice_status_pinned_publish;
DROP INDEX IF EXISTS idx_msg_notice_status_kind_publish;

DELETE FROM sys_config WHERE config_key IN (
  'AUDIT_ALERT_ENABLED',
  'AUDIT_ALERT_RULE_BRUTE_FORCE',
  'AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS',
  'AUDIT_ALERT_BRUTE_FORCE_THRESHOLD',
  'AUDIT_ALERT_ALERT_COOLDOWN_SECONDS',
  'AUDIT_ALERT_NOTIFY_PUSH',
  'AUDIT_ALERT_NOTIFY_EMAIL',
  'AUDIT_ALERT_NOTIFY_EMAIL_TO',
  'AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK',
  'AUDIT_ALERT_WEBHOOK_URL',
  'AUTH_REGISTER_PORTAL_ENABLED',
  'AUTH_REGISTER_PORTAL_REQUIRE_EMAIL',
  'AUTH_REGISTER_PORTAL_REQUIRE_PHONE',
  'AUTH_DEFAULT_PASSWORD'
);
