-- +goose Up
-- Bootstrap superadmin (password: 123456)

INSERT INTO sys_account (id, password_hash, account_type, account_status, created_at, updated_at)
VALUES ('1', '$2a$10$sOiaxK4ALKezK5lkt5zV6eDsUjqQmMygASlR3VRUuqX6HKnwLTGzW', 'ADMIN', 'ENABLED', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_account_identity (id, account_id, identity_type, identifier, verified, is_primary, bind_status, created_at, updated_at)
VALUES ('1', '1', 'ACCOUNT', 'superadmin', TRUE, TRUE, 'BOUND', NOW(), NOW())
ON CONFLICT (identity_type, identifier) DO NOTHING;

INSERT INTO admin_user_profile (account_id, name, nickname, created_at, updated_at)
VALUES ('1', 'Super Admin', 'superadmin', NOW(), NOW())
ON CONFLICT (account_id) DO NOTHING;

INSERT INTO sys_role (id, code, name, category, scope_type, sort, status, is_builtin, extra, created_at, updated_at)
VALUES ('1', 'SUPER_ADMIN', '超级管理员', 'SYSTEM', 'ALL', 0, 'ENABLED', TRUE, '{}', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO sys_iam_relation (
  id, subject_type, subject_id, account_type, relation_type,
  target_type, target_id, target_key, grant_mode, data_scope,
  custom_scope_dept_ids, is_primary, sort, status, extra, created_at, updated_at
) VALUES (
  '1', 'ACCOUNT', '1', 'ADMIN', 'ACCOUNT_ROLE',
  'ROLE', '1', 'SUPER_ADMIN', 'ALLOW', 'ALL',
  '[]', TRUE, 0, 'ENABLED', '{}', NOW(), NOW()
) ON CONFLICT DO NOTHING;

INSERT INTO sys_iam_relation (
  id, subject_type, subject_id, account_type, relation_type,
  target_type, target_id, target_key, grant_mode, data_scope,
  custom_scope_dept_ids, is_primary, sort, status, extra, created_at, updated_at
) VALUES (
  '2', 'ROLE', '1', 'ADMIN', 'ROLE_PERMISSION',
  'PERMISSION', '*', '*:*:*', 'ALLOW', 'ALL',
  '[]', FALSE, 0, 'ENABLED', '{}', NOW(), NOW()
) ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM sys_iam_relation WHERE id IN ('1', '2');
DELETE FROM sys_role WHERE id = '1';
DELETE FROM admin_user_profile WHERE account_id = '1';
DELETE FROM sys_account_identity WHERE id = '1';
DELETE FROM sys_account WHERE id = '1';
