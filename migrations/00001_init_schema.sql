-- +goose Up
-- Align with hei-fastapi initial schema (27c193fc4b22)

CREATE TABLE IF NOT EXISTS admin_user_profile (
  account_id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NULL,
  nickname VARCHAR(64) NULL,
  avatar TEXT NULL,
  signature TEXT NULL,
  phone VARCHAR(32) NULL,
  email VARCHAR(128) NULL,
  remark TEXT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_admin_user_profile PRIMARY KEY (account_id)
);

CREATE TABLE IF NOT EXISTS cg_test_activity (
  id VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(120) NOT NULL,
  category VARCHAR(32) NULL,
  type VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  cover_url VARCHAR(512) NULL,
  description TEXT NULL,
  start_at TIMESTAMPTZ NOT NULL,
  end_at TIMESTAMPTZ NULL,
  max_participants INTEGER NOT NULL,
  price NUMERIC NOT NULL,
  is_public BOOLEAN NOT NULL,
  need_approval BOOLEAN NOT NULL,
  rule_config JSONB NOT NULL,
  extra JSONB NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  owner_dept_id VARCHAR(64) NULL,
  CONSTRAINT pk_cg_test_activity PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cg_test_catalog (
  id VARCHAR(64) NOT NULL,
  parent_id VARCHAR(64) NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(120) NOT NULL,
  category VARCHAR(32) NULL,
  status VARCHAR(32) NOT NULL,
  sort INTEGER NOT NULL,
  is_visible BOOLEAN NOT NULL,
  icon VARCHAR(128) NULL,
  description TEXT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  owner_dept_id VARCHAR(64) NULL,
  CONSTRAINT pk_cg_test_catalog PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cg_test_knowledge_category (
  id VARCHAR(64) NOT NULL,
  parent_id VARCHAR(64) NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(120) NOT NULL,
  status VARCHAR(32) NOT NULL,
  sort INTEGER NOT NULL,
  is_visible BOOLEAN NOT NULL,
  description TEXT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  owner_dept_id VARCHAR(64) NULL,
  CONSTRAINT pk_cg_test_knowledge_category PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cg_test_knowledge_doc (
  id VARCHAR(64) NOT NULL,
  category_id VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  title VARCHAR(160) NOT NULL,
  type VARCHAR(32) NOT NULL,
  status VARCHAR(32) NOT NULL,
  summary VARCHAR(512) NULL,
  content TEXT NULL,
  author VARCHAR(64) NULL,
  published_at TIMESTAMPTZ NULL,
  view_count INTEGER NOT NULL,
  sort INTEGER NOT NULL,
  is_top BOOLEAN NOT NULL,
  settings JSONB NOT NULL,
  extra JSONB NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_cg_test_knowledge_doc PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cg_test_order (
  id VARCHAR(64) NOT NULL,
  order_no VARCHAR(64) NOT NULL,
  name VARCHAR(120) NOT NULL,
  customer_name VARCHAR(120) NOT NULL,
  customer_phone VARCHAR(32) NULL,
  status VARCHAR(32) NOT NULL,
  type VARCHAR(32) NOT NULL,
  ordered_at TIMESTAMPTZ NOT NULL,
  paid_at TIMESTAMPTZ NULL,
  total_amount NUMERIC NOT NULL,
  item_count INTEGER NOT NULL,
  need_invoice BOOLEAN NOT NULL,
  invoice_config JSONB NOT NULL,
  remark TEXT NULL,
  extra JSONB NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  owner_dept_id VARCHAR(64) NULL,
  CONSTRAINT pk_cg_test_order PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cg_test_order_item (
  id VARCHAR(64) NOT NULL,
  order_id VARCHAR(64) NOT NULL,
  sku_code VARCHAR(64) NOT NULL,
  name VARCHAR(120) NOT NULL,
  category VARCHAR(32) NULL,
  status VARCHAR(32) NOT NULL,
  quantity INTEGER NOT NULL,
  unit_price NUMERIC NOT NULL,
  shipped_at TIMESTAMPTZ NULL,
  is_gift BOOLEAN NOT NULL,
  item_config JSONB NOT NULL,
  remark TEXT NULL,
  extra JSONB NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_cg_test_order_item PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS msg_feedback (
  id VARCHAR(64) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  category VARCHAR(64) NOT NULL,
  contact VARCHAR(255) NULL,
  attach_object_names JSONB NOT NULL,
  status VARCHAR(32) NOT NULL,
  reply TEXT NULL,
  replied_by VARCHAR(64) NULL,
  replied_at TIMESTAMPTZ NULL,
  submitter_account_type VARCHAR(32) NOT NULL,
  submitter_account_id VARCHAR(64) NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_msg_feedback PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS msg_notice (
  id VARCHAR(64) NOT NULL,
  kind VARCHAR(32) NOT NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT NOT NULL,
  content_type VARCHAR(32) NOT NULL,
  category VARCHAR(32) NULL,
  severity VARCHAR(32) NOT NULL,
  target_scope VARCHAR(32) NOT NULL,
  target_account_types JSONB NOT NULL,
  target_account_ids JSONB NOT NULL,
  target_dept_ids JSONB NOT NULL,
  target_role_ids JSONB NOT NULL,
  publish_locations JSONB NOT NULL,
  is_pinned BOOLEAN NOT NULL,
  pinned_until TIMESTAMPTZ NULL,
  sender_account_type VARCHAR(32) NULL,
  sender_account_id VARCHAR(64) NULL,
  source_type VARCHAR(64) NULL,
  source_id VARCHAR(64) NULL,
  status VARCHAR(32) NOT NULL,
  publish_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL,
  expire_at TIMESTAMPTZ NULL,
  view_count INTEGER NOT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_msg_notice PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS msg_notice_read (
  id VARCHAR(64) NOT NULL,
  notice_id VARCHAR(64) NOT NULL,
  account_type VARCHAR(32) NOT NULL,
  account_id VARCHAR(64) NOT NULL,
  read_at TIMESTAMPTZ NULL,
  CONSTRAINT pk_msg_notice_read PRIMARY KEY (id),
  CONSTRAINT uq_msg_notice_read_account UNIQUE (notice_id, account_type, account_id)
);

CREATE TABLE IF NOT EXISTS portal_user_profile (
  account_id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NULL,
  nickname VARCHAR(64) NULL,
  avatar TEXT NULL,
  signature TEXT NULL,
  phone VARCHAR(32) NULL,
  email VARCHAR(128) NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_portal_user_profile PRIMARY KEY (account_id)
);

CREATE TABLE IF NOT EXISTS sys_account (
  id VARCHAR(64) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  account_type VARCHAR(32) NOT NULL,
  account_status VARCHAR(32) NOT NULL,
  cancelled_at TIMESTAMPTZ NULL,
  cancelled_by VARCHAR(64) NULL,
  cancel_reason TEXT NULL,
  cancel_notify_email VARCHAR(128) NULL,
  cancel_notify_phone VARCHAR(32) NULL,
  last_login_ip VARCHAR(64) NULL,
  last_login_address VARCHAR(255) NULL,
  last_login_time TIMESTAMPTZ NULL,
  last_login_device TEXT NULL,
  latest_login_ip VARCHAR(64) NULL,
  latest_login_address VARCHAR(255) NULL,
  latest_login_time TIMESTAMPTZ NULL,
  latest_login_device TEXT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_account PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_account_identity (
  id VARCHAR(64) NOT NULL,
  account_id VARCHAR(64) NOT NULL,
  identity_type VARCHAR(32) NOT NULL,
  identifier VARCHAR(128) NOT NULL,
  verified BOOLEAN NOT NULL,
  is_primary BOOLEAN NOT NULL,
  bind_status VARCHAR(32) NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_account_identity PRIMARY KEY (id),
  CONSTRAINT uq_sys_account_identity_type_identifier UNIQUE (identity_type, identifier)
);

CREATE TABLE IF NOT EXISTS sys_account_password_history (
  id VARCHAR(64) NOT NULL,
  account_id VARCHAR(64) NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  changed_by VARCHAR(64) NULL,
  change_reason VARCHAR(64) NULL,
  created_at TIMESTAMPTZ NULL,
  CONSTRAINT pk_sys_account_password_history PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_alert_log (
  id VARCHAR(64) NOT NULL,
  rule_name VARCHAR(64) NOT NULL,
  severity VARCHAR(16) NOT NULL,
  summary VARCHAR(255) NOT NULL,
  details JSONB NULL,
  notified_via VARCHAR(64) NULL,
  created_at TIMESTAMPTZ NULL,
  CONSTRAINT pk_sys_alert_log PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_banner (
  id VARCHAR(64) NOT NULL,
  title VARCHAR(255) NOT NULL,
  image VARCHAR(500) NOT NULL,
  url VARCHAR(500) NULL,
  link_type VARCHAR(16) NOT NULL,
  summary VARCHAR(500) NULL,
  description TEXT NULL,
  category VARCHAR(32) NOT NULL,
  type VARCHAR(32) NOT NULL,
  position VARCHAR(32) NOT NULL,
  target_account_types JSONB NOT NULL,
  sort INTEGER NOT NULL,
  interaction_count BIGINT NOT NULL,
  status VARCHAR(32) NOT NULL,
  start_at TIMESTAMPTZ NULL,
  end_at TIMESTAMPTZ NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_banner PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_client_module (
  id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  account_type VARCHAR(32) NOT NULL,
  icon VARCHAR(255) NULL,
  color VARCHAR(32) NULL,
  sort INTEGER NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_client_module PRIMARY KEY (id),
  CONSTRAINT uq_sys_client_module_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS sys_client_resource (
  id VARCHAR(64) NOT NULL,
  parent_id VARCHAR(64) NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  resource_type VARCHAR(32) NOT NULL,
  module_id VARCHAR(64) NULL,
  path VARCHAR(255) NULL,
  component VARCHAR(255) NULL,
  redirect VARCHAR(255) NULL,
  icon VARCHAR(255) NULL,
  color VARCHAR(32) NULL,
  href VARCHAR(255) NULL,
  sort INTEGER NOT NULL,
  is_visible BOOLEAN NOT NULL,
  is_cache BOOLEAN NOT NULL,
  is_affix BOOLEAN NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT NULL,
  layout VARCHAR(255) NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_client_resource PRIMARY KEY (id),
  CONSTRAINT uq_sys_client_resource_module_id_code UNIQUE (module_id, code)
);

CREATE TABLE IF NOT EXISTS sys_codegen_field (
  id VARCHAR(64) NOT NULL,
  plan_id VARCHAR(64) NOT NULL,
  table_role VARCHAR(16) NOT NULL,
  column_name VARCHAR(128) NOT NULL,
  column_comment VARCHAR(255) NULL,
  db_type VARCHAR(128) NOT NULL,
  python_type VARCHAR(64) NOT NULL,
  typescript_type VARCHAR(64) NOT NULL,
  form_widget VARCHAR(32) NOT NULL,
  dict_code VARCHAR(128) NULL,
  query_operator VARCHAR(32) NULL,
  show_in_table BOOLEAN NOT NULL,
  show_in_form BOOLEAN NOT NULL,
  show_in_detail BOOLEAN NOT NULL,
  show_in_query BOOLEAN NOT NULL,
  is_primary_key BOOLEAN NOT NULL,
  is_required BOOLEAN NOT NULL,
  is_unique BOOLEAN NOT NULL,
  is_nullable BOOLEAN NOT NULL,
  max_length INTEGER NULL,
  sort INTEGER NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_codegen_field PRIMARY KEY (id),
  CONSTRAINT uq_sys_codegen_field_plan_role_column UNIQUE (plan_id, table_role, column_name)
);

CREATE TABLE IF NOT EXISTS sys_codegen_plan (
  id VARCHAR(64) NOT NULL,
  name VARCHAR(128) NOT NULL,
  gen_type VARCHAR(32) NOT NULL,
  author VARCHAR(64) NOT NULL,
  description TEXT NULL,
  main_table VARCHAR(128) NOT NULL,
  main_pk VARCHAR(128) NOT NULL,
  main_entity_name VARCHAR(128) NOT NULL,
  main_module_path VARCHAR(255) NOT NULL,
  main_business_name VARCHAR(128) NOT NULL,
  api_prefix VARCHAR(255) NOT NULL,
  permission_prefix VARCHAR(128) NOT NULL,
  resource_module_id VARCHAR(64) NULL,
  parent_resource_id VARCHAR(64) NULL,
  menu_name VARCHAR(64) NOT NULL,
  menu_path VARCHAR(255) NOT NULL,
  component_path VARCHAR(255) NOT NULL,
  icon VARCHAR(255) NULL,
  sort INTEGER NOT NULL,
  tree_parent_field VARCHAR(128) NULL,
  tree_label_field VARCHAR(128) NULL,
  sub_table VARCHAR(128) NULL,
  sub_pk VARCHAR(128) NULL,
  sub_foreign_key VARCHAR(128) NULL,
  sub_entity_name VARCHAR(128) NULL,
  sub_business_name VARCHAR(128) NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_codegen_plan PRIMARY KEY (id),
  CONSTRAINT uq_sys_codegen_plan_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS sys_config (
  id VARCHAR(64) NOT NULL,
  config_key VARCHAR(255) NOT NULL,
  config_value TEXT NULL,
  category VARCHAR(255) NULL,
  remark VARCHAR(255) NULL,
  sort_code INTEGER NOT NULL,
  value_type VARCHAR(32) NOT NULL,
  label VARCHAR(128) NULL,
  scope VARCHAR(32) NULL,
  scene VARCHAR(64) NULL,
  is_builtin BOOLEAN NOT NULL,
  ext_json JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_config PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_dept (
  id VARCHAR(64) NOT NULL,
  parent_id VARCHAR(64) NULL,
  master_id VARCHAR(64) NULL,
  deputy_master_id VARCHAR(64) NULL,
  name VARCHAR(64) NOT NULL,
  category VARCHAR(64) NOT NULL,
  sort INTEGER NOT NULL,
  is_virtual BOOLEAN NOT NULL,
  status VARCHAR(32) NOT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_dept PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_dict (
  id VARCHAR(32) NOT NULL,
  code VARCHAR(50) NOT NULL,
  label VARCHAR(255) NULL,
  value VARCHAR(255) NULL,
  color VARCHAR(32) NULL,
  category VARCHAR(64) NULL,
  parent_id VARCHAR(32) NULL,
  status VARCHAR(16) NOT NULL,
  sort INTEGER NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_dict PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_file (
  id VARCHAR(64) NOT NULL,
  object_name VARCHAR(255) NOT NULL,
  original_name VARCHAR(255) NOT NULL,
  storage_provider VARCHAR(32) NOT NULL,
  bucket VARCHAR(255) NULL,
  content_type VARCHAR(128) NOT NULL,
  size BIGINT NOT NULL,
  url VARCHAR(1024) NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_file PRIMARY KEY (id),
  CONSTRAINT uq_sys_file_object_name UNIQUE (object_name)
);

CREATE TABLE IF NOT EXISTS sys_group (
  id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  owner_dept_id VARCHAR(64) NULL,
  description TEXT NULL,
  status VARCHAR(32) NOT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_group PRIMARY KEY (id),
  CONSTRAINT uq_sys_group_name UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS sys_iam_relation (
  id VARCHAR(64) NOT NULL,
  subject_type VARCHAR(32) NOT NULL,
  subject_id VARCHAR(64) NOT NULL,
  account_type VARCHAR(32) NOT NULL,
  relation_type VARCHAR(64) NOT NULL,
  target_type VARCHAR(32) NOT NULL,
  target_id VARCHAR(64) NOT NULL,
  target_key VARCHAR(128) NOT NULL,
  grant_mode VARCHAR(32) NOT NULL,
  data_scope VARCHAR(32) NOT NULL,
  custom_scope_dept_ids JSONB NOT NULL,
  is_primary BOOLEAN NOT NULL,
  sort INTEGER NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT NULL,
  reason TEXT NULL,
  expired_at TIMESTAMPTZ NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_iam_relation PRIMARY KEY (id),
  CONSTRAINT uq_sys_iam_relation_subject_relation_target UNIQUE (subject_type, subject_id, relation_type, target_type, target_id, target_key, account_type)
);

CREATE TABLE IF NOT EXISTS sys_operation_audit_log (
  id VARCHAR(64) NOT NULL,
  module VARCHAR(64) NOT NULL,
  resource_type VARCHAR(128) NULL,
  resource_id VARCHAR(128) NULL,
  action VARCHAR(64) NOT NULL,
  summary VARCHAR(255) NULL,
  before_data JSONB NULL,
  after_data JSONB NULL,
  account_id VARCHAR(64) NULL,
  account_type VARCHAR(32) NULL,
  request_id VARCHAR(64) NULL,
  ip VARCHAR(64) NULL,
  user_agent VARCHAR(512) NULL,
  success BOOLEAN NOT NULL,
  error_message TEXT NULL,
  created_at TIMESTAMPTZ NULL,
  CONSTRAINT pk_sys_operation_audit_log PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_operation_audit_outbox (
  id VARCHAR(64) NOT NULL,
  payload TEXT NOT NULL,
  status VARCHAR(32) NOT NULL,
  attempts INTEGER NOT NULL,
  created_at TIMESTAMPTZ NULL,
  claimed_at TIMESTAMPTZ NULL,
  CONSTRAINT pk_sys_operation_audit_outbox PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_position (
  id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  category VARCHAR(32) NOT NULL,
  owner_dept_id VARCHAR(64) NULL,
  sort INTEGER NOT NULL,
  is_virtual BOOLEAN NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_position PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sys_resource (
  id VARCHAR(64) NOT NULL,
  parent_id VARCHAR(64) NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  resource_type VARCHAR(32) NOT NULL,
  module_id VARCHAR(64) NULL,
  path VARCHAR(255) NULL,
  component VARCHAR(255) NULL,
  redirect VARCHAR(255) NULL,
  icon VARCHAR(255) NULL,
  color VARCHAR(32) NULL,
  href VARCHAR(255) NULL,
  sort INTEGER NOT NULL,
  is_visible BOOLEAN NOT NULL,
  is_cache BOOLEAN NOT NULL,
  is_affix BOOLEAN NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT NULL,
  layout VARCHAR(255) NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_resource PRIMARY KEY (id),
  CONSTRAINT uq_sys_resource_module_id_code UNIQUE (module_id, code)
);

CREATE TABLE IF NOT EXISTS sys_resource_module (
  id VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  client VARCHAR(32) NOT NULL,
  icon VARCHAR(255) NULL,
  color VARCHAR(32) NULL,
  sort INTEGER NOT NULL,
  status VARCHAR(32) NOT NULL,
  description TEXT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_resource_module PRIMARY KEY (id),
  CONSTRAINT uq_sys_resource_module_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS sys_role (
  id VARCHAR(64) NOT NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(64) NOT NULL,
  category VARCHAR(64) NOT NULL,
  scope_type VARCHAR(32) NOT NULL,
  owner_dept_id VARCHAR(64) NULL,
  sort INTEGER NOT NULL,
  status VARCHAR(32) NOT NULL,
  is_builtin BOOLEAN NOT NULL,
  description TEXT NULL,
  extra JSONB NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_role PRIMARY KEY (id),
  CONSTRAINT uq_sys_role_code UNIQUE (code)
);

CREATE TABLE IF NOT EXISTS sys_weak_password (
  id VARCHAR(64) NOT NULL,
  password VARCHAR(255) NOT NULL,
  created_at TIMESTAMPTZ NULL,
  created_by VARCHAR(64) NULL,
  updated_at TIMESTAMPTZ NULL,
  updated_by VARCHAR(64) NULL,
  CONSTRAINT pk_sys_weak_password PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS ix_cg_test_activity_owner_dept_id ON cg_test_activity (owner_dept_id);
CREATE INDEX IF NOT EXISTS ix_cg_test_catalog_owner_dept_id ON cg_test_catalog (owner_dept_id);
CREATE INDEX IF NOT EXISTS ix_cg_test_knowledge_category_owner_dept_id ON cg_test_knowledge_category (owner_dept_id);
CREATE INDEX IF NOT EXISTS ix_cg_test_order_owner_dept_id ON cg_test_order (owner_dept_id);
CREATE INDEX IF NOT EXISTS idx_pwd_history_account_created ON sys_account_password_history (account_id, created_at);
CREATE INDEX IF NOT EXISTS ix_sys_banner_position_status_sort ON sys_banner (position, status, sort);
CREATE INDEX IF NOT EXISTS ix_sys_codegen_field_plan_role_sort ON sys_codegen_field (plan_id, table_role, sort);
CREATE INDEX IF NOT EXISTS ix_sys_codegen_plan_gen_type ON sys_codegen_plan (gen_type);
CREATE INDEX IF NOT EXISTS ix_sys_codegen_plan_main_table ON sys_codegen_plan (main_table);
CREATE INDEX IF NOT EXISTS idx_sys_config_category ON sys_config (category);
CREATE INDEX IF NOT EXISTS idx_sys_config_category_scope_scene ON sys_config (category, scope, scene);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_config_key ON sys_config (config_key);
CREATE INDEX IF NOT EXISTS idx_sys_dict_category ON sys_dict (category);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_dict_code ON sys_dict (code);
CREATE INDEX IF NOT EXISTS idx_sys_dict_parent_id ON sys_dict (parent_id);
CREATE INDEX IF NOT EXISTS ix_sys_iam_relation_account_type_relation ON sys_iam_relation (account_type, relation_type);
CREATE INDEX IF NOT EXISTS ix_sys_iam_relation_subject ON sys_iam_relation (subject_type, subject_id, relation_type);
CREATE INDEX IF NOT EXISTS ix_sys_iam_relation_target ON sys_iam_relation (target_type, target_id, target_key);
CREATE INDEX IF NOT EXISTS idx_sys_operation_audit_account_id ON sys_operation_audit_log (account_id);
CREATE INDEX IF NOT EXISTS idx_sys_operation_audit_created_at ON sys_operation_audit_log (created_at);
CREATE INDEX IF NOT EXISTS idx_sys_operation_audit_module_action ON sys_operation_audit_log (module, action);
CREATE INDEX IF NOT EXISTS idx_sys_operation_audit_resource ON sys_operation_audit_log (resource_type, resource_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_sys_weak_password_password ON sys_weak_password (password);

-- +goose Down
DROP TABLE IF EXISTS sys_weak_password CASCADE;
DROP TABLE IF EXISTS sys_role CASCADE;
DROP TABLE IF EXISTS sys_resource_module CASCADE;
DROP TABLE IF EXISTS sys_resource CASCADE;
DROP TABLE IF EXISTS sys_position CASCADE;
DROP TABLE IF EXISTS sys_operation_audit_outbox CASCADE;
DROP TABLE IF EXISTS sys_operation_audit_log CASCADE;
DROP TABLE IF EXISTS sys_iam_relation CASCADE;
DROP TABLE IF EXISTS sys_group CASCADE;
DROP TABLE IF EXISTS sys_file CASCADE;
DROP TABLE IF EXISTS sys_dict CASCADE;
DROP TABLE IF EXISTS sys_dept CASCADE;
DROP TABLE IF EXISTS sys_config CASCADE;
DROP TABLE IF EXISTS sys_codegen_plan CASCADE;
DROP TABLE IF EXISTS sys_codegen_field CASCADE;
DROP TABLE IF EXISTS sys_client_resource CASCADE;
DROP TABLE IF EXISTS sys_client_module CASCADE;
DROP TABLE IF EXISTS sys_banner CASCADE;
DROP TABLE IF EXISTS sys_alert_log CASCADE;
DROP TABLE IF EXISTS sys_account_password_history CASCADE;
DROP TABLE IF EXISTS sys_account_identity CASCADE;
DROP TABLE IF EXISTS sys_account CASCADE;
DROP TABLE IF EXISTS portal_user_profile CASCADE;
DROP TABLE IF EXISTS msg_notice_read CASCADE;
DROP TABLE IF EXISTS msg_notice CASCADE;
DROP TABLE IF EXISTS msg_feedback CASCADE;
DROP TABLE IF EXISTS cg_test_order_item CASCADE;
DROP TABLE IF EXISTS cg_test_order CASCADE;
DROP TABLE IF EXISTS cg_test_knowledge_doc CASCADE;
DROP TABLE IF EXISTS cg_test_knowledge_category CASCADE;
DROP TABLE IF EXISTS cg_test_catalog CASCADE;
DROP TABLE IF EXISTS cg_test_activity CASCADE;
DROP TABLE IF EXISTS admin_user_profile CASCADE;
