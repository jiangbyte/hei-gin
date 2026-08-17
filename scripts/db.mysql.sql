/*
 MySQL schema converted from scripts/db.sql (PostgreSQL dump).
 Charset: utf8mb4

 Usage:
   mysql -u root -p -e "CREATE DATABASE hei_gin DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   mysql -u root -p hei_gin < scripts/db.mysql.sql
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cg_test_activity
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_activity`;
CREATE TABLE `cg_test_activity` (
 `id` varchar(64) NOT NULL,
 `code` varchar(64) NOT NULL,
 `name` varchar(120) NOT NULL,
 `category` varchar(32),
 `type` varchar(32) NOT NULL,
 `status` varchar(32) NOT NULL,
 `cover_url` varchar(512),
 `description` text,
 `start_at` datetime(6) NOT NULL,
 `end_at` datetime(6),
 `max_participants` int NOT NULL,
 `price` decimal(20,6) NOT NULL,
 `is_public` tinyint(1) NOT NULL,
 `need_approval` tinyint(1) NOT NULL,
 `rule_config` json NOT NULL,
 `extra` json,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64),
 `owner_dept_id` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of cg_test_activity
-- ----------------------------
INSERT INTO `cg_test_activity` VALUES ('900000000000000001', 'ACT-BOOTCAMP', '暑期训练营', 'TRAINING', 'OFFLINE', 'ENABLED', 'https://example.com/activity/bootcamp.png', '覆盖文本域、时间、金额、开关、JSON 的 CRUD 测试数据。', '2026-07-19 01:00:00', '2026-07-19 09:00:00', 120, 199.00, 0, 0, '{"limit": {"daily": 3}, "checkin": true}', '{"tags": ["codegen", "crud"]}', '2026-08-08 13:09:50.554189', '1', '2026-08-08 13:09:50.554189', '1', NULL);

-- ----------------------------
-- Table structure for cg_test_catalog
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_catalog`;
CREATE TABLE `cg_test_catalog` (
 `id` varchar(64) NOT NULL,
 `parent_id` varchar(64),
 `code` varchar(64) NOT NULL,
 `name` varchar(120) NOT NULL,
 `category` varchar(32),
 `status` varchar(32) NOT NULL,
 `sort` int NOT NULL,
 `is_visible` tinyint(1) NOT NULL,
 `icon` varchar(128),
 `description` text,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64),
 `owner_dept_id` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of cg_test_catalog
-- ----------------------------
INSERT INTO `cg_test_catalog` VALUES ('900000000000000101', NULL, 'ROOT', '根目录', 'SYSTEM', 'ENABLED', 1, 0, 'folder', '一级节点', '{"level": 1}', '2026-08-08 13:09:50.667031', '1', '2026-08-08 13:09:50.667031', '1', NULL);
INSERT INTO `cg_test_catalog` VALUES ('900000000000000102', '900000000000000101', 'CHILD-A', '子目录A', 'SYSTEM', 'ENABLED', 10, 0, 'folder-open', '二级节点', '{"level": 2}', '2026-08-08 13:09:50.667031', '1', '2026-08-08 13:09:50.667031', '1', NULL);
INSERT INTO `cg_test_catalog` VALUES ('900000000000000103', '900000000000000101', 'CHILD-B', '子目录B', 'BUSINESS', 'DISABLED', 20, 0, 'folder-open', '二级节点', '{"level": 2}', '2026-08-08 13:09:50.667031', '1', '2026-08-08 13:09:50.667031', '1', NULL);

-- ----------------------------
-- Table structure for cg_test_knowledge_category
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_knowledge_category`;
CREATE TABLE `cg_test_knowledge_category` (
 `id` varchar(64) NOT NULL,
 `parent_id` varchar(64),
 `code` varchar(64) NOT NULL,
 `name` varchar(120) NOT NULL,
 `status` varchar(32) NOT NULL,
 `sort` int NOT NULL,
 `is_visible` tinyint(1) NOT NULL,
 `description` text,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64),
 `owner_dept_id` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of cg_test_knowledge_category
-- ----------------------------
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000301', NULL, 'KB', '知识库', 'ENABLED', 1, 0, '根分类', '{"level": 1}', '2026-08-08 13:09:50.841963', '1', '2026-08-08 13:09:50.841963', '1', NULL);
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000302', '900000000000000301', 'KB-DEV', '研发文档', 'ENABLED', 10, 0, '研发相关文档', '{"level": 2}', '2026-08-08 13:09:50.841963', '1', '2026-08-08 13:09:50.841963', '1', NULL);

-- ----------------------------
-- Table structure for cg_test_knowledge_doc
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_knowledge_doc`;
CREATE TABLE `cg_test_knowledge_doc` (
 `id` varchar(64) NOT NULL,
 `category_id` varchar(64) NOT NULL,
 `code` varchar(64) NOT NULL,
 `title` varchar(160) NOT NULL,
 `type` varchar(32) NOT NULL,
 `status` varchar(32) NOT NULL,
 `summary` varchar(512),
 `content` text,
 `author` varchar(64),
 `published_at` datetime(6),
 `view_count` int NOT NULL,
 `sort` int NOT NULL,
 `is_top` tinyint(1) NOT NULL,
 `settings` json NOT NULL,
 `extra` json,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of cg_test_knowledge_doc
-- ----------------------------
INSERT INTO `cg_test_knowledge_doc` VALUES ('900000000000000311', '900000000000000302', 'DOC-CODEGEN-001', '代码生成测试文档', 'ARTICLE', 'PUBLISHED', '用于测试左树右表生成。', '正文内容用于触发 textarea。', 'tester', '2026-07-19 01:19:18', 88, 1, 0, '{"theme": "default", "showToc": true}', '{"tags": ["tree", "table"]}', '2026-08-08 13:09:50.841963', '1', '2026-08-08 13:09:50.841963', '1');

-- ----------------------------
-- Table structure for cg_test_order
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_order`;
CREATE TABLE `cg_test_order` (
 `id` varchar(64) NOT NULL,
 `order_no` varchar(64) NOT NULL,
 `name` varchar(120) NOT NULL,
 `customer_name` varchar(120) NOT NULL,
 `customer_phone` varchar(32),
 `status` varchar(32) NOT NULL,
 `type` varchar(32) NOT NULL,
 `ordered_at` datetime(6) NOT NULL,
 `paid_at` datetime(6),
 `total_amount` decimal(20,6) NOT NULL,
 `item_count` int NOT NULL,
 `need_invoice` tinyint(1) NOT NULL,
 `invoice_config` json NOT NULL,
 `remark` text,
 `extra` json,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64),
 `owner_dept_id` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of cg_test_order
-- ----------------------------
INSERT INTO `cg_test_order` VALUES ('900000000000000201', 'CG-ORDER-001', '测试订单001', '张三', '13800000000', 'PAID', 'NORMAL', '2026-07-19 01:10:00', '2026-07-19 01:20:00', 399.00, 2, 0, '{"taxNo": "91300000000000000X", "title": "张三"}', '主子表生成测试订单', '{"source": "codegen"}', '2026-08-08 13:09:50.732542', '1', '2026-08-08 13:09:50.732542', '1', NULL);

-- ----------------------------
-- Table structure for cg_test_order_item
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_order_item`;
CREATE TABLE `cg_test_order_item` (
 `id` varchar(64) NOT NULL,
 `order_id` varchar(64) NOT NULL,
 `sku_code` varchar(64) NOT NULL,
 `name` varchar(120) NOT NULL,
 `category` varchar(32),
 `status` varchar(32) NOT NULL,
 `quantity` int NOT NULL,
 `unit_price` decimal(20,6) NOT NULL,
 `shipped_at` datetime(6),
 `is_gift` tinyint(1) NOT NULL,
 `item_config` json NOT NULL,
 `remark` text,
 `extra` json,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of cg_test_order_item
-- ----------------------------
INSERT INTO `cg_test_order_item` VALUES ('900000000000000211', '900000000000000201', 'SKU-001', '测试商品A', 'BOOK', 'ENABLED', 1, 199.00, NULL, 0, '{"color": "red"}', '普通明细', '{"line": 1}', '2026-08-08 13:09:50.732542', '1', '2026-08-08 13:09:50.732542', '1');
INSERT INTO `cg_test_order_item` VALUES ('900000000000000212', '900000000000000201', 'SKU-002', '测试商品B', 'COURSE', 'ENABLED', 1, 200.00, '2026-07-19 02:30:00', 0, '{"duration": 30}', '赠品明细', '{"line": 2}', '2026-08-08 13:09:50.732542', '1', '2026-08-08 13:09:50.732542', '1');

-- ----------------------------
-- Table structure for profile_user_admin
-- ----------------------------
DROP TABLE IF EXISTS `profile_user_admin`;
CREATE TABLE `profile_user_admin` (
 `account_id` varchar(64) NOT NULL,
 `name` varchar(64),
 `nickname` varchar(64),
 `avatar` text,
 `signature` text,
 `phone` varchar(32),
 `email` varchar(128),
 `remark` text,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of profile_user_admin
-- ----------------------------
INSERT INTO `profile_user_admin` VALUES ('1', '超级管理员', '超管', 'uploads/2026/08/09/02acc3dee5454d34913b07f49fe59cac.png', NULL, NULL, 'jiangbytebb@163.com', '系统内置超管账户', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 13:17:41.018249', '1');

-- ----------------------------
-- Table structure for profile_user_portal
-- ----------------------------
DROP TABLE IF EXISTS `profile_user_portal`;
CREATE TABLE `profile_user_portal` (
 `account_id` varchar(64) NOT NULL,
 `name` varchar(64),
 `nickname` varchar(64),
 `avatar` text,
 `signature` text,
 `phone` varchar(32),
 `email` varchar(128),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of profile_user_portal
-- ----------------------------
INSERT INTO `profile_user_portal` VALUES ('7491872891940786176', NULL, 'user-a527e592', NULL, NULL, '17286916074', '3317229064@qq.com', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `profile_user_portal` VALUES ('7491847383584804864', '', 'user-171fd244', 'uploads/2026/08/09/85e1b98acfc9465abbbba86ef3b4fec8.jpg', NULL, NULL, 'jiangbyte@163.com', '2026-08-08 13:26:48.032837', NULL, '2026-08-08 13:48:45.931196', '7491847383584804864');

-- ----------------------------
-- Table structure for sys_account
-- ----------------------------
DROP TABLE IF EXISTS `sys_account`;
CREATE TABLE `sys_account` (
 `id` varchar(64) NOT NULL,
 `password_hash` varchar(255) NOT NULL,
 `account_type` varchar(32) NOT NULL,
 `account_status` varchar(32) NOT NULL,
 `cancelled_at` datetime(6),
 `cancelled_by` varchar(64),
 `cancel_reason` text,
 `cancel_notify_email` varchar(128),
 `cancel_notify_phone` varchar(32),
 `last_login_ip` varchar(64),
 `last_login_address` varchar(255),
 `last_login_time` datetime(6),
 `last_login_device` text,
 `latest_login_ip` varchar(64),
 `latest_login_address` varchar(255),
 `latest_login_time` datetime(6),
 `latest_login_device` text,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_account
-- ----------------------------
INSERT INTO `sys_account` VALUES ('7491872891940786176', '$2b$12$kghdYhio.WATOZvDNdTLe.ACM2ibhvP.v88NudZcvjroz/H8M8.z.', 'PORTAL', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account` VALUES ('7491847383584804864', '$2a$10$ZvgY90jMCQpobPlmaCXqie6rCzii8JEciVkXUVM.Kc2DkQHc639xy', 'PORTAL', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '127.0.0.1', NULL, '2026-08-13 00:05:54.418751', 'Desktop', '2026-08-08 13:26:48.032837', NULL, '2026-08-08 22:35:11.342559', '7491847383584804864');
INSERT INTO `sys_account` VALUES ('1', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '127.0.0.1', NULL, '2026-08-13 11:13:03.970991', 'Desktop', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 11:56:13.747886', '1');

-- ----------------------------
-- Table structure for sys_account_identity
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_identity`;
CREATE TABLE `sys_account_identity` (
 `id` varchar(64) NOT NULL,
 `account_id` varchar(64) NOT NULL,
 `identity_type` varchar(32) NOT NULL,
 `identifier` varchar(128) NOT NULL,
 `verified` tinyint(1) NOT NULL,
 `is_primary` tinyint(1) NOT NULL,
 `bind_status` varchar(32) NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_account_identity
-- ----------------------------
INSERT INTO `sys_account_identity` VALUES ('7491872891999506432', '7491872891940786176', 'ACCOUNT', 'usera', 0, 0, 'BOUND', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account_identity` VALUES ('7491872891999506433', '7491872891940786176', 'EMAIL', '3317229064@qq.com', 0, 0, 'BOUND', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account_identity` VALUES ('7491872891999506434', '7491872891940786176', 'PHONE', '17286916074', 0, 0, 'BOUND', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account_identity` VALUES ('2086415885489872897', '1', 'ACCOUNT', 'superadmin', 0, 0, 'BOUND', '2026-08-09 11:34:45.504657', '1', '2026-08-09 11:34:45.504657', '1');
INSERT INTO `sys_account_identity` VALUES ('2086415885552787457', '1', 'EMAIL', 'jiangbytebb@163.com', 0, 0, 'BOUND', '2026-08-09 11:34:45.513351', '1', '2026-08-09 11:34:45.513351', '1');
INSERT INTO `sys_account_identity` VALUES ('2087538500522606594', '7491847383584804864', 'ACCOUNT', 'user', 0, 0, 'BOUND', '2026-08-12 13:55:37.773544', '1', '2026-08-12 13:55:37.773544', '1');
INSERT INTO `sys_account_identity` VALUES ('2087538500522606595', '7491847383584804864', 'EMAIL', 'jiangbyte@163.com', 0, 0, 'BOUND', '2026-08-12 13:55:37.784348', '1', '2026-08-12 13:55:37.784348', '1');

-- ----------------------------
-- Table structure for sys_account_oauth_binding
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_oauth_binding`;
CREATE TABLE `sys_account_oauth_binding` (
 `id` varchar(64) NOT NULL,
 `account_id` varchar(64) NOT NULL,
 `provider` varchar(32) NOT NULL,
 `open_id` varchar(128) NOT NULL,
 `union_id` varchar(128),
 `nickname` varchar(128),
 `avatar` text,
 `raw_profile` json NOT NULL,
 `bound_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_account_oauth_binding
-- ----------------------------

-- ----------------------------
-- Table structure for sys_account_password_history
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_password_history`;
CREATE TABLE `sys_account_password_history` (
 `id` varchar(64) NOT NULL,
 `account_id` varchar(64) NOT NULL,
 `password_hash` varchar(255) NOT NULL,
 `changed_by` varchar(64),
 `change_reason` varchar(64),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_account_password_history
-- ----------------------------
INSERT INTO `sys_account_password_history` VALUES ('7491847383723216896', '7491847383584804864', '$2b$12$wWwMojZp.MT59xARucHu8ODN7EjIHUvc4dwSuKeEnNxgjQZydiKE.', '7491847383584804864', 'register', '2026-08-08 13:26:48.032837');
INSERT INTO `sys_account_password_history` VALUES ('7491872892125335552', '7491872891940786176', '$2b$12$bgU/uMSYlZ.9sYZDHToB4.X6J3EheYp7lpRYheGSQDKrUm8EVp9zC', '7491872891940786176', 'register', '2026-08-08 15:08:09.699685');
INSERT INTO `sys_account_password_history` VALUES ('7491985391344615424', '7491847383584804864', '$2b$12$562T0duxv9fT5lEOWRMMjey8MwEZeXOuyQuP705mJznyOCnbvsxGu', '7491847383584804864', 'self_reset', '2026-08-08 22:35:11.342559');

-- ----------------------------
-- Table structure for sys_alert_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_alert_log`;
CREATE TABLE `sys_alert_log` (
 `id` varchar(64) NOT NULL,
 `rule_name` varchar(64) NOT NULL,
 `severity` varchar(16) NOT NULL,
 `summary` varchar(255) NOT NULL,
 `details` json,
 `notified_via` varchar(64),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_alert_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_banner
-- ----------------------------
DROP TABLE IF EXISTS `sys_banner`;
CREATE TABLE `sys_banner` (
 `id` varchar(64) NOT NULL,
 `title` varchar(255) NOT NULL,
 `image` varchar(500) NOT NULL,
 `url` varchar(500),
 `link_type` varchar(16) NOT NULL,
 `summary` varchar(500),
 `description` text,
 `category` varchar(32) NOT NULL,
 `type` varchar(32) NOT NULL,
 `position` varchar(32) NOT NULL,
 `target_account_types` json NOT NULL,
 `sort` int NOT NULL,
 `interaction_count` bigint NOT NULL,
 `status` varchar(32) NOT NULL,
 `start_at` datetime(6),
 `end_at` datetime(6),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_banner
-- ----------------------------
INSERT INTO `sys_banner` VALUES ('7491889345134235648', '最新公告速递', 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1600&h=700&fit=crop', '/announcements', 'ROUTE', '及时了解平台动态与重要通知。', '门户首页轮播示例图二', 'HOME', 'CAROUSEL', 'HOME_TOP', '["PORTAL"]', 20, 0, 'ENABLED', NULL, NULL, '2026-08-08 16:13:32.393714', NULL, '2026-08-08 16:13:32.393714', NULL);
INSERT INTO `sys_banner` VALUES ('7491889345142624256', '完善个人资料', 'https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=1600&h=700&fit=crop', '/usercenter', 'ROUTE', '前往个人中心完善资料，获得更好的使用体验。', '门户首页轮播示例图三', 'HOME', 'CAROUSEL', 'HOME_TOP', '["PORTAL"]', 30, 0, 'ENABLED', NULL, NULL, '2026-08-08 16:13:32.393714', NULL, '2026-08-08 16:13:32.393714', NULL);
INSERT INTO `sys_banner` VALUES ('7491889345146818560', '管理端运营位示例', 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=1600&h=700&fit=crop', NULL, 'NONE', '仅面向管理端账号类型，不在门户轮播出现。', '管理端展示图示例', 'ADMIN_DASHBOARD', 'CARD', 'ADMIN_TOP', '["ADMIN"]', 10, 0, 'ENABLED', NULL, NULL, '2026-08-08 16:13:32.393714', NULL, '2026-08-08 16:13:32.393714', NULL);

-- ----------------------------
-- Table structure for sys_client_module
-- ----------------------------
DROP TABLE IF EXISTS `sys_client_module`;
CREATE TABLE `sys_client_module` (
 `id` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `code` varchar(64) NOT NULL,
 `account_type` varchar(32) NOT NULL,
 `icon` varchar(255),
 `color` varchar(32),
 `sort` int NOT NULL,
 `status` varchar(32) NOT NULL,
 `description` text,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_client_module
-- ----------------------------
INSERT INTO `sys_client_module` VALUES ('221001', '管理端默认模块', 'admin-default', 'ADMIN', 'icon-park-outline:application-one', NULL, 1, 'ENABLED', '管理端默认客户端模块', '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_client_module` VALUES ('221002', '门户端默认模块', 'portal-default', 'PORTAL', 'icon-park-outline:application-one', NULL, 1, 'ENABLED', '门户端默认客户端模块', '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);

-- ----------------------------
-- Table structure for sys_client_resource
-- ----------------------------
DROP TABLE IF EXISTS `sys_client_resource`;
CREATE TABLE `sys_client_resource` (
 `id` varchar(64) NOT NULL,
 `parent_id` varchar(64),
 `code` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `resource_type` varchar(32) NOT NULL,
 `module_id` varchar(64),
 `path` varchar(255),
 `component` varchar(255),
 `redirect` varchar(255),
 `icon` varchar(255),
 `color` varchar(32),
 `href` varchar(255),
 `sort` int NOT NULL,
 `is_visible` tinyint(1) NOT NULL,
 `is_cache` tinyint(1) NOT NULL,
 `is_affix` tinyint(1) NOT NULL,
 `status` varchar(32) NOT NULL,
 `description` text,
 `layout` varchar(255),
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_client_resource
-- ----------------------------
INSERT INTO `sys_client_resource` VALUES ('222001', NULL, 'home', '首页', 'MENU', '221001', '/home', '/home/index.vue', NULL, 'icon-park-outline:home', NULL, NULL, 1, 0, 0, 0, 'ENABLED', '管理端客户端样例菜单', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_client_resource` VALUES ('222002', NULL, 'home', '首页', 'MENU', '221002', '/home', '/home/index.vue', NULL, 'icon-park-outline:home', NULL, NULL, 1, 0, 0, 0, 'ENABLED', '门户端客户端样例菜单', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);

-- ----------------------------
-- Table structure for sys_codegen_field
-- ----------------------------
DROP TABLE IF EXISTS `sys_codegen_field`;
CREATE TABLE `sys_codegen_field` (
 `id` varchar(64) NOT NULL,
 `plan_id` varchar(64) NOT NULL,
 `table_role` varchar(16) NOT NULL,
 `column_name` varchar(128) NOT NULL,
 `column_comment` varchar(255),
 `db_type` varchar(128) NOT NULL,
 `python_type` varchar(64) NOT NULL,
 `typescript_type` varchar(64) NOT NULL,
 `form_widget` varchar(32) NOT NULL,
 `dict_code` varchar(128),
 `query_operator` varchar(32),
 `show_in_table` tinyint(1) NOT NULL,
 `show_in_form` tinyint(1) NOT NULL,
 `show_in_detail` tinyint(1) NOT NULL,
 `show_in_query` tinyint(1) NOT NULL,
 `is_primary_key` tinyint(1) NOT NULL,
 `is_required` tinyint(1) NOT NULL,
 `is_unique` tinyint(1) NOT NULL,
 `is_nullable` tinyint(1) NOT NULL,
 `max_length` int,
 `sort` int NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_codegen_field
-- ----------------------------
INSERT INTO `sys_codegen_field` VALUES ('2086441951986069505', '2086441951554056193', 'MAIN', 'id', '主键', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 1, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952053178370', '2086441951554056193', 'MAIN', 'code', '活动编码', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 64, 2, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952120287234', '2086441951554056193', 'MAIN', 'name', '活动名称', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 120, 3, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952187396098', '2086441951554056193', 'MAIN', 'category', '活动分类', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 32, 4, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952254504961', '2086441951554056193', 'MAIN', 'type', '活动类型', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 32, 5, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952321613826', '2086441951554056193', 'MAIN', 'status', '状态', 'varchar', 'str', 'string', 'dict', 'COMMON_STATUS', 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, 32, 6, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952388722689', '2086441951554056193', 'MAIN', 'cover_url', '封面地址', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 512, 7, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952451637250', '2086441951554056193', 'MAIN', 'description', '活动描述', 'text', 'str', 'string', 'textarea', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 8, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952556494849', '2086441951554056193', 'MAIN', 'start_at', '开始时间', 'timestamptz', 'datetime', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 9, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952556494850', '2086441951554056193', 'MAIN', 'end_at', '结束时间', 'timestamptz', 'datetime', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 10, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952623603714', '2086441951554056193', 'MAIN', 'max_participants', '最大参与人数', 'int', 'int', 'number', 'number', NULL, 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, NULL, 11, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952690712577', '2086441951554056193', 'MAIN', 'price', '报名费用', 'decimal(20,6)', 'float', 'number', 'number', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 12, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952757821442', '2086441951554056193', 'MAIN', 'is_public', '是否公开', 'tinyint(1)', 'tinyint(1)', 'boolean', 'switch', NULL, 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, NULL, 13, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952829124610', '2086441951554056193', 'MAIN', 'need_approval', '是否需要审批', 'tinyint(1)', 'tinyint(1)', 'boolean', 'switch', NULL, 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, NULL, 14, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952829124611', '2086441951554056193', 'MAIN', 'rule_config', '规则配置', 'json', 'dict', 'Record<string, any>', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 15, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952892039169', '2086441951554056193', 'MAIN', 'extra', '扩展信息', 'json', 'dict', 'Record<string, any>', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 16, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441952959148033', '2086441951554056193', 'MAIN', 'created_at', '创建时间', 'timestamptz', 'datetime', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 17, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441953026256897', '2086441951554056193', 'MAIN', 'created_by', '创建人', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 18, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441953026256898', '2086441951554056193', 'MAIN', 'updated_at', '更新时间', 'timestamptz', 'datetime', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 19, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441953089171457', '2086441951554056193', 'MAIN', 'updated_by', '更新人', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 20, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086441953156280321', '2086441951554056193', 'MAIN', 'owner_dept_id', '所属部门ID（数据范围）', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 21, '2026-08-09 13:18:20.21131', '1', '2026-08-09 13:18:20.21131', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497010528257', '2086449496716926978', 'MAIN', 'id', '主键', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 1, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497056665601', '2086449496716926978', 'MAIN', 'parent_id', '父级ID', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 2, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497115385858', '2086449496716926978', 'MAIN', 'code', '目录编码', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 64, 3, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497161523201', '2086449496716926978', 'MAIN', 'name', '目录名称', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 120, 4, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497195077634', '2086449496716926978', 'MAIN', 'category', '目录分类', 'varchar', 'str', 'string', 'input', NULL, 'LIKE', 0, 0, 0, 0, 0, 0, 0, 0, 32, 5, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497232826369', '2086449496716926978', 'MAIN', 'status', '状态', 'varchar', 'str', 'string', 'dict', 'COMMON_STATUS', 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, 32, 6, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497295740930', '2086449496716926978', 'MAIN', 'sort', '排序', 'int', 'int', 'number', 'number', NULL, 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, NULL, 7, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497497067522', '2086449496716926978', 'MAIN', 'is_visible', '是否显示', 'tinyint(1)', 'tinyint(1)', 'boolean', 'switch', NULL, 'EQ', 0, 0, 0, 0, 0, 0, 0, 0, NULL, 8, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497509650433', '2086449496716926978', 'MAIN', 'icon', '图标', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 128, 9, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497572564994', '2086449496716926978', 'MAIN', 'description', '描述', 'text', 'str', 'string', 'textarea', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 10, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497627090945', '2086449496716926978', 'MAIN', 'extra', '扩展信息', 'json', 'dict', 'Record<string, any>', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 11, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497673228289', '2086449496716926978', 'MAIN', 'created_at', '创建时间', 'timestamptz', 'datetime', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 12, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497719365633', '2086449496716926978', 'MAIN', 'created_by', '创建人', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 13, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497786474498', '2086449496716926978', 'MAIN', 'updated_at', '更新时间', 'timestamptz', 'datetime', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, NULL, 14, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497786474499', '2086449496716926978', 'MAIN', 'updated_by', '更新人', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 15, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');
INSERT INTO `sys_codegen_field` VALUES ('2086449497853583362', '2086449496716926978', 'MAIN', 'owner_dept_id', '所属部门ID（数据范围）', 'varchar', 'str', 'string', 'input', NULL, NULL, 0, 0, 0, 0, 0, 0, 0, 0, 64, 16, '2026-08-09 13:48:19.093754', '1', '2026-08-09 13:48:19.093754', '1');

-- ----------------------------
-- Table structure for sys_codegen_plan
-- ----------------------------
DROP TABLE IF EXISTS `sys_codegen_plan`;
CREATE TABLE `sys_codegen_plan` (
 `id` varchar(64) NOT NULL,
 `name` varchar(128) NOT NULL,
 `gen_type` varchar(32) NOT NULL,
 `author` varchar(64) NOT NULL,
 `description` text,
 `main_table` varchar(128) NOT NULL,
 `main_pk` varchar(128) NOT NULL,
 `main_entity_name` varchar(128) NOT NULL,
 `main_module_path` varchar(255) NOT NULL,
 `main_business_name` varchar(128) NOT NULL,
 `api_prefix` varchar(255) NOT NULL,
 `permission_prefix` varchar(128) NOT NULL,
 `resource_module_id` varchar(64),
 `parent_resource_id` varchar(64),
 `menu_name` varchar(64) NOT NULL,
 `menu_path` varchar(255) NOT NULL,
 `component_path` varchar(255) NOT NULL,
 `icon` varchar(255),
 `sort` int NOT NULL,
 `tree_parent_field` varchar(128),
 `tree_label_field` varchar(128),
 `sub_table` varchar(128),
 `sub_pk` varchar(128),
 `sub_foreign_key` varchar(128),
 `sub_entity_name` varchar(128),
 `sub_business_name` varchar(128),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_codegen_plan
-- ----------------------------
INSERT INTO `sys_codegen_plan` VALUES ('2086441951554056193', '11', 'TABLE', '11', '11111111', 'cg_test_activity', 'id', 'CgTestActivity', 'biz/cg_test_activity', 'CgTestActivity', '/biz/cg-test-activity', 'biz:cgtestactivity', '210001', '202030', 'CgTestActivity', '/biz/cg-test-activity', '/biz/cg-test-activity/index.vue', 'icon-park-outline:code', 99, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 13:18:20.129202', '1', '2026-08-09 13:18:20.129202', '1');
INSERT INTO `sys_codegen_plan` VALUES ('2086449496716926978', '4444', 'TABLE', '444', '', 'cg_test_catalog', 'id', 'CgTestCatalog', 'biz/cg_test_catalog', 'CgTestCatalog', '/biz/cg-test-catalog', 'biz:cgtestcatalog', '210001', '202030', 'CgTestCatalog', '/biz/cg-test-catalog', 'biz/cg-test-catalog/index.vue', 'icon-park-outline:code', 99, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 13:48:19.037199', '1', '2026-08-09 13:48:19.037199', '1');

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config` (
 `id` varchar(64) NOT NULL,
 `config_key` varchar(255) NOT NULL,
 `config_value` text,
 `category` varchar(255),
 `remark` varchar(255),
 `sort_code` int NOT NULL,
 `value_type` varchar(32) NOT NULL,
 `label` varchar(128),
 `scope` varchar(32),
 `scene` varchar(64),
 `is_builtin` tinyint(1) NOT NULL,
 `ext_json` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES ('7491869125225193480', 'STORAGE_MINIO_ACCESS_KEY', 'gAAAAABqd2UjQzg7UyUYFbmdQe6DHLXzJI7dO2Ql7IH_dmCaWvHkCfsmkFOfhGyXG_Q3kAWCsYEoWWZz0kaqgn9XSHGZDoZB8Q==', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 17:19:31.423597', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193481', 'STORAGE_MINIO_SECRET_KEY', 'gAAAAABqd2Uj8lcg6wI7znRSZmoU7WRxiPFY1ZQ9nU5O2kIyQHwbp0xPbRkrP7ww153nsWr-szThKe07RGmeicdbkHLFl4eSmA==', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 17:19:31.423597', '1');
INSERT INTO `sys_config` VALUES ('7878222758007273827', 'COPYRIGHT_TEXT', 'hei-fastapi', 'SYS', '版权文案', 1, 'STRING', '版权文案', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7310793430810554097', 'PUSH_WECHAT_WORK_WEBHOOK', '', 'PUSH', '企业微信 Webhook', 30, 'STRING', '企业微信 Webhook', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7131883348553291227', 'PUSH_LARK_SECRET', '', 'PUSH', '飞书加签密钥', 21, 'STRING', '飞书加签密钥', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7384500216009438994', 'PASSWORD_CUSTOM_WEAK_WORDS', '', 'AUTH_PASSWORD', '自定义弱密码词（逗号分隔）', 20, 'STRING', '自定义弱密码词（逗号分隔）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7791661612876719423', 'PASSWORD_HISTORY_CHECK_COUNT', '5', 'AUTH_PASSWORD', '历史密码检查条数', 16, 'INT', '历史密码检查条数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7296603632888365225', 'AUDIT_ALERT_WEBHOOK_URL', '', 'AUDIT_ALERT', 'Webhook 地址', 5, 'STRING', 'Webhook 地址', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7128087998001730719', 'PASSWORD_MIN_LENGTH', '8', 'AUTH_PASSWORD', '密码最小长度', 10, 'INT', '密码最小长度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7695727871068877789', 'AUDIT_ALERT_NOTIFY_EMAIL', 'TRUE', 'AUDIT_ALERT', '邮件通知', 2, 'tinyint(1)', '邮件通知', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7545090679203476556', 'AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS', '60', 'AUDIT_ALERT', '分析周期(秒)', 7, 'INT', '分析周期(秒)', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7127821217542267803', 'AUTH_REGISTER_ADMIN_DEFAULT_ROLE_ID', '', 'AUTH_REGISTER', 'ADMIN 注册默认角色', 4, 'STRING', 'ADMIN 注册默认角色', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7932896058923796364', 'MAIL_LOCAL_USE_STARTTLS', 'FALSE', 'MAIL', 'SMTP 使用 STARTTLS', 18, 'tinyint(1)', 'SMTP 使用 STARTTLS', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7182623752582531304', 'AUDIT_ALERT_WEBHOOK_SECRET', '', 'AUDIT_ALERT', 'Webhook 签名密钥', 6, 'STRING', 'Webhook 签名密钥', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7121863095719606243', 'AUTH_REGISTER_ADMIN_DEFAULT_DEPT_ID', '', 'AUTH_REGISTER', 'ADMIN 注册默认部门', 5, 'STRING', 'ADMIN 注册默认部门', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7557049262663315054', 'MAIL_LOCAL_HOST', 'localhost', 'MAIL', 'SMTP 服务器地址', 10, 'STRING', 'SMTP 服务器地址', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7345028346771097677', 'PASSWORD_MAX_LENGTH', '128', 'AUTH_PASSWORD', '密码最大长度', 11, 'INT', '密码最大长度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7337967471264658438', 'AUDIT_ALERT_ALERT_COOLDOWN_SECONDS', '1800', 'AUDIT_ALERT', '告警冷却(秒)', 8, 'INT', '告警冷却(秒)', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7915486803474335428', 'PASSWORD_MAX_CONSECUTIVE_CHARS', '3', 'AUTH_PASSWORD', '最大连续相同字符数', 13, 'INT', '最大连续相同字符数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7116199988527416451', 'AUTH_PASSWORD_RESET_URL_ADMIN', 'http://localhost:5173/auth/forgot-password', 'AUTH_TOKEN', 'ADMIN 密码重置页完整 URL', 3, 'STRING', 'ADMIN 密码重置页完整 URL', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7331249362091683364', 'MAIL_TENCENT_SECRET_KEY', '', 'MAIL', '腾讯云邮件 SecretKey', 31, 'STRING', '腾讯云邮件 SecretKey', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7667316617799595990', 'AUDIT_ALERT_RULE_BULK_DELETE', 'TRUE', 'AUDIT_ALERT', '批量删除检测', 13, 'tinyint(1)', '批量删除检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7486049791614620352', 'AUDIT_ALERT_IP_ANOMALY_THRESHOLD', '3', 'AUDIT_ALERT', 'IP异常阈值', 22, 'INT', 'IP异常阈值', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7564883844348625200', 'AUDIT_ALERT_RULE_IP_ANOMALY', 'TRUE', 'AUDIT_ALERT', 'IP 异常检测', 14, 'tinyint(1)', 'IP 异常检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7965270553477718376', 'AUTH_LOGIN_ADMIN_ALLOW_PHONE', 'TRUE', 'AUTH_LOGIN', 'ADMIN 允许手机号登录', 13, 'tinyint(1)', 'ADMIN 允许手机号登录', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7165407484130097877', 'AUTH_LOGIN_PORTAL_MAX_FAILURES', '5', 'AUTH_LOGIN', 'PORTAL 最大失败次数', 19, 'INT', 'PORTAL 最大失败次数', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7707546743105654318', 'AUTH_LOGIN_IP_MAX_FAILURES', '30', 'AUTH_LOGIN', '单 IP 最大登录失败次数', 3, 'INT', '单 IP 最大登录失败次数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7645847924523510728', 'AUTH_REGISTER_PORTAL_REQUIRE_EMAIL', 'TRUE', 'AUTH_REGISTER', 'PORTAL 注册要求邮箱', 8, 'tinyint(1)', 'PORTAL 注册要求邮箱', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7110705537461490506', 'AUTH_REGISTER_PORTAL_DEFAULT_DEPT_ID', '', 'AUTH_REGISTER', 'PORTAL 注册默认部门', 10, 'STRING', 'PORTAL 注册默认部门', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7607784543080994141', 'SMS_ALIYUN_ACCESS_KEY_ID', '', 'SMS', '阿里云短信 AccessKeyId', 10, 'STRING', '阿里云短信 AccessKeyId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7350757244774876723', 'DEFAULT_MESSAGE_PUSH_ENGINE', 'DINGTALK', 'PUSH', '默认消息推送引擎', 1, 'STRING', '默认消息推送引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7858449177900745217', 'DEFAULT_EMAIL_ENGINE', 'LOCAL', 'MAIL', '默认邮件引擎', 1, 'STRING', '默认邮件引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7204221884041359188', 'AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK', 'FALSE', 'AUDIT_ALERT', '自定义 Webhook 通知', 4, 'tinyint(1)', '自定义 Webhook 通知', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7436928202175303081', 'PASSWORD_VALIDITY_DAYS', '90', 'AUTH_PASSWORD', '密码有效期（天）', 18, 'INT', '密码有效期（天）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7395157445765927617', 'SMS_TENCENT_SECRET_KEY', '', 'SMS', '腾讯云短信 SecretKey', 21, 'STRING', '腾讯云短信 SecretKey', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7440731969037919022', 'MAIL_LOCAL_FROM_EMAIL', 'test@hei-fastapi.local', 'MAIL', '发件人邮箱', 14, 'STRING', '发件人邮箱', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7774327802304558750', 'PUSH_DINGTALK_SECRET', '', 'PUSH', '钉钉加签密钥', 11, 'STRING', '钉钉加签密钥', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7840109703376148568', 'MAIL_LOCAL_PASSWORD', '', 'MAIL', 'SMTP 密码', 13, 'STRING', 'SMTP 密码', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7143992105573126791', 'MAIL_TENCENT_SECRET_ID', '', 'MAIL', '腾讯云邮件 SecretId', 30, 'STRING', '腾讯云邮件 SecretId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7159687456156883406', 'PASSWORD_FORBID_WEAK_LIST', 'TRUE', 'AUTH_PASSWORD', '禁止弱密码库命中', 17, 'tinyint(1)', '禁止弱密码库命中', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7107236720876061841', 'AUTH_REGISTER_ADMIN_ENABLED', 'FALSE', 'AUTH_REGISTER', 'ADMIN 开放注册', 1, 'tinyint(1)', 'ADMIN 开放注册', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7312697174587677983', 'MAIL_LOCAL_USE_SSL', 'FALSE', 'MAIL', 'SMTP 使用 SSL', 17, 'tinyint(1)', 'SMTP 使用 SSL', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7279087550089826228', 'MAIL_LOCAL_USERNAME', '', 'MAIL', 'SMTP 用户名', 12, 'STRING', 'SMTP 用户名', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7817231488316128319', 'DEFAULT_SMS_ENGINE', 'ALIYUN', 'SMS', '默认短信引擎', 1, 'STRING', '默认短信引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7841576503746820287', 'MAIL_LOCAL_PORT', '1025', 'MAIL', 'SMTP 端口', 11, 'INT', 'SMTP 端口', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7717547996356753442', 'PUSH_LARK_WEBHOOK', '', 'PUSH', '飞书 Webhook', 20, 'STRING', '飞书 Webhook', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7148106912230055136', 'PASSWORD_COMPLEXITY', 'DIGITS_UPPER_LOWER_SPECIAL', 'AUTH_PASSWORD', '密码复杂度', 12, 'STRING', '密码复杂度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7569398987799485767', 'PUSH_DINGTALK_WEBHOOK', '', 'PUSH', '钉钉 Webhook', 10, 'STRING', '钉钉 Webhook', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7743096612050880755', 'PASSWORD_FORBID_USER_INFO', 'TRUE', 'AUTH_PASSWORD', '禁止包含用户信息', 14, 'tinyint(1)', '禁止包含用户信息', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7297229307635338316', 'MAIL_ALIYUN_ACCOUNT_NAME', '', 'MAIL', '阿里云发信地址', 22, 'STRING', '阿里云发信地址', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7156018694924792840', 'SMS_ALIYUN_SIGN_NAME', '', 'SMS', '阿里云短信签名', 12, 'STRING', '阿里云短信签名', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7961221992910853426', 'PASSWORD_CHANGE_VERIFY_METHOD', 'OLD_PASSWORD', 'AUTH_PASSWORD', '自助改密验证方式', 2, 'STRING', '自助改密验证方式', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7223837583572900295', 'SMS_TEMPLATE_CHANGE_PASSWORD_CODE', '{"code": "", "content": "改密验证码 {{code}}"}', 'SMS_TEMPLATE', '修改密码短信模板', 2, 'JSON', '修改密码短信模板', NULL, 'CHANGE_PASSWORD_CODE', 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7943155149240026436', 'AUTH_PASSWORD_RESET_TOKEN_TTL_SECONDS', '600', 'AUTH_TOKEN', '密码重置 Token 有效期（秒）', 2, 'INT', '密码重置 Token 有效期（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7640258913173401645', 'AUDIT_ALERT_ENABLED', 'TRUE', 'AUDIT_ALERT', '审计告警总开关', 1, 'tinyint(1)', '审计告警总开关', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7259830455098789920', 'AUTH_REGISTER_PORTAL_ENABLED', 'TRUE', 'AUTH_REGISTER', 'PORTAL 开放注册', 6, 'tinyint(1)', 'PORTAL 开放注册', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7925585644534920178', 'AUTH_LOGIN_ADMIN_FAILURE_WINDOW_SECONDS', '300', 'AUTH_LOGIN', 'ADMIN 登录失败窗口（秒）', 10, 'INT', 'ADMIN 登录失败窗口（秒）', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7371141348954395191', 'AUTH_LOGIN_ADMIN_MAX_FAILURES', '5', 'AUTH_LOGIN', 'ADMIN 最大失败次数', 11, 'INT', 'ADMIN 最大失败次数', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7616410287911469795', 'AUTH_REGISTER_PORTAL_REQUIRE_PHONE', 'FALSE', 'AUTH_REGISTER', 'PORTAL 注册要求手机号', 7, 'tinyint(1)', 'PORTAL 注册要求手机号', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7389279063701289685', 'AUTH_REGISTER_ADMIN_REQUIRE_PHONE', 'FALSE', 'AUTH_REGISTER', 'ADMIN 注册要求手机号', 2, 'tinyint(1)', 'ADMIN 注册要求手机号', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7680761152547899501', 'AUTH_REGISTER_ADMIN_REQUIRE_EMAIL', 'FALSE', 'AUTH_REGISTER', 'ADMIN 注册要求邮箱', 3, 'tinyint(1)', 'ADMIN 注册要求邮箱', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7123144141092942155', 'AUTH_OAUTH_PORTAL_GITEE_CLIENT_SECRET', '', 'AUTH_OAUTH', '门户 Gitee ClientSecret', 13, 'STRING', 'Client Secret', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.427549', NULL, '2026-08-12 15:57:59.427549', NULL);
INSERT INTO `sys_config` VALUES ('7734399746267501494', 'AUTH_OAUTH_PORTAL_QQ_CLIENT_SECRET', '', 'AUTH_OAUTH', '门户 QQ ClientSecret', 23, 'STRING', 'Client Secret', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.443762', NULL, '2026-08-12 15:57:59.443762', NULL);
INSERT INTO `sys_config` VALUES ('7744551651175282801', 'AUTH_OAUTH_PORTAL_GITEE_CLIENT_ID', '', 'AUTH_OAUTH', '门户 Gitee ClientId', 12, 'STRING', 'Client ID', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.421381', NULL, '2026-08-12 15:57:59.421381', '1');
INSERT INTO `sys_config` VALUES ('7125885979202993561', 'AUTH_OAUTH_PORTAL_GITEE_REDIRECT_URI', '', 'AUTH_OAUTH', '门户 Gitee 回调', 14, 'STRING', 'Redirect URI', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.431808', NULL, '2026-08-12 15:57:59.431808', '1');
INSERT INTO `sys_config` VALUES ('7691150031767579943', 'AUTH_OAUTH_PORTAL_QQ_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户 QQ 登录', 21, 'tinyint(1)', '启用', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.435873', NULL, '2026-08-12 15:57:59.435873', '1');
INSERT INTO `sys_config` VALUES ('7149430200646252861', 'AUTH_OAUTH_PORTAL_QQ_CLIENT_ID', '', 'AUTH_OAUTH', '门户 QQ ClientId', 22, 'STRING', 'Client ID', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.439704', NULL, '2026-08-12 15:57:59.439704', '1');
INSERT INTO `sys_config` VALUES ('7462877140739332791', 'MAIL_TEMPLATE_REGISTER_SUCCESS', '{"subject": "欢迎注册 {{app_name}}", "body": "账号 {{account}} 注册成功。"}', 'MAIL_TEMPLATE', '注册成功邮件模板', 4, 'JSON', '注册成功邮件模板', NULL, 'REGISTER_SUCCESS', 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7833544856736466882', 'ACCOUNT_CANCEL_RETENTION_DAYS', '15', 'OTHER', '注销账号保留天数', 10, 'INT', '注销账号保留天数', NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7884464514418971420', 'AUDIT_ALERT_BRUTE_FORCE_THRESHOLD', '10', 'AUDIT_ALERT', '暴力破解阈值', 20, 'INT', '暴力破解阈值', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7780925862931981131', 'MAIL_LOCAL_AUTH_REQUIRED', 'FALSE', 'MAIL', 'SMTP 是否需要认证', 16, 'tinyint(1)', 'SMTP 是否需要认证', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7168944809004050220', 'SMS_TENCENT_SECRET_ID', '', 'SMS', '腾讯云短信 SecretId', 20, 'STRING', '腾讯云短信 SecretId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7650241926616235479', 'COPYRIGHT_URL', '', 'SYS', '版权链接', 2, 'STRING', '版权链接', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7642474443917569970', 'MAIL_ALIYUN_ACCESS_KEY_ID', '', 'MAIL', '阿里云邮件 AccessKeyId', 20, 'STRING', '阿里云邮件 AccessKeyId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7377151897664869996', 'AUDIT_ALERT_NOTIFY_PUSH', 'TRUE', 'AUDIT_ALERT', '推送通知', 3, 'tinyint(1)', '推送通知', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7387934960432348080', 'MAIL_TEMPLATE_RESET_PASSWORD_CODE', '{"subject": "{{app_name}} 密码重置", "body": "请点击以下链接重置密码，该链接将在 {{expire_minutes}} 分钟内有效。\n\n{{reset_link}}"}', 'MAIL_TEMPLATE', '重置密码邮件模板', 1, 'JSON', '重置密码邮件模板', NULL, 'RESET_PASSWORD_CODE', 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7639638206762559428', 'AUTH_PASSWORD_RESET_URL_PORTAL', 'http://localhost:5174/auth/forgot-password', 'AUTH_TOKEN', 'PORTAL 密码重置页完整 URL', 4, 'STRING', 'PORTAL 密码重置页完整 URL', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7491104580893387509', 'SMS_TEMPLATE_LOGIN_CODE', '{"code": "", "content": "登录验证码 {{code}}"}', 'SMS_TEMPLATE', '登录验证码短信模板', 1, 'JSON', '登录验证码短信模板', NULL, 'LOGIN_CODE', 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7971733830500536612', 'PASSWORD_FORBID_HISTORICAL', 'TRUE', 'AUTH_PASSWORD', '禁止复用历史密码', 15, 'tinyint(1)', '禁止复用历史密码', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7231060999533820222', 'MAIL_TEMPLATE_LOGIN_CODE', '{"subject": "{{app_name}} 登录验证码", "body": "您的登录验证码是 {{code}}，{{expire_minutes}} 分钟内有效。"}', 'MAIL_TEMPLATE', '登录验证码邮件模板', 2, 'JSON', '登录验证码邮件模板', NULL, 'LOGIN_CODE', 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7799955389648191999', 'PASSWORD_EXPIRY_WARNING_DAYS', '7', 'AUTH_PASSWORD', '密码过期提前提醒（天）', 19, 'INT', '密码过期提前提醒（天）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7137492378552518589', 'AUTH_DEFAULT_PASSWORD', '', 'AUTH_PASSWORD', '新建账户默认密码', 1, 'STRING', '新建账户默认密码', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7815227162417732636', 'MAIL_ALIYUN_ACCESS_KEY_SECRET', '', 'MAIL', '阿里云邮件 AccessKeySecret', 21, 'STRING', '阿里云邮件 AccessKeySecret', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7386527519961798374', 'AUTH_TOKEN_TTL_SECONDS', '2592000', 'AUTH_TOKEN', 'Token 过期时间（秒），默认 30 天', 1, 'INT', 'Token 过期时间（秒），默认 30 天', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7254222471397278270', 'MAIL_LOCAL_FROM_NAME', 'hei-fastapi', 'MAIL', '发件人显示名称', 15, 'STRING', '发件人显示名称', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7577237128270583039', 'SMS_ALIYUN_ACCESS_KEY_SECRET', '', 'SMS', '阿里云短信 AccessKeySecret', 11, 'STRING', '阿里云短信 AccessKeySecret', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7755198709300048806', 'SMS_TENCENT_SIGN_NAME', '', 'SMS', '腾讯云短信签名', 23, 'STRING', '腾讯云短信签名', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7403276917562314417', 'SMS_TENCENT_REGION', '', 'SMS', '腾讯云短信区域', 90, 'STRING', '腾讯云短信区域', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7111122815413902706', 'MAIL_TENCENT_FROM_EMAIL', '', 'MAIL', '腾讯云发件邮箱', 32, 'STRING', '腾讯云发件邮箱', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7975896582003472913', 'MAIL_TENCENT_REGION', '', 'MAIL', '腾讯云邮件区域', 90, 'STRING', '腾讯云邮件区域', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7509581770976973374', 'SMS_TENCENT_SDK_APP_ID', '', 'SMS', '腾讯云短信 SdkAppId', 22, 'STRING', '腾讯云短信 SdkAppId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7317566250917157196', 'MAIL_TEMPLATE_CHANGE_PASSWORD_CODE', '{"subject": "{{app_name}} 修改密码验证码", "body": "验证码 {{code}}，{{expire_minutes}} 分钟内有效。"}', 'MAIL_TEMPLATE', '修改密码邮件模板', 3, 'JSON', '修改密码邮件模板', NULL, 'CHANGE_PASSWORD_CODE', 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7214189881608199926', 'AUDIT_ALERT_RULE_BRUTE_FORCE', 'TRUE', 'AUDIT_ALERT', '暴力破解检测', 10, 'tinyint(1)', '暴力破解检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7316084524318600148', 'MAIL_TEMPLATE_ACCOUNT_CANCELLED', '{"subject": "{{app_name}} 账号注销确认", "body": "您好，您的账号已申请注销。\n\n我们将在 {{retention_days}} 天内保留账号数据；到期且期间未再登录使用后，系统将彻底删除账号及相关数据。\n\n预计清理时间：{{purge_at}}\n如非本人操作，请尽快联系管理员。"}', 'MAIL_TEMPLATE', '账号注销确认邮件模板', 20, 'JSON', '账号注销确认邮件模板', NULL, 'ACCOUNT_CANCELLED', 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7440617985015094217', 'MAIL_TEMPLATE_ACCOUNT_PURGED', '{"subject": "{{app_name}} 账号已彻底删除", "body": "您好，您此前注销的账号已完成保留期清理，账号及相关个人数据已彻底删除。\n\n清理时间：{{purged_at}}\n感谢您曾使用 {{app_name}}。"}', 'MAIL_TEMPLATE', '账号彻底删除邮件模板', 21, 'JSON', '账号彻底删除邮件模板', NULL, 'ACCOUNT_PURGED', 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7932451368437798893', 'SMS_TEMPLATE_ACCOUNT_CANCELLED', '{"code": "", "content": "账号已申请注销，将于{{retention_days}}天后彻底删除。"}', 'SMS_TEMPLATE', '账号注销确认短信模板', 20, 'JSON', '账号注销确认短信模板', NULL, 'ACCOUNT_CANCELLED', 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7490352982435195735', 'AUTH_REGISTER_PORTAL_DEFAULT_ROLE_ID', '', 'AUTH_REGISTER', 'PORTAL 注册默认角色', 9, 'STRING', 'PORTAL 注册默认角色', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7924936151245798814', 'AUTH_LOGIN_ADMIN_EMAIL_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'ADMIN 邮箱无用户策略', 16, 'STRING', 'ADMIN 邮箱无用户策略', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7725294497454302580', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_CLIENT_SECRET', '', 'AUTH_OAUTH', '门户微信开放平台 Secret', 33, 'STRING', 'AppSecret', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.458015', NULL, '2026-08-12 15:57:59.458015', NULL);
INSERT INTO `sys_config` VALUES ('7805678264902064100', 'AUTH_OAUTH_PORTAL_QQ_REDIRECT_URI', '', 'AUTH_OAUTH', '门户 QQ 回调', 24, 'STRING', 'Redirect URI', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.447608', NULL, '2026-08-12 15:57:59.447608', '1');
INSERT INTO `sys_config` VALUES ('7402885978532192492', 'SMS_TEMPLATE_ACCOUNT_PURGED', '{"code": "", "content": "您的账号已完成注销清理并彻底删除。"}', 'SMS_TEMPLATE', '账号彻底删除短信模板', 21, 'JSON', '账号彻底删除短信模板', NULL, 'ACCOUNT_PURGED', 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7141145699218587845', 'STORAGE_UPLOAD_MAX_BYTES', '10485760', 'UPLOAD', '上传文件大小上限（字节）', 1, 'INT', '上传文件大小上限（字节）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7634546496447633922', 'STORAGE_PRESIGN_EXPIRE_SECONDS', '3600', 'UPLOAD', '预签名 URL 有效期（秒）', 3, 'INT', '预签名 URL 有效期（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7424366524857713971', 'DEFAULT_FILE_ENGINE', 'RUSTFS', 'STORAGE', '默认文件引擎', 1, 'STRING', '默认文件引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 17:19:31.423597', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999168', 'STORAGE_ALIYUN_ENDPOINT', 'oss-cn-hangzhou.aliyuncs.com', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999169', 'STORAGE_ALIYUN_BUCKET', 'defaultbucket', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999170', 'STORAGE_ALIYUN_REGION', 'cn-hangzhou', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999171', 'STORAGE_ALIYUN_USE_SSL', 'TRUE', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193472', 'STORAGE_ALIYUN_BASE_URL', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7264846649644584871', 'STORAGE_ALIYUN_BUCKET_PUBLIC', 'FALSE', 'STORAGE', '阿里云桶是否公开', 14, 'tinyint(1)', NULL, NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193474', 'STORAGE_TENCENT_ENDPOINT', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193475', 'STORAGE_TENCENT_BUCKET', 'defaultbucket', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193476', 'STORAGE_TENCENT_REGION', 'ap-beijing', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193477', 'STORAGE_TENCENT_USE_SSL', 'TRUE', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193478', 'STORAGE_TENCENT_BASE_URL', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7362777511165276641', 'STORAGE_TENCENT_BUCKET_PUBLIC', 'FALSE', 'STORAGE', '腾讯云桶是否公开', 20, 'tinyint(1)', NULL, NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193482', 'STORAGE_MINIO_ENDPOINT', 'http://127.0.0.1:9000', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:50.574307', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193483', 'STORAGE_MINIO_BUCKET', 'vms', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193484', 'STORAGE_MINIO_REGION', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193485', 'STORAGE_MINIO_USE_SSL', 'FALSE', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193486', 'STORAGE_MINIO_BASE_URL', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7525309778220488671', 'STORAGE_MINIO_BUCKET_PUBLIC', 'FALSE', 'STORAGE', 'MinIO 桶是否公开', 26, 'tinyint(1)', NULL, NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7247260109501730845', 'STORAGE_RUSTFS_ACCESS_KEY', 'gAAAAABqeFtjtOr0lzjyRl8TPLDP5be0LWRCOizoe7P6RlqHaxasJvILz2P27NauLfjSKM71tWtwpBPZiltAP2aT5Zi5Jzr1lA', 'STORAGE', 'RustFS Access Key', 42, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 10:49:34.896059', '1');
INSERT INTO `sys_config` VALUES ('7159166510497302996', 'STORAGE_UPLOAD_CATEGORY_MAX_LENGTH', '64', 'UPLOAD', '上传分类名最大长度', 7, 'INT', '上传分类名最大长度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7236608062114386965', 'STORAGE_UPLOAD_ALLOWED_CONTENT_TYPES', '["image/jpeg","image/png","image/webp","application/pdf","text/plain","application/octet-stream"]', 'UPLOAD', '允许的 MIME 类型列表（JSON 数组）', 4, 'JSON', '允许的 MIME 类型列表（JSON 数组）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7280596390016858524', 'STORAGE_UPLOAD_ALLOWED_EXTENSIONS', '[".jpg",".jpeg",".png",".webp",".pdf",".txt",".ini",".xlsx"]', 'UPLOAD', '允许的文件扩展名列表（JSON 数组）', 5, 'JSON', '允许的文件扩展名列表（JSON 数组）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7325262355928690941', 'STORAGE_UPLOAD_DENIED_EXTENSIONS', '[".exe",".bat",".cmd",".sh",".js",".html",".php",".py",".jar"]', 'UPLOAD', '禁止上传的扩展名列表（JSON 数组）', 6, 'JSON', '禁止上传的扩展名列表（JSON 数组）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7952653754833581211', 'STORAGE_RUSTFS_SECRET_KEY', 'gAAAAABqeFtjH6R_u459TAGXEzeQOa9GRHtG_8xkDrnY8hwiZYyxeRtHh0YqhXsuzfkY3_Nn3OWId3rStJkCQgQUdwfkLtuAyA', 'STORAGE', 'RustFS Secret Key', 43, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 10:49:34.902513', '1');
INSERT INTO `sys_config` VALUES ('7950285006814595256', 'STORAGE_RUSTFS_ENDPOINT', 'http://127.0.0.1:9002', 'STORAGE', 'RustFS S3 API 端点', 41, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-08 17:17:14.441312', '1');
INSERT INTO `sys_config` VALUES ('7334700335238107691', 'STORAGE_RUSTFS_BUCKET', 'defaultbucket', 'STORAGE', 'RustFS 存储桶', 40, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7175913351351122695', 'STORAGE_RUSTFS_REGION', 'us-east-1', 'STORAGE', 'RustFS Region', 44, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7618985031080401037', 'STORAGE_RUSTFS_USE_SSL', 'FALSE', 'STORAGE', 'RustFS 是否 SSL', 45, 'tinyint(1)', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7748406326636979011', 'STORAGE_RUSTFS_BASE_URL', '', 'STORAGE', 'RustFS 自定义基础 URL', 46, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7507807560036605420', 'STORAGE_RUSTFS_BUCKET_PUBLIC', 'FALSE', 'STORAGE', 'RustFS 桶是否公开', 47, 'tinyint(1)', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7388788716511880208', 'AUDIT_ALERT_NOTIFY_EMAIL_TO', '', 'AUDIT_ALERT', '审计告警收件邮箱', 2, 'STRING', '告警收件邮箱', NULL, NULL, 0, '{}', '2026-08-12 14:10:48.638877', NULL, '2026-08-12 14:10:48.638877', NULL);
INSERT INTO `sys_config` VALUES ('7172687233882333188', 'AUDIT_ALERT_RULE_SENSITIVE_OPS', 'TRUE', 'AUDIT_ALERT', '敏感操作监控', 12, 'tinyint(1)', '敏感操作监控', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7260823451582098683', 'AUDIT_ALERT_RULE_UNUSUAL_HOURS', 'TRUE', 'AUDIT_ALERT', '异常时间操作检测', 11, 'tinyint(1)', '异常时间操作检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7879475119693769672', 'AUDIT_ALERT_BULK_DELETE_THRESHOLD', '20', 'AUDIT_ALERT', '批量删除阈值', 21, 'INT', '批量删除阈值', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_config` VALUES ('7957493498148551921', 'AUTH_LOGIN_ADMIN_LOCK_SECONDS', '300', 'AUTH_LOGIN', 'ADMIN 锁定时间（秒）', 12, 'INT', 'ADMIN 锁定时间（秒）', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7641214881822118472', 'AUTH_LOGIN_ADMIN_PHONE_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'ADMIN 手机号无用户策略', 14, 'STRING', 'ADMIN 手机号无用户策略', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7154674032244377737', 'AUTH_LOGIN_ADMIN_ALLOW_EMAIL', 'TRUE', 'AUTH_LOGIN', 'ADMIN 允许邮箱登录', 15, 'tinyint(1)', 'ADMIN 允许邮箱登录', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7189190642894026372', 'AUTH_LOGIN_ADMIN_ALLOW_OTP', 'TRUE', 'AUTH_LOGIN', 'ADMIN 允许 OTP 登录', 17, 'tinyint(1)', 'ADMIN 允许 OTP 登录', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7286296755297320815', 'AUTH_LOGIN_PORTAL_FAILURE_WINDOW_SECONDS', '300', 'AUTH_LOGIN', 'PORTAL 登录失败窗口（秒）', 18, 'INT', 'PORTAL 登录失败窗口（秒）', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7120119863717980260', 'AUTH_LOGIN_PORTAL_LOCK_SECONDS', '300', 'AUTH_LOGIN', 'PORTAL 锁定时间（秒）', 20, 'INT', 'PORTAL 锁定时间（秒）', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7961921801574207800', 'AUTH_LOGIN_PORTAL_ALLOW_PHONE', 'TRUE', 'AUTH_LOGIN', 'PORTAL 允许手机号登录', 21, 'tinyint(1)', 'PORTAL 允许手机号登录', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7553572003476288661', 'AUTH_LOGIN_PORTAL_PHONE_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'PORTAL 手机号无用户策略', 22, 'STRING', 'PORTAL 手机号无用户策略', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7643758115191934361', 'AUTH_LOGIN_PORTAL_ALLOW_EMAIL', 'TRUE', 'AUTH_LOGIN', 'PORTAL 允许邮箱登录', 23, 'tinyint(1)', 'PORTAL 允许邮箱登录', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7152439802015666741', 'AUTH_LOGIN_PORTAL_EMAIL_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'PORTAL 邮箱无用户策略', 24, 'STRING', 'PORTAL 邮箱无用户策略', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7219205172797610985', 'AUTH_LOGIN_PORTAL_ALLOW_OTP', 'TRUE', 'AUTH_LOGIN', 'PORTAL 允许 OTP 登录', 25, 'tinyint(1)', 'PORTAL 允许 OTP 登录', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7528718642174231668', 'AUTH_LOGIN_FAILURE_WINDOW_SECONDS', '900', 'AUTH_LOGIN', '登录失败统计窗口（秒）', 1, 'INT', '登录失败统计窗口（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7140248531257076308', 'AUTH_LOGIN_ACCOUNT_MAX_FAILURES', '5', 'AUTH_LOGIN', '单账号最大登录失败次数', 2, 'INT', '单账号最大登录失败次数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7224651178645144555', 'AUTH_LOGIN_LOCK_SECONDS', '900', 'AUTH_LOGIN', '登录锁定时间（秒）', 4, 'INT', '登录锁定时间（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', '1');
INSERT INTO `sys_config` VALUES ('7585813464464377276', 'AUTH_FORCE_BIND_ADMIN_EMAIL', 'FALSE', 'AUTH_FORCE_BIND', 'ADMIN 强制绑定邮箱', 3, 'tinyint(1)', '强制绑定邮箱', 'ADMIN', NULL, 0, '{}', '2026-08-12 14:36:35.252523', NULL, '2026-08-12 14:36:35.252523', NULL);
INSERT INTO `sys_config` VALUES ('7137737529370775030', 'AUTH_FORCE_BIND_ADMIN_PHONE', 'FALSE', 'AUTH_FORCE_BIND', 'ADMIN 强制绑定手机', 4, 'tinyint(1)', '强制绑定手机', 'ADMIN', NULL, 0, '{}', '2026-08-12 14:36:35.255718', NULL, '2026-08-12 14:36:35.255718', NULL);
INSERT INTO `sys_config` VALUES ('7926454604208823829', 'MAIL_TEMPLATE_BIND_EMAIL_CODE', '{"subject": "{{app_name}} 绑定邮箱验证码", "body": "您的绑定验证码是 {{code}}，{{expire_minutes}} 分钟内有效。"}', 'MAIL_TEMPLATE', '绑定邮箱验证码邮件模板', 20, 'JSON', '绑定邮箱验证码邮件模板', NULL, 'BIND_EMAIL_CODE', 0, '{}', '2026-08-12 14:36:35.270099', NULL, '2026-08-12 14:36:35.270099', NULL);
INSERT INTO `sys_config` VALUES ('7191040278958713174', 'SMS_TEMPLATE_BIND_PHONE_CODE', '{"code": "", "content": "绑定验证码 {{code}}，{{expire_minutes}} 分钟内有效"}', 'SMS_TEMPLATE', '绑定手机验证码短信模板', 20, 'JSON', '绑定手机验证码短信模板', NULL, 'BIND_PHONE_CODE', 0, '{}', '2026-08-12 14:36:35.274517', NULL, '2026-08-12 14:36:35.274517', NULL);
INSERT INTO `sys_config` VALUES ('7231113269178162861', 'AUTH_REGISTER_PORTAL_ALLOW_ACCOUNT', 'TRUE', 'AUTH_REGISTER', 'PORTAL 允许用户名注册', 11, 'tinyint(1)', '允许账号注册', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.228596', NULL, '2026-08-12 14:36:35.228596', '1');
INSERT INTO `sys_config` VALUES ('7285230910149567688', 'AUTH_REGISTER_PORTAL_ALLOW_EMAIL', 'TRUE', 'AUTH_REGISTER', 'PORTAL 允许邮箱注册', 12, 'tinyint(1)', '允许邮箱注册', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.234457', NULL, '2026-08-12 14:36:35.234457', '1');
INSERT INTO `sys_config` VALUES ('7800328101940188056', 'AUTH_OAUTH_PORTAL_GITHUB_CLIENT_ID', 'superadmin', 'AUTH_OAUTH', '门户 GitHub ClientId', 2, 'STRING', 'Client ID', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.403535', NULL, '2026-08-12 15:57:59.403535', '1');
INSERT INTO `sys_config` VALUES ('7925472518021370098', 'AUTH_OAUTH_PORTAL_GITHUB_CLIENT_SECRET', 'gAAAAABqfJrd_yTDy0DkEmej80F2frWaLzP6NMbZMBDyOzUwUOTnjRHn7UGI2ACEkt4EzA9zS9q1dK4T0B0yRAP_LDtVMcNs3w', 'AUTH_OAUTH', '门户 GitHub ClientSecret', 3, 'STRING', 'Client Secret', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.407874', NULL, '2026-08-12 15:57:59.407874', '1');
INSERT INTO `sys_config` VALUES ('7407734083761585458', 'AUTH_OAUTH_PORTAL_GITEE_ENABLED', 'TRUE', 'AUTH_OAUTH', '门户 Gitee 登录', 11, 'tinyint(1)', '启用', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.416643', NULL, '2026-08-12 15:57:59.416643', '1');
INSERT INTO `sys_config` VALUES ('7712189083486710611', 'AUTH_OAUTH_PORTAL_WECHAT_MP_APP_SECRET', '', 'AUTH_OAUTH', '门户小程序 AppSecret', 43, 'STRING', 'AppSecret', 'PORTAL', 'WECHAT_MP', 0, '{}', '2026-08-12 15:57:59.472798', NULL, '2026-08-12 15:57:59.472798', NULL);
INSERT INTO `sys_config` VALUES ('7163571438509411369', 'AUTH_OAUTH_ADMIN_GITEE_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端 Gitee ClientSecret', 113, 'STRING', 'Client Secret', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.499737', NULL, '2026-08-12 15:57:59.499737', NULL);
INSERT INTO `sys_config` VALUES ('7905539512046799183', 'AUTH_OAUTH_ADMIN_QQ_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端 QQ ClientSecret', 123, 'STRING', 'Client Secret', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.515748', NULL, '2026-08-12 15:57:59.515748', NULL);
INSERT INTO `sys_config` VALUES ('7374162031239208376', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端微信开放平台 Secret', 133, 'STRING', 'AppSecret', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.5331', NULL, '2026-08-12 15:57:59.5331', NULL);
INSERT INTO `sys_config` VALUES ('7115927954405030347', 'AUTH_OAUTH_ADMIN_GITHUB_ENABLED', 'FALSE', 'AUTH_OAUTH', '管理端 GitHub 登录', 101, 'tinyint(1)', '启用', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.476901', NULL, '2026-08-12 15:57:59.476901', '1');
INSERT INTO `sys_config` VALUES ('7192344291042963470', 'AUTH_OAUTH_ADMIN_GITHUB_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端 GitHub ClientSecret', 103, 'STRING', 'Client Secret', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.483943', NULL, '2026-08-12 15:57:59.483943', NULL);
INSERT INTO `sys_config` VALUES ('7203040505764618021', 'AUTH_OAUTH_ADMIN_GITHUB_CLIENT_ID', '', 'AUTH_OAUTH', '管理端 GitHub ClientId', 102, 'STRING', 'Client ID', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.480274', NULL, '2026-08-12 15:57:59.480274', '1');
INSERT INTO `sys_config` VALUES ('7125526233449487784', 'AUTH_OAUTH_ADMIN_GITHUB_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端 GitHub 回调', 104, 'STRING', 'Redirect URI', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.487901', NULL, '2026-08-12 15:57:59.487901', '1');
INSERT INTO `sys_config` VALUES ('7196439739936279174', 'AUTH_REGISTER_PORTAL_ALLOW_PHONE', 'TRUE', 'AUTH_REGISTER', 'PORTAL 允许手机注册', 13, 'tinyint(1)', '允许手机注册', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.237264', NULL, '2026-08-12 14:36:35.237264', '1');
INSERT INTO `sys_config` VALUES ('7147573270237877567', 'AUTH_FORCE_BIND_PORTAL_EMAIL', 'FALSE', 'AUTH_FORCE_BIND', 'PORTAL 强制绑定邮箱', 1, 'tinyint(1)', '强制绑定邮箱', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.239972', NULL, '2026-08-12 14:36:35.239972', '1');
INSERT INTO `sys_config` VALUES ('7317628235066225421', 'AUTH_FORCE_BIND_PORTAL_PHONE', 'FALSE', 'AUTH_FORCE_BIND', 'PORTAL 强制绑定手机', 2, 'tinyint(1)', '强制绑定手机', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.245505', NULL, '2026-08-12 14:36:35.245505', '1');
INSERT INTO `sys_config` VALUES ('7179891447845803919', 'AUTH_OAUTH_ADMIN_GITEE_ENABLED', 'TRUE', 'AUTH_OAUTH', '管理端 Gitee 登录', 111, 'tinyint(1)', '启用', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.491641', NULL, '2026-08-12 15:57:59.491641', '1');
INSERT INTO `sys_config` VALUES ('7820162979363132264', 'AUTH_OAUTH_ADMIN_GITEE_CLIENT_ID', '', 'AUTH_OAUTH', '管理端 Gitee ClientId', 112, 'STRING', 'Client ID', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.495992', NULL, '2026-08-12 15:57:59.495992', '1');
INSERT INTO `sys_config` VALUES ('7277178859125997111', 'AUTH_OAUTH_ADMIN_GITEE_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端 Gitee 回调', 114, 'STRING', 'Redirect URI', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.503669', NULL, '2026-08-12 15:57:59.503669', '1');
INSERT INTO `sys_config` VALUES ('7983845910800985649', 'AUTH_OAUTH_ADMIN_QQ_ENABLED', 'FALSE', 'AUTH_OAUTH', '管理端 QQ 登录', 121, 'tinyint(1)', '启用', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.507724', NULL, '2026-08-12 15:57:59.507724', '1');
INSERT INTO `sys_config` VALUES ('7953467271865167705', 'AUTH_OAUTH_ADMIN_QQ_CLIENT_ID', '', 'AUTH_OAUTH', '管理端 QQ ClientId', 122, 'STRING', 'Client ID', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.512033', NULL, '2026-08-12 15:57:59.512033', '1');
INSERT INTO `sys_config` VALUES ('7997330304022587335', 'AUTH_OAUTH_ADMIN_QQ_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端 QQ 回调', 124, 'STRING', 'Redirect URI', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.520333', NULL, '2026-08-12 15:57:59.520333', '1');
INSERT INTO `sys_config` VALUES ('7566333830794456163', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_ENABLED', 'FALSE', 'AUTH_OAUTH', '管理端微信网页登录', 131, 'tinyint(1)', '启用', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.525073', NULL, '2026-08-12 15:57:59.525073', '1');
INSERT INTO `sys_config` VALUES ('7709767367597654277', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_CLIENT_ID', '', 'AUTH_OAUTH', '管理端微信开放平台 AppId', 132, 'STRING', 'AppId', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.529053', NULL, '2026-08-12 15:57:59.529053', '1');
INSERT INTO `sys_config` VALUES ('7878128811204573416', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端微信开放平台回调', 134, 'STRING', 'Redirect URI', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.53685', NULL, '2026-08-12 15:57:59.53685', '1');
INSERT INTO `sys_config` VALUES ('7319284294792948722', 'AUTH_OAUTH_PORTAL_GITHUB_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户 GitHub 登录', 1, 'tinyint(1)', '启用', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.397127', NULL, '2026-08-12 15:57:59.397127', '1');
INSERT INTO `sys_config` VALUES ('7964584885414219711', 'AUTH_OAUTH_PORTAL_GITHUB_REDIRECT_URI', '', 'AUTH_OAUTH', '门户 GitHub 回调', 4, 'STRING', 'Redirect URI', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.412339', NULL, '2026-08-12 15:57:59.412339', '1');
INSERT INTO `sys_config` VALUES ('7346432777171636629', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户微信网页登录', 31, 'tinyint(1)', '启用', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.450854', NULL, '2026-08-12 15:57:59.450854', '1');
INSERT INTO `sys_config` VALUES ('7112014435091478940', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_CLIENT_ID', '', 'AUTH_OAUTH', '门户微信开放平台 AppId', 32, 'STRING', 'AppId', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.454644', NULL, '2026-08-12 15:57:59.454644', '1');
INSERT INTO `sys_config` VALUES ('7535593666554670318', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_REDIRECT_URI', '', 'AUTH_OAUTH', '门户微信开放平台回调', 34, 'STRING', 'Redirect URI', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.461291', NULL, '2026-08-12 15:57:59.461291', '1');
INSERT INTO `sys_config` VALUES ('7408452631572683423', 'AUTH_OAUTH_PORTAL_WECHAT_MP_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户微信小程序登录', 41, 'tinyint(1)', '启用', 'PORTAL', 'WECHAT_MP', 0, '{}', '2026-08-12 15:57:59.464887', NULL, '2026-08-12 15:57:59.464887', '1');
INSERT INTO `sys_config` VALUES ('7143068116163573568', 'AUTH_OAUTH_PORTAL_WECHAT_MP_APP_ID', '', 'AUTH_OAUTH', '门户小程序 AppId', 42, 'STRING', 'AppId', 'PORTAL', 'WECHAT_MP', 0, '{}', '2026-08-12 15:57:59.468946', NULL, '2026-08-12 15:57:59.468946', '1');
INSERT INTO `sys_config` VALUES ('7499926412524933665', 'AUTH_OAUTH_FRONTEND_CALLBACK_PORTAL', 'http://localhost:5174/auth/oauth/callback', 'AUTH_OAUTH', '门户 OAuth 前端回调页（空则用默认）', 200, 'STRING', '门户前端回调', NULL, NULL, 0, '{}', '2026-08-12 15:57:59.540869', NULL, '2026-08-16 10:27:33.694774', '1');
INSERT INTO `sys_config` VALUES ('7814530538364155449', 'AUTH_OAUTH_FRONTEND_CALLBACK_ADMIN', 'http://localhost:5173/auth/oauth/callback', 'AUTH_OAUTH', '管理端 OAuth 前端回调页（空则用默认）', 201, 'STRING', '管理端前端回调', NULL, NULL, 0, '{}', '2026-08-12 15:57:59.544738', NULL, '2026-08-16 10:27:33.715381', '1');

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept` (
 `id` varchar(64) NOT NULL,
 `parent_id` varchar(64),
 `master_id` varchar(64),
 `deputy_master_id` varchar(64),
 `name` varchar(64) NOT NULL,
 `category` varchar(64) NOT NULL,
 `sort` int NOT NULL,
 `is_virtual` tinyint(1) NOT NULL,
 `status` varchar(32) NOT NULL,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_dept
-- ----------------------------

-- ----------------------------
-- Table structure for sys_dict
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict`;
CREATE TABLE `sys_dict` (
 `id` varchar(32) NOT NULL,
 `code` varchar(50) NOT NULL,
 `label` varchar(255),
 `value` varchar(255),
 `color` varchar(32),
 `category` varchar(64),
 `parent_id` varchar(32),
 `status` varchar(16) NOT NULL,
 `sort` int NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_dict
-- ----------------------------
INSERT INTO `sys_dict` VALUES ('100001', 'COMMON_STATUS', '状态', 'COMMON_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100002', 'COMMON_STATUS_ENABLED', '启用', 'ENABLED', '#18a058', 'SYS', '100001', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100003', 'COMMON_STATUS_DISABLED', '禁用', 'DISABLED', '#d03050', 'SYS', '100001', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100004', 'SYS_BIZ_CATEGORY', '系统/业务分类', 'SYS_BIZ_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100005', 'SYS_BIZ_CATEGORY_SYS', '系统', 'SYS', '#2080f0', 'SYS', '100004', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100006', 'SYS_BIZ_CATEGORY_BIZ', '业务', 'BIZ', '#f0a020', 'SYS', '100004', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100010', 'ACCOUNT_STATUS', '账号状态', 'ACCOUNT_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100011', 'ACCOUNT_STATUS_ENABLED', '启用', 'ENABLED', '#18a058', 'SYS', '100010', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100012', 'ACCOUNT_STATUS_DISABLED', '禁用', 'DISABLED', '#d03050', 'SYS', '100010', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100013', 'ACCOUNT_STATUS_CANCELLED', '已注销', 'CANCELLED', '#909399', 'SYS', '100010', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100014', 'ROLE_SCOPE_TYPE', '角色范围类型', 'ROLE_SCOPE_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100015', 'ROLE_SCOPE_TYPE_PLATFORM', '平台', 'PLATFORM', '#2080f0', 'SYS', '100014', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100016', 'ROLE_SCOPE_TYPE_DEPT', '部门', 'DEPT', '#18a058', 'SYS', '100014', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100017', 'RESOURCE_TYPE', '资源类型', 'RESOURCE_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100018', 'RESOURCE_TYPE_CATALOG', '目录', 'CATALOG', '#722ed1', 'SYS', '100017', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100019', 'RESOURCE_TYPE_MENU', '菜单', 'MENU', '#2080f0', 'SYS', '100017', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100020', 'RESOURCE_TYPE_PAGE', '页面', 'PAGE', '#18a058', 'SYS', '100017', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100021', 'RESOURCE_TYPE_BUTTON', '按钮', 'BUTTON', '#f0a020', 'SYS', '100017', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100022', 'RESOURCE_TYPE_ACTION', '操作', 'ACTION', '#d03050', 'SYS', '100017', 'ENABLED', 5, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100023', 'RESOURCE_TYPE_API_GROUP', '接口组', 'API_GROUP', '#1677ff', 'SYS', '100017', 'ENABLED', 6, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100024', 'DATA_SCOPE', '数据范围', 'DATA_SCOPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100025', 'DATA_SCOPE_ALL', '全部', 'ALL', '#18a058', 'SYS', '100024', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100026', 'DATA_SCOPE_DEPT_AND_CHILD', '本部门及子部门', 'DEPT_AND_CHILD', '#2080f0', 'SYS', '100024', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100027', 'DATA_SCOPE_DEPT', '本部门', 'DEPT', '#2db7f5', 'SYS', '100024', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100028', 'DATA_SCOPE_SELF', '本人', 'SELF', '#f0a020', 'SYS', '100024', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100029', 'DATA_SCOPE_CUSTOM', '自定义部门', 'CUSTOM', '#722ed1', 'SYS', '100024', 'ENABLED', 5, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100040', 'DEPT_CATEGORY', '部门分类', 'DEPT_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100041', 'DEPT_CATEGORY_COMPANY', '公司', 'COMPANY', '#2080f0', 'SYS', '100040', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100042', 'DEPT_CATEGORY_DEPARTMENT', '部门', 'DEPARTMENT', '#18a058', 'SYS', '100040', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100043', 'DEPT_CATEGORY_TEAM', '团队', 'TEAM', '#f0a020', 'SYS', '100040', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100044', 'DEPT_CATEGORY_VIRTUAL', '虚拟组织', 'VIRTUAL', '#909399', 'SYS', '100040', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100045', 'POSITION_CATEGORY', '岗位分类', 'POSITION_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100046', 'POSITION_CATEGORY_MANAGEMENT', '管理', 'MANAGEMENT', '#2080f0', 'SYS', '100045', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100047', 'POSITION_CATEGORY_TECHNICAL', '技术', 'TECHNICAL', '#18a058', 'SYS', '100045', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100048', 'POSITION_CATEGORY_OPERATION', '运营', 'OPERATION', '#f0a020', 'SYS', '100045', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100049', 'POSITION_CATEGORY_SUPPORT', '支持', 'SUPPORT', '#909399', 'SYS', '100045', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100054', 'BANNER_CATEGORY', '展示图分类', 'BANNER_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100055', 'BANNER_CATEGORY_HOME', '首页', 'HOME', '#18a058', 'SYS', '100054', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100056', 'BANNER_CATEGORY_LOGIN', '登录', 'LOGIN', '#2080f0', 'SYS', '100054', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100057', 'BANNER_CATEGORY_WORKPLACE', '工作台', 'WORKPLACE', '#722ed1', 'SYS', '100054', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100058', 'BANNER_CATEGORY_NOTICE', '公告', 'NOTICE', '#f0a020', 'SYS', '100054', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100059', 'BANNER_CATEGORY_ADMIN_DASHBOARD', '管理端仪表盘', 'ADMIN_DASHBOARD', '#2080f0', 'SYS', '100054', 'ENABLED', 5, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100060', 'BANNER_CATEGORY_SYSTEM_UPGRADE', '系统升级', 'SYSTEM_UPGRADE', '#d03050', 'SYS', '100054', 'ENABLED', 6, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100061', 'BANNER_TYPE', '展示图类型', 'BANNER_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100062', 'BANNER_TYPE_CAROUSEL', '轮播图', 'CAROUSEL', '#18a058', 'SYS', '100061', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100063', 'BANNER_TYPE_HERO', '主视觉', 'HERO', '#2080f0', 'SYS', '100061', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100064', 'BANNER_TYPE_NOTICE', '公告', 'NOTICE', '#f0a020', 'SYS', '100061', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100065', 'BANNER_TYPE_CARD', '卡片', 'CARD', '#722ed1', 'SYS', '100061', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100066', 'BANNER_TYPE_POPUP', '弹窗', 'POPUP', '#d03050', 'SYS', '100061', 'ENABLED', 5, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100067', 'BANNER_TYPE_SIDEBAR', '侧边栏', 'SIDEBAR', '#2080f0', 'SYS', '100061', 'ENABLED', 6, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100068', 'BANNER_POSITION', '展示图位置', 'BANNER_POSITION', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100069', 'BANNER_POSITION_HOME_TOP', '首页顶部', 'HOME_TOP', '#18a058', 'SYS', '100068', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100070', 'BANNER_POSITION_HOME_MIDDLE', '首页中部', 'HOME_MIDDLE', '#18a058', 'SYS', '100068', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100071', 'BANNER_POSITION_HOME_BOTTOM', '首页底部', 'HOME_BOTTOM', '#18a058', 'SYS', '100068', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100072', 'BANNER_POSITION_LOGIN_SIDE', '登录侧边', 'LOGIN_SIDE', '#2080f0', 'SYS', '100068', 'ENABLED', 4, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100073', 'BANNER_POSITION_WORKPLACE_TOP', '工作台顶部', 'WORKPLACE_TOP', '#722ed1', 'SYS', '100068', 'ENABLED', 5, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100074', 'BANNER_POSITION_NOTICE_AREA', '公告区域', 'NOTICE_AREA', '#f0a020', 'SYS', '100068', 'ENABLED', 6, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100075', 'BANNER_POSITION_ADMIN_TOP', '管理端顶部', 'ADMIN_TOP', '#2080f0', 'SYS', '100068', 'ENABLED', 7, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100076', 'BANNER_POSITION_ADMIN_SIDEBAR', '管理端侧边栏', 'ADMIN_SIDEBAR', '#2080f0', 'SYS', '100068', 'ENABLED', 8, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100077', 'BANNER_LINK_TYPE', '展示图链接类型', 'BANNER_LINK_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100078', 'BANNER_LINK_TYPE_URL', '外部链接', 'URL', '#18a058', 'SYS', '100077', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100079', 'BANNER_LINK_TYPE_ROUTE', '路由', 'ROUTE', '#2080f0', 'SYS', '100077', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100080', 'BANNER_LINK_TYPE_NONE', '无链接', 'NONE', '#909399', 'SYS', '100077', 'ENABLED', 3, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100085', 'ACCOUNT_IDENTITY_BIND_STATUS', '账号身份绑定状态', 'ACCOUNT_IDENTITY_BIND_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100086', 'ACCOUNT_IDENTITY_BIND_STATUS_BOUND', '已绑定', 'BOUND', '#18a058', 'SYS', '100085', 'ENABLED', 1, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100087', 'ACCOUNT_IDENTITY_BIND_STATUS_UNBOUND', '未绑定', 'UNBOUND', '#909399', 'SYS', '100085', 'ENABLED', 2, '2026-06-29 00:00:00', NULL, '2026-06-29 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100095', 'NOTIFICATION_SEVERITY', '通知严重级别', 'NOTIFICATION_SEVERITY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100096', 'NOTIFICATION_SEVERITY_INFO', '信息', 'INFO', '#2080f0', 'SYS', '100095', 'ENABLED', 1, '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100097', 'NOTIFICATION_SEVERITY_SUCCESS', '成功', 'SUCCESS', '#18a058', 'SYS', '100095', 'ENABLED', 2, '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100098', 'NOTIFICATION_SEVERITY_WARNING', '警告', 'WARNING', '#f0a020', 'SYS', '100095', 'ENABLED', 3, '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100099', 'NOTIFICATION_SEVERITY_ERROR', '错误', 'ERROR', '#d03050', 'SYS', '100095', 'ENABLED', 4, '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100126', 'NOTIFICATION_SEVERITY_URGENT', '紧急', 'URGENT', '#d03050', 'SYS', '100095', 'ENABLED', 5, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100127', 'CONTENT_TYPE', '内容格式', 'CONTENT_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100128', 'CONTENT_TYPE_TEXT', '纯文本', 'text', '#909399', 'SYS', '100127', 'ENABLED', 1, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100129', 'CONTENT_TYPE_HTML', '富文本', 'html', '#18a058', 'SYS', '100127', 'ENABLED', 2, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100130', 'CONTENT_TYPE_MARKDOWN', 'Markdown', 'markdown', '#722ed1', 'SYS', '100127', 'ENABLED', 3, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100131', 'TARGET_SCOPE', '目标范围', 'TARGET_SCOPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100132', 'TARGET_SCOPE_ALL', '全部', 'ALL', '#2080f0', 'SYS', '100131', 'ENABLED', 1, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100135', 'TARGET_SCOPE_SPECIFIC', '指定用户', 'SPECIFIC', '#d03050', 'SYS', '100131', 'ENABLED', 3, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100136', 'NOTIFY_LOCATION', '通知位置', 'NOTIFY_LOCATION', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100137', 'NOTIFY_LOCATION_CENTER', '通知中心', 'center', '#2080f0', 'SYS', '100136', 'ENABLED', 1, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100138', 'NOTIFY_LOCATION_POPUP', '弹窗', 'popup', '#f0a020', 'SYS', '100136', 'ENABLED', 2, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100139', 'NOTIFY_LOCATION_DASHBOARD', '工作台公告区', 'dashboard', '#722ed1', 'SYS', '100136', 'ENABLED', 3, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100140', 'FEEDBACK_CATEGORY', '反馈分类', 'FEEDBACK_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100141', 'FEEDBACK_CATEGORY_SUGGESTION', '功能建议', 'SUGGESTION', '#18a058', 'SYS', '100140', 'ENABLED', 1, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100142', 'FEEDBACK_CATEGORY_BUG', '问题反馈', 'BUG', '#d03050', 'SYS', '100140', 'ENABLED', 2, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100143', 'FEEDBACK_CATEGORY_OTHER', '其他', 'OTHER', '#909399', 'SYS', '100140', 'ENABLED', 3, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100144', 'FEEDBACK_STATUS', '反馈状态', 'FEEDBACK_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100145', 'FEEDBACK_STATUS_PENDING', '待处理', 'PENDING', '#f0a020', 'SYS', '100144', 'ENABLED', 1, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100146', 'FEEDBACK_STATUS_REVIEWED', '已查看', 'REVIEWED', '#2080f0', 'SYS', '100144', 'ENABLED', 2, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100147', 'FEEDBACK_STATUS_RESOLVED', '已解决', 'RESOLVED', '#18a058', 'SYS', '100144', 'ENABLED', 3, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100148', 'FEEDBACK_STATUS_CLOSED', '已关闭', 'CLOSED', '#909399', 'SYS', '100144', 'ENABLED', 4, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100149', 'PUBLISH_STATUS', '发布状态', 'PUBLISH_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100150', 'PUBLISH_STATUS_DRAFT', '草稿', 'DRAFT', '#909399', 'SYS', '100149', 'ENABLED', 1, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100151', 'PUBLISH_STATUS_PUBLISHED', '已发布', 'PUBLISHED', '#18a058', 'SYS', '100149', 'ENABLED', 2, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100152', 'PUBLISH_STATUS_REVOKED', '已撤回', 'REVOKED', '#d03050', 'SYS', '100149', 'ENABLED', 3, '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100206', 'TARGET_SCOPE_ACCOUNT_TYPE', '按账号类型', 'ACCOUNT_TYPE', '#722ed1', 'SYS', '100131', 'ENABLED', 2, '2026-08-08 04:14:19.198462', NULL, '2026-08-08 04:14:19.198462', NULL);
INSERT INTO `sys_dict` VALUES ('100210', 'NOTIFICATION_CATEGORY', '通知分类', 'NOTIFICATION_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100211', 'NOTIFICATION_CATEGORY_ORDER', '订单', 'ORDER', '#2080f0', 'SYS', '100210', 'ENABLED', 1, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100212', 'NOTIFICATION_CATEGORY_APPROVAL', '审批', 'APPROVAL', '#722ed1', 'SYS', '100210', 'ENABLED', 2, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100213', 'NOTIFICATION_CATEGORY_SYSTEM', '系统', 'SYSTEM', '#18a058', 'SYS', '100210', 'ENABLED', 3, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100214', 'NOTIFICATION_CATEGORY_SECURITY', '安全', 'SECURITY', '#d03050', 'SYS', '100210', 'ENABLED', 4, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('100215', 'NOTIFICATION_CATEGORY_BIZ', '业务', 'BIZ', '#f0a020', 'SYS', '100210', 'ENABLED', 5, '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_dict` VALUES ('7656105493975688393', 'OAUTH_PROVIDER', '三方登录提供商', 'OAUTH_PROVIDER', NULL, 'SYS', NULL, 'ENABLED', 90, '2026-08-12 15:57:59.360135', NULL, '2026-08-12 15:57:59.360135', NULL);
INSERT INTO `sys_dict` VALUES ('7587538474363101234', 'GITHUB', 'GitHub', 'GITHUB', NULL, 'SYS', '7656105493975688393', 'ENABLED', 1, '2026-08-12 15:57:59.376072', NULL, '2026-08-12 15:57:59.376072', NULL);
INSERT INTO `sys_dict` VALUES ('7351536453246655198', 'GITEE', 'Gitee', 'GITEE', NULL, 'SYS', '7656105493975688393', 'ENABLED', 2, '2026-08-12 15:57:59.380323', NULL, '2026-08-12 15:57:59.380323', NULL);
INSERT INTO `sys_dict` VALUES ('7399371307733482375', 'QQ', 'QQ', 'QQ', NULL, 'SYS', '7656105493975688393', 'ENABLED', 3, '2026-08-12 15:57:59.384444', NULL, '2026-08-12 15:57:59.384444', NULL);
INSERT INTO `sys_dict` VALUES ('7211108434702344271', 'WECHAT_OPEN', '微信开放平台', 'WECHAT_OPEN', NULL, 'SYS', '7656105493975688393', 'ENABLED', 4, '2026-08-12 15:57:59.388141', NULL, '2026-08-12 15:57:59.388141', NULL);
INSERT INTO `sys_dict` VALUES ('7260287663522585806', 'WECHAT_MP', '微信小程序', 'WECHAT_MP', NULL, 'SYS', '7656105493975688393', 'ENABLED', 5, '2026-08-12 15:57:59.392548', NULL, '2026-08-12 15:57:59.392548', NULL);

-- ----------------------------
-- Table structure for sys_feedback
-- ----------------------------
DROP TABLE IF EXISTS `sys_feedback`;
CREATE TABLE `sys_feedback` (
 `id` varchar(64) NOT NULL,
 `title` varchar(255) NOT NULL,
 `content` text NOT NULL,
 `category` varchar(64) NOT NULL,
 `contact` varchar(255),
 `attach_object_names` json NOT NULL,
 `status` varchar(32) NOT NULL,
 `reply` text,
 `replied_by` varchar(64),
 `replied_at` datetime(6),
 `submitter_account_type` varchar(32) NOT NULL,
 `submitter_account_id` varchar(64) NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_feedback
-- ----------------------------
INSERT INTO `sys_feedback` VALUES ('7491849375090675712', '哈哈哈', '撒擦三次', 'SUGGESTION', '擦拭擦拭', '["uploads/2026/08/08/be3515e142974cf08b46e348d1c3d8d3.png"]', 'RESOLVED', 'ok', '1', '2026-08-08 13:40:49.757811', 'PORTAL', '7491847383584804864', '2026-08-08 13:34:42.831539', '7491847383584804864', '2026-08-08 13:40:49.742207', '1');

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS `sys_file`;
CREATE TABLE `sys_file` (
 `id` varchar(64) NOT NULL,
 `object_name` varchar(255) NOT NULL,
 `original_name` varchar(255) NOT NULL,
 `storage_provider` varchar(32) NOT NULL,
 `bucket` varchar(255),
 `content_type` varchar(128) NOT NULL,
 `size` bigint NOT NULL,
 `url` varchar(1024) NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_file
-- ----------------------------
INSERT INTO `sys_file` VALUES ('2086404723008208898', 'uploads/2026/08/09/48af08ca31e346d3a75139fa84eb0282.png', 'QR2026080700024_1786146767761.png', 'rustfs', 'defaultbucket', 'image/png', 38636, 'uploads/2026/08/09/48af08ca31e346d3a75139fa84eb0282.png', '2026-08-09 10:50:24.153462', '1', '2026-08-09 10:50:24.153462', '1');
INSERT INTO `sys_file` VALUES ('7491869364283744256', 'uploads/2026/08/08/93c4b0dba86f483bb1905ba079550521.png', 'QR2026080700024_1786146767761.png', 'minio', 'vms', 'image/png', 38636, 'uploads/2026/08/08/93c4b0dba86f483bb1905ba079550521.png', '2026-08-08 14:54:08.592047', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('7491906012023353344', 'uploads/2026/08/08/0e535c4dc69241eab526c5e94d9eb19b.png', 'QR2026080700024_1786146767761.png', 'rustfs', 'defaultbucket', 'image/png', 38636, 'uploads/2026/08/08/0e535c4dc69241eab526c5e94d9eb19b.png', '2026-08-08 17:19:46.108576', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2086408061170970625', 'uploads/2026/08/09/3dbf7f27b66f4023a0736bdccaa9596b.png', 'QR2026080700024_1786146767761.png', 'rustfs', 'defaultbucket', 'image/png', 38636, 'uploads/2026/08/09/3dbf7f27b66f4023a0736bdccaa9596b.png', '2026-08-09 11:03:40.030025', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2086410328867565570', 'uploads/2026/08/09/02acc3dee5454d34913b07f49fe59cac.png', 'avatar.png', 'rustfs', 'defaultbucket', 'image/png', 193387, 'uploads/2026/08/09/02acc3dee5454d34913b07f49fe59cac.png', '2026-08-09 11:12:40.692242', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2086415187620601857', 'uploads/2026/08/09/85e1b98acfc9465abbbba86ef3b4fec8.jpg', '120153703_touxiang_bobopic (1).jpg', 'rustfs', 'defaultbucket', 'image/jpeg', 65451, 'uploads/2026/08/09/85e1b98acfc9465abbbba86ef3b4fec8.jpg', '2026-08-09 11:31:59.117567', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2088928636523130882', 'uploads/2026/08/16/947e530849e245d3bcb873354da8f113.txt', 'PyPI-Recovery-Codes-charliebyte-2026-06-14T07_19_21.473015.txt', 'rustfs', 'defaultbucket', 'text/plain', 135, 'uploads/2026/08/16/947e530849e245d3bcb873354da8f113.txt', '2026-08-16 09:59:32.021932', '1', '2026-08-16 10:27:45.912235', '1');

-- ----------------------------
-- Table structure for sys_group
-- ----------------------------
DROP TABLE IF EXISTS `sys_group`;
CREATE TABLE `sys_group` (
 `id` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `owner_dept_id` varchar(64),
 `description` text,
 `status` varchar(32) NOT NULL,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_group
-- ----------------------------

-- ----------------------------
-- Table structure for sys_iam_relation
-- ----------------------------
DROP TABLE IF EXISTS `sys_iam_relation`;
CREATE TABLE `sys_iam_relation` (
 `id` varchar(64) NOT NULL,
 `subject_type` varchar(32) NOT NULL,
 `subject_id` varchar(64) NOT NULL,
 `account_type` varchar(32) NOT NULL,
 `relation_type` varchar(64) NOT NULL,
 `target_type` varchar(32) NOT NULL,
 `target_id` varchar(64) NOT NULL,
 `target_key` varchar(128) NOT NULL,
 `grant_mode` varchar(32) NOT NULL,
 `data_scope` varchar(32) NOT NULL,
 `custom_scope_dept_ids` json NOT NULL,
 `is_primary` tinyint(1) NOT NULL,
 `sort` int NOT NULL,
 `status` varchar(32) NOT NULL,
 `description` text,
 `reason` text,
 `expired_at` datetime(6),
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_iam_relation
-- ----------------------------
INSERT INTO `sys_iam_relation` VALUES ('1', 'ACCOUNT', '1', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '1', '', 'CASCADE', 'SELF', '[]', 0, 99, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 11:56:13.747886', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7859474578876774469', 'RESOURCE', '201061', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7109538496802851524', 'RESOURCE', '201060', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7606481288132251344', 'RESOURCE', '200029', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7232506669573029115', 'RESOURCE', '200028', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7639169641875772298', 'RESOURCE', '203011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:page', 'CASCADE', 'ALL', '[]', 0, 10, 'ENABLED', '分页活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7611870955633752502', 'RESOURCE', '203012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:create', 'CASCADE', 'ALL', '[]', 0, 20, 'ENABLED', '新增活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7654259910312149696', 'RESOURCE', '203013', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:detail', 'CASCADE', 'ALL', '[]', 0, 30, 'ENABLED', '详情活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7223792660518235132', 'RESOURCE', '203014', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:update', 'CASCADE', 'ALL', '[]', 0, 40, 'ENABLED', '编辑活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7567716937788130247', 'RESOURCE', '203015', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:delete', 'CASCADE', 'ALL', '[]', 0, 50, 'ENABLED', '删除活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7986838434768267433', 'RESOURCE', '203021', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:page', 'CASCADE', 'ALL', '[]', 0, 10, 'ENABLED', '分页目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7768912188750692632', 'RESOURCE', '203022', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:create', 'CASCADE', 'ALL', '[]', 0, 20, 'ENABLED', '新增目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7159667080064467923', 'RESOURCE', '203023', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:detail', 'CASCADE', 'ALL', '[]', 0, 30, 'ENABLED', '详情目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7660302211516474641', 'RESOURCE', '203024', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:update', 'CASCADE', 'ALL', '[]', 0, 40, 'ENABLED', '编辑目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7624773106003991812', 'RESOURCE', '203025', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:delete', 'CASCADE', 'ALL', '[]', 0, 50, 'ENABLED', '删除目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7822968206741092129', 'RESOURCE', '203026', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:list', 'CASCADE', 'ALL', '[]', 0, 90, 'ENABLED', '树列表目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7837316234393882458', 'RESOURCE', '203031', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:page', 'CASCADE', 'ALL', '[]', 0, 10, 'ENABLED', '分页订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7518444451967536602', 'RESOURCE', '203032', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:create', 'CASCADE', 'ALL', '[]', 0, 20, 'ENABLED', '新增订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7112731414735234196', 'RESOURCE', '203033', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:detail', 'CASCADE', 'ALL', '[]', 0, 30, 'ENABLED', '详情订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7304374671311844075', 'RESOURCE', '203034', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:update', 'CASCADE', 'ALL', '[]', 0, 40, 'ENABLED', '编辑订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7334894311555546200', 'RESOURCE', '203035', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:delete', 'CASCADE', 'ALL', '[]', 0, 50, 'ENABLED', '删除订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7453524161449865528', 'RESOURCE', '203041', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:page', 'CASCADE', 'ALL', '[]', 0, 10, 'ENABLED', '分页知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7601824633419714671', 'RESOURCE', '203042', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:create', 'CASCADE', 'ALL', '[]', 0, 20, 'ENABLED', '新增知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7574875650561833761', 'RESOURCE', '203043', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:detail', 'CASCADE', 'ALL', '[]', 0, 30, 'ENABLED', '详情知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7990168508290078017', 'RESOURCE', '203044', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:update', 'CASCADE', 'ALL', '[]', 0, 40, 'ENABLED', '编辑知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7845666016635732956', 'RESOURCE', '203045', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:delete', 'CASCADE', 'ALL', '[]', 0, 50, 'ENABLED', '删除知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7133323711340623180', 'RESOURCE', '203046', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:list', 'CASCADE', 'ALL', '[]', 0, 90, 'ENABLED', '树列表知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7569914124743592951', 'RESOURCE', '202201', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分页消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7575458644059959564', 'RESOURCE', '202202', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7927293661818445174', 'RESOURCE', '202203', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '详情消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7742840047784749526', 'RESOURCE', '202204', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7174771647983316441', 'RESOURCE', '202205', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7624067801880049144', 'RESOURCE', '202209', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:publish', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '发布消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7122813088692955083', 'RESOURCE', '202240', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:revoke', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '撤回消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7380222627904177407', 'RESOURCE', '202241', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:pin', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '置顶消息', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7284368855990246834', 'RESOURCE', '204001', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分页任务', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7939060249732762857', 'RESOURCE', '204011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:create', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '新增任务', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7212116468775288981', 'RESOURCE', '204012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:update', 'CASCADE', 'ALL', '[]', 0, 2, 'ENABLED', '编辑任务', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7560408972191285564', 'RESOURCE', '204013', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:delete', 'CASCADE', 'ALL', '[]', 0, 3, 'ENABLED', '删除任务', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7873407408257995473', 'RESOURCE', '204014', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:detail', 'CASCADE', 'ALL', '[]', 0, 4, 'ENABLED', '任务详情', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7792386902249912041', 'RESOURCE', '204015', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:run', 'CASCADE', 'ALL', '[]', 0, 5, 'ENABLED', '立即执行', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7740028803530587951', 'RESOURCE', '204016', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:joblog:page', 'CASCADE', 'ALL', '[]', 0, 6, 'ENABLED', '执行日志', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);

-- ----------------------------
-- Table structure for sys_job
-- ----------------------------
DROP TABLE IF EXISTS `sys_job`;
CREATE TABLE `sys_job` (
 `id` varchar(64) NOT NULL,
 `job_name` varchar(128) NOT NULL,
 `execute_class` varchar(255) NOT NULL,
 `execute_type` varchar(16) NOT NULL,
 `trigger_config` varchar(255) NOT NULL,
 `execute_param` json,
 `last_run_time` datetime(6),
 `next_run_time` datetime(6) NOT NULL,
 `last_execute_result` varchar(500),
 `enabled` tinyint(1) NOT NULL,
 `description` varchar(500),
 `sort` int NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_job
-- ----------------------------
INSERT INTO `sys_job` VALUES ('7541000000000000007', '任务执行日志清理', 'sys_job_log_cleanup', 'CRON', '0 30 3 * * *', '{"retentionDays": 30, "batchSize": 1000}', '2026-08-16 09:53:27.104743', '2026-08-16 19:30:00', 'deleted=0,retentionDays=30,batchSize=1000', 0, '按保留天数批量清理过期 sys_job_log', 7, '2026-08-16 09:53:26.903237', NULL, '2026-08-16 09:53:27.146058', NULL);
INSERT INTO `sys_job` VALUES ('7541000000000000006', '注销账号清理', 'iam_account_purge_cancelled', 'CRON', '0 0 3 * * *', '{"retentionDays": 15}', '2026-08-16 09:53:27.108262', '2026-08-16 19:00:00', 'purged=0', 0, '每日清理已取消且超过保留期的账号数据', 6, '2026-08-16 09:53:26.897926', NULL, '2026-08-16 09:53:27.152576', NULL);
INSERT INTO `sys_job` VALUES ('7541000000000000004', '审计告警', 'sys_audit_alert', 'FIXED', '300', '{}', '2026-08-16 10:25:56.244268', '2026-08-16 10:30:56.283976', 'done fired=0', 0, '按配置规则扫描审计日志并发送告警', 4, '2026-08-16 09:53:26.892521', NULL, '2026-08-16 10:25:56.289499', NULL);
INSERT INTO `sys_job` VALUES ('7541000000000000001', '示例任务', 'sys_job_sample', 'FIXED', '60', '{}', '2026-08-16 10:27:57.95913', '2026-08-16 10:28:57.960251', 'echo: (无参数)', 0, '演示调度链路：回显执行参数', 1, '2026-08-16 09:53:26.876522', NULL, '2026-08-16 10:27:57.960251', NULL);
INSERT INTO `sys_job` VALUES ('7541000000000000003', 'Banner 互动计数刷库', 'sys_banner_flush_interactions', 'FIXED', '60', '{}', '2026-08-16 10:27:57.95913', '2026-08-16 10:28:57.962143', 'flushed=0', 0, '将 Redis 互动增量写入 sys_banner.interaction_count', 3, '2026-08-16 09:53:26.887732', NULL, '2026-08-16 10:27:57.963134', NULL);
INSERT INTO `sys_job` VALUES ('7541000000000000002', 'Banner 状态同步', 'sys_banner_status_sync', 'FIXED', '60', '{}', '2026-08-16 10:27:57.95913', '2026-08-16 10:28:57.970142', 'expired=0,activated=0', 0, '按 start_at / end_at 激活或过期 Banner', 2, '2026-08-16 09:53:26.882777', NULL, '2026-08-16 10:27:57.970142', NULL);

-- ----------------------------
-- Table structure for sys_job_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_job_log`;
CREATE TABLE `sys_job_log` (
 `id` varchar(64) NOT NULL,
 `job_id` varchar(64) NOT NULL,
 `job_name` varchar(128) NOT NULL,
 `execute_param` json,
 `execute_time` datetime(6) NOT NULL,
 `execute_duration_ms` bigint,
 `success` tinyint(1) NOT NULL,
 `execute_result` text,
 `executor` varchar(64),
 `ip` varchar(64),
 `process_id` varchar(32),
 `app_dir` varchar(500),
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_job_log
-- ----------------------------
INSERT INTO `sys_job_log` VALUES ('2088927105711222785', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:53:27.057194', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:53:27.063199', NULL, '2026-08-16 09:53:27.063199', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927105841246210', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:53:27.057194', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:53:27.08054', NULL, '2026-08-16 09:53:27.08054', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927106109681666', '7541000000000000004', '审计告警', '{}', '2026-08-16 09:53:27.059189', 78, 0, 'done fired=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:53:27.143059', NULL, '2026-08-16 09:53:27.143059', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927106109681668', '7541000000000000006', '注销账号清理', '{"retentionDays":15}', '2026-08-16 09:53:27.108262', 43, 0, 'purged=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:53:27.156571', NULL, '2026-08-16 09:53:27.156571', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927106109681667', '7541000000000000007', '任务执行日志清理', '{"retentionDays":30,"batchSize":1000}', '2026-08-16 09:53:27.104743', 24, 0, 'deleted=0,retentionDays=30,batchSize=1000', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:53:27.149059', NULL, '2026-08-16 09:53:27.149059', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927107363778562', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:53:27.055208', 393, 0, '执行失败: 
### Error updating database. Cause: org.postgresql.util.PSQLException: ERROR: relation `sys_banner` does not exist
 Position: 9
### The error may exist in github/jiangbyte/io/sys/modules/banner/mapper/SysBannerMapper.java (best guess)
### The error may involve github.jiangbyte.io.sys.modules.banner.mapper.SysBannerMapper.update-Inline
### The error occurred while setting parameters
### SQL: UPDATE sys_banner SET status=?,updated_at=? WHERE (status = ? AND end_at IS NOT NULL AND end_at < ?)
### Cause: org.postgresql.util.PSQLException: ERROR: relation `sys_banner` does not exist
 Position: 9
; bad SQL grammar []', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:53:27.452823', NULL, '2026-08-16 09:53:27.452823', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927360938815489', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:54:27.891935', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:54:27.900684', NULL, '2026-08-16 09:54:27.900684', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927360938815490', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:54:27.895819', 3, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:54:27.904202', NULL, '2026-08-16 09:54:27.904202', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927360984952833', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:54:27.895819', 12, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:54:27.912183', NULL, '2026-08-16 09:54:27.912183', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927616212545538', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:55:28.752201', 3, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:55:28.762236', NULL, '2026-08-16 09:55:28.762236', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927616212545537', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:55:28.752201', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:55:28.761184', NULL, '2026-08-16 09:55:28.761184', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927616246099970', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:55:28.750203', 12, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:55:28.7684', NULL, '2026-08-16 09:55:28.7684', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927870076989441', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:56:29.281856', 3, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:56:29.288857', NULL, '2026-08-16 09:56:29.288857', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927870076989442', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:56:29.285852', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:56:29.289371', NULL, '2026-08-16 09:56:29.289371', NULL);
INSERT INTO `sys_job_log` VALUES ('2088927870076989443', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:56:29.284857', 9, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:56:29.297382', NULL, '2026-08-16 09:56:29.297382', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928124180508674', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:57:29.881393', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:57:29.884391', NULL, '2026-08-16 09:57:29.884391', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928124247617537', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:57:29.881393', 0, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:57:29.885615', NULL, '2026-08-16 09:57:29.885615', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928124247617538', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:57:29.879395', 4, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:57:29.887615', NULL, '2026-08-16 09:57:29.887615', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928365785001986', '7541000000000000004', '审计告警', '{}', '2026-08-16 09:58:27.458677', 8, 0, 'done fired=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:58:27.476045', NULL, '2026-08-16 09:58:27.476045', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928378409857025', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:58:30.49088', 0, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:58:30.494351', NULL, '2026-08-16 09:58:30.494351', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928378460188673', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:58:30.49188', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:58:30.494351', NULL, '2026-08-16 09:58:30.494351', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928378464382977', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:58:30.49188', 3, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:58:30.498381', NULL, '2026-08-16 09:58:30.498381', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928633280933889', '7541000000000000001', '示例任务', '{}', '2026-08-16 09:59:31.244234', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:59:31.248923', NULL, '2026-08-16 09:59:31.248923', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928633280933890', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 09:59:31.244234', 3, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:59:31.248923', NULL, '2026-08-16 09:59:31.248923', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928633280933891', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 09:59:31.242237', 6, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 09:59:31.257338', NULL, '2026-08-16 09:59:31.257338', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928887933906946', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:00:31.961377', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:00:31.966934', NULL, '2026-08-16 10:00:31.966934', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928887933906947', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:00:31.964457', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:00:31.970929', NULL, '2026-08-16 10:00:31.970929', NULL);
INSERT INTO `sys_job_log` VALUES ('2088928888001015810', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:00:31.964457', 7, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:00:31.978208', NULL, '2026-08-16 10:00:31.978208', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929142687543297', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:01:32.701165', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:01:32.708191', NULL, '2026-08-16 10:01:32.708191', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929142687543298', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:01:32.702203', 1, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:01:32.709678', NULL, '2026-08-16 10:01:32.709678', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929142750457858', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:01:32.703456', 8, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:01:32.717355', NULL, '2026-08-16 10:01:32.717355', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929397172744194', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:02:33.383666', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:02:33.388188', NULL, '2026-08-16 10:02:33.388188', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929397239853058', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:02:33.384973', 4, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:02:33.393205', NULL, '2026-08-16 10:02:33.393205', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929397239853059', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:02:33.384973', 8, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:02:33.395612', NULL, '2026-08-16 10:02:33.395612', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929626609561601', '7541000000000000004', '审计告警', '{}', '2026-08-16 10:03:28.055669', 18, 0, 'done fired=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:03:28.079276', NULL, '2026-08-16 10:03:28.079276', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929652035432451', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:03:34.140859', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:03:34.147752', NULL, '2026-08-16 10:03:34.147752', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929906495467523', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:04:34.8138', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:04:34.820209', NULL, '2026-08-16 10:04:34.820209', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930161400098818', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:05:35.576894', 7, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:05:35.589373', NULL, '2026-08-16 10:05:35.589373', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930416074043393', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:06:36.305621', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:06:36.306939', NULL, '2026-08-16 10:06:36.306939', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930670785736707', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:07:37.023242', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:07:37.028469', NULL, '2026-08-16 10:07:37.028469', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930925421932547', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:08:37.737872', 6, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:08:37.749604', NULL, '2026-08-16 10:08:37.749604', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929652035432452', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:03:34.13769', 9, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:03:34.150764', NULL, '2026-08-16 10:03:34.150764', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929906495467522', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:04:34.807801', 1, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:04:34.816192', NULL, '2026-08-16 10:04:34.816192', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930161400098819', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:05:35.580523', 0, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:05:35.589373', NULL, '2026-08-16 10:05:35.589373', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930416074043395', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:06:36.306939', 5, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:06:36.315414', NULL, '2026-08-16 10:06:36.315414', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930670785736706', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:07:37.022249', 0, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:07:37.028469', NULL, '2026-08-16 10:07:37.028469', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930925421932546', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:08:37.740066', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:08:37.744113', NULL, '2026-08-16 10:08:37.744113', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931180049739777', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:09:38.440697', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:09:38.447166', NULL, '2026-08-16 10:09:38.447166', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931434568495107', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:10:39.12234', 4, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:10:39.130045', NULL, '2026-08-16 10:10:39.130045', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931689120804866', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:11:39.828285', 1, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:11:39.831652', NULL, '2026-08-16 10:11:39.831652', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931944155459588', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:12:40.625177', 8, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:12:40.637277', NULL, '2026-08-16 10:12:40.637277', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932147981856770', '7541000000000000004', '审计告警', '{}', '2026-08-16 10:13:29.213709', 14, 0, 'done fired=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:13:29.23146', NULL, '2026-08-16 10:13:29.23146', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932198909095937', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:13:41.35542', 3, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:13:41.362956', NULL, '2026-08-16 10:13:41.362956', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932453742424067', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:14:42.125253', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:14:42.127254', NULL, '2026-08-16 10:14:42.127254', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932708277956610', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:15:42.797756', 6, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:15:42.80824', NULL, '2026-08-16 10:15:42.80824', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932962696048641', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:16:43.460235', 1, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:16:43.466772', NULL, '2026-08-16 10:16:43.466772', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933217357410305', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:17:44.175377', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:17:44.183342', NULL, '2026-08-16 10:17:44.183342', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933471783890946', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:18:44.847336', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:18:44.852872', NULL, '2026-08-16 10:18:44.852872', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933726340395011', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:19:45.529981', 1, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:19:45.536951', NULL, '2026-08-16 10:19:45.536951', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931180049739778', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:09:38.440697', 4, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:09:38.448163', NULL, '2026-08-16 10:09:38.448163', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931434568495106', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:10:39.12234', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:10:39.128031', NULL, '2026-08-16 10:10:39.128031', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931689120804867', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:11:39.826288', 3, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:11:39.83265', NULL, '2026-08-16 10:11:39.83265', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931944155459586', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:12:40.622833', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:12:40.628181', NULL, '2026-08-16 10:12:40.628181', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932198909095939', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:13:41.358507', 9, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:13:41.372395', NULL, '2026-08-16 10:13:41.372395', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932453742424068', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:14:42.125253', 2, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:14:42.130252', NULL, '2026-08-16 10:14:42.130252', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932708215042051', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:15:42.797756', 3, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:15:42.803956', NULL, '2026-08-16 10:15:42.803956', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932962696048642', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:16:43.463127', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:16:43.467775', NULL, '2026-08-16 10:16:43.467775', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933217357410306', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:17:44.177343', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:17:44.183342', NULL, '2026-08-16 10:17:44.183342', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933471846805505', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:18:44.847336', 5, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:18:44.852872', NULL, '2026-08-16 10:18:44.852872', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933726340395012', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:19:45.527979', 8, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:19:45.540749', NULL, '2026-08-16 10:19:45.540749', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931180049739779', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:09:38.439683', 8, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:09:38.45217', NULL, '2026-08-16 10:09:38.45217', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931434505580545', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:10:39.12133', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:10:39.125343', NULL, '2026-08-16 10:10:39.125343', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931689120804865', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:11:39.828285', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:11:39.829485', NULL, '2026-08-16 10:11:39.829485', NULL);
INSERT INTO `sys_job_log` VALUES ('2088931944155459587', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:12:40.629287', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:12:40.63219', NULL, '2026-08-16 10:12:40.63219', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932198909095938', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:13:41.357445', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:13:41.362956', NULL, '2026-08-16 10:13:41.362956', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932453742424066', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:14:42.124254', 0, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:14:42.127254', NULL, '2026-08-16 10:14:42.127254', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932708215042050', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:15:42.797756', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:15:42.803956', NULL, '2026-08-16 10:15:42.803956', NULL);
INSERT INTO `sys_job_log` VALUES ('2088932962696048643', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:16:43.463127', 7, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:16:43.47477', NULL, '2026-08-16 10:16:43.47477', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933217357410307', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:17:44.177343', 5, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:17:44.185609', NULL, '2026-08-16 10:17:44.185609', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933408345042945', '7541000000000000004', '审计告警', '{}', '2026-08-16 10:18:29.691565', 26, 0, 'done fired=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:18:29.722212', NULL, '2026-08-16 10:18:29.722212', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933471783890945', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:18:44.847336', 0, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:18:44.852872', NULL, '2026-08-16 10:18:44.852872', NULL);
INSERT INTO `sys_job_log` VALUES ('2088933726340395010', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:19:45.528984', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:19:45.530984', NULL, '2026-08-16 10:19:45.530984', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935281412177921', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:25:56.241705', 38, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:25:56.286976', NULL, '2026-08-16 10:25:56.286976', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935281286348801', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:25:56.241705', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:25:56.257657', NULL, '2026-08-16 10:25:56.257657', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935281437343746', '7541000000000000004', '审计告警', '{}', '2026-08-16 10:25:56.244268', 39, 0, 'done fired=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:25:56.293494', NULL, '2026-08-16 10:25:56.293494', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935536581050370', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:26:57.11642', 2, 0, 'echo: (无参数)', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:26:57.129211', NULL, '2026-08-16 10:26:57.129211', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935536581050371', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:26:57.119242', 2, 0, 'flushed=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:26:57.132257', NULL, '2026-08-16 10:26:57.132257', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935536648159234', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:26:57.11642', 18, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:26:57.141369', NULL, '2026-08-16 10:26:57.141369', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935791720562689', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:27:57.95913', 1, 0, 'echo: (无参数)', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:27:57.965126', NULL, '2026-08-16 10:27:57.965126', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935791720562690', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:27:57.95913', 3, 0, 'flushed=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:27:57.968132', NULL, '2026-08-16 10:27:57.968132', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935791791865857', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:27:57.95913', 11, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:27:57.976625', NULL, '2026-08-16 10:27:57.976625', NULL);
INSERT INTO `sys_job_log` VALUES ('2088935281286348802', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:25:56.244268', 3, 0, 'flushed=0', 'system', '192.168.10.4', '8828', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:25:56.257657', NULL, '2026-08-16 10:25:56.257657', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929652035432450', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:03:34.140859', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:03:34.147752', NULL, '2026-08-16 10:03:34.147752', NULL);
INSERT INTO `sys_job_log` VALUES ('2088929906558382081', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:04:34.810807', 7, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:04:34.823206', NULL, '2026-08-16 10:04:34.823206', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930161328795649', '7541000000000000001', '示例任务', '{}', '2026-08-16 10:05:35.575834', 0, 0, 'echo: (无参数)', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:05:35.57943', NULL, '2026-08-16 10:05:35.57943', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930416074043394', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:06:36.305621', 1, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:06:36.310887', NULL, '2026-08-16 10:06:36.310887', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930670785736708', '7541000000000000002', 'Banner 状态同步', '{}', '2026-08-16 10:07:37.023242', 8, 0, 'expired=0,activated=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:07:37.035507', NULL, '2026-08-16 10:07:37.035507', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930887169880066', '7541000000000000004', '审计告警', '{}', '2026-08-16 10:08:28.610841', 8, 0, 'done fired=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:08:28.621044', NULL, '2026-08-16 10:08:28.621044', NULL);
INSERT INTO `sys_job_log` VALUES ('2088930925421932545', '7541000000000000003', 'Banner 互动计数刷库', '{}', '2026-08-16 10:08:37.736858', 2, 0, 'flushed=0', 'system', '192.168.10.4', '1464', 'E:\projects\mine\hei\hei-boot', '2026-08-16 10:08:37.743081', NULL, '2026-08-16 10:08:37.743081', NULL);

-- ----------------------------
-- Table structure for sys_notice
-- ----------------------------
DROP TABLE IF EXISTS `sys_notice`;
CREATE TABLE `sys_notice` (
 `id` varchar(64) NOT NULL,
 `kind` varchar(32) NOT NULL,
 `title` varchar(255) NOT NULL,
 `content` text NOT NULL,
 `content_type` varchar(32) NOT NULL,
 `category` varchar(32),
 `severity` varchar(32) NOT NULL,
 `target_scope` varchar(32) NOT NULL,
 `target_account_types` json NOT NULL,
 `target_account_ids` json NOT NULL,
 `target_dept_ids` json NOT NULL,
 `target_role_ids` json NOT NULL,
 `publish_locations` json NOT NULL,
 `is_pinned` tinyint(1) NOT NULL,
 `pinned_until` datetime(6),
 `sender_account_type` varchar(32),
 `sender_account_id` varchar(64),
 `source_type` varchar(64),
 `source_id` varchar(64),
 `status` varchar(32) NOT NULL,
 `publish_at` datetime(6),
 `revoked_at` datetime(6),
 `expire_at` datetime(6),
 `view_count` int NOT NULL,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_notice
-- ----------------------------
INSERT INTO `sys_notice` VALUES ('7491842112464527360', 'NOTIFICATION', '急急急', '哈哈哈', 'text', 'SYSTEM', 'INFO', 'ALL', '["PORTAL"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 13:05:46', NULL, NULL, 0, '{}', '2026-08-08 13:05:51.295643', '1', '2026-08-08 13:05:51.295643', '1');
INSERT INTO `sys_notice` VALUES ('7491842211315884032', 'NOTIFICATION', '哈哈哈', '哈哈哈哈哈', 'text', 'SYSTEM', 'INFO', 'ALL', '["ADMIN"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 13:06:10', NULL, NULL, 0, '{}', '2026-08-08 13:06:14.871274', '1', '2026-08-08 13:06:14.871274', '1');
INSERT INTO `sys_notice` VALUES ('7491853809015291905', 'ANNOUNCEMENT', '系统维护预告', '本周日 02:00-04:00 将进行例行维护，期间门户可能短暂不可用，请提前做好安排。', 'markdown', 'SYSTEM', 'WARNING', 'ALL', '["PORTAL", "ADMIN"]', '[]', '[]', '[]', '{"center": true, "popup": true, "dashboard": true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 12:52:19.890496', NULL, '2026-08-22 13:52:19.890496', 1, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 14:11:23.438139', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809015291906', 'ANNOUNCEMENT', '意见反馈功能上线', '现已支持在线提交意见反馈并查看处理进度。登录后打开用户菜单中的「我的反馈」即可使用。', 'text', 'SYSTEM', 'SUCCESS', 'ACCOUNT_TYPE', '["PORTAL"]', '[]', '[]', '[]', '{"center": true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 13:22:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809044652032', 'NOTIFICATION', '账号安全提醒', '建议定期修改密码，并确保绑定的手机号与邮箱可用，以便找回账号。', 'text', 'SECURITY', 'WARNING', 'ALL', '["PORTAL"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, 'SYSTEM', NULL, 'PUBLISHED', '2026-08-08 13:32:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809044652033', 'NOTIFICATION', '新功能提示：消息中心', '右上角铃铛可查看未读通知与公告，支持一键全部已读。', 'text', 'SYSTEM', 'INFO', 'ALL', '["PORTAL"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, 'SYSTEM', NULL, 'PUBLISHED', '2026-08-08 13:42:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809044652034', 'NOTIFICATION', '管理端测试通知', '这是一条仅面向 ADMIN 的通知，用于验证账户类型过滤。', 'text', 'SYSTEM', 'INFO', 'ACCOUNT_TYPE', '["ADMIN"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, 'SYSTEM', NULL, 'PUBLISHED', '2026-08-08 13:47:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809015291904', 'ANNOUNCEMENT', '欢迎使用 HEI 门户', '门户账号体系、个人中心与消息中心已就绪。如有问题可通过「我的反馈」提交。', 'text', 'SYSTEM', 'INFO', 'ALL', '["PORTAL"]', '[]', '[]', '[]', '{"center":true,"dashboard":true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 11:52:19.890496', NULL, '2026-11-06 13:52:19.890496', 6, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 14:56:51.531823', '7491847383584804864');

-- ----------------------------
-- Table structure for sys_notice_read
-- ----------------------------
DROP TABLE IF EXISTS `sys_notice_read`;
CREATE TABLE `sys_notice_read` (
 `id` varchar(64) NOT NULL,
 `notice_id` varchar(64) NOT NULL,
 `account_type` varchar(32) NOT NULL,
 `account_id` varchar(64) NOT NULL,
 `read_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_notice_read
-- ----------------------------
INSERT INTO `sys_notice_read` VALUES ('7491843548694908928', '7491842211315884032', 'ADMIN', '1', '2026-08-08 13:11:33.720396');
INSERT INTO `sys_notice_read` VALUES ('7491847557782638592', '7491842112464527360', 'PORTAL', '7491847383584804864', '2026-08-08 13:27:29.554552');
INSERT INTO `sys_notice_read` VALUES ('7491853876019257344', '7491853809015291904', 'PORTAL', '7491847383584804864', '2026-08-08 13:52:35.928972');
INSERT INTO `sys_notice_read` VALUES ('7491853918088126464', '7491853809044652033', 'PORTAL', '7491847383584804864', '2026-08-08 13:52:45.966126');
INSERT INTO `sys_notice_read` VALUES ('7491853936312377344', '7491853809015291905', 'PORTAL', '7491847383584804864', '2026-08-08 13:52:50.326041');
INSERT INTO `sys_notice_read` VALUES ('7491853936312377345', '7491853809015291906', 'PORTAL', '7491847383584804864', '2026-08-08 13:52:50.326041');
INSERT INTO `sys_notice_read` VALUES ('7491853936312377346', '7491853809044652032', 'PORTAL', '7491847383584804864', '2026-08-08 13:52:50.326041');
INSERT INTO `sys_notice_read` VALUES ('7491856130134695936', '7491853809015291905', 'ADMIN', '1', '2026-08-08 14:01:33.364809');
INSERT INTO `sys_notice_read` VALUES ('7491856130134695937', '7491853809044652034', 'ADMIN', '1', '2026-08-08 14:01:33.364809');

-- ----------------------------
-- Table structure for sys_operation_audit_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_audit_log`;
CREATE TABLE `sys_operation_audit_log` (
 `id` varchar(64) NOT NULL,
 `module` varchar(64) NOT NULL,
 `resource_type` varchar(128),
 `resource_id` varchar(128),
 `action` varchar(64) NOT NULL,
 `summary` varchar(255),
 `before_data` json,
 `after_data` json,
 `account_id` varchar(64),
 `account_type` varchar(32),
 `request_id` varchar(64),
 `ip` varchar(64),
 `user_agent` varchar(512),
 `success` tinyint(1) NOT NULL,
 `error_message` text,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_operation_audit_log
-- ----------------------------
INSERT INTO `sys_operation_audit_log` VALUES ('7491824755243462656', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', NULL, NULL, NULL, 0, NULL, '2026-08-08 11:56:52.549136');
INSERT INTO `sys_operation_audit_log` VALUES ('7491826689375354880', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '997f4d2c30b34d2fab35091b27f20dc1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 12:04:33.755576');
INSERT INTO `sys_operation_audit_log` VALUES ('7491842112552607744', 'iam', 'create', NULL, 'post', 'POST /api/v1/admin/sys/notices/create', 'null', 'null', '1', 'ADMIN', '6c8b3cc0d20a4a60866b37c2e8a98baf', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:05:51.329499');
INSERT INTO `sys_operation_audit_log` VALUES ('7491842211382992896', 'iam', 'create', NULL, 'post', 'POST /api/v1/admin/sys/notices/create', 'null', 'null', '1', 'ADMIN', '606a87d34ba3408db97735ddeaabe10a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:06:14.892188');
INSERT INTO `sys_operation_audit_log` VALUES ('7491845090663665664', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/profile/avatar/upload', 'null', 'null', '1', 'ADMIN', 'e779ab8350f7430ca6ab3d3c5689a822', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:17:41.365887');
INSERT INTO `sys_operation_audit_log` VALUES ('7491847385635819520', 'auth', 'account', '7491847383584804864', 'register', 'Portal account registered', 'null', 'null', '7491847383584804864', 'PORTAL', '178b180693724ed4bbdef49d6411609c', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:26:48.530169');
INSERT INTO `sys_operation_audit_log` VALUES ('7491847464409042944', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'ef3e093461f64d63a33ceff7fa281b98', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:27:06.809483');
INSERT INTO `sys_operation_audit_log` VALUES ('7491849365728989184', 'iam', 'upload', NULL, 'post', 'POST /api/v1/portal/sys/file/upload', 'null', 'null', '7491847383584804864', 'PORTAL', 'ef951dc41e29426689d40f63b2650a2c', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:34:40.621796');
INSERT INTO `sys_operation_audit_log` VALUES ('7491849375170367488', 'iam', 'submit', NULL, 'post', 'POST /api/v1/portal/sys/feedbacks/submit', 'null', 'null', '7491847383584804864', 'PORTAL', '45261bcedc5c4900a5daa0e6b51ba112', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:34:42.871686');
INSERT INTO `sys_operation_audit_log` VALUES ('7491850767205326848', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '7491847383584804864', 'PORTAL', '9db3060fba264773ae39314a7a2a3e3d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '403', '2026-08-08 13:40:14.75986');
INSERT INTO `sys_operation_audit_log` VALUES ('7491850787824525312', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', 'cbef755d2cf74399bbf78075b1b4a569', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:40:19.210215');
INSERT INTO `sys_operation_audit_log` VALUES ('7491850914173739008', 'iam', 'update', NULL, 'post', 'POST /api/v1/admin/sys/feedbacks/update', 'null', 'null', '1', 'ADMIN', '10c752edb8ee4d45b7a5868a9a5c0bee', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:40:49.799448');
INSERT INTO `sys_operation_audit_log` VALUES ('7491851004623904768', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', '8121b706c0ab43a59cf412d3d4a889d1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:41:10.980972');
INSERT INTO `sys_operation_audit_log` VALUES ('7491852911862013952', 'iam', 'upload', NULL, 'post', 'POST /api/v1/portal/profile/avatar/upload', 'null', 'null', '7491847383584804864', 'PORTAL', 'a340f1b7b3174d1c8b5918818666e40d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:48:46.085676');
INSERT INTO `sys_operation_audit_log` VALUES ('7491853936371097600', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/portal/sys/notices/read-all', 'null', 'null', '7491847383584804864', 'PORTAL', '4a07b49c65ab41aa8f6eb409aeb81f44', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 13:52:50.347166');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856130252136448', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'a38dac89228e43c4809a2d0ccddade60', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:01:33.410073');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856145267744768', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'e4a86183e35b4211b5d81c2d78ac128a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:01:36.98952');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856147004186624', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', '36b62e3e1dc54b438d1c0a694de8b974', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:01:37.403846');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856149260722176', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'f7d8fe609757442c9f42a9557ce49a4b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:01:37.941479');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856211504193536', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', '3b4a2d936c0f44d78d8c9dc24c2cd748', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:01:52.781516');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856220572278784', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'e8b49b5611db4c498bb8fcb2f0ddfe12', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:01:54.9441');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856359823171584', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'e95abd971e9b45f2a0e930e50c233117', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:02:28.143287');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856638173962240', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'a78eed380e434144a1b1e74729d6e54f', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:03:34.506926');
INSERT INTO `sys_operation_audit_log` VALUES ('7491867956285251584', 'auth', 'account', 'k-oPNnS5p0ls7b7JnuEO8yRR80EZ3FR4CCStRQvGdLo', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', '60488ac9dcec4a59bf287de53e535a92', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:48:32.938124');
INSERT INTO `sys_operation_audit_log` VALUES ('7491867956377526272', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', '60488ac9dcec4a59bf287de53e535a92', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:48:32.976961');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868063009316864', 'auth', 'account', 'P48I4wXyp2T1K6Gws0WwWa-48PnEE1AEQshPdhtbqko', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', '8aee2686706648ac8febf62481191666', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:48:58.392732');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868063068037120', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', '8aee2686706648ac8febf62481191666', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:48:58.413543');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868407806271488', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'd249fc3c14104900b1c5ab59acc4ebdc', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:50:20.14571');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868707275382784', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '3436832124ae49b7bd8885920146ed05', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:51:31.644571');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868946300379136', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', 'd2b5c5b182d844089cde10f6802c702d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:52:28.992926');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869125392965632', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', 'bfad8619bcb04e929dbe3b8f4ceeb34b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:53:11.69228');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869202706571264', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '54a3eb743686490089f5f64617d88f23', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-08 14:53:30.125444');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869288723357696', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', '94a81d39bd6d4179b431c24ba933c36b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:53:50.63331');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869364359241728', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '36a39cc7c1b04ed3bd3f7be81c619a76', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:54:08.665933');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869492054827008', 'auth', 'account', 'Lr0IMLzxN8NytPPOtKb-7axzFWyr4kB2Y1fiu1Ofk7w', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'f7e05c09efd5437c9fa4e07971765789', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:54:39.098545');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869492138713088', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'f7e05c09efd5437c9fa4e07971765789', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:54:39.131162');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869580353314816', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'ed81b9adb9254c798c5389e9304fae3e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:54:59.8293');
INSERT INTO `sys_operation_audit_log` VALUES ('7491870806759415808', 'auth', 'account', '8GXa7mbl8xxMrkS8e50Kcnc6wjnIDWwHCEZ9RZkFi6Y', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', '5d92a1b221184b848bcad88547b77ec2', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:59:52.539321');
INSERT INTO `sys_operation_audit_log` VALUES ('7491870806889439232', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', '5d92a1b221184b848bcad88547b77ec2', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:59:52.592863');
INSERT INTO `sys_operation_audit_log` VALUES ('7491870836610277376', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '9f84348d2ea048569a529eb0db336177', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 14:59:59.115155');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871812532523008', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', '6505ddfea04149a39f336a1a62857bde', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:03:52.356699');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871872020336640', 'auth', 'account', '-1hh-Gk-I2IVhH5QvuGxD6aKn71kIqmWf2LHo6gW7tA', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'ae74d9a084604aa28c8f3c54af94de95', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:04:06.528253');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871872095834112', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'ae74d9a084604aa28c8f3c54af94de95', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:04:06.55724');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871928467279872', 'auth', 'account', '7491847383584804864', 'forgot_password', 'PORTAL password reset requested', 'null', 'null', '7491847383584804864', 'PORTAL', '57b11ecf39894f2e8e40bcdd2f636fe6', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:04:19.399391');
INSERT INTO `sys_operation_audit_log` VALUES ('7491872608439390208', 'auth', 'account', '7491847383584804864', 'forgot_password', 'PORTAL password reset requested', 'null', 'null', '7491847383584804864', 'PORTAL', '6e791c69c0cc4d858d75af582256e3f4', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:07:01.437344');
INSERT INTO `sys_operation_audit_log` VALUES ('7491872894360899584', 'auth', 'account', '7491872891940786176', 'register', 'Portal account registered', 'null', 'null', '7491872891940786176', 'PORTAL', 'bd310a77940d4c1d8e3ad519c3cc09c0', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:08:10.284656');
INSERT INTO `sys_operation_audit_log` VALUES ('7491875322690949120', 'iam', 'send-login-code', NULL, 'post', 'POST /api/v1/portal/send-login-code', 'null', 'null', NULL, NULL, 'a270b9c191dd4879bbfbb36ba9ea9cd4', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '400', '2026-08-08 15:17:49.244173');
INSERT INTO `sys_operation_audit_log` VALUES ('7491875361068830720', 'iam', 'send-login-code', NULL, 'post', 'POST /api/v1/portal/send-login-code', 'null', 'null', NULL, NULL, '2eb366121b4a49989fa62fa361f605f2', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:17:58.394406');
INSERT INTO `sys_operation_audit_log` VALUES ('7491875444590006272', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'c2b3d7256b6e4c54a489859978feb288', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:18:18.2083');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876150378123264', 'auth', 'account', 'cKRn3a7CVntuJV6gMUKhYxVdj66PajD-FmVnQE8WHGQ', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'c81c429632dd4d6cb6c5e71ea025465b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:21:06.560599');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876150457815040', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'c81c429632dd4d6cb6c5e71ea025465b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:21:06.597915');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876284402913280', 'auth', 'account', '218B1W9aw378shi2KJwdvkHL6dy4KOxhy0RThJIwBkU', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', '4c93ad1119ae46fda36dfcbce10e964a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:21:38.515325');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876284495187968', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', '4c93ad1119ae46fda36dfcbce10e964a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 15:21:38.555219');
INSERT INTO `sys_operation_audit_log` VALUES ('7491886508237017088', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '06ffe05bd308431a85c804da1deff708', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 16:02:15.699453');
INSERT INTO `sys_operation_audit_log` VALUES ('7491890281995022336', 'iam', 'interaction', NULL, 'post', 'POST /api/v1/portal/sys/banners/interaction', 'null', 'null', NULL, NULL, '557c0d9c0e7e4a0186aff815a7655afb', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 16:17:15.820139');
INSERT INTO `sys_operation_audit_log` VALUES ('7491905950954287104', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', 'fa3e17a02e6b4df3b99c9d7e4ebbf6ee', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 17:19:31.591092');
INSERT INTO `sys_operation_audit_log` VALUES ('7491906012094656512', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', 'ad21a931262a45928ee617d6a1425ece', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 17:19:46.1667');
INSERT INTO `sys_operation_audit_log` VALUES ('7491978456146837504', 'iam', 'interaction', NULL, 'post', 'POST /api/v1/portal/sys/banners/interaction', 'null', 'null', NULL, NULL, '605ec219cc7c4aa9988ea4a65a790f1e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:07:38.175092');
INSERT INTO `sys_operation_audit_log` VALUES ('7491978848712720384', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '2353489a36fc450ba6fc61e02c043450', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:09:10.963818');
INSERT INTO `sys_operation_audit_log` VALUES ('7491978930061246464', 'iam', 'delete', NULL, 'post', 'POST /api/v1/admin/sys/banners/delete', 'null', 'null', '1', 'ADMIN', '2f55e8ff734c4795bb34ee2d53d21f54', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:09:31.165058');
INSERT INTO `sys_operation_audit_log` VALUES ('7491984102703431680', 'auth', 'account', 'isGNdSOJtqTMMLd5tE89-zWE9E-gZTI2e9QtAN9O5IY', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', 'e91dd3e589f04d1fbd2847d159404e59', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:30:04.376178');
INSERT INTO `sys_operation_audit_log` VALUES ('7491984102850232320', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', 'e91dd3e589f04d1fbd2847d159404e59', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:30:04.454205');
INSERT INTO `sys_operation_audit_log` VALUES ('7491984343368400896', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '565285e8a92d4db183eb55f22bf5ebd5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:31:01.16333');
INSERT INTO `sys_operation_audit_log` VALUES ('7491985237854060544', 'auth', 'account', '7491847383584804864', 'forgot_password', 'PORTAL password reset requested', 'null', 'null', '7491847383584804864', 'PORTAL', 'e2553b09785a4da580824d763b076a86', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:34:34.641992');
INSERT INTO `sys_operation_audit_log` VALUES ('7491985393911529472', 'auth', 'account', '7491847383584804864', 'reset_password', 'PORTAL password reset', 'null', 'null', '7491847383584804864', 'PORTAL', 'ef552be360104587bf1e3160d613aeb1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:35:11.342559');
INSERT INTO `sys_operation_audit_log` VALUES ('7491985465474744320', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'b2cf5db7a39c42d792ee141212e7892e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-08 22:35:29.00935');
INSERT INTO `sys_operation_audit_log` VALUES ('7492049918010503168', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', 'bb7ce0ff89594460ac78a96c5f680e69', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 02:51:35.362724');
INSERT INTO `sys_operation_audit_log` VALUES ('7492050216368123904', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '45af2d5c4eb3416e953db56f294eb386', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 02:52:47.145857');
INSERT INTO `sys_operation_audit_log` VALUES ('7492050342729920512', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '3104ce6d274243619c356adc1f04279b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 02:53:17.271997');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070224125145088', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '187c6bf3c132445e87e101755a4a236b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 04:12:16.604286');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070372163104768', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', 'bc76b2dc3b9e45108466763edffa9ba0', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 04:12:52.661485');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070399157645312', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '8f9d09659bae4815806de838cf2699f5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 04:12:59.097782');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070449363464192', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '9ce941c98b414b3cbc84e53cf7b25849', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 04:13:11.067709');
INSERT INTO `sys_operation_audit_log` VALUES ('7492074203068411904', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '69519ea4cf5542ccabebc8c6182125f6', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 04:28:06.020572');
INSERT INTO `sys_operation_audit_log` VALUES ('7492074211410882560', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '4228001a67784ec7a6201e518de4fd92', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', '2026-08-09 04:28:08.009335');
INSERT INTO `sys_operation_audit_log` VALUES ('7492078890945560576', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', 'b4da294ce5f746989d1ae6bebad11c65', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 04:46:43.055912');
INSERT INTO `sys_operation_audit_log` VALUES ('7492079018624368640', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', '11a3c39a52e84f90bdd30204d1f90a90', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 04:47:13.496031');
INSERT INTO `sys_operation_audit_log` VALUES ('2086348203017056258', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', 'aab6d330ea5648ac8609cedf04598930', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 07:05:48.561987');
INSERT INTO `sys_operation_audit_log` VALUES ('2086404723817709569', 'iam', 'sys_file', NULL, 'upload', 'POST /api/v1/admin/sys/file/upload', NULL, NULL, '1', 'admin', '9680f40a4f99443ca7ad392b4a8bcc01', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 10:50:24.191888');
INSERT INTO `sys_operation_audit_log` VALUES ('2086408062643171330', 'iam', 'sys_file', NULL, 'upload', 'POST /api/v1/admin/sys/file/upload', NULL, NULL, '1', 'admin', '0b0147e7fc334c46bd0ce20b4bd3e8e5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:03:40.112575');
INSERT INTO `sys_operation_audit_log` VALUES ('2086413233146265601', 'iam', 'auth', NULL, 'login', 'POST /api/v1/portal/login', NULL, NULL, '7491847383584804864', 'portal', '346ae88b022e4dcca715f45c9d1acd99', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:24:12.949234');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414670496473090', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/portal/logout', NULL, NULL, NULL, NULL, 'f69f316b0d1d456c96267ff44156cc38', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:29:55.529704');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414718772912129', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, 'b57f1bbb39e14d19ab83073979cb9ad9', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:30:07.299008');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414768404111361', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '0c33636455a84a968a602e8ca91efd15', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:30:19.117349');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414793913868289', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, '16ca7bf3218641f199759626faaba64d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:30:25.202419');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414860225814529', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '3f32732033144eeeab4d6ab4fa098853', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:30:41.022353');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415187943563265', 'iam', 'sys_file', NULL, 'upload', 'POST /api/v1/admin/sys/file/upload', NULL, NULL, '1', 'admin', '94e68d53f1c54a39899e57e8777303a4', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:31:59.153041');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415233657282562', 'iam', 'iam_account', NULL, 'update', 'POST /api/v1/admin/sys/accounts/update', NULL, NULL, '1', 'admin', '74e6acb5b526402da0cb5ae3a6bfe261', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:32:10.050871');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415589262958594', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, '61422f3e916c4389b835b42cd7ee89a1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:33:34.827482');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415657378455554', 'iam', 'auth', NULL, 'forgot_password', 'POST /api/v1/admin/forgot-password', NULL, NULL, NULL, NULL, '7d41b61b6e01477490fb699253770a64', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:33:51.061495');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415774965768194', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '964907d0682a4dfea90556655625386e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:34:19.122222');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415885812834306', 'iam', 'iam_account', NULL, 'update', 'POST /api/v1/admin/sys/accounts/update', NULL, NULL, '1', 'admin', '6627a2127a38469ebf9ea00ed3a6a4f9', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:34:45.539416');
INSERT INTO `sys_operation_audit_log` VALUES ('2086416001558847490', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, 'd022212e771f473aa82cf2f1de58c3b0', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:35:13.139818');
INSERT INTO `sys_operation_audit_log` VALUES ('2086416040901419009', 'iam', 'auth', NULL, 'forgot_password', 'POST /api/v1/admin/forgot-password', NULL, NULL, NULL, NULL, '65956d3abbe443db8be96f482cdf7d60', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:35:22.522867');
INSERT INTO `sys_operation_audit_log` VALUES ('2086418125776633857', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '45ce46c4357e44de944d457b6c5f6bad', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, '2026-08-09 11:43:39.491727');

-- ----------------------------
-- Table structure for sys_operation_audit_outbox
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_audit_outbox`;
CREATE TABLE `sys_operation_audit_outbox` (
 `id` varchar(64) NOT NULL,
 `payload` text NOT NULL,
 `status` varchar(32) NOT NULL,
 `attempts` int NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `claimed_at` datetime(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_operation_audit_outbox
-- ----------------------------

-- ----------------------------
-- Table structure for sys_position
-- ----------------------------
DROP TABLE IF EXISTS `sys_position`;
CREATE TABLE `sys_position` (
 `id` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `category` varchar(32) NOT NULL,
 `owner_dept_id` varchar(64),
 `sort` int NOT NULL,
 `is_virtual` tinyint(1) NOT NULL,
 `status` varchar(32) NOT NULL,
 `description` text,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_position
-- ----------------------------

-- ----------------------------
-- Table structure for sys_resource
-- ----------------------------
DROP TABLE IF EXISTS `sys_resource`;
CREATE TABLE `sys_resource` (
 `id` varchar(64) NOT NULL,
 `parent_id` varchar(64),
 `code` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `resource_type` varchar(32) NOT NULL,
 `module_id` varchar(64),
 `path` varchar(255),
 `component` varchar(255),
 `redirect` varchar(255),
 `icon` varchar(255),
 `color` varchar(32),
 `href` varchar(255),
 `sort` int NOT NULL,
 `is_visible` tinyint(1) NOT NULL,
 `is_cache` tinyint(1) NOT NULL,
 `is_affix` tinyint(1) NOT NULL,
 `status` varchar(32) NOT NULL,
 `description` text,
 `layout` varchar(255),
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_resource
-- ----------------------------
INSERT INTO `sys_resource` VALUES ('200001', NULL, 'dashboard', '运营工作台', 'MENU', '210001', '/dashboard', '/dashboard/index.vue', NULL, 'icon-park-outline:analysis', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200003', NULL, 'ops', '系统运维', 'CATALOG', '210001', '/sys', NULL, NULL, 'icon-park-outline:setting-two', NULL, NULL, 25, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200004', '200003', 'sys-dict', '字典管理', 'MENU', '210001', '/sys/dict', '/sys/dict/index.vue', NULL, 'icon-park-outline:file-search', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200005', '200019', 'content-banner', '展示图管理', 'MENU', '210001', '/sys/banner', '/sys/banner/index.vue', NULL, 'icon-park-outline:ad-product', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200006', NULL, 'org', '组织权限', 'CATALOG', '210001', '/iam', NULL, NULL, 'icon-park-outline:people', NULL, NULL, 10, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200007', '200006', 'iam-account', '账号管理', 'MENU', '210001', '/iam/account', '/iam/account/index.vue', NULL, 'icon-park-outline:people', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200008', '200006', 'iam-dept', '部门管理', 'MENU', '210001', '/iam/dept', '/iam/dept/index.vue', NULL, 'icon-park-outline:tree-diagram', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200009', '200006', 'iam-group', '用户组管理', 'MENU', '210001', '/iam/group', '/iam/group/index.vue', NULL, 'icon-park-outline:group', NULL, NULL, 3, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200010', '200006', 'iam-position', '岗位管理', 'MENU', '210001', '/iam/position', '/iam/position/index.vue', NULL, 'icon-park-outline:people-bottom', NULL, NULL, 4, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200011', '200006', 'iam-role', '角色管理', 'MENU', '210001', '/iam/role', '/iam/role/index.vue', NULL, 'icon-park-outline:peoples', NULL, NULL, 5, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200012', '200040', 'iam-resource', '资源管理', 'MENU', '210001', '/iam/resource', '/iam/resource/index.vue', NULL, 'icon-park-outline:all-application', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200018', '200040', 'iam-resourcemodule', '资源模块管理', 'MENU', '210001', '/iam/resource_module', '/iam/resource_module/index.vue', NULL, 'icon-park-outline:blocks-and-arrows', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200019', NULL, 'content', '内容运营', 'CATALOG', '210001', '/content', NULL, '/sys/notice', 'icon-park-outline:picture-album', NULL, NULL, 20, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200023', '200003', 'sys-file', '文件管理', 'MENU', '210001', '/sys/file', '/sys/file/index.vue', NULL, 'icon-park-outline:file-code', NULL, NULL, 3, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200025', '200003', 'sys-session', '在线会话', 'MENU', '210001', '/sys/session', '/auth/session/index.vue', NULL, 'icon-park-outline:connection', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200027', '200003', 'sys-audit-api', '操作审计接口', 'API_GROUP', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', '操作审计后端权限组', NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200028', '200003', 'sys-login-log', '登录日志', 'MENU', '210001', '/sys/login-log', '/sys/login-log/index.vue', NULL, 'icon-park-outline:log', NULL, NULL, 5, 1, 0, 0, 'ENABLED', '登录成功/失败历史记录', NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200029', '200003', 'sys-audit', '操作审计', 'MENU', '210001', '/sys/audit', '/sys/audit/index.vue', NULL, 'icon-park-outline:audit', NULL, NULL, 7, 1, 0, 0, 'ENABLED', '系统操作审计日志', NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200031', '200041', 'iam-clientmodule', '客户端模块管理', 'MENU', '210001', '/iam/client_module', '/iam/client_module/index.vue', NULL, 'icon-park-outline:application-one', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200032', '200041', 'iam-clientresource', '客户端资源管理', 'MENU', '210001', '/iam/client_resource', '/iam/client_resource/index.vue', NULL, 'icon-park-outline:page-template', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200040', NULL, 'resource-auth', '资源授权', 'CATALOG', '210001', '/resource-auth', NULL, NULL, 'icon-park-outline:all-application', NULL, NULL, 15, 1, 0, 0, 'ENABLED', '菜单资源与资源模块授权配置', NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('200041', NULL, 'client-resource-auth', '客户端资源授权', 'CATALOG', '210001', '/client-resource-auth', NULL, NULL, 'icon-park-outline:application-one', NULL, NULL, 16, 1, 0, 0, 'ENABLED', '客户端模块与客户端资源授权配置', NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201011', '200004', 'sys-dict-create', '新增字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201012', '200004', 'sys-dict-detail', '查看字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201013', '200004', 'sys-dict-update', '编辑字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201014', '200004', 'sys-dict-delete', '删除字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201021', '200005', 'sys-banner-create', '新增展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201022', '200005', 'sys-banner-detail', '查看展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201023', '200005', 'sys-banner-update', '编辑展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201024', '200005', 'sys-banner-delete', '删除展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201025', '200005', 'sys-banner-create-page', '新增展示图页', 'PAGE', '210001', '/sys/banner/create', '/sys/banner/form.vue', NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201026', '200005', 'sys-banner-edit-page', '编辑展示图页', 'PAGE', '210001', '/sys/banner/edit', '/sys/banner/form.vue', NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201027', '200005', 'sys-banner-detail-page', '展示图详情页', 'PAGE', '210001', '/sys/banner/detail', '/sys/banner/detail.vue', NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201031', '200023', 'sys-file-upload', '上传文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201032', '200023', 'sys-file-detail', '查看文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201033', '200023', 'sys-file-update', '编辑文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201034', '200023', 'sys-file-url', '打开文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201035', '200023', 'sys-file-delete', '删除文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201041', '200025', 'sys-session-tokenlist', '查看令牌', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201042', '200025', 'sys-session-exit', '强退账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201043', '200025', 'sys-session-tokenexit', '强退令牌', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201051', '202015', 'sys-codegen-create', '新增生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201052', '202015', 'sys-codegen-detail', '查看生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201053', '202015', 'sys-codegen-update', '编辑生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201054', '202015', 'sys-codegen-delete', '删除生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201055', '202015', 'sys-codegen-tables', '读取数据库表', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201056', '202015', 'sys-codegen-preview', '预览代码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 60, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201057', '202015', 'sys-codegen-download', '下载代码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 70, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-07-18 16:10:45', NULL);
INSERT INTO `sys_resource` VALUES ('201060', '200028', 'sys-login-log-detail', '查看登录日志', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201061', '200029', 'sys-audit-detail', '查看审计详情', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201101', '200007', 'iam-account-create', '新增账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201102', '200007', 'iam-account-detail', '查看账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201103', '200007', 'iam-account-update', '编辑账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201104', '200007', 'iam-account-delete', '删除账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201105', '200007', 'iam-account-grant-role', '分配角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201106', '200007', 'iam-account-grant-group', '分配用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201107', '200007', 'iam-account-grant-dept', '分配部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201108', '200007', 'iam-account-grant-resource', '分配资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201109', '200007', 'iam-account-grant-client-resource', '分配客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 9, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201121', '200008', 'iam-dept-create', '新增部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201122', '200008', 'iam-dept-detail', '查看部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201123', '200008', 'iam-dept-update', '编辑部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201124', '200008', 'iam-dept-delete', '删除部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201131', '200009', 'iam-group-create', '新增用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201132', '200009', 'iam-group-detail', '查看用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201133', '200009', 'iam-group-update', '编辑用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201134', '200009', 'iam-group-delete', '删除用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201135', '200009', 'iam-group-grant-user', '分配用户', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201136', '200009', 'iam-group-grant-role', '分配角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201137', '200009', 'iam-group-grant-resource', '分配资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201138', '200009', 'iam-group-grant-client-resource', '分配客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201151', '200010', 'iam-position-create', '新增岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201152', '200010', 'iam-position-detail', '查看岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201153', '200010', 'iam-position-update', '编辑岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201154', '200010', 'iam-position-delete', '删除岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201161', '200011', 'iam-role-create', '新增角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201162', '200011', 'iam-role-detail', '查看角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201163', '200011', 'iam-role-update', '编辑角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201164', '200011', 'iam-role-delete', '删除角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201165', '200011', 'iam-role-grant-resource', '分配资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201167', '200011', 'iam-role-grant-user', '分配用户', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201168', '200011', 'iam-role-grant-client-resource', '分配客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201181', '200012', 'iam-resource-create', '新增资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201182', '200012', 'iam-resource-detail', '查看资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201183', '200012', 'iam-resource-update', '编辑资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201184', '200012', 'iam-resource-delete', '删除资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201185', '200012', 'iam-resource-grant', '绑定权限', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201191', '200018', 'iam-resourcemodule-create', '新增资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201192', '200018', 'iam-resourcemodule-detail', '查看资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201193', '200018', 'iam-resourcemodule-update', '编辑资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201194', '200018', 'iam-resourcemodule-delete', '删除资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00', NULL, '2026-07-03 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201311', '200031', 'iam-clientmodule-create', '新增客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201312', '200031', 'iam-clientmodule-detail', '查看客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201313', '200031', 'iam-clientmodule-update', '编辑客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201314', '200031', 'iam-clientmodule-delete', '删除客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201321', '200032', 'iam-clientresource-create', '新增客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201322', '200032', 'iam-clientresource-detail', '查看客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201323', '200032', 'iam-clientresource-update', '编辑客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201324', '200032', 'iam-clientresource-delete', '删除客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201325', '200032', 'iam-clientresource-list', '客户端资源树', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('201326', '200032', 'iam-clientresource-grant', '绑定客户端资源权限', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202001', NULL, 'devtools', '开发工具', 'CATALOG', '210001', '/test', NULL, '/sys/codegen', 'icon-park-outline:code', NULL, NULL, 90, 1, 0, 0, 'ENABLED', '系统模块测试页面目录', NULL, '{}', '2026-07-18 12:39:16', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202002', '202001', 'system-test-editor', '编辑器测试', 'MENU', '210001', '/test/editor', '/test/editor/index.vue', NULL, 'icon-park-outline:edit', NULL, NULL, 2, 1, 0, 0, 'ENABLED', 'Markdown、富文本和代码编辑器组件测试页面', NULL, '{}', '2026-07-18 12:39:16', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202003', '202001', 'system-test-icon', '图标选择器测试', 'MENU', '210001', '/test/icon', '/test/icon/index.vue', NULL, 'icon-park-outline:all-application', NULL, NULL, 3, 1, 0, 0, 'ENABLED', 'Iconify 离线图标选择器测试页面', NULL, '{}', '2026-07-18 12:49:42', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202004', '202030', 'biz-cgtestactivity', '代码生成测试-活动', 'MENU', '210001', '/biz/cg-test-activity', '/biz/cg-test-activity/index.vue', NULL, 'icon-park-outline:calendar', NULL, NULL, 1, 1, 0, 0, 'ENABLED', '代码生成 CRUD 样例', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202005', '202030', 'biz-cgtestcatalog', '代码生成测试-目录树', 'MENU', '210001', '/biz/cg-test-catalog', '/biz/cg-test-catalog/index.vue', NULL, 'icon-park-outline:tree-list', NULL, NULL, 2, 1, 0, 0, 'ENABLED', '代码生成树表样例', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202006', '202030', 'biz-cgtestorder', '代码生成测试-订单', 'MENU', '210001', '/biz/cg-test-order', '/biz/cg-test-order/index.vue', NULL, 'icon-park-outline:transaction-order', NULL, NULL, 3, 1, 0, 0, 'ENABLED', '代码生成主子表样例', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202007', '202030', 'biz-cgtestknowledgecategory', '代码生成测试-知识分类', 'MENU', '210001', '/biz/cg-test-knowledge-category', '/biz/cg-test-knowledge-category/index.vue', NULL, 'icon-park-outline:book-open', NULL, NULL, 4, 1, 0, 0, 'ENABLED', '代码生成左树右表样例', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202010', '200003', 'system-config', '系统配置', 'MENU', '210001', '/sys/config', '/sys/config/index.vue', NULL, 'icon-park-outline:setting-config', NULL, NULL, 4, 1, 0, 0, 'ENABLED', '系统配置管理页面', NULL, '{}', '2026-07-18 14:07:48', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202011', '202010', 'sys:config:create', '新增系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48', NULL, '2026-07-18 14:07:48', NULL);
INSERT INTO `sys_resource` VALUES ('202012', '202010', 'sys:config:detail', '查看系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48', NULL, '2026-07-18 14:07:48', NULL);
INSERT INTO `sys_resource` VALUES ('202013', '202010', 'sys:config:update', '编辑系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48', NULL, '2026-07-18 14:07:48', NULL);
INSERT INTO `sys_resource` VALUES ('202014', '202010', 'sys:config:delete', '删除系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48', NULL, '2026-07-18 14:07:48', NULL);
INSERT INTO `sys_resource` VALUES ('202015', '202001', 'sys-codegen', '代码生成', 'MENU', '210001', '/sys/codegen', '/sys/codegen/index.vue', NULL, 'icon-park-outline:code', NULL, NULL, 1, 1, 0, 0, 'ENABLED', '代码生成管理', NULL, '{}', '2026-07-18 16:10:45', NULL, '2026-08-09 00:00:00', '1');
INSERT INTO `sys_resource` VALUES ('202030', NULL, 'biz-demo', '业务示例', 'CATALOG', '210001', '/biz', NULL, NULL, 'icon-park-outline:application-one', NULL, NULL, 40, 1, 0, 0, 'ENABLED', '代码生成业务示例页面', NULL, '{}', '2026-08-09 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202200', '200019', 'content-notice', '通知消息', 'MENU', '210001', '/sys/notice', '/sys/notice/index.vue', NULL, 'icon-park-outline:message', NULL, NULL, 2, 1, 0, 0, 'ENABLED', '消息管理', NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202201', '202200', 'sys-notice-page', '分页消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202202', '202200', 'sys-notice-create', '新增消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202203', '202200', 'sys-notice-detail', '详情消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202204', '202200', 'sys-notice-update', '编辑消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202205', '202200', 'sys-notice-delete', '删除消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202206', '202200', 'sys-notice-create-page', '新增消息页', 'PAGE', '210001', '/sys/notice/create', '/sys/notice/form.vue', NULL, NULL, NULL, NULL, 60, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202207', '202200', 'sys-notice-edit-page', '编辑消息页', 'PAGE', '210001', '/sys/notice/edit', '/sys/notice/form.vue', NULL, NULL, NULL, NULL, 70, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202208', '202200', 'sys-notice-detail-page', '消息详情页', 'PAGE', '210001', '/sys/notice/detail', '/sys/notice/detail.vue', NULL, NULL, NULL, NULL, 80, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202209', '202200', 'sys-notice-publish', '发布消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 55, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202220', '200019', 'content-feedback', '反馈管理', 'MENU', '210001', '/sys/feedback', '/sys/feedback/index.vue', NULL, 'icon-park-outline:write', NULL, NULL, 3, 1, 0, 0, 'ENABLED', '意见反馈管理', NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-08-09 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202221', '202220', 'sys-feedback-page', '分页反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202222', '202220', 'sys-feedback-detail', '查看反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202223', '202220', 'sys-feedback-update', '处理反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202224', '202220', 'sys-feedback-delete', '删除反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202225', '202220', 'sys-feedback-edit-page', '处理反馈页', 'PAGE', '210001', '/sys/feedback/edit', '/sys/feedback/form.vue', NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202226', '202220', 'sys-feedback-detail-page', '反馈详情页', 'PAGE', '210001', '/sys/feedback/detail', '/sys/feedback/detail.vue', NULL, NULL, NULL, NULL, 60, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00', NULL, '2026-07-24 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202240', '202200', 'sys-notice-revoke', '撤回消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 56, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('202241', '202200', 'sys-notice-pin', '置顶消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 57, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203011', '202004', 'biz-cgtestactivity-page', '分页活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203012', '202004', 'biz-cgtestactivity-create', '新增活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203013', '202004', 'biz-cgtestactivity-detail', '详情活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203014', '202004', 'biz-cgtestactivity-update', '编辑活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203015', '202004', 'biz-cgtestactivity-delete', '删除活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203021', '202005', 'biz-cgtestcatalog-page', '分页目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203022', '202005', 'biz-cgtestcatalog-create', '新增目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203023', '202005', 'biz-cgtestcatalog-detail', '详情目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203024', '202005', 'biz-cgtestcatalog-update', '编辑目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203025', '202005', 'biz-cgtestcatalog-delete', '删除目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203026', '202005', 'biz-cgtestcatalog-list', '树列表目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 90, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203031', '202006', 'biz-cgtestorder-page', '分页订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203032', '202006', 'biz-cgtestorder-create', '新增订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203033', '202006', 'biz-cgtestorder-detail', '详情订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203034', '202006', 'biz-cgtestorder-update', '编辑订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203035', '202006', 'biz-cgtestorder-delete', '删除订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203041', '202007', 'biz-cgtestknowledgecategory-page', '分页知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203042', '202007', 'biz-cgtestknowledgecategory-create', '新增知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203043', '202007', 'biz-cgtestknowledgecategory-detail', '详情知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203044', '202007', 'biz-cgtestknowledgecategory-update', '编辑知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203045', '202007', 'biz-cgtestknowledgecategory-delete', '删除知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('203046', '202007', 'biz-cgtestknowledgecategory-list', '树列表知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 90, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204001', '200003', 'sys-job', '任务管理', 'MENU', '210001', '/sys/job', '/sys/job/index.vue', NULL, 'icon-park-outline:timer', NULL, NULL, 4, 1, 0, 0, 'ENABLED', '任务调度管理（CRON / 固定间隔，Redis 锁防多实例重复执行）', NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204011', '204001', 'sys-job-create', '新增任务', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204012', '204001', 'sys-job-update', '编辑任务', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204013', '204001', 'sys-job-delete', '删除任务', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204014', '204001', 'sys-job-detail', '任务详情', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204015', '204001', 'sys-job-run', '立即执行', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204016', '204001', 'sys-job-log', '执行日志', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204021', '204001', 'sys-job-create-page', '新增任务页', 'PAGE', '210001', '/sys/job/create', '/sys/job/form.vue', NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204022', '204001', 'sys-job-edit-page', '编辑任务页', 'PAGE', '210001', '/sys/job/edit', '/sys/job/form.vue', NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204023', '204001', 'sys-job-detail-page', '任务详情页', 'PAGE', '210001', '/sys/job/detail', '/sys/job/detail.vue', NULL, NULL, NULL, NULL, 9, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);
INSERT INTO `sys_resource` VALUES ('204024', '204001', 'sys-job-log-page', '任务执行记录页', 'PAGE', '210001', '/sys/job/log', '/sys/job/log.vue', NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00', NULL, '2026-08-16 00:00:00', NULL);

-- ----------------------------
-- Table structure for sys_resource_module
-- ----------------------------
DROP TABLE IF EXISTS `sys_resource_module`;
CREATE TABLE `sys_resource_module` (
 `id` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `code` varchar(64) NOT NULL,
 `client` varchar(32) NOT NULL,
 `icon` varchar(255),
 `color` varchar(32),
 `sort` int NOT NULL,
 `status` varchar(32) NOT NULL,
 `description` text,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_resource_module
-- ----------------------------
INSERT INTO `sys_resource_module` VALUES ('210001', '管理端', 'admin', 'ADMIN', 'icon-park-outline:all-application', '#2080f0', 1, 'ENABLED', '管理端菜单与权限资源模块', '{}', '2026-06-30 00:00:00', NULL, '2026-06-30 00:00:00', NULL);

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
 `id` varchar(64) NOT NULL,
 `code` varchar(64) NOT NULL,
 `name` varchar(64) NOT NULL,
 `category` varchar(64) NOT NULL,
 `scope_type` varchar(32) NOT NULL,
 `owner_dept_id` varchar(64),
 `sort` int NOT NULL,
 `status` varchar(32) NOT NULL,
 `is_builtin` tinyint(1) NOT NULL,
 `description` text,
 `extra` json NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES ('1', 'SUPER_ADMIN', '超级管理员', 'SYS', 'PLATFORM', NULL, 1, 'ENABLED', 0, '系统内置超级管理员角色', '{}', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 11:56:13.747886', NULL);

-- ----------------------------
-- Table structure for sys_weak_password
-- ----------------------------
DROP TABLE IF EXISTS `sys_weak_password`;
CREATE TABLE `sys_weak_password` (
 `id` varchar(64) NOT NULL,
 `password` varchar(255) NOT NULL,
 `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `created_by` varchar(64),
 `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
 `updated_by` varchar(64)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of sys_weak_password
-- ----------------------------
INSERT INTO `sys_weak_password` VALUES ('7267772085910261393', '123456', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_weak_password` VALUES ('7661438436788304682', 'password', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_weak_password` VALUES ('7404805181363764417', 'admin123', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_weak_password` VALUES ('7142597855705676121', 'qwerty', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);
INSERT INTO `sys_weak_password` VALUES ('7411256677926569870', '111111', '2026-08-08 00:00:00', NULL, '2026-08-08 00:00:00', NULL);

-- ----------------------------
-- Indexes structure for table cg_test_activity
-- ----------------------------
CREATE INDEX `ix_cg_test_activity_owner_dept_id` ON `cg_test_activity` (
 `owner_dept_id`
);

-- ----------------------------
-- Primary Key structure for table cg_test_activity
-- ----------------------------
ALTER TABLE `cg_test_activity` ADD CONSTRAINT `pk_cg_test_activity` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table cg_test_catalog
-- ----------------------------
CREATE INDEX `ix_cg_test_catalog_owner_dept_id` ON `cg_test_catalog` (
 `owner_dept_id`
);

-- ----------------------------
-- Primary Key structure for table cg_test_catalog
-- ----------------------------
ALTER TABLE `cg_test_catalog` ADD CONSTRAINT `pk_cg_test_catalog` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table cg_test_knowledge_category
-- ----------------------------
CREATE INDEX `ix_cg_test_knowledge_category_owner_dept_id` ON `cg_test_knowledge_category` (
 `owner_dept_id`
);

-- ----------------------------
-- Primary Key structure for table cg_test_knowledge_category
-- ----------------------------
ALTER TABLE `cg_test_knowledge_category` ADD CONSTRAINT `pk_cg_test_knowledge_category` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table cg_test_knowledge_doc
-- ----------------------------
ALTER TABLE `cg_test_knowledge_doc` ADD CONSTRAINT `pk_cg_test_knowledge_doc` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table cg_test_order
-- ----------------------------
CREATE INDEX `ix_cg_test_order_owner_dept_id` ON `cg_test_order` (
 `owner_dept_id`
);

-- ----------------------------
-- Primary Key structure for table cg_test_order
-- ----------------------------
ALTER TABLE `cg_test_order` ADD CONSTRAINT `pk_cg_test_order` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table cg_test_order_item
-- ----------------------------
ALTER TABLE `cg_test_order_item` ADD CONSTRAINT `pk_cg_test_order_item` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table profile_user_admin
-- ----------------------------
ALTER TABLE `profile_user_admin` ADD CONSTRAINT `pk_profile_user_admin` PRIMARY KEY (`account_id`);

-- ----------------------------
-- Primary Key structure for table profile_user_portal
-- ----------------------------
ALTER TABLE `profile_user_portal` ADD CONSTRAINT `pk_profile_user_portal` PRIMARY KEY (`account_id`);

-- ----------------------------
-- Primary Key structure for table sys_account
-- ----------------------------
ALTER TABLE `sys_account` ADD CONSTRAINT `pk_sys_account` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_account_identity
-- ----------------------------
CREATE INDEX `idx_sys_account_identity_account_id` ON `sys_account_identity` (
 `account_id`
);

-- ----------------------------
-- Uniques structure for table sys_account_identity
-- ----------------------------
ALTER TABLE `sys_account_identity` ADD CONSTRAINT `uq_sys_account_identity_type_identifier` UNIQUE (`identity_type`, `identifier`);

-- ----------------------------
-- Primary Key structure for table sys_account_identity
-- ----------------------------
ALTER TABLE `sys_account_identity` ADD CONSTRAINT `pk_sys_account_identity` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_account_oauth_binding
-- ----------------------------
CREATE INDEX `idx_oauth_binding_account` ON `sys_account_oauth_binding` (
 `account_id`
);

-- ----------------------------
-- Uniques structure for table sys_account_oauth_binding
-- ----------------------------
ALTER TABLE `sys_account_oauth_binding` ADD CONSTRAINT `uq_oauth_provider_open_id` UNIQUE (`provider`, `open_id`);
ALTER TABLE `sys_account_oauth_binding` ADD CONSTRAINT `uq_oauth_account_provider` UNIQUE (`account_id`, `provider`);

-- ----------------------------
-- Primary Key structure for table sys_account_oauth_binding
-- ----------------------------
ALTER TABLE `sys_account_oauth_binding` ADD CONSTRAINT `sys_account_oauth_binding_pkey` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_account_password_history
-- ----------------------------
CREATE INDEX `idx_pwd_history_account_created` ON `sys_account_password_history` (
 `account_id`,
 `created_at`
);

-- ----------------------------
-- Primary Key structure for table sys_account_password_history
-- ----------------------------
ALTER TABLE `sys_account_password_history` ADD CONSTRAINT `pk_sys_account_password_history` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table sys_alert_log
-- ----------------------------
ALTER TABLE `sys_alert_log` ADD CONSTRAINT `pk_sys_alert_log` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_banner
-- ----------------------------
CREATE INDEX `ix_sys_banner_position_status_sort` ON `sys_banner` (
 `position`,
 `status`,
 `sort`
);

-- ----------------------------
-- Primary Key structure for table sys_banner
-- ----------------------------
ALTER TABLE `sys_banner` ADD CONSTRAINT `pk_sys_banner` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_client_module
-- ----------------------------
ALTER TABLE `sys_client_module` ADD CONSTRAINT `uq_sys_client_module_code` UNIQUE (`code`);

-- ----------------------------
-- Primary Key structure for table sys_client_module
-- ----------------------------
ALTER TABLE `sys_client_module` ADD CONSTRAINT `pk_sys_client_module` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_client_resource
-- ----------------------------
ALTER TABLE `sys_client_resource` ADD CONSTRAINT `uq_sys_client_resource_module_id_code` UNIQUE (`module_id`, `code`);

-- ----------------------------
-- Primary Key structure for table sys_client_resource
-- ----------------------------
ALTER TABLE `sys_client_resource` ADD CONSTRAINT `pk_sys_client_resource` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_codegen_field
-- ----------------------------
CREATE INDEX `ix_sys_codegen_field_plan_role_sort` ON `sys_codegen_field` (
 `plan_id`,
 `table_role`,
 `sort`
);

-- ----------------------------
-- Uniques structure for table sys_codegen_field
-- ----------------------------
ALTER TABLE `sys_codegen_field` ADD CONSTRAINT `uq_sys_codegen_field_plan_role_column` UNIQUE (`plan_id`, `table_role`, `column_name`);

-- ----------------------------
-- Primary Key structure for table sys_codegen_field
-- ----------------------------
ALTER TABLE `sys_codegen_field` ADD CONSTRAINT `pk_sys_codegen_field` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_codegen_plan
-- ----------------------------
CREATE INDEX `ix_sys_codegen_plan_gen_type` ON `sys_codegen_plan` (
 `gen_type`
);
CREATE INDEX `ix_sys_codegen_plan_main_table` ON `sys_codegen_plan` (
 `main_table`
);

-- ----------------------------
-- Uniques structure for table sys_codegen_plan
-- ----------------------------
ALTER TABLE `sys_codegen_plan` ADD CONSTRAINT `uq_sys_codegen_plan_name` UNIQUE (`name`);

-- ----------------------------
-- Primary Key structure for table sys_codegen_plan
-- ----------------------------
ALTER TABLE `sys_codegen_plan` ADD CONSTRAINT `pk_sys_codegen_plan` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_config
-- ----------------------------
CREATE INDEX `idx_sys_config_category` ON `sys_config` (
 `category`
);
CREATE INDEX `idx_sys_config_category_scope_scene` ON `sys_config` (
 `category`,
 `scope`,
 `scene`
);
CREATE UNIQUE INDEX `idx_sys_config_key` ON `sys_config` (
 `config_key`
);

-- ----------------------------
-- Primary Key structure for table sys_config
-- ----------------------------
ALTER TABLE `sys_config` ADD CONSTRAINT `pk_sys_config` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table sys_dept
-- ----------------------------
ALTER TABLE `sys_dept` ADD CONSTRAINT `pk_sys_dept` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_dict
-- ----------------------------
CREATE INDEX `idx_sys_dict_category` ON `sys_dict` (
 `category`
);
CREATE UNIQUE INDEX `idx_sys_dict_code` ON `sys_dict` (
 `code`
);
CREATE INDEX `idx_sys_dict_parent_id` ON `sys_dict` (
 `parent_id`
);

-- ----------------------------
-- Primary Key structure for table sys_dict
-- ----------------------------
ALTER TABLE `sys_dict` ADD CONSTRAINT `pk_sys_dict` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table sys_feedback
-- ----------------------------
ALTER TABLE `sys_feedback` ADD CONSTRAINT `pk_sys_feedback` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_file
-- ----------------------------
ALTER TABLE `sys_file` ADD CONSTRAINT `uq_sys_file_object_name` UNIQUE (`object_name`);

-- ----------------------------
-- Primary Key structure for table sys_file
-- ----------------------------
ALTER TABLE `sys_file` ADD CONSTRAINT `pk_sys_file` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_group
-- ----------------------------
ALTER TABLE `sys_group` ADD CONSTRAINT `uq_sys_group_name` UNIQUE (`name`);

-- ----------------------------
-- Primary Key structure for table sys_group
-- ----------------------------
ALTER TABLE `sys_group` ADD CONSTRAINT `pk_sys_group` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_iam_relation
-- ----------------------------
CREATE INDEX `ix_sys_iam_relation_account_type_relation` ON `sys_iam_relation` (
 `account_type`,
 `relation_type`
);
CREATE INDEX `ix_sys_iam_relation_subject` ON `sys_iam_relation` (
 `subject_type`,
 `subject_id`,
 `relation_type`
);
CREATE INDEX `ix_sys_iam_relation_target` ON `sys_iam_relation` (
 `target_type`,
 `target_id`,
 `target_key`
);

-- ----------------------------
-- Uniques structure for table sys_iam_relation
-- ----------------------------
ALTER TABLE `sys_iam_relation` ADD CONSTRAINT `uq_sys_iam_relation_subject_relation_target` UNIQUE (`subject_type`, `subject_id`, `relation_type`, `target_type`, `target_id`, `target_key`, `account_type`);

-- ----------------------------
-- Primary Key structure for table sys_iam_relation
-- ----------------------------
ALTER TABLE `sys_iam_relation` ADD CONSTRAINT `pk_sys_iam_relation` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_job_log
-- ----------------------------
CREATE INDEX `idx_sys_job_log_execute_time` ON `sys_job_log` (
 `execute_time`
);

-- ----------------------------
-- Indexes structure for table sys_notice
-- ----------------------------
CREATE INDEX `idx_sys_notice_status_kind_publish` ON `sys_notice` (
 `status`,
 `kind`,
 `publish_at`
);
CREATE INDEX `idx_sys_notice_status_pinned_publish` ON `sys_notice` (
 `status`,
 `is_pinned`,
 `publish_at` DESC
);

-- ----------------------------
-- Primary Key structure for table sys_notice
-- ----------------------------
ALTER TABLE `sys_notice` ADD CONSTRAINT `pk_sys_notice` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_notice_read
-- ----------------------------
CREATE INDEX `idx_sys_notice_read_account` ON `sys_notice_read` (
 `account_type`,
 `account_id`
);

-- ----------------------------
-- Uniques structure for table sys_notice_read
-- ----------------------------
ALTER TABLE `sys_notice_read` ADD CONSTRAINT `uq_sys_notice_read_account` UNIQUE (`notice_id`, `account_type`, `account_id`);

-- ----------------------------
-- Primary Key structure for table sys_notice_read
-- ----------------------------
ALTER TABLE `sys_notice_read` ADD CONSTRAINT `pk_sys_notice_read` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_operation_audit_log
-- ----------------------------
CREATE INDEX `idx_sys_operation_audit_account_id` ON `sys_operation_audit_log` (
 `account_id`
);
CREATE INDEX `idx_sys_operation_audit_created_at` ON `sys_operation_audit_log` (
 `created_at`
);
CREATE INDEX `idx_sys_operation_audit_module_action` ON `sys_operation_audit_log` (
 `module`,
 `action`
);
CREATE INDEX `idx_sys_operation_audit_resource` ON `sys_operation_audit_log` (
 `resource_type`,
 `resource_id`
);

-- ----------------------------
-- Primary Key structure for table sys_operation_audit_log
-- ----------------------------
ALTER TABLE `sys_operation_audit_log` ADD CONSTRAINT `pk_sys_operation_audit_log` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table sys_operation_audit_outbox
-- ----------------------------
ALTER TABLE `sys_operation_audit_outbox` ADD CONSTRAINT `pk_sys_operation_audit_outbox` PRIMARY KEY (`id`);

-- ----------------------------
-- Primary Key structure for table sys_position
-- ----------------------------
ALTER TABLE `sys_position` ADD CONSTRAINT `pk_sys_position` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_resource
-- ----------------------------
ALTER TABLE `sys_resource` ADD CONSTRAINT `uq_sys_resource_module_id_code` UNIQUE (`module_id`, `code`);

-- ----------------------------
-- Primary Key structure for table sys_resource
-- ----------------------------
ALTER TABLE `sys_resource` ADD CONSTRAINT `pk_sys_resource` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_resource_module
-- ----------------------------
ALTER TABLE `sys_resource_module` ADD CONSTRAINT `uq_sys_resource_module_code` UNIQUE (`code`);

-- ----------------------------
-- Primary Key structure for table sys_resource_module
-- ----------------------------
ALTER TABLE `sys_resource_module` ADD CONSTRAINT `pk_sys_resource_module` PRIMARY KEY (`id`);

-- ----------------------------
-- Uniques structure for table sys_role
-- ----------------------------
ALTER TABLE `sys_role` ADD CONSTRAINT `uq_sys_role_code` UNIQUE (`code`);

-- ----------------------------
-- Primary Key structure for table sys_role
-- ----------------------------
ALTER TABLE `sys_role` ADD CONSTRAINT `pk_sys_role` PRIMARY KEY (`id`);

-- ----------------------------
-- Indexes structure for table sys_weak_password
-- ----------------------------
CREATE UNIQUE INDEX `idx_sys_weak_password_password` ON `sys_weak_password` (
 `password`
);

-- ----------------------------
-- Primary Key structure for table sys_weak_password
-- ----------------------------
ALTER TABLE `sys_weak_password` ADD CONSTRAINT `pk_sys_weak_password` PRIMARY KEY (`id`);

SET FOREIGN_KEY_CHECKS = 1;
