/*
 Navicat Premium Dump SQL

 Source Server         : mysql-dev
 Source Server Type    : MySQL
 Source Server Version : 90600 (9.6.0)
 Source Host           : 127.0.0.1:3306
 Source Schema         : hei_gin

 Target Server Type    : MySQL
 Target Server Version : 90600 (9.6.0)
 File Encoding         : 65001

 Date: 23/08/2026 17:56:44
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for cg_test_activity
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_activity`;
CREATE TABLE `cg_test_activity`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '分类',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态',
  `cover_url` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '封面地址',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '描述说明',
  `start_at` datetime(6) NOT NULL COMMENT '开始时间',
  `end_at` datetime(6) NULL DEFAULT NULL COMMENT '结束时间',
  `max_participants` int NOT NULL COMMENT '最大参与人数',
  `price` decimal(10, 0) NOT NULL COMMENT '报名费用',
  `is_public` tinyint(1) NOT NULL COMMENT '是否公开：1 公开 / 0 不公开',
  `need_approval` tinyint(1) NOT NULL COMMENT '是否需要审批：1 需要 / 0 不需要',
  `rule_config` json NOT NULL COMMENT '规则配置（JSON）',
  `extra` json NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据范围）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `ix_cg_test_activity_owner_dept_id`(`owner_dept_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成测试-活动' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cg_test_activity
-- ----------------------------
INSERT INTO `cg_test_activity` VALUES ('900000000000000001', 'ACT-BOOTCAMP', '暑期训练营', 'TRAINING', 'OFFLINE', 'ENABLED', 'https://example.com/activity/bootcamp.png', '覆盖文本域、时间、金额、开关、JSON 的 CRUD 测试数据。', '2026-07-19 01:00:00.000000', '2026-07-19 09:00:00.000000', 120, 199, 0, 0, '{\"limit\": {\"daily\": 3}, \"checkin\": true}', '{\"tags\": [\"codegen\", \"crud\"]}', '2026-08-08 13:09:50.554189', '1', '2026-08-08 13:09:50.554189', '1', NULL);
INSERT INTO `cg_test_activity` VALUES ('900000000000000002', 'ACT-RD-01', '研发部活动', 'TRAINING', 'OFFLINE', 'ENABLED', NULL, '研发部 ALL 样例', '2026-08-01 09:00:00.000000', '2026-08-01 18:00:00.000000', 50, 0, 1, 0, '{}', '{}', '2026-08-23 08:00:00.000000', '8200000000000202', '2026-08-23 08:00:00.000000', '8200000000000202', '8200000000000102');
INSERT INTO `cg_test_activity` VALUES ('900000000000000003', 'ACT-MKT-01', '市场部活动', 'MARKETING', 'ONLINE', 'ENABLED', NULL, '市场部 ALL 样例', '2026-08-02 09:00:00.000000', '2026-08-02 18:00:00.000000', 100, 0, 1, 0, '{}', '{}', '2026-08-23 08:00:00.000000', '8200000000000202', '2026-08-23 08:00:00.000000', '8200000000000202', '8200000000000103');

-- ----------------------------
-- Table structure for cg_test_catalog
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_catalog`;
CREATE TABLE `cg_test_catalog`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `parent_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '父级ID',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '分类',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_visible` tinyint(1) NOT NULL COMMENT '是否可见：1 可见 / 0 隐藏',
  `icon` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标标识',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '描述说明',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据范围）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `ix_cg_test_catalog_owner_dept_id`(`owner_dept_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成测试-目录' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cg_test_catalog
-- ----------------------------
INSERT INTO `cg_test_catalog` VALUES ('900000000000000101', NULL, 'ROOT', '根目录', 'SYSTEM', 'ENABLED', 1, 0, 'folder', '一级节点', '{\"level\": 1}', '2026-08-08 13:09:50.667031', '1', '2026-08-08 13:09:50.667031', '1', NULL);
INSERT INTO `cg_test_catalog` VALUES ('900000000000000102', '900000000000000101', 'CHILD-A', '子目录A', 'SYSTEM', 'ENABLED', 10, 0, 'folder-open', '二级节点', '{\"level\": 2}', '2026-08-08 13:09:50.667031', '1', '2026-08-08 13:09:50.667031', '1', NULL);
INSERT INTO `cg_test_catalog` VALUES ('900000000000000103', '900000000000000101', 'CHILD-B', '子目录B', 'BUSINESS', 'DISABLED', 20, 0, 'folder-open', '二级节点', '{\"level\": 2}', '2026-08-08 13:09:50.667031', '1', '2026-08-08 13:09:50.667031', '1', NULL);
INSERT INTO `cg_test_catalog` VALUES ('900000000000000104', '900000000000000101', 'RD-CAT', '研发目录', 'SYSTEM', 'ENABLED', 11, 1, 'folder', '研发部目录', '{}', '2026-08-23 08:00:00.000000', '8200000000000203', '2026-08-23 08:00:00.000000', '8200000000000203', '8200000000000102');
INSERT INTO `cg_test_catalog` VALUES ('900000000000000105', '900000000000000101', 'MKT-CAT', '市场目录', 'BUSINESS', 'ENABLED', 21, 1, 'folder', '市场部目录', '{}', '2026-08-23 08:00:00.000000', '8200000000000203', '2026-08-23 08:00:00.000000', '8200000000000203', '8200000000000103');
INSERT INTO `cg_test_catalog` VALUES ('900000000000000106', '900000000000000101', 'BE-CAT', '后端目录', 'SYSTEM', 'ENABLED', 12, 1, 'folder', '后端组目录', '{}', '2026-08-23 08:00:00.000000', '8200000000000208', '2026-08-23 08:00:00.000000', '8200000000000208', '8200000000000105');

-- ----------------------------
-- Table structure for cg_test_knowledge_category
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_knowledge_category`;
CREATE TABLE `cg_test_knowledge_category`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `parent_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '父级ID',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_visible` tinyint(1) NOT NULL COMMENT '是否可见：1 可见 / 0 隐藏',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '描述说明',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据范围）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `ix_cg_test_knowledge_category_owner_dept_id`(`owner_dept_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成测试-知识分类' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cg_test_knowledge_category
-- ----------------------------
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000301', NULL, 'KB', '知识库', 'ENABLED', 1, 0, '根分类', '{\"level\": 1}', '2026-08-08 13:09:50.841963', '1', '2026-08-08 13:09:50.841963', '1', NULL);
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000302', '900000000000000301', 'KB-DEV', '研发文档', 'ENABLED', 10, 0, '研发相关文档', '{\"level\": 2}', '2026-08-08 13:09:50.841963', '1', '2026-08-08 13:09:50.841963', '1', NULL);
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000303', NULL, 'RD-KB', '研发知识库', 'ENABLED', 2, 1, '研发根分类', '{}', '2026-08-23 08:00:00.000000', '8200000000000205', '2026-08-23 08:00:00.000000', '8200000000000205', '8200000000000102');
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000304', '900000000000000303', 'RD-FE-KB', '前端知识库', 'ENABLED', 1, 1, '前端子分类', '{}', '2026-08-23 08:00:00.000000', '8200000000000205', '2026-08-23 08:00:00.000000', '8200000000000205', '8200000000000104');
INSERT INTO `cg_test_knowledge_category` VALUES ('900000000000000305', NULL, 'MKT-KB', '市场知识库', 'ENABLED', 3, 1, '市场根分类', '{}', '2026-08-23 08:00:00.000000', '8200000000000205', '2026-08-23 08:00:00.000000', '8200000000000205', '8200000000000103');

-- ----------------------------
-- Table structure for cg_test_knowledge_doc
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_knowledge_doc`;
CREATE TABLE `cg_test_knowledge_doc`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `category_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类ID',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `title` varchar(160) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态',
  `summary` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '摘要',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '内容',
  `author` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '作者',
  `published_at` datetime(6) NULL DEFAULT NULL COMMENT '发布时间',
  `view_count` int NOT NULL COMMENT '浏览/查看次数',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_top` tinyint(1) NOT NULL COMMENT '是否置顶：1 置顶 / 0 不置顶',
  `settings` json NOT NULL COMMENT '展示设置（JSON）',
  `extra` json NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成测试-知识文档' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cg_test_knowledge_doc
-- ----------------------------
INSERT INTO `cg_test_knowledge_doc` VALUES ('900000000000000311', '900000000000000302', 'DOC-CODEGEN-001', '代码生成测试文档', 'ARTICLE', 'PUBLISHED', '用于测试左树右表生成。', '正文内容用于触发 textarea。', 'tester', '2026-07-19 01:19:18.000000', 88, 1, 0, '{\"theme\": \"default\", \"showToc\": true}', '{\"tags\": [\"tree\", \"table\"]}', '2026-08-08 13:09:50.841963', '1', '2026-08-08 13:09:50.841963', '1');

-- ----------------------------
-- Table structure for cg_test_order
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_order`;
CREATE TABLE `cg_test_order`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `order_no` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单号',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `customer_name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '客户名称',
  `customer_phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '客户手机号',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '类型',
  `ordered_at` datetime(6) NOT NULL COMMENT '下单时间',
  `paid_at` datetime(6) NULL DEFAULT NULL COMMENT '支付时间',
  `total_amount` decimal(10, 0) NOT NULL COMMENT '订单金额',
  `item_count` int NOT NULL COMMENT '商品数量',
  `need_invoice` tinyint(1) NOT NULL COMMENT '是否开票：1 需要 / 0 不需要',
  `invoice_config` json NOT NULL COMMENT '发票配置（JSON）',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注说明',
  `extra` json NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据范围）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `ix_cg_test_order_owner_dept_id`(`owner_dept_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成测试-订单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cg_test_order
-- ----------------------------
INSERT INTO `cg_test_order` VALUES ('900000000000000201', 'CG-ORDER-001', '测试订单001', '张三', '13800000000', 'PAID', 'NORMAL', '2026-07-19 01:10:00.000000', '2026-07-19 01:20:00.000000', 399, 2, 0, '{\"taxNo\": \"91300000000000000X\", \"title\": \"张三\"}', '主子表生成测试订单', '{\"source\": \"codegen\"}', '2026-08-08 13:09:50.732542', '1', '2026-08-08 13:09:50.732542', '1', NULL);
INSERT INTO `cg_test_order` VALUES ('900000000000000202', 'CG-ORDER-RD', '研发订单', '张三', '13800002001', 'PAID', 'NORMAL', '2026-08-10 10:00:00.000000', NULL, 100, 1, 0, '{}', '研发部订单', '{}', '2026-08-23 08:00:00.000000', '8200000000000203', '2026-08-23 08:00:00.000000', '8200000000000203', '8200000000000102');
INSERT INTO `cg_test_order` VALUES ('900000000000000203', 'CG-ORDER-SELF', '本人订单', '李四', '13800002002', 'PAID', 'NORMAL', '2026-08-11 10:00:00.000000', NULL, 200, 1, 0, '{}', '仅本人可见', '{}', '2026-08-23 08:00:00.000000', '8200000000000204', '2026-08-23 08:00:00.000000', '8200000000000204', NULL);
INSERT INTO `cg_test_order` VALUES ('900000000000000204', 'CG-ORDER-OTHER', '他人订单', '王五', '13800002003', 'PAID', 'NORMAL', '2026-08-12 10:00:00.000000', NULL, 300, 1, 0, '{}', '他人创建', '{}', '2026-08-23 08:00:00.000000', '8200000000000202', '2026-08-23 08:00:00.000000', '8200000000000202', NULL);

-- ----------------------------
-- Table structure for cg_test_order_item
-- ----------------------------
DROP TABLE IF EXISTS `cg_test_order_item`;
CREATE TABLE `cg_test_order_item`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `order_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '订单ID',
  `sku_code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'SKU编码',
  `name` varchar(120) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '分类',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态',
  `quantity` int NOT NULL COMMENT '数量',
  `unit_price` decimal(10, 0) NOT NULL COMMENT '单价',
  `shipped_at` datetime(6) NULL DEFAULT NULL COMMENT '发货时间',
  `is_gift` tinyint(1) NOT NULL COMMENT '是否赠品：1 是 / 0 否',
  `item_config` json NOT NULL COMMENT '明细配置（JSON）',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注说明',
  `extra` json NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成测试-订单明细' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of cg_test_order_item
-- ----------------------------
INSERT INTO `cg_test_order_item` VALUES ('900000000000000211', '900000000000000201', 'SKU-001', '测试商品A', 'BOOK', 'ENABLED', 1, 199, NULL, 0, '{\"color\": \"red\"}', '普通明细', '{\"line\": 1}', '2026-08-08 13:09:50.732542', '1', '2026-08-08 13:09:50.732542', '1');
INSERT INTO `cg_test_order_item` VALUES ('900000000000000212', '900000000000000201', 'SKU-002', '测试商品B', 'COURSE', 'ENABLED', 1, 200, '2026-07-19 02:30:00.000000', 0, '{\"duration\": 30}', '赠品明细', '{\"line\": 2}', '2026-08-08 13:09:50.732542', '1', '2026-08-08 13:09:50.732542', '1');

-- ----------------------------
-- Table structure for profile_identity
-- ----------------------------
DROP TABLE IF EXISTS `profile_identity`;
CREATE TABLE `profile_identity`  (
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联系统账号ID（主键）',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'UNVERIFIED' COMMENT '认证状态：UNVERIFIED/PENDING/VERIFIED/REJECTED',
  `document_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '证件类型：ID_CARD/PASSPORT 等',
  `real_name_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '真实姓名密文（加密存储）',
  `document_no_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '证件号码密文（加密存储）',
  `document_no_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '证件号码哈希（用于脱敏检索）',
  `verify_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '认证通道：THIRD_PARTY/MANUAL',
  `provider` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方服务提供方',
  `provider_order_no` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方业务订单号',
  `verified_at` datetime(6) NULL DEFAULT NULL COMMENT '实名认证通过时间',
  `source_case_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '来源实名工单ID',
  `revoked_at` datetime(6) NULL DEFAULT NULL COMMENT '实名认证撤销时间',
  `revoked_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '撤销操作人账户ID',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`account_id`) USING BTREE,
  UNIQUE INDEX `uk_profile_identity_document_hash`(`document_no_hash` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '账号实名认证快照' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of profile_identity
-- ----------------------------

-- ----------------------------
-- Table structure for profile_user_admin
-- ----------------------------
DROP TABLE IF EXISTS `profile_user_admin`;
CREATE TABLE `profile_user_admin`  (
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联系统账号ID（主键）',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '管理端显示昵称',
  `avatar` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '头像 object_name 或 URL',
  `signature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '个性签名',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '绑定手机号',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '绑定邮箱',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '备注说明',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`account_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '管理端用户资料' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of profile_user_admin
-- ----------------------------
INSERT INTO `profile_user_admin` VALUES ('1', '超管', 'uploads/2026/08/09/02acc3dee5454d34913b07f49fe59cac.png', NULL, NULL, 'jiangbytebb@163.com', '系统内置超管账户', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 13:17:41.018249', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000201', 'IAM管理员', NULL, NULL, NULL, 'iam-admin@demo.local', 'IAM 管理员', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000202', '业务全量', NULL, NULL, NULL, 'biz-all@demo.local', '活动 ALL', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000203', '研发部员', NULL, NULL, NULL, 'biz-dept@demo.local', '目录 DEPT', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000204', '仅本人', NULL, NULL, NULL, 'biz-self@demo.local', '订单 SELF', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000205', '研发含子部门', NULL, NULL, NULL, 'biz-child@demo.local', '知识分类 DEPT_AND_CHILD', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000206', '只读账号', NULL, NULL, NULL, 'readonly@demo.local', '账号只读', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000207', '组授权', NULL, NULL, NULL, 'group-rd@demo.local', '用户组继承角色', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000208', '后端工程师', NULL, NULL, NULL, 'biz-be@demo.local', '后端组，目录 DEPT', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_admin` VALUES ('8200000000000209', '测试工程师', NULL, NULL, NULL, 'biz-qa@demo.local', '测试组，知识 DEPT_AND_CHILD', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for profile_user_portal
-- ----------------------------
DROP TABLE IF EXISTS `profile_user_portal`;
CREATE TABLE `profile_user_portal`  (
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联系统账号ID（主键）',
  `nickname` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '门户端显示昵称',
  `avatar` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '头像 object_name 或 URL',
  `signature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '个性签名',
  `phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '绑定手机号',
  `email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '绑定邮箱',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`account_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '门户用户资料' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of profile_user_portal
-- ----------------------------
INSERT INTO `profile_user_portal` VALUES ('7491847383584804864', 'user-171fd244', 'uploads/2026/08/09/85e1b98acfc9465abbbba86ef3b4fec8.jpg', NULL, NULL, 'jiangbyte@163.com', '2026-08-08 13:26:48.032837', NULL, '2026-08-08 13:48:45.931196', '7491847383584804864');
INSERT INTO `profile_user_portal` VALUES ('7491872891940786176', 'user-a527e592', NULL, NULL, '17286916074', '3317229064@qq.com', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `profile_user_portal` VALUES ('8200000000000211', 'Bob', NULL, NULL, '13800001001', 'bob@demo.local', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `profile_user_portal` VALUES ('8200000000000212', 'Alice', NULL, NULL, '13800001002', 'alice@demo.local', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for real_name_case
-- ----------------------------
DROP TABLE IF EXISTS `real_name_case`;
CREATE TABLE `real_name_case`  (
  `case_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '实名工单ID（主键）',
  `business_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '业务类型：ACCOUNT_VERIFY/ACCOUNT_RECOVERY',
  `verify_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '认证通道：THIRD_PARTY/MANUAL',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '工单状态：PENDING/APPROVED/REJECTED 等',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '账户ID',
  `target_account_hint_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '目标账户提示信息密文',
  `applicant_contact_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '申请人联系方式密文',
  `document_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '证件类型：ID_CARD/PASSPORT 等',
  `real_name_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '真实姓名密文',
  `document_no_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '证件号码密文',
  `document_no_hash` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '证件号码哈希',
  `attachment_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '附件ID列表（JSON数组）',
  `payload_cipher` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '扩展业务载荷密文',
  `handler_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '受理部门ID',
  `provider` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方服务提供方',
  `provider_order_no` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方业务订单号',
  `submitter_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '提交人账户ID',
  `reviewer_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '审核人账户ID',
  `reviewed_at` datetime(6) NULL DEFAULT NULL COMMENT '审核完成时间',
  `reject_reason` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '审核驳回原因',
  `expire_at` datetime(6) NULL DEFAULT NULL COMMENT '过期时间',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`case_id`) USING BTREE,
  INDEX `idx_real_name_case_account`(`account_id` ASC) USING BTREE,
  INDEX `idx_real_name_case_status`(`business_type` ASC, `status` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '实名业务工单' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of real_name_case
-- ----------------------------

-- ----------------------------
-- Table structure for real_name_case_record
-- ----------------------------
DROP TABLE IF EXISTS `real_name_case_record`;
CREATE TABLE `real_name_case_record`  (
  `record_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '工单流水ID（主键）',
  `case_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联实名工单ID',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '账户ID',
  `business_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '业务类型',
  `action` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '流水动作：SUBMIT/APPROVE/REJECT/REVOKE 等',
  `status_before` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '动作前工单状态',
  `status_after` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '动作后工单状态',
  `verify_channel` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '认证通道：THIRD_PARTY（三方）/ MANUAL（人工）',
  `provider` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方服务提供方',
  `operator_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作人账户ID',
  `dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作所属部门ID',
  `remark` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注说明',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  PRIMARY KEY (`record_id`) USING BTREE,
  INDEX `idx_real_name_case_record_case`(`case_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '实名业务工单流水' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of real_name_case_record
-- ----------------------------

-- ----------------------------
-- Table structure for sys_account
-- ----------------------------
DROP TABLE IF EXISTS `sys_account`;
CREATE TABLE `sys_account`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录密码哈希值（bcrypt 等）',
  `account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户类型：ADMIN（管理端）/ PORTAL（门户端）',
  `account_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户状态：ACTIVE（正常）/ LOCKED（锁定）/ CANCELLED（已注销）',
  `cancelled_at` datetime(6) NULL DEFAULT NULL COMMENT '账号注销完成时间',
  `cancelled_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '执行注销的操作人账户ID',
  `cancel_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '注销原因说明',
  `cancel_notify_email` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '注销前快照：通知邮箱（身份清理前保留）',
  `cancel_notify_phone` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '注销前快照：通知手机号（身份清理前保留）',
  `last_login_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '上一次成功登录 IP',
  `last_login_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '上一次成功登录地理位置',
  `last_login_time` datetime(6) NULL DEFAULT NULL COMMENT '上一次成功登录时间',
  `last_login_device` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '上一次成功登录设备标识',
  `latest_login_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '最近一次成功登录 IP',
  `latest_login_address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '最近一次成功登录地理位置',
  `latest_login_time` datetime(6) NULL DEFAULT NULL COMMENT '最近一次成功登录时间',
  `latest_login_device` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '最近一次成功登录设备标识',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统账号' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_account
-- ----------------------------
INSERT INTO `sys_account` VALUES ('1', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '127.0.0.1', NULL, '2026-08-23 08:43:29.340653', 'Desktop', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 11:56:13.747886', '1');
INSERT INTO `sys_account` VALUES ('7491847383584804864', '$2a$10$ZvgY90jMCQpobPlmaCXqie6rCzii8JEciVkXUVM.Kc2DkQHc639xy', 'PORTAL', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '127.0.0.1', NULL, '2026-08-23 09:34:13.894237', 'Desktop', '2026-08-08 13:26:48.032837', NULL, '2026-08-08 22:35:11.342559', '7491847383584804864');
INSERT INTO `sys_account` VALUES ('7491872891940786176', '$2b$12$kghdYhio.WATOZvDNdTLe.ACM2ibhvP.v88NudZcvjroz/H8M8.z.', 'PORTAL', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account` VALUES ('8200000000000201', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000202', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000203', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000204', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000205', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000206', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000207', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000208', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000209', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'ADMIN', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000211', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'PORTAL', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account` VALUES ('8200000000000212', '$2b$12$LLjMHXpuX5MuFwevZiU11OP4OQYVIDUKBBqNeK6VvDqfTWeqs8BJi', 'PORTAL', 'ENABLED', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_account_identity
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_identity`;
CREATE TABLE `sys_account_identity`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户ID',
  `identity_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '身份类型：USERNAME/EMAIL/PHONE 等',
  `identifier` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录标识值（用户名/邮箱/手机号）',
  `verified` tinyint(1) NOT NULL COMMENT '标识是否已完成验证：1 是 / 0 否',
  `is_primary` tinyint(1) NOT NULL COMMENT '是否主登录标识：1 主标识 / 0 次标识',
  `bind_status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '绑定状态：BOUND/UNBOUND/PENDING 等',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_account_identity_type_identifier`(`identity_type` ASC, `identifier` ASC) USING BTREE,
  INDEX `idx_sys_account_identity_account_id`(`account_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '账号身份标识（手机/邮箱/登录名）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_account_identity
-- ----------------------------
INSERT INTO `sys_account_identity` VALUES ('2086415885489872897', '1', 'ACCOUNT', 'superadmin', 0, 0, 'BOUND', '2026-08-09 11:34:45.504657', '1', '2026-08-09 11:34:45.504657', '1');
INSERT INTO `sys_account_identity` VALUES ('2086415885552787457', '1', 'EMAIL', 'jiangbytebb@163.com', 0, 0, 'BOUND', '2026-08-09 11:34:45.513351', '1', '2026-08-09 11:34:45.513351', '1');
INSERT INTO `sys_account_identity` VALUES ('2087538500522606594', '7491847383584804864', 'ACCOUNT', 'user', 0, 0, 'BOUND', '2026-08-12 13:55:37.773544', '1', '2026-08-12 13:55:37.773544', '1');
INSERT INTO `sys_account_identity` VALUES ('2087538500522606595', '7491847383584804864', 'EMAIL', 'jiangbyte@163.com', 0, 0, 'BOUND', '2026-08-12 13:55:37.784348', '1', '2026-08-12 13:55:37.784348', '1');
INSERT INTO `sys_account_identity` VALUES ('7491872891999506432', '7491872891940786176', 'ACCOUNT', 'usera', 0, 0, 'BOUND', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account_identity` VALUES ('7491872891999506433', '7491872891940786176', 'EMAIL', '3317229064@qq.com', 0, 0, 'BOUND', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account_identity` VALUES ('7491872891999506434', '7491872891940786176', 'PHONE', '17286916074', 0, 0, 'BOUND', '2026-08-08 15:08:09.699685', NULL, '2026-08-08 15:08:09.699685', NULL);
INSERT INTO `sys_account_identity` VALUES ('820000000000020101', '8200000000000201', 'ACCOUNT', 'admin_iam', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020102', '8200000000000201', 'EMAIL', 'iam-admin@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020201', '8200000000000202', 'ACCOUNT', 'admin_all', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020202', '8200000000000202', 'EMAIL', 'biz-all@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020301', '8200000000000203', 'ACCOUNT', 'admin_dept', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020302', '8200000000000203', 'EMAIL', 'biz-dept@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020401', '8200000000000204', 'ACCOUNT', 'admin_self', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020402', '8200000000000204', 'EMAIL', 'biz-self@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020501', '8200000000000205', 'ACCOUNT', 'admin_child', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020502', '8200000000000205', 'EMAIL', 'biz-child@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020601', '8200000000000206', 'ACCOUNT', 'admin_readonly', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020602', '8200000000000206', 'EMAIL', 'readonly@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020701', '8200000000000207', 'ACCOUNT', 'admin_group', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020702', '8200000000000207', 'EMAIL', 'group-rd@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020801', '8200000000000208', 'ACCOUNT', 'admin_be', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020802', '8200000000000208', 'EMAIL', 'biz-be@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020901', '8200000000000209', 'ACCOUNT', 'admin_qa', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000020902', '8200000000000209', 'EMAIL', 'biz-qa@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000021101', '8200000000000211', 'ACCOUNT', 'portal_bob', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000021102', '8200000000000211', 'EMAIL', 'bob@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000021201', '8200000000000212', 'ACCOUNT', 'portal_alice', 1, 1, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_account_identity` VALUES ('820000000000021202', '8200000000000212', 'EMAIL', 'alice@demo.local', 1, 0, 'BOUND', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_account_oauth_binding
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_oauth_binding`;
CREATE TABLE `sys_account_oauth_binding`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户ID',
  `provider` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'OAuth 提供方：wechat/github/google 等',
  `open_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '第三方平台 OpenID',
  `union_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方平台 UnionID（跨应用统一标识）',
  `nickname` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '第三方账号昵称快照',
  `avatar` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '第三方账号头像 URL 快照',
  `raw_profile` json NOT NULL COMMENT '第三方原始资料（JSON）',
  `bound_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '与本系统账号绑定时间',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_oauth_provider_open_id`(`provider` ASC, `open_id` ASC) USING BTREE,
  UNIQUE INDEX `uq_oauth_account_provider`(`account_id` ASC, `provider` ASC) USING BTREE,
  INDEX `idx_oauth_binding_account`(`account_id` ASC) USING BTREE,
  INDEX `idx_oauth_binding_union`(`union_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '三方登录绑定' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_account_oauth_binding
-- ----------------------------

-- ----------------------------
-- Table structure for sys_account_password_history
-- ----------------------------
DROP TABLE IF EXISTS `sys_account_password_history`;
CREATE TABLE `sys_account_password_history`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户ID',
  `password_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码哈希值（不可逆）',
  `changed_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '密码变更操作人（账户ID 或 system）',
  `change_reason` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '变更原因：register/admin_reset/self_reset/password_expired',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '密码写入历史时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_pwd_history_account_created`(`account_id` ASC, `created_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '账号密码历史' ROW_FORMAT = Dynamic;

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
CREATE TABLE `sys_alert_log`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `rule_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '触发告警的规则名称',
  `severity` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '严重级别：INFO/WARNING/CRITICAL',
  `summary` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '告警摘要（展示用）',
  `details` json NULL COMMENT '告警详情上下文（JSON）',
  `notified_via` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '通知渠道：email/webhook 等',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '告警产生/通知时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统告警日志' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_alert_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_banner
-- ----------------------------
DROP TABLE IF EXISTS `sys_banner`;
CREATE TABLE `sys_banner`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `image` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Banner 图片 object_name（由服务层解析访问 URL）',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '点击跳转链接地址',
  `link_type` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '链接类型（字典 BANNER_LINK_TYPE）',
  `summary` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '摘要',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '描述说明',
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Banner 分类（字典 BANNER_CATEGORY）',
  `type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Banner 类型（字典 BANNER_TYPE）',
  `position` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '展示位置（字典 BANNER_POSITION）',
  `target_account_types` json NOT NULL COMMENT '可见账户类型：ADMIN/PORTAL（JSON 数组）',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `interaction_count` bigint NOT NULL COMMENT '用户交互次数统计',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Banner 状态：ENABLED/DISABLED 等',
  `start_at` datetime(6) NULL DEFAULT NULL COMMENT '开始展示时间',
  `end_at` datetime(6) NULL DEFAULT NULL COMMENT '结束展示时间',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `ix_sys_banner_position_status_sort`(`position` ASC, `status` ASC, `sort` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'Banner 轮播' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_banner
-- ----------------------------
INSERT INTO `sys_banner` VALUES ('7491889345134235648', '最新公告速递', 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1600&h=700&fit=crop', '/announcements', 'ROUTE', '及时了解平台动态与重要通知。', '门户首页轮播示例图二', 'HOME', 'CAROUSEL', 'HOME_TOP', '[\"PORTAL\"]', 20, 0, 'ENABLED', NULL, NULL, '2026-08-08 16:13:32.393714', NULL, '2026-08-08 16:13:32.393714', NULL);
INSERT INTO `sys_banner` VALUES ('7491889345142624256', '完善个人资料', 'https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=1600&h=700&fit=crop', '/usercenter', 'ROUTE', '前往个人中心完善资料，获得更好的使用体验。', '门户首页轮播示例图三', 'HOME', 'CAROUSEL', 'HOME_TOP', '[\"PORTAL\"]', 30, 0, 'ENABLED', NULL, NULL, '2026-08-08 16:13:32.393714', NULL, '2026-08-08 16:13:32.393714', NULL);
INSERT INTO `sys_banner` VALUES ('7491889345146818560', '管理端运营位示例', 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=1600&h=700&fit=crop', NULL, 'NONE', '仅面向管理端账号类型，不在门户轮播出现。', '管理端展示图示例', 'ADMIN_DASHBOARD', 'CARD', 'ADMIN_TOP', '[\"ADMIN\"]', 10, 0, 'ENABLED', NULL, NULL, '2026-08-08 16:13:32.393714', NULL, '2026-08-08 16:13:32.393714', NULL);
INSERT INTO `sys_banner` VALUES ('8300000000000101', '门户首页欢迎', 'https://images.unsplash.com/photo-1497366216548-37526070297c?w=1600&h=700&fit=crop', '/', 'ROUTE', '欢迎使用 HEI 门户', '门户首页顶部轮播', 'HOME', 'CAROUSEL', 'HOME_TOP', '[\"PORTAL\"]', 5, 0, 'ENABLED', NULL, NULL, '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_banner` VALUES ('8300000000000102', '活动中心', 'https://images.unsplash.com/photo-1451187580459-43490279c0fa?w=1600&h=700&fit=crop', '/activities', 'ROUTE', '查看最新活动与报名', '门户活动推广位', 'HOME', 'CAROUSEL', 'HOME_TOP', '[\"PORTAL\"]', 15, 0, 'ENABLED', NULL, NULL, '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_banner` VALUES ('8300000000000103', '安全合规提示', 'https://images.unsplash.com/photo-1563986768609-322da13575f3?w=1600&h=700&fit=crop', NULL, 'NONE', '请妥善保管账号密码', '门户安全宣传卡片', 'HOME', 'CARD', 'HOME_MIDDLE', '[\"PORTAL\"]', 25, 0, 'ENABLED', NULL, NULL, '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_banner` VALUES ('8300000000000104', '管理端工作台', 'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=1600&h=700&fit=crop', '/workspace', 'ROUTE', '快捷进入工作台', '管理端顶部运营位', 'ADMIN_DASHBOARD', 'CAROUSEL', 'ADMIN_TOP', '[\"ADMIN\"]', 5, 0, 'ENABLED', NULL, NULL, '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_banner` VALUES ('8300000000000105', 'IAM 权限指引', 'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=1600&h=700&fit=crop', '/iam/account', 'ROUTE', '组织权限配置入口', '管理端 IAM 引导', 'ADMIN_DASHBOARD', 'CARD', 'ADMIN_SIDE', '[\"ADMIN\"]', 15, 0, 'ENABLED', NULL, NULL, '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_banner` VALUES ('8300000000000106', '双端同步公告', 'https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=1600&h=700&fit=crop', '/sys/notice', 'ROUTE', '管理端与门户同步消息', '双端可见展示图', 'HOME', 'BANNER', 'GLOBAL_TOP', '[\"ADMIN\", \"PORTAL\"]', 20, 0, 'ENABLED', NULL, NULL, '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_client_module
-- ----------------------------
DROP TABLE IF EXISTS `sys_client_module`;
CREATE TABLE `sys_client_module`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '适用账户体系：ADMIN/PORTAL',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标标识',
  `color` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '颜色值',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模块状态：ENABLED/DISABLED',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '客户端模块描述',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_client_module_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '客户端模块' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_client_module
-- ----------------------------
INSERT INTO `sys_client_module` VALUES ('221001', '管理端默认模块', 'admin-default', 'ADMIN', 'icon-park-outline:application-one', NULL, 1, 'ENABLED', '管理端默认客户端模块', '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_client_module` VALUES ('221002', '门户端默认模块', 'portal-default', 'PORTAL', 'icon-park-outline:application-one', NULL, 1, 'ENABLED', '门户端默认客户端模块', '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);

-- ----------------------------
-- Table structure for sys_client_resource
-- ----------------------------
DROP TABLE IF EXISTS `sys_client_resource`;
CREATE TABLE `sys_client_resource`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `parent_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '父级客户端资源ID',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `resource_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型：MENU/BUTTON/API 等',
  `module_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属客户端模块ID',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端路由路径',
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端组件路径',
  `redirect` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '路由重定向地址',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标标识',
  `color` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '颜色值',
  `href` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '外链跳转地址',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_visible` tinyint(1) NOT NULL COMMENT '是否可见：1 可见 / 0 隐藏',
  `is_cache` tinyint(1) NOT NULL COMMENT '是否缓存路由：1 缓存 / 0 不缓存',
  `is_affix` tinyint(1) NOT NULL COMMENT '是否固定标签页：1 固定 / 0 不固定',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源状态：ENABLED/DISABLED',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '客户端资源描述',
  `layout` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '页面布局类型',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_client_resource_module_id_code`(`module_id` ASC, `code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '客户端资源' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_client_resource
-- ----------------------------
INSERT INTO `sys_client_resource` VALUES ('222001', NULL, 'home', '首页', 'MENU', '221001', '/home', '/home/index.vue', NULL, 'icon-park-outline:home', NULL, NULL, 1, 0, 0, 0, 'ENABLED', '管理端客户端样例菜单', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_client_resource` VALUES ('222002', NULL, 'home', '首页', 'MENU', '221002', '/home', '/home/index.vue', NULL, 'icon-park-outline:home', NULL, NULL, 1, 0, 0, 0, 'ENABLED', '门户端客户端样例菜单', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);

-- ----------------------------
-- Table structure for sys_codegen_field
-- ----------------------------
DROP TABLE IF EXISTS `sys_codegen_field`;
CREATE TABLE `sys_codegen_field`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `plan_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属代码生成方案ID',
  `table_role` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '表角色：MASTER/SUB 等',
  `column_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据库列名',
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字段展示标签（通常取自表注释）',
  `db_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据库物理类型',
  `value_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '语义值类型：str/int/bool/datetime/dict 等',
  `ui_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '前端 UI 类型：string/number/boolean 等',
  `widget` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '表单控件类型',
  `dict_code` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '关联数据字典编码',
  `query_operator` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '列表查询运算符：eq/like/between 等',
  `in_table` tinyint(1) NOT NULL COMMENT '是否在表格列展示：1 是 / 0 否',
  `in_form` tinyint(1) NOT NULL COMMENT '是否在表单展示：1 是 / 0 否',
  `in_detail` tinyint(1) NOT NULL COMMENT '是否在详情展示：1 是 / 0 否',
  `in_query` tinyint(1) NOT NULL COMMENT '是否作为查询条件：1 是 / 0 否',
  `primary_key` tinyint(1) NOT NULL COMMENT '是否主键列：1 是 / 0 否',
  `required` tinyint(1) NOT NULL COMMENT '是否必填：1 是 / 0 否',
  `unique_flag` tinyint(1) NOT NULL COMMENT '是否唯一：1 是 / 0 否',
  `nullable` tinyint(1) NOT NULL COMMENT '是否允许为空：1 可空 / 0 非空',
  `max_length` int NULL DEFAULT NULL COMMENT '字段最大长度限制',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_codegen_field_plan_role_column`(`plan_id` ASC, `table_role` ASC, `column_name` ASC) USING BTREE,
  INDEX `ix_sys_codegen_field_plan_role_sort`(`plan_id` ASC, `table_role` ASC, `sort` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成字段配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_codegen_field
-- ----------------------------

-- ----------------------------
-- Table structure for sys_codegen_plan
-- ----------------------------
DROP TABLE IF EXISTS `sys_codegen_plan`;
CREATE TABLE `sys_codegen_plan`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `gen_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '生成类型：CRUD/TREE/SUB_TABLE 等',
  `author` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '作者',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '代码生成方案描述',
  `table_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主表数据库名',
  `pk_column` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主表主键列名',
  `entity_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '生成的主实体类名',
  `module_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '后端模块包路径',
  `business_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主业务中文名',
  `api_prefix` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'REST API 路径前缀',
  `permission_prefix` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限标识前缀',
  `resource_module_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '挂载的资源模块ID',
  `parent_resource_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '挂载的父菜单资源ID',
  `menu_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '生成菜单名称',
  `menu_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '生成菜单路由路径',
  `component_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '生成前端组件路径',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标标识',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `tree_parent_field` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '树表父级字段名',
  `tree_label_field` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '树节点展示字段名',
  `sub_table` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '子表数据库名',
  `sub_pk` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '子表主键列名',
  `sub_foreign_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '子表外键列名',
  `sub_entity_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '子实体类名',
  `sub_business_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '子业务中文名',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_codegen_plan_name`(`name` ASC) USING BTREE,
  INDEX `ix_sys_codegen_plan_gen_type`(`gen_type` ASC) USING BTREE,
  INDEX `ix_sys_codegen_plan_table_name`(`table_name` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '代码生成方案' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_codegen_plan
-- ----------------------------

-- ----------------------------
-- Table structure for sys_config
-- ----------------------------
DROP TABLE IF EXISTS `sys_config`;
CREATE TABLE `sys_config`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `config_key` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '配置项唯一键',
  `config_value` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '配置项值（按 value_type 解析）',
  `category` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置分类/分组',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '备注说明',
  `sort_code` int NOT NULL COMMENT '同分类下排序码',
  `value_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '值类型：STRING/JSON/BOOL/NUMBER',
  `label` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '配置项展示名称',
  `scope` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '作用域账户类型：GLOBAL/ADMIN/PORTAL',
  `scene` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '业务场景编码',
  `is_builtin` tinyint(1) NOT NULL COMMENT '是否内置配置：1 内置不可删 / 0 可维护',
  `ext_json` json NOT NULL COMMENT '扩展配置（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_config_key`(`config_key` ASC) USING BTREE,
  INDEX `idx_sys_config_category`(`category` ASC) USING BTREE,
  INDEX `idx_sys_config_category_scope_scene`(`category` ASC, `scope` ASC, `scene` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统动态配置' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_config
-- ----------------------------
INSERT INTO `sys_config` VALUES ('7107236720876061841', 'AUTH_REGISTER_ADMIN_ENABLED', 'FALSE', 'AUTH_REGISTER', 'ADMIN 开放注册', 1, 'BOOL', 'ADMIN 开放注册', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7110705537461490506', 'AUTH_REGISTER_PORTAL_DEFAULT_DEPT_ID', '', 'AUTH_REGISTER', 'PORTAL 注册默认部门', 10, 'STRING', 'PORTAL 注册默认部门', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7111122815413902706', 'MAIL_TENCENT_FROM_EMAIL', '', 'MAIL', '腾讯云发件邮箱', 32, 'STRING', '腾讯云发件邮箱', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7112014435091478940', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_CLIENT_ID', '', 'AUTH_OAUTH', '门户微信开放平台 AppId', 32, 'STRING', 'AppId', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.454644', NULL, '2026-08-12 15:57:59.454644', '1');
INSERT INTO `sys_config` VALUES ('7115927954405030347', 'AUTH_OAUTH_ADMIN_GITHUB_ENABLED', 'FALSE', 'AUTH_OAUTH', '管理端 GitHub 登录', 101, 'BOOL', '启用', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.476901', NULL, '2026-08-12 15:57:59.476901', '1');
INSERT INTO `sys_config` VALUES ('7116199988527416451', 'AUTH_PASSWORD_RESET_URL_ADMIN', 'http://localhost:5173/auth/forgot-password', 'AUTH_TOKEN', 'ADMIN 密码重置页完整 URL', 3, 'STRING', 'ADMIN 密码重置页完整 URL', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7120119863717980260', 'AUTH_LOGIN_PORTAL_LOCK_SECONDS', '300', 'AUTH_LOGIN', 'PORTAL 锁定时间（秒）', 20, 'INT', 'PORTAL 锁定时间（秒）', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7121863095719606243', 'AUTH_REGISTER_ADMIN_DEFAULT_DEPT_ID', '', 'AUTH_REGISTER', 'ADMIN 注册默认部门', 5, 'STRING', 'ADMIN 注册默认部门', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7123144141092942155', 'AUTH_OAUTH_PORTAL_GITEE_CLIENT_SECRET', '', 'AUTH_OAUTH', '门户 Gitee ClientSecret', 13, 'STRING', 'Client Secret', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.427549', NULL, '2026-08-12 15:57:59.427549', NULL);
INSERT INTO `sys_config` VALUES ('7125526233449487784', 'AUTH_OAUTH_ADMIN_GITHUB_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端 GitHub 回调', 104, 'STRING', 'Redirect URI', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.487901', NULL, '2026-08-12 15:57:59.487901', '1');
INSERT INTO `sys_config` VALUES ('7125885979202993561', 'AUTH_OAUTH_PORTAL_GITEE_REDIRECT_URI', '', 'AUTH_OAUTH', '门户 Gitee 回调', 14, 'STRING', 'Redirect URI', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.431808', NULL, '2026-08-12 15:57:59.431808', '1');
INSERT INTO `sys_config` VALUES ('7127821217542267803', 'AUTH_REGISTER_ADMIN_DEFAULT_ROLE_ID', '', 'AUTH_REGISTER', 'ADMIN 注册默认角色', 4, 'STRING', 'ADMIN 注册默认角色', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7128087998001730719', 'PASSWORD_MIN_LENGTH', '8', 'AUTH_PASSWORD', '密码最小长度', 10, 'INT', '密码最小长度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7131883348553291227', 'PUSH_LARK_SECRET', '', 'PUSH', '飞书加签密钥', 21, 'STRING', '飞书加签密钥', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7137492378552518589', 'AUTH_DEFAULT_PASSWORD', '', 'AUTH_PASSWORD', '新建账户默认密码', 1, 'STRING', '新建账户默认密码', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7137737529370775030', 'AUTH_FORCE_BIND_ADMIN_PHONE', 'FALSE', 'AUTH_FORCE_BIND', 'ADMIN 强制绑定手机', 4, 'BOOL', '强制绑定手机', 'ADMIN', NULL, 0, '{}', '2026-08-12 14:36:35.255718', NULL, '2026-08-12 14:36:35.255718', NULL);
INSERT INTO `sys_config` VALUES ('7140248531257076308', 'AUTH_LOGIN_ACCOUNT_MAX_FAILURES', '5', 'AUTH_LOGIN', '单账号最大登录失败次数', 2, 'INT', '单账号最大登录失败次数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7141145699218587845', 'STORAGE_UPLOAD_MAX_BYTES', '10485760', 'UPLOAD', '上传文件大小上限（字节）', 1, 'INT', '上传文件大小上限（字节）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7143068116163573568', 'AUTH_OAUTH_PORTAL_WECHAT_MP_APP_ID', '', 'AUTH_OAUTH', '门户小程序 AppId', 42, 'STRING', 'AppId', 'PORTAL', 'WECHAT_MP', 0, '{}', '2026-08-12 15:57:59.468946', NULL, '2026-08-12 15:57:59.468946', '1');
INSERT INTO `sys_config` VALUES ('7143992105573126791', 'MAIL_TENCENT_SECRET_ID', '', 'MAIL', '腾讯云邮件 SecretId', 30, 'STRING', '腾讯云邮件 SecretId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7147573270237877567', 'AUTH_FORCE_BIND_PORTAL_EMAIL', 'FALSE', 'AUTH_FORCE_BIND', 'PORTAL 强制绑定邮箱', 1, 'BOOL', '强制绑定邮箱', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.239972', NULL, '2026-08-12 14:36:35.239972', '1');
INSERT INTO `sys_config` VALUES ('7148106912230055136', 'PASSWORD_COMPLEXITY', 'DIGITS_UPPER_LOWER_SPECIAL', 'AUTH_PASSWORD', '密码复杂度', 12, 'STRING', '密码复杂度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7149430200646252861', 'AUTH_OAUTH_PORTAL_QQ_CLIENT_ID', '', 'AUTH_OAUTH', '门户 QQ ClientId', 22, 'STRING', 'Client ID', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.439704', NULL, '2026-08-12 15:57:59.439704', '1');
INSERT INTO `sys_config` VALUES ('7152439802015666741', 'AUTH_LOGIN_PORTAL_EMAIL_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'PORTAL 邮箱无用户策略', 24, 'STRING', 'PORTAL 邮箱无用户策略', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7154674032244377737', 'AUTH_LOGIN_ADMIN_ALLOW_EMAIL', 'TRUE', 'AUTH_LOGIN', 'ADMIN 允许邮箱登录', 15, 'BOOL', 'ADMIN 允许邮箱登录', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7156018694924792840', 'SMS_ALIYUN_SIGN_NAME', '', 'SMS', '阿里云短信签名', 12, 'STRING', '阿里云短信签名', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7159166510497302996', 'STORAGE_UPLOAD_CATEGORY_MAX_LENGTH', '64', 'UPLOAD', '上传分类名最大长度', 7, 'INT', '上传分类名最大长度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7159687456156883406', 'PASSWORD_FORBID_WEAK_LIST', 'TRUE', 'AUTH_PASSWORD', '禁止弱密码库命中', 17, 'BOOL', '禁止弱密码库命中', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7163571438509411369', 'AUTH_OAUTH_ADMIN_GITEE_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端 Gitee ClientSecret', 113, 'STRING', 'Client Secret', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.499737', NULL, '2026-08-12 15:57:59.499737', NULL);
INSERT INTO `sys_config` VALUES ('7165407484130097877', 'AUTH_LOGIN_PORTAL_MAX_FAILURES', '5', 'AUTH_LOGIN', 'PORTAL 最大失败次数', 19, 'INT', 'PORTAL 最大失败次数', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7168944809004050220', 'SMS_TENCENT_SECRET_ID', '', 'SMS', '腾讯云短信 SecretId', 20, 'STRING', '腾讯云短信 SecretId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7172687233882333188', 'AUDIT_ALERT_RULE_SENSITIVE_OPS', 'TRUE', 'AUDIT_ALERT', '敏感操作监控', 12, 'BOOL', '敏感操作监控', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7175913351351122695', 'STORAGE_RUSTFS_REGION', 'us-east-1', 'STORAGE', 'RustFS Region', 44, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7179891447845803919', 'AUTH_OAUTH_ADMIN_GITEE_ENABLED', 'TRUE', 'AUTH_OAUTH', '管理端 Gitee 登录', 111, 'BOOL', '启用', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.491641', NULL, '2026-08-12 15:57:59.491641', '1');
INSERT INTO `sys_config` VALUES ('7182623752582531304', 'AUDIT_ALERT_WEBHOOK_SECRET', '', 'AUDIT_ALERT', 'Webhook 签名密钥', 6, 'STRING', 'Webhook 签名密钥', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7189190642894026372', 'AUTH_LOGIN_ADMIN_ALLOW_OTP', 'TRUE', 'AUTH_LOGIN', 'ADMIN 允许 OTP 登录', 17, 'BOOL', 'ADMIN 允许 OTP 登录', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7191040278958713174', 'SMS_TEMPLATE_BIND_PHONE_CODE', '{\"code\": \"\", \"content\": \"绑定验证码 {{code}}，{{expire_minutes}} 分钟内有效\"}', 'SMS_TEMPLATE', '绑定手机验证码短信模板', 20, 'JSON', '绑定手机验证码短信模板', NULL, 'BIND_PHONE_CODE', 0, '{}', '2026-08-12 14:36:35.274517', NULL, '2026-08-12 14:36:35.274517', NULL);
INSERT INTO `sys_config` VALUES ('7192344291042963470', 'AUTH_OAUTH_ADMIN_GITHUB_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端 GitHub ClientSecret', 103, 'STRING', 'Client Secret', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.483943', NULL, '2026-08-12 15:57:59.483943', NULL);
INSERT INTO `sys_config` VALUES ('7196439739936279174', 'AUTH_REGISTER_PORTAL_ALLOW_PHONE', 'TRUE', 'AUTH_REGISTER', 'PORTAL 允许手机注册', 13, 'BOOL', '允许手机注册', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.237264', NULL, '2026-08-12 14:36:35.237264', '1');
INSERT INTO `sys_config` VALUES ('7203040505764618021', 'AUTH_OAUTH_ADMIN_GITHUB_CLIENT_ID', '', 'AUTH_OAUTH', '管理端 GitHub ClientId', 102, 'STRING', 'Client ID', 'ADMIN', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.480274', NULL, '2026-08-12 15:57:59.480274', '1');
INSERT INTO `sys_config` VALUES ('7204221884041359188', 'AUDIT_ALERT_NOTIFY_CUSTOM_WEBHOOK', 'FALSE', 'AUDIT_ALERT', '自定义 Webhook 通知', 4, 'BOOL', '自定义 Webhook 通知', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7214189881608199926', 'AUDIT_ALERT_RULE_BRUTE_FORCE', 'TRUE', 'AUDIT_ALERT', '暴力破解检测', 10, 'BOOL', '暴力破解检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7219205172797610985', 'AUTH_LOGIN_PORTAL_ALLOW_OTP', 'TRUE', 'AUTH_LOGIN', 'PORTAL 允许 OTP 登录', 25, 'BOOL', 'PORTAL 允许 OTP 登录', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7223837583572900295', 'SMS_TEMPLATE_CHANGE_PASSWORD_CODE', '{\"code\": \"\", \"content\": \"改密验证码 {{code}}\"}', 'SMS_TEMPLATE', '修改密码短信模板', 2, 'JSON', '修改密码短信模板', NULL, 'CHANGE_PASSWORD_CODE', 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7224651178645144555', 'AUTH_LOGIN_LOCK_SECONDS', '900', 'AUTH_LOGIN', '登录锁定时间（秒）', 4, 'INT', '登录锁定时间（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7231060999533820222', 'MAIL_TEMPLATE_LOGIN_CODE', '{\"subject\": \"{{app_name}} 登录验证码\", \"body\": \"您的登录验证码是 {{code}}，{{expire_minutes}} 分钟内有效。\"}', 'MAIL_TEMPLATE', '登录验证码邮件模板', 2, 'JSON', '登录验证码邮件模板', NULL, 'LOGIN_CODE', 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7231113269178162861', 'AUTH_REGISTER_PORTAL_ALLOW_ACCOUNT', 'TRUE', 'AUTH_REGISTER', 'PORTAL 允许用户名注册', 11, 'BOOL', '允许账号注册', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.228596', NULL, '2026-08-12 14:36:35.228596', '1');
INSERT INTO `sys_config` VALUES ('7236608062114386965', 'STORAGE_UPLOAD_ALLOWED_CONTENT_TYPES', '[\"image/jpeg\",\"image/png\",\"image/webp\",\"application/pdf\",\"text/plain\",\"application/octet-stream\"]', 'UPLOAD', '允许的 MIME 类型列表（JSON 数组）', 4, 'JSON', '允许的 MIME 类型列表（JSON 数组）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7247260109501730845', 'STORAGE_RUSTFS_ACCESS_KEY', 'gAAAAABqeFtjtOr0lzjyRl8TPLDP5be0LWRCOizoe7P6RlqHaxasJvILz2P27NauLfjSKM71tWtwpBPZiltAP2aT5Zi5Jzr1lA', 'STORAGE', 'RustFS Access Key', 42, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 10:49:34.896059', '1');
INSERT INTO `sys_config` VALUES ('7254222471397278270', 'MAIL_LOCAL_FROM_NAME', 'hei-fastapi', 'MAIL', '发件人显示名称', 15, 'STRING', '发件人显示名称', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7259830455098789920', 'AUTH_REGISTER_PORTAL_ENABLED', 'TRUE', 'AUTH_REGISTER', 'PORTAL 开放注册', 6, 'BOOL', 'PORTAL 开放注册', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7260823451582098683', 'AUDIT_ALERT_RULE_UNUSUAL_HOURS', 'TRUE', 'AUDIT_ALERT', '异常时间操作检测', 11, 'BOOL', '异常时间操作检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7264846649644584871', 'STORAGE_ALIYUN_BUCKET_PUBLIC', 'FALSE', 'STORAGE', '阿里云桶是否公开', 14, 'BOOL', NULL, NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7277178859125997111', 'AUTH_OAUTH_ADMIN_GITEE_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端 Gitee 回调', 114, 'STRING', 'Redirect URI', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.503669', NULL, '2026-08-12 15:57:59.503669', '1');
INSERT INTO `sys_config` VALUES ('7279087550089826228', 'MAIL_LOCAL_USERNAME', '', 'MAIL', 'SMTP 用户名', 12, 'STRING', 'SMTP 用户名', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7280596390016858524', 'STORAGE_UPLOAD_ALLOWED_EXTENSIONS', '[\".jpg\",\".jpeg\",\".png\",\".webp\",\".pdf\",\".txt\",\".ini\",\".xlsx\"]', 'UPLOAD', '允许的文件扩展名列表（JSON 数组）', 5, 'JSON', '允许的文件扩展名列表（JSON 数组）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7285230910149567688', 'AUTH_REGISTER_PORTAL_ALLOW_EMAIL', 'TRUE', 'AUTH_REGISTER', 'PORTAL 允许邮箱注册', 12, 'BOOL', '允许邮箱注册', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.234457', NULL, '2026-08-12 14:36:35.234457', '1');
INSERT INTO `sys_config` VALUES ('7286296755297320815', 'AUTH_LOGIN_PORTAL_FAILURE_WINDOW_SECONDS', '300', 'AUTH_LOGIN', 'PORTAL 登录失败窗口（秒）', 18, 'INT', 'PORTAL 登录失败窗口（秒）', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7296603632888365225', 'AUDIT_ALERT_WEBHOOK_URL', '', 'AUDIT_ALERT', 'Webhook 地址', 5, 'STRING', 'Webhook 地址', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7297229307635338316', 'MAIL_ALIYUN_ACCOUNT_NAME', '', 'MAIL', '阿里云发信地址', 22, 'STRING', '阿里云发信地址', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7310793430810554097', 'PUSH_WECHAT_WORK_WEBHOOK', '', 'PUSH', '企业微信 Webhook', 30, 'STRING', '企业微信 Webhook', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7312697174587677983', 'MAIL_LOCAL_USE_SSL', 'FALSE', 'MAIL', 'SMTP 使用 SSL', 17, 'BOOL', 'SMTP 使用 SSL', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7316084524318600148', 'MAIL_TEMPLATE_ACCOUNT_CANCELLED', '{\"subject\": \"{{app_name}} 账号注销确认\", \"body\": \"您好，您的账号已申请注销。\\n\\n我们将在 {{retention_days}} 天内保留账号数据；到期且期间未再登录使用后，系统将彻底删除账号及相关数据。\\n\\n预计清理时间：{{purge_at}}\\n如非本人操作，请尽快联系管理员。\"}', 'MAIL_TEMPLATE', '账号注销确认邮件模板', 20, 'JSON', '账号注销确认邮件模板', NULL, 'ACCOUNT_CANCELLED', 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7317566250917157196', 'MAIL_TEMPLATE_CHANGE_PASSWORD_CODE', '{\"subject\": \"{{app_name}} 修改密码验证码\", \"body\": \"验证码 {{code}}，{{expire_minutes}} 分钟内有效。\"}', 'MAIL_TEMPLATE', '修改密码邮件模板', 3, 'JSON', '修改密码邮件模板', NULL, 'CHANGE_PASSWORD_CODE', 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7317628235066225421', 'AUTH_FORCE_BIND_PORTAL_PHONE', 'FALSE', 'AUTH_FORCE_BIND', 'PORTAL 强制绑定手机', 2, 'BOOL', '强制绑定手机', 'PORTAL', NULL, 0, '{}', '2026-08-12 14:36:35.245505', NULL, '2026-08-12 14:36:35.245505', '1');
INSERT INTO `sys_config` VALUES ('7319284294792948722', 'AUTH_OAUTH_PORTAL_GITHUB_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户 GitHub 登录', 1, 'BOOL', '启用', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.397127', NULL, '2026-08-12 15:57:59.397127', '1');
INSERT INTO `sys_config` VALUES ('7325262355928690941', 'STORAGE_UPLOAD_DENIED_EXTENSIONS', '[\".exe\",\".bat\",\".cmd\",\".sh\",\".js\",\".html\",\".php\",\".py\",\".jar\"]', 'UPLOAD', '禁止上传的扩展名列表（JSON 数组）', 6, 'JSON', '禁止上传的扩展名列表（JSON 数组）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7331249362091683364', 'MAIL_TENCENT_SECRET_KEY', '', 'MAIL', '腾讯云邮件 SecretKey', 31, 'STRING', '腾讯云邮件 SecretKey', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7334700335238107691', 'STORAGE_RUSTFS_BUCKET', 'defaultbucket', 'STORAGE', 'RustFS 存储桶', 40, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7337967471264658438', 'AUDIT_ALERT_ALERT_COOLDOWN_SECONDS', '1800', 'AUDIT_ALERT', '告警冷却(秒)', 8, 'INT', '告警冷却(秒)', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7345028346771097677', 'PASSWORD_MAX_LENGTH', '128', 'AUTH_PASSWORD', '密码最大长度', 11, 'INT', '密码最大长度', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7346432777171636629', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户微信网页登录', 31, 'BOOL', '启用', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.450854', NULL, '2026-08-12 15:57:59.450854', '1');
INSERT INTO `sys_config` VALUES ('7350757244774876723', 'DEFAULT_MESSAGE_PUSH_ENGINE', 'DINGTALK', 'PUSH', '默认消息推送引擎', 1, 'STRING', '默认消息推送引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7362777511165276641', 'STORAGE_TENCENT_BUCKET_PUBLIC', 'FALSE', 'STORAGE', '腾讯云桶是否公开', 20, 'BOOL', NULL, NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7371141348954395191', 'AUTH_LOGIN_ADMIN_MAX_FAILURES', '5', 'AUTH_LOGIN', 'ADMIN 最大失败次数', 11, 'INT', 'ADMIN 最大失败次数', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7374162031239208376', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端微信开放平台 Secret', 133, 'STRING', 'AppSecret', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.533100', NULL, '2026-08-12 15:57:59.533100', NULL);
INSERT INTO `sys_config` VALUES ('7377151897664869996', 'AUDIT_ALERT_NOTIFY_PUSH', 'TRUE', 'AUDIT_ALERT', '推送通知', 3, 'BOOL', '推送通知', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7384500216009438994', 'PASSWORD_CUSTOM_WEAK_WORDS', '', 'AUTH_PASSWORD', '自定义弱密码词（逗号分隔）', 20, 'STRING', '自定义弱密码词（逗号分隔）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7386527519961798374', 'AUTH_TOKEN_TTL_SECONDS', '2592000', 'AUTH_TOKEN', 'Token 过期时间（秒），默认 30 天', 1, 'INT', 'Token 过期时间（秒），默认 30 天', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7387934960432348080', 'MAIL_TEMPLATE_RESET_PASSWORD_CODE', '{\"subject\": \"{{app_name}} 密码重置\", \"body\": \"请点击以下链接重置密码，该链接将在 {{expire_minutes}} 分钟内有效。\\n\\n{{reset_link}}\"}', 'MAIL_TEMPLATE', '重置密码邮件模板', 1, 'JSON', '重置密码邮件模板', NULL, 'RESET_PASSWORD_CODE', 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7388788716511880208', 'AUDIT_ALERT_NOTIFY_EMAIL_TO', '', 'AUDIT_ALERT', '审计告警收件邮箱', 2, 'STRING', '告警收件邮箱', NULL, NULL, 0, '{}', '2026-08-12 14:10:48.638877', NULL, '2026-08-12 14:10:48.638877', NULL);
INSERT INTO `sys_config` VALUES ('7389279063701289685', 'AUTH_REGISTER_ADMIN_REQUIRE_PHONE', 'FALSE', 'AUTH_REGISTER', 'ADMIN 注册要求手机号', 2, 'BOOL', 'ADMIN 注册要求手机号', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7395157445765927617', 'SMS_TENCENT_SECRET_KEY', '', 'SMS', '腾讯云短信 SecretKey', 21, 'STRING', '腾讯云短信 SecretKey', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7402885978532192492', 'SMS_TEMPLATE_ACCOUNT_PURGED', '{\"code\": \"\", \"content\": \"您的账号已完成注销清理并彻底删除。\"}', 'SMS_TEMPLATE', '账号彻底删除短信模板', 21, 'JSON', '账号彻底删除短信模板', NULL, 'ACCOUNT_PURGED', 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7403276917562314417', 'SMS_TENCENT_REGION', '', 'SMS', '腾讯云短信区域', 90, 'STRING', '腾讯云短信区域', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7407734083761585458', 'AUTH_OAUTH_PORTAL_GITEE_ENABLED', 'TRUE', 'AUTH_OAUTH', '门户 Gitee 登录', 11, 'BOOL', '启用', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.416643', NULL, '2026-08-12 15:57:59.416643', '1');
INSERT INTO `sys_config` VALUES ('7408452631572683423', 'AUTH_OAUTH_PORTAL_WECHAT_MP_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户微信小程序登录', 41, 'BOOL', '启用', 'PORTAL', 'WECHAT_MP', 0, '{}', '2026-08-12 15:57:59.464887', NULL, '2026-08-12 15:57:59.464887', '1');
INSERT INTO `sys_config` VALUES ('7424366524857713971', 'DEFAULT_FILE_ENGINE', 'RUSTFS', 'STORAGE', '默认文件引擎', 1, 'STRING', '默认文件引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 17:19:31.423597', '1');
INSERT INTO `sys_config` VALUES ('7436928202175303081', 'PASSWORD_VALIDITY_DAYS', '90', 'AUTH_PASSWORD', '密码有效期（天）', 18, 'INT', '密码有效期（天）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7440617985015094217', 'MAIL_TEMPLATE_ACCOUNT_PURGED', '{\"subject\": \"{{app_name}} 账号已彻底删除\", \"body\": \"您好，您此前注销的账号已完成保留期清理，账号及相关个人数据已彻底删除。\\n\\n清理时间：{{purged_at}}\\n感谢您曾使用 {{app_name}}。\"}', 'MAIL_TEMPLATE', '账号彻底删除邮件模板', 21, 'JSON', '账号彻底删除邮件模板', NULL, 'ACCOUNT_PURGED', 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7440731969037919022', 'MAIL_LOCAL_FROM_EMAIL', 'test@hei-fastapi.local', 'MAIL', '发件人邮箱', 14, 'STRING', '发件人邮箱', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7462877140739332791', 'MAIL_TEMPLATE_REGISTER_SUCCESS', '{\"subject\": \"欢迎注册 {{app_name}}\", \"body\": \"账号 {{account}} 注册成功。\"}', 'MAIL_TEMPLATE', '注册成功邮件模板', 4, 'JSON', '注册成功邮件模板', NULL, 'REGISTER_SUCCESS', 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7486049791614620352', 'AUDIT_ALERT_IP_ANOMALY_THRESHOLD', '3', 'AUDIT_ALERT', 'IP异常阈值', 22, 'INT', 'IP异常阈值', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7490352982435195735', 'AUTH_REGISTER_PORTAL_DEFAULT_ROLE_ID', '', 'AUTH_REGISTER', 'PORTAL 注册默认角色', 9, 'STRING', 'PORTAL 注册默认角色', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7491104580893387509', 'SMS_TEMPLATE_LOGIN_CODE', '{\"code\": \"\", \"content\": \"登录验证码 {{code}}\"}', 'SMS_TEMPLATE', '登录验证码短信模板', 1, 'JSON', '登录验证码短信模板', NULL, 'LOGIN_CODE', 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7491869125220999168', 'STORAGE_ALIYUN_ENDPOINT', 'oss-cn-hangzhou.aliyuncs.com', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999169', 'STORAGE_ALIYUN_BUCKET', 'defaultbucket', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999170', 'STORAGE_ALIYUN_REGION', 'cn-hangzhou', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125220999171', 'STORAGE_ALIYUN_USE_SSL', 'TRUE', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193472', 'STORAGE_ALIYUN_BASE_URL', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193474', 'STORAGE_TENCENT_ENDPOINT', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193475', 'STORAGE_TENCENT_BUCKET', 'defaultbucket', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193476', 'STORAGE_TENCENT_REGION', 'ap-beijing', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193477', 'STORAGE_TENCENT_USE_SSL', 'TRUE', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193478', 'STORAGE_TENCENT_BASE_URL', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193480', 'STORAGE_MINIO_ACCESS_KEY', 'gAAAAABqd2UjQzg7UyUYFbmdQe6DHLXzJI7dO2Ql7IH_dmCaWvHkCfsmkFOfhGyXG_Q3kAWCsYEoWWZz0kaqgn9XSHGZDoZB8Q==', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 17:19:31.423597', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193481', 'STORAGE_MINIO_SECRET_KEY', 'gAAAAABqd2Uj8lcg6wI7znRSZmoU7WRxiPFY1ZQ9nU5O2kIyQHwbp0xPbRkrP7ww153nsWr-szThKe07RGmeicdbkHLFl4eSmA==', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 17:19:31.423597', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193482', 'STORAGE_MINIO_ENDPOINT', 'http://127.0.0.1:9000', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:50.574307', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193483', 'STORAGE_MINIO_BUCKET', 'vms', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193484', 'STORAGE_MINIO_REGION', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193485', 'STORAGE_MINIO_USE_SSL', 'FALSE', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7491869125225193486', 'STORAGE_MINIO_BASE_URL', '', 'STORAGE', NULL, 0, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-08 14:53:11.633342', '1', '2026-08-08 14:53:11.633342', '1');
INSERT INTO `sys_config` VALUES ('7499926412524933665', 'AUTH_OAUTH_FRONTEND_CALLBACK_PORTAL', 'http://localhost:5174/auth/oauth/callback', 'AUTH_OAUTH', '门户 OAuth 前端回调页（空则用默认）', 200, 'STRING', '门户前端回调', NULL, NULL, 0, '{}', '2026-08-12 15:57:59.540869', NULL, '2026-08-16 10:27:33.694774', '1');
INSERT INTO `sys_config` VALUES ('7507807560036605420', 'STORAGE_RUSTFS_BUCKET_PUBLIC', 'FALSE', 'STORAGE', 'RustFS 桶是否公开', 47, 'BOOL', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7509581770976973374', 'SMS_TENCENT_SDK_APP_ID', '', 'SMS', '腾讯云短信 SdkAppId', 22, 'STRING', '腾讯云短信 SdkAppId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7525309778220488671', 'STORAGE_MINIO_BUCKET_PUBLIC', 'FALSE', 'STORAGE', 'MinIO 桶是否公开', 26, 'BOOL', NULL, NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7528718642174231668', 'AUTH_LOGIN_FAILURE_WINDOW_SECONDS', '900', 'AUTH_LOGIN', '登录失败统计窗口（秒）', 1, 'INT', '登录失败统计窗口（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7535593666554670318', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_REDIRECT_URI', '', 'AUTH_OAUTH', '门户微信开放平台回调', 34, 'STRING', 'Redirect URI', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.461291', NULL, '2026-08-12 15:57:59.461291', '1');
INSERT INTO `sys_config` VALUES ('7545090679203476556', 'AUDIT_ALERT_ANALYSIS_INTERVAL_SECONDS', '60', 'AUDIT_ALERT', '分析周期(秒)', 7, 'INT', '分析周期(秒)', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7553572003476288661', 'AUTH_LOGIN_PORTAL_PHONE_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'PORTAL 手机号无用户策略', 22, 'STRING', 'PORTAL 手机号无用户策略', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7557049262663315054', 'MAIL_LOCAL_HOST', 'localhost', 'MAIL', 'SMTP 服务器地址', 10, 'STRING', 'SMTP 服务器地址', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7564883844348625200', 'AUDIT_ALERT_RULE_IP_ANOMALY', 'TRUE', 'AUDIT_ALERT', 'IP 异常检测', 14, 'BOOL', 'IP 异常检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7566333830794456163', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_ENABLED', 'FALSE', 'AUTH_OAUTH', '管理端微信网页登录', 131, 'BOOL', '启用', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.525073', NULL, '2026-08-12 15:57:59.525073', '1');
INSERT INTO `sys_config` VALUES ('7569398987799485767', 'PUSH_DINGTALK_WEBHOOK', '', 'PUSH', '钉钉 Webhook', 10, 'STRING', '钉钉 Webhook', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7577237128270583039', 'SMS_ALIYUN_ACCESS_KEY_SECRET', '', 'SMS', '阿里云短信 AccessKeySecret', 11, 'STRING', '阿里云短信 AccessKeySecret', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7585813464464377276', 'AUTH_FORCE_BIND_ADMIN_EMAIL', 'FALSE', 'AUTH_FORCE_BIND', 'ADMIN 强制绑定邮箱', 3, 'BOOL', '强制绑定邮箱', 'ADMIN', NULL, 0, '{}', '2026-08-12 14:36:35.252523', NULL, '2026-08-12 14:36:35.252523', NULL);
INSERT INTO `sys_config` VALUES ('7607784543080994141', 'SMS_ALIYUN_ACCESS_KEY_ID', '', 'SMS', '阿里云短信 AccessKeyId', 10, 'STRING', '阿里云短信 AccessKeyId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7616410287911469795', 'AUTH_REGISTER_PORTAL_REQUIRE_PHONE', 'FALSE', 'AUTH_REGISTER', 'PORTAL 注册要求手机号', 7, 'BOOL', 'PORTAL 注册要求手机号', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7618985031080401037', 'STORAGE_RUSTFS_USE_SSL', 'FALSE', 'STORAGE', 'RustFS 是否 SSL', 45, 'BOOL', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7634546496447633922', 'STORAGE_PRESIGN_EXPIRE_SECONDS', '3600', 'UPLOAD', '预签名 URL 有效期（秒）', 3, 'INT', '预签名 URL 有效期（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7639638206762559428', 'AUTH_PASSWORD_RESET_URL_PORTAL', 'http://localhost:5174/auth/forgot-password', 'AUTH_TOKEN', 'PORTAL 密码重置页完整 URL', 4, 'STRING', 'PORTAL 密码重置页完整 URL', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7640258913173401645', 'AUDIT_ALERT_ENABLED', 'TRUE', 'AUDIT_ALERT', '审计告警总开关', 1, 'BOOL', '审计告警总开关', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7641214881822118472', 'AUTH_LOGIN_ADMIN_PHONE_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'ADMIN 手机号无用户策略', 14, 'STRING', 'ADMIN 手机号无用户策略', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7642474443917569970', 'MAIL_ALIYUN_ACCESS_KEY_ID', '', 'MAIL', '阿里云邮件 AccessKeyId', 20, 'STRING', '阿里云邮件 AccessKeyId', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7643758115191934361', 'AUTH_LOGIN_PORTAL_ALLOW_EMAIL', 'TRUE', 'AUTH_LOGIN', 'PORTAL 允许邮箱登录', 23, 'BOOL', 'PORTAL 允许邮箱登录', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7645847924523510728', 'AUTH_REGISTER_PORTAL_REQUIRE_EMAIL', 'TRUE', 'AUTH_REGISTER', 'PORTAL 注册要求邮箱', 8, 'BOOL', 'PORTAL 注册要求邮箱', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7650241926616235479', 'COPYRIGHT_URL', '', 'SYS', '版权链接', 2, 'STRING', '版权链接', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7667316617799595990', 'AUDIT_ALERT_RULE_BULK_DELETE', 'TRUE', 'AUDIT_ALERT', '批量删除检测', 13, 'BOOL', '批量删除检测', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7680761152547899501', 'AUTH_REGISTER_ADMIN_REQUIRE_EMAIL', 'FALSE', 'AUTH_REGISTER', 'ADMIN 注册要求邮箱', 3, 'BOOL', 'ADMIN 注册要求邮箱', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7691150031767579943', 'AUTH_OAUTH_PORTAL_QQ_ENABLED', 'FALSE', 'AUTH_OAUTH', '门户 QQ 登录', 21, 'BOOL', '启用', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.435873', NULL, '2026-08-12 15:57:59.435873', '1');
INSERT INTO `sys_config` VALUES ('7695727871068877789', 'AUDIT_ALERT_NOTIFY_EMAIL', 'TRUE', 'AUDIT_ALERT', '邮件通知', 2, 'BOOL', '邮件通知', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7707546743105654318', 'AUTH_LOGIN_IP_MAX_FAILURES', '30', 'AUTH_LOGIN', '单 IP 最大登录失败次数', 3, 'INT', '单 IP 最大登录失败次数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7709767367597654277', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_CLIENT_ID', '', 'AUTH_OAUTH', '管理端微信开放平台 AppId', 132, 'STRING', 'AppId', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.529053', NULL, '2026-08-12 15:57:59.529053', '1');
INSERT INTO `sys_config` VALUES ('7712189083486710611', 'AUTH_OAUTH_PORTAL_WECHAT_MP_APP_SECRET', '', 'AUTH_OAUTH', '门户小程序 AppSecret', 43, 'STRING', 'AppSecret', 'PORTAL', 'WECHAT_MP', 0, '{}', '2026-08-12 15:57:59.472798', NULL, '2026-08-12 15:57:59.472798', NULL);
INSERT INTO `sys_config` VALUES ('7717547996356753442', 'PUSH_LARK_WEBHOOK', '', 'PUSH', '飞书 Webhook', 20, 'STRING', '飞书 Webhook', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7725294497454302580', 'AUTH_OAUTH_PORTAL_WECHAT_OPEN_CLIENT_SECRET', '', 'AUTH_OAUTH', '门户微信开放平台 Secret', 33, 'STRING', 'AppSecret', 'PORTAL', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.458015', NULL, '2026-08-12 15:57:59.458015', NULL);
INSERT INTO `sys_config` VALUES ('7734399746267501494', 'AUTH_OAUTH_PORTAL_QQ_CLIENT_SECRET', '', 'AUTH_OAUTH', '门户 QQ ClientSecret', 23, 'STRING', 'Client Secret', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.443762', NULL, '2026-08-12 15:57:59.443762', NULL);
INSERT INTO `sys_config` VALUES ('7743096612050880755', 'PASSWORD_FORBID_USER_INFO', 'TRUE', 'AUTH_PASSWORD', '禁止包含用户信息', 14, 'BOOL', '禁止包含用户信息', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7744551651175282801', 'AUTH_OAUTH_PORTAL_GITEE_CLIENT_ID', '', 'AUTH_OAUTH', '门户 Gitee ClientId', 12, 'STRING', 'Client ID', 'PORTAL', 'GITEE', 0, '{}', '2026-08-12 15:57:59.421381', NULL, '2026-08-12 15:57:59.421381', '1');
INSERT INTO `sys_config` VALUES ('7748406326636979011', 'STORAGE_RUSTFS_BASE_URL', '', 'STORAGE', 'RustFS 自定义基础 URL', 46, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7755198709300048806', 'SMS_TENCENT_SIGN_NAME', '', 'SMS', '腾讯云短信签名', 23, 'STRING', '腾讯云短信签名', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7774327802304558750', 'PUSH_DINGTALK_SECRET', '', 'PUSH', '钉钉加签密钥', 11, 'STRING', '钉钉加签密钥', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7780925862931981131', 'MAIL_LOCAL_AUTH_REQUIRED', 'FALSE', 'MAIL', 'SMTP 是否需要认证', 16, 'BOOL', 'SMTP 是否需要认证', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7791661612876719423', 'PASSWORD_HISTORY_CHECK_COUNT', '5', 'AUTH_PASSWORD', '历史密码检查条数', 16, 'INT', '历史密码检查条数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7799955389648191999', 'PASSWORD_EXPIRY_WARNING_DAYS', '7', 'AUTH_PASSWORD', '密码过期提前提醒（天）', 19, 'INT', '密码过期提前提醒（天）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7800328101940188056', 'AUTH_OAUTH_PORTAL_GITHUB_CLIENT_ID', 'superadmin', 'AUTH_OAUTH', '门户 GitHub ClientId', 2, 'STRING', 'Client ID', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.403535', NULL, '2026-08-12 15:57:59.403535', '1');
INSERT INTO `sys_config` VALUES ('7805678264902064100', 'AUTH_OAUTH_PORTAL_QQ_REDIRECT_URI', '', 'AUTH_OAUTH', '门户 QQ 回调', 24, 'STRING', 'Redirect URI', 'PORTAL', 'QQ', 0, '{}', '2026-08-12 15:57:59.447608', NULL, '2026-08-12 15:57:59.447608', '1');
INSERT INTO `sys_config` VALUES ('7812345678901234501', 'APP_NAME', 'HEI', 'SYS', '应用名称', 3, 'STRING', '邮件/短信模板中的应用名称', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234502', 'AUTH_OTP_TTL_SECONDS', '300', 'AUTH_TOKEN', '验证码有效期（秒）', 5, 'INT', '登录/注册/绑定等 OTP 有效期', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234503', 'AUDIT_ALERT_BULK_DELETE_WINDOW_SECONDS', '300', 'AUDIT_ALERT', '批量删除检测窗口（秒）', 16, 'INT', '批量删除检测窗口（秒）', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234504', 'AUDIT_ALERT_IP_ANOMALY_WINDOW_SECONDS', '900', 'AUDIT_ALERT', '异地 IP 检测窗口（秒）', 17, 'INT', '异地 IP 检测窗口（秒）', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234510', 'MAIL_TEMPLATE_RESET_PASSWORD_SUCCESS', '{\"subject\": \"{{app_name}} 密码已重置\", \"body\": \"您好，账号 {{account}} 的密码已成功重置。如非本人操作，请尽快联系管理员。\"}', 'MAIL_TEMPLATE', '重置密码成功邮件模板', 5, 'JSON', '重置密码成功邮件模板', NULL, 'RESET_PASSWORD_SUCCESS', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234511', 'MAIL_TEMPLATE_PASSWORD_EXPIRING', '{\"subject\": \"{{app_name}} 密码即将过期\", \"body\": \"您好，账号 {{account}} 的密码将在 {{remaining_days}} 天后过期，请尽快修改密码。\"}', 'MAIL_TEMPLATE', '密码即将过期邮件模板', 6, 'JSON', '密码即将过期邮件模板', NULL, 'PASSWORD_EXPIRING', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234512', 'MAIL_TEMPLATE_CHANGE_EMAIL_CODE', '{\"subject\": \"{{app_name}} 修改邮箱验证码\", \"body\": \"您的修改邮箱验证码是 {{code}}，{{expire_minutes}} 分钟内有效。\"}', 'MAIL_TEMPLATE', '修改邮箱验证码邮件模板', 21, 'JSON', '修改邮箱验证码邮件模板', NULL, 'CHANGE_EMAIL_CODE', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234520', 'SMS_TEMPLATE_REGISTER_SUCCESS', '{\"code\": \"\", \"content\": \"欢迎注册 {{app_name}}，账号 {{account}} 注册成功\"}', 'SMS_TEMPLATE', '注册成功短信模板', 3, 'JSON', '注册成功短信模板', NULL, 'REGISTER_SUCCESS', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234521', 'SMS_TEMPLATE_RESET_PASSWORD_CODE', '{\"code\": \"\", \"content\": \"重置密码验证码 {{code}}，{{expire_minutes}} 分钟内有效\"}', 'SMS_TEMPLATE', '重置密码验证码短信模板', 4, 'JSON', '重置密码验证码短信模板', NULL, 'RESET_PASSWORD_CODE', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234522', 'SMS_TEMPLATE_RESET_PASSWORD_SUCCESS', '{\"code\": \"\", \"content\": \"账号 {{account}} 密码已重置，如非本人操作请联系管理员\"}', 'SMS_TEMPLATE', '重置密码成功短信模板', 5, 'JSON', '重置密码成功短信模板', NULL, 'RESET_PASSWORD_SUCCESS', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234523', 'SMS_TEMPLATE_PASSWORD_EXPIRING', '{\"code\": \"\", \"content\": \"账号 {{account}} 密码将在 {{remaining_days}} 天后过期，请尽快修改\"}', 'SMS_TEMPLATE', '密码即将过期短信模板', 6, 'JSON', '密码即将过期短信模板', NULL, 'PASSWORD_EXPIRING', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234524', 'SMS_TEMPLATE_CHANGE_PHONE_CODE', '{\"code\": \"\", \"content\": \"修改手机验证码 {{code}}，{{expire_minutes}} 分钟内有效\"}', 'SMS_TEMPLATE', '修改手机验证码短信模板', 22, 'JSON', '修改手机验证码短信模板', NULL, 'CHANGE_PHONE_CODE', 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234551', 'SITE_ICP_NUMBER', '', 'SYS', 'ICP 备案号', 4, 'STRING', 'ICP 备案号', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234552', 'SITE_ICP_URL', 'https://beian.miit.gov.cn/', 'SYS', 'ICP 备案链接', 5, 'STRING', 'ICP 备案查询链接', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234553', 'SITE_PSB_NUMBER', '', 'SYS', '公安备案号', 6, 'STRING', '公安备案号', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7812345678901234554', 'SITE_PSB_URL', '', 'SYS', '公安备案链接', 7, 'STRING', '公安备案查询链接', NULL, NULL, 0, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7814530538364155449', 'AUTH_OAUTH_FRONTEND_CALLBACK_ADMIN', 'http://localhost:5173/auth/oauth/callback', 'AUTH_OAUTH', '管理端 OAuth 前端回调页（空则用默认）', 201, 'STRING', '管理端前端回调', NULL, NULL, 0, '{}', '2026-08-12 15:57:59.544738', NULL, '2026-08-16 10:27:33.715381', '1');
INSERT INTO `sys_config` VALUES ('7815227162417732636', 'MAIL_ALIYUN_ACCESS_KEY_SECRET', '', 'MAIL', '阿里云邮件 AccessKeySecret', 21, 'STRING', '阿里云邮件 AccessKeySecret', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7817231488316128319', 'DEFAULT_SMS_ENGINE', 'ALIYUN', 'SMS', '默认短信引擎', 1, 'STRING', '默认短信引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7820162979363132264', 'AUTH_OAUTH_ADMIN_GITEE_CLIENT_ID', '', 'AUTH_OAUTH', '管理端 Gitee ClientId', 112, 'STRING', 'Client ID', 'ADMIN', 'GITEE', 0, '{}', '2026-08-12 15:57:59.495992', NULL, '2026-08-12 15:57:59.495992', '1');
INSERT INTO `sys_config` VALUES ('7833544856736466882', 'ACCOUNT_CANCEL_RETENTION_DAYS', '15', 'AUTH_PASSWORD', '注销账号保留天数', 25, 'INT', '注销账号保留天数', NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7840109703376148568', 'MAIL_LOCAL_PASSWORD', '', 'MAIL', 'SMTP 密码', 13, 'STRING', 'SMTP 密码', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7841576503746820287', 'MAIL_LOCAL_PORT', '1025', 'MAIL', 'SMTP 端口', 11, 'INT', 'SMTP 端口', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7858449177900745217', 'DEFAULT_EMAIL_ENGINE', 'LOCAL', 'MAIL', '默认邮件引擎', 1, 'STRING', '默认邮件引擎', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7878128811204573416', 'AUTH_OAUTH_ADMIN_WECHAT_OPEN_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端微信开放平台回调', 134, 'STRING', 'Redirect URI', 'ADMIN', 'WECHAT_OPEN', 0, '{}', '2026-08-12 15:57:59.536850', NULL, '2026-08-12 15:57:59.536850', '1');
INSERT INTO `sys_config` VALUES ('7878222758007273827', 'COPYRIGHT_TEXT', 'hei-fastapi', 'SYS', '版权文案', 1, 'STRING', '版权文案', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7879475119693769672', 'AUDIT_ALERT_BULK_DELETE_THRESHOLD', '20', 'AUDIT_ALERT', '批量删除阈值', 21, 'INT', '批量删除阈值', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7884464514418971420', 'AUDIT_ALERT_BRUTE_FORCE_THRESHOLD', '10', 'AUDIT_ALERT', '暴力破解阈值', 20, 'INT', '暴力破解阈值', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7905539512046799183', 'AUTH_OAUTH_ADMIN_QQ_CLIENT_SECRET', '', 'AUTH_OAUTH', '管理端 QQ ClientSecret', 123, 'STRING', 'Client Secret', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.515748', NULL, '2026-08-12 15:57:59.515748', NULL);
INSERT INTO `sys_config` VALUES ('7915486803474335428', 'PASSWORD_MAX_CONSECUTIVE_CHARS', '3', 'AUTH_PASSWORD', '最大连续相同字符数', 13, 'INT', '最大连续相同字符数', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7924936151245798814', 'AUTH_LOGIN_ADMIN_EMAIL_NO_USER_POLICY', 'DENY', 'AUTH_LOGIN', 'ADMIN 邮箱无用户策略', 16, 'STRING', 'ADMIN 邮箱无用户策略', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7925472518021370098', 'AUTH_OAUTH_PORTAL_GITHUB_CLIENT_SECRET', 'gAAAAABqfJrd_yTDy0DkEmej80F2frWaLzP6NMbZMBDyOzUwUOTnjRHn7UGI2ACEkt4EzA9zS9q1dK4T0B0yRAP_LDtVMcNs3w', 'AUTH_OAUTH', '门户 GitHub ClientSecret', 3, 'STRING', 'Client Secret', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.407874', NULL, '2026-08-12 15:57:59.407874', '1');
INSERT INTO `sys_config` VALUES ('7925585644534920178', 'AUTH_LOGIN_ADMIN_FAILURE_WINDOW_SECONDS', '300', 'AUTH_LOGIN', 'ADMIN 登录失败窗口（秒）', 10, 'INT', 'ADMIN 登录失败窗口（秒）', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7926454604208823829', 'MAIL_TEMPLATE_BIND_EMAIL_CODE', '{\"subject\": \"{{app_name}} 绑定邮箱验证码\", \"body\": \"您的绑定验证码是 {{code}}，{{expire_minutes}} 分钟内有效。\"}', 'MAIL_TEMPLATE', '绑定邮箱验证码邮件模板', 20, 'JSON', '绑定邮箱验证码邮件模板', NULL, 'BIND_EMAIL_CODE', 0, '{}', '2026-08-12 14:36:35.270099', NULL, '2026-08-12 14:36:35.270099', NULL);
INSERT INTO `sys_config` VALUES ('7932451368437798893', 'SMS_TEMPLATE_ACCOUNT_CANCELLED', '{\"code\": \"\", \"content\": \"账号已申请注销，将于{{retention_days}}天后彻底删除。\"}', 'SMS_TEMPLATE', '账号注销确认短信模板', 20, 'JSON', '账号注销确认短信模板', NULL, 'ACCOUNT_CANCELLED', 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7932896058923796364', 'MAIL_LOCAL_USE_STARTTLS', 'FALSE', 'MAIL', 'SMTP 使用 STARTTLS', 18, 'BOOL', 'SMTP 使用 STARTTLS', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7943155149240026436', 'AUTH_PASSWORD_RESET_TOKEN_TTL_SECONDS', '600', 'AUTH_TOKEN', '密码重置 Token 有效期（秒）', 2, 'INT', '密码重置 Token 有效期（秒）', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7950285006814595256', 'STORAGE_RUSTFS_ENDPOINT', 'http://127.0.0.1:9002', 'STORAGE', 'RustFS S3 API 端点', 41, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-08 17:17:14.441312', '1');
INSERT INTO `sys_config` VALUES ('7952653754833581211', 'STORAGE_RUSTFS_SECRET_KEY', 'gAAAAABqeFtjH6R_u459TAGXEzeQOa9GRHtG_8xkDrnY8hwiZYyxeRtHh0YqhXsuzfkY3_Nn3OWId3rStJkCQgQUdwfkLtuAyA', 'STORAGE', 'RustFS Secret Key', 43, 'STRING', NULL, NULL, NULL, 0, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 10:49:34.902513', '1');
INSERT INTO `sys_config` VALUES ('7953467271865167705', 'AUTH_OAUTH_ADMIN_QQ_CLIENT_ID', '', 'AUTH_OAUTH', '管理端 QQ ClientId', 122, 'STRING', 'Client ID', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.512033', NULL, '2026-08-12 15:57:59.512033', '1');
INSERT INTO `sys_config` VALUES ('7957493498148551921', 'AUTH_LOGIN_ADMIN_LOCK_SECONDS', '300', 'AUTH_LOGIN', 'ADMIN 锁定时间（秒）', 12, 'INT', 'ADMIN 锁定时间（秒）', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7961221992910853426', 'PASSWORD_CHANGE_VERIFY_METHOD', 'OLD_PASSWORD', 'AUTH_PASSWORD', '自助改密验证方式', 2, 'STRING', '自助改密验证方式', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7961921801574207800', 'AUTH_LOGIN_PORTAL_ALLOW_PHONE', 'TRUE', 'AUTH_LOGIN', 'PORTAL 允许手机号登录', 21, 'BOOL', 'PORTAL 允许手机号登录', 'PORTAL', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7964584885414219711', 'AUTH_OAUTH_PORTAL_GITHUB_REDIRECT_URI', '', 'AUTH_OAUTH', '门户 GitHub 回调', 4, 'STRING', 'Redirect URI', 'PORTAL', 'GITHUB', 0, '{}', '2026-08-12 15:57:59.412339', NULL, '2026-08-12 15:57:59.412339', '1');
INSERT INTO `sys_config` VALUES ('7965270553477718376', 'AUTH_LOGIN_ADMIN_ALLOW_PHONE', 'TRUE', 'AUTH_LOGIN', 'ADMIN 允许手机号登录', 13, 'BOOL', 'ADMIN 允许手机号登录', 'ADMIN', NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', '1');
INSERT INTO `sys_config` VALUES ('7971733830500536612', 'PASSWORD_FORBID_HISTORICAL', 'TRUE', 'AUTH_PASSWORD', '禁止复用历史密码', 15, 'BOOL', '禁止复用历史密码', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7975896582003472913', 'MAIL_TENCENT_REGION', '', 'MAIL', '腾讯云邮件区域', 90, 'STRING', '腾讯云邮件区域', NULL, NULL, 0, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_config` VALUES ('7983845910800985649', 'AUTH_OAUTH_ADMIN_QQ_ENABLED', 'FALSE', 'AUTH_OAUTH', '管理端 QQ 登录', 121, 'BOOL', '启用', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.507724', NULL, '2026-08-12 15:57:59.507724', '1');
INSERT INTO `sys_config` VALUES ('7997330304022587335', 'AUTH_OAUTH_ADMIN_QQ_REDIRECT_URI', '', 'AUTH_OAUTH', '管理端 QQ 回调', 124, 'STRING', 'Redirect URI', 'ADMIN', 'QQ', 0, '{}', '2026-08-12 15:57:59.520333', NULL, '2026-08-12 15:57:59.520333', '1');

-- ----------------------------
-- Table structure for sys_dept
-- ----------------------------
DROP TABLE IF EXISTS `sys_dept`;
CREATE TABLE `sys_dept`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `parent_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '父级ID',
  `master_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '部门主管账户ID',
  `deputy_master_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '部门副主管账户ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门名称',
  `category` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门类别/层级类型',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_virtual` tinyint(1) NOT NULL COMMENT '是否虚拟组织：1 虚拟 / 0 实体',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '部门状态：ENABLED/DISABLED',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '部门' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dept
-- ----------------------------
INSERT INTO `sys_dept` VALUES ('8200000000000101', NULL, NULL, NULL, '总部', 'ORG', 1, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_dept` VALUES ('8200000000000102', '8200000000000101', '8200000000000201', '8200000000000205', '研发部', 'ORG', 1, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_dept` VALUES ('8200000000000103', '8200000000000101', '8200000000000202', NULL, '市场部', 'ORG', 2, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_dept` VALUES ('8200000000000104', '8200000000000102', '8200000000000203', NULL, '前端组', 'ORG', 1, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_dept` VALUES ('8200000000000105', '8200000000000102', '8200000000000208', NULL, '后端组', 'ORG', 2, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_dept` VALUES ('8200000000000106', '8200000000000102', '8200000000000209', NULL, '测试组', 'ORG', 3, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_dept` VALUES ('8200000000000107', '8200000000000101', '8200000000000206', NULL, '人事行政部', 'ORG', 3, 0, 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_dict
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict`;
CREATE TABLE `sys_dict`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典项编码（同父级下唯一）',
  `label` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典项展示标签',
  `value` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典项实际值',
  `color` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端展示颜色',
  `category` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '字典分类：SYSTEM（系统）/ BUSINESS（业务）',
  `parent_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '父级字典项ID',
  `status` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '字典项状态：ENABLED/DISABLED',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_dict_code`(`code` ASC) USING BTREE,
  INDEX `idx_sys_dict_category`(`category` ASC) USING BTREE,
  INDEX `idx_sys_dict_parent_id`(`parent_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '数据字典' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_dict
-- ----------------------------
INSERT INTO `sys_dict` VALUES ('100001', 'COMMON_STATUS', '状态', 'COMMON_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100002', 'COMMON_STATUS_ENABLED', '启用', 'ENABLED', '#18a058', 'SYS', '100001', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100003', 'COMMON_STATUS_DISABLED', '禁用', 'DISABLED', '#d03050', 'SYS', '100001', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100004', 'SYS_BIZ_CATEGORY', '系统/业务分类', 'SYS_BIZ_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100005', 'SYS_BIZ_CATEGORY_SYS', '系统', 'SYS', '#2080f0', 'SYS', '100004', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100006', 'SYS_BIZ_CATEGORY_BIZ', '业务', 'BIZ', '#f0a020', 'SYS', '100004', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100010', 'ACCOUNT_STATUS', '账号状态', 'ACCOUNT_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100011', 'ACCOUNT_STATUS_ENABLED', '启用', 'ENABLED', '#18a058', 'SYS', '100010', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100012', 'ACCOUNT_STATUS_DISABLED', '禁用', 'DISABLED', '#d03050', 'SYS', '100010', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100013', 'ACCOUNT_STATUS_CANCELLED', '已注销', 'CANCELLED', '#909399', 'SYS', '100010', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100014', 'ROLE_SCOPE_TYPE', '角色范围类型', 'ROLE_SCOPE_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100015', 'ROLE_SCOPE_TYPE_PLATFORM', '平台', 'PLATFORM', '#2080f0', 'SYS', '100014', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100016', 'ROLE_SCOPE_TYPE_DEPT', '部门', 'DEPT', '#18a058', 'SYS', '100014', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100017', 'RESOURCE_TYPE', '资源类型', 'RESOURCE_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100018', 'RESOURCE_TYPE_CATALOG', '目录', 'CATALOG', '#722ed1', 'SYS', '100017', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100019', 'RESOURCE_TYPE_MENU', '菜单', 'MENU', '#2080f0', 'SYS', '100017', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100020', 'RESOURCE_TYPE_PAGE', '页面', 'PAGE', '#18a058', 'SYS', '100017', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100021', 'RESOURCE_TYPE_BUTTON', '按钮', 'BUTTON', '#f0a020', 'SYS', '100017', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100022', 'RESOURCE_TYPE_ACTION', '操作', 'ACTION', '#d03050', 'SYS', '100017', 'ENABLED', 5, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100023', 'RESOURCE_TYPE_API_GROUP', '接口组', 'API_GROUP', '#1677ff', 'SYS', '100017', 'ENABLED', 6, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100024', 'DATA_SCOPE', '数据范围', 'DATA_SCOPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100025', 'DATA_SCOPE_ALL', '全部', 'ALL', '#18a058', 'SYS', '100024', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100026', 'DATA_SCOPE_DEPT_AND_CHILD', '本部门及子部门', 'DEPT_AND_CHILD', '#2080f0', 'SYS', '100024', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100027', 'DATA_SCOPE_DEPT', '本部门', 'DEPT', '#2db7f5', 'SYS', '100024', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100028', 'DATA_SCOPE_SELF', '本人', 'SELF', '#f0a020', 'SYS', '100024', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100029', 'DATA_SCOPE_CUSTOM', '自定义部门', 'CUSTOM', '#722ed1', 'SYS', '100024', 'ENABLED', 5, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100040', 'DEPT_CATEGORY', '部门分类', 'DEPT_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100041', 'DEPT_CATEGORY_COMPANY', '公司', 'COMPANY', '#2080f0', 'SYS', '100040', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100042', 'DEPT_CATEGORY_DEPARTMENT', '部门', 'DEPARTMENT', '#18a058', 'SYS', '100040', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100043', 'DEPT_CATEGORY_TEAM', '团队', 'TEAM', '#f0a020', 'SYS', '100040', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100044', 'DEPT_CATEGORY_VIRTUAL', '虚拟组织', 'VIRTUAL', '#909399', 'SYS', '100040', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100045', 'POSITION_CATEGORY', '岗位分类', 'POSITION_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100046', 'POSITION_CATEGORY_MANAGEMENT', '管理', 'MANAGEMENT', '#2080f0', 'SYS', '100045', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100047', 'POSITION_CATEGORY_TECHNICAL', '技术', 'TECHNICAL', '#18a058', 'SYS', '100045', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100048', 'POSITION_CATEGORY_OPERATION', '运营', 'OPERATION', '#f0a020', 'SYS', '100045', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100049', 'POSITION_CATEGORY_SUPPORT', '支持', 'SUPPORT', '#909399', 'SYS', '100045', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100054', 'BANNER_CATEGORY', '展示图分类', 'BANNER_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100055', 'BANNER_CATEGORY_HOME', '首页', 'HOME', '#18a058', 'SYS', '100054', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100056', 'BANNER_CATEGORY_LOGIN', '登录', 'LOGIN', '#2080f0', 'SYS', '100054', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100057', 'BANNER_CATEGORY_WORKPLACE', '工作台', 'WORKPLACE', '#722ed1', 'SYS', '100054', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100058', 'BANNER_CATEGORY_NOTICE', '公告', 'NOTICE', '#f0a020', 'SYS', '100054', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100059', 'BANNER_CATEGORY_ADMIN_DASHBOARD', '管理端仪表盘', 'ADMIN_DASHBOARD', '#2080f0', 'SYS', '100054', 'ENABLED', 5, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100060', 'BANNER_CATEGORY_SYSTEM_UPGRADE', '系统升级', 'SYSTEM_UPGRADE', '#d03050', 'SYS', '100054', 'ENABLED', 6, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100061', 'BANNER_TYPE', '展示图类型', 'BANNER_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100062', 'BANNER_TYPE_CAROUSEL', '轮播图', 'CAROUSEL', '#18a058', 'SYS', '100061', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100063', 'BANNER_TYPE_HERO', '主视觉', 'HERO', '#2080f0', 'SYS', '100061', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100064', 'BANNER_TYPE_NOTICE', '公告', 'NOTICE', '#f0a020', 'SYS', '100061', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100065', 'BANNER_TYPE_CARD', '卡片', 'CARD', '#722ed1', 'SYS', '100061', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100066', 'BANNER_TYPE_POPUP', '弹窗', 'POPUP', '#d03050', 'SYS', '100061', 'ENABLED', 5, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100067', 'BANNER_TYPE_SIDEBAR', '侧边栏', 'SIDEBAR', '#2080f0', 'SYS', '100061', 'ENABLED', 6, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100068', 'BANNER_POSITION', '展示图位置', 'BANNER_POSITION', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100069', 'BANNER_POSITION_HOME_TOP', '首页顶部', 'HOME_TOP', '#18a058', 'SYS', '100068', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100070', 'BANNER_POSITION_HOME_MIDDLE', '首页中部', 'HOME_MIDDLE', '#18a058', 'SYS', '100068', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100071', 'BANNER_POSITION_HOME_BOTTOM', '首页底部', 'HOME_BOTTOM', '#18a058', 'SYS', '100068', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100072', 'BANNER_POSITION_LOGIN_SIDE', '登录侧边', 'LOGIN_SIDE', '#2080f0', 'SYS', '100068', 'ENABLED', 4, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100073', 'BANNER_POSITION_WORKPLACE_TOP', '工作台顶部', 'WORKPLACE_TOP', '#722ed1', 'SYS', '100068', 'ENABLED', 5, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100074', 'BANNER_POSITION_NOTICE_AREA', '公告区域', 'NOTICE_AREA', '#f0a020', 'SYS', '100068', 'ENABLED', 6, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100075', 'BANNER_POSITION_ADMIN_TOP', '管理端顶部', 'ADMIN_TOP', '#2080f0', 'SYS', '100068', 'ENABLED', 7, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100076', 'BANNER_POSITION_ADMIN_SIDEBAR', '管理端侧边栏', 'ADMIN_SIDEBAR', '#2080f0', 'SYS', '100068', 'ENABLED', 8, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100077', 'BANNER_LINK_TYPE', '展示图链接类型', 'BANNER_LINK_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100078', 'BANNER_LINK_TYPE_URL', '外部链接', 'URL', '#18a058', 'SYS', '100077', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100079', 'BANNER_LINK_TYPE_ROUTE', '路由', 'ROUTE', '#2080f0', 'SYS', '100077', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100080', 'BANNER_LINK_TYPE_NONE', '无链接', 'NONE', '#909399', 'SYS', '100077', 'ENABLED', 3, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100085', 'ACCOUNT_IDENTITY_BIND_STATUS', '账号身份绑定状态', 'ACCOUNT_IDENTITY_BIND_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100086', 'ACCOUNT_IDENTITY_BIND_STATUS_BOUND', '已绑定', 'BOUND', '#18a058', 'SYS', '100085', 'ENABLED', 1, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100087', 'ACCOUNT_IDENTITY_BIND_STATUS_UNBOUND', '未绑定', 'UNBOUND', '#909399', 'SYS', '100085', 'ENABLED', 2, '2026-06-29 00:00:00.000000', NULL, '2026-06-29 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100095', 'NOTIFICATION_SEVERITY', '通知严重级别', 'NOTIFICATION_SEVERITY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100096', 'NOTIFICATION_SEVERITY_INFO', '信息', 'INFO', '#2080f0', 'SYS', '100095', 'ENABLED', 1, '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100097', 'NOTIFICATION_SEVERITY_SUCCESS', '成功', 'SUCCESS', '#18a058', 'SYS', '100095', 'ENABLED', 2, '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100098', 'NOTIFICATION_SEVERITY_WARNING', '警告', 'WARNING', '#f0a020', 'SYS', '100095', 'ENABLED', 3, '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100099', 'NOTIFICATION_SEVERITY_ERROR', '错误', 'ERROR', '#d03050', 'SYS', '100095', 'ENABLED', 4, '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100126', 'NOTIFICATION_SEVERITY_URGENT', '紧急', 'URGENT', '#d03050', 'SYS', '100095', 'ENABLED', 5, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100127', 'CONTENT_TYPE', '内容格式', 'CONTENT_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100128', 'CONTENT_TYPE_TEXT', '纯文本', 'text', '#909399', 'SYS', '100127', 'ENABLED', 1, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100129', 'CONTENT_TYPE_HTML', '富文本', 'html', '#18a058', 'SYS', '100127', 'ENABLED', 2, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100130', 'CONTENT_TYPE_MARKDOWN', 'Markdown', 'markdown', '#722ed1', 'SYS', '100127', 'ENABLED', 3, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100131', 'TARGET_SCOPE', '目标范围', 'TARGET_SCOPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100132', 'TARGET_SCOPE_ALL', '全部', 'ALL', '#2080f0', 'SYS', '100131', 'ENABLED', 1, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100135', 'TARGET_SCOPE_SPECIFIC', '指定用户', 'SPECIFIC', '#d03050', 'SYS', '100131', 'ENABLED', 3, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100136', 'NOTIFY_LOCATION', '通知位置', 'NOTIFY_LOCATION', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100137', 'NOTIFY_LOCATION_CENTER', '通知中心', 'center', '#2080f0', 'SYS', '100136', 'ENABLED', 1, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100138', 'NOTIFY_LOCATION_POPUP', '弹窗', 'popup', '#f0a020', 'SYS', '100136', 'ENABLED', 2, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100139', 'NOTIFY_LOCATION_WORKSPACE', '工作台公告区', 'workspace', '#722ed1', 'SYS', '100136', 'ENABLED', 3, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100140', 'FEEDBACK_CATEGORY', '反馈分类', 'FEEDBACK_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100141', 'FEEDBACK_CATEGORY_SUGGESTION', '功能建议', 'SUGGESTION', '#18a058', 'SYS', '100140', 'ENABLED', 1, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100142', 'FEEDBACK_CATEGORY_BUG', '问题反馈', 'BUG', '#d03050', 'SYS', '100140', 'ENABLED', 2, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100143', 'FEEDBACK_CATEGORY_OTHER', '其他', 'OTHER', '#909399', 'SYS', '100140', 'ENABLED', 3, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100144', 'FEEDBACK_STATUS', '反馈状态', 'FEEDBACK_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100145', 'FEEDBACK_STATUS_PENDING', '待处理', 'PENDING', '#f0a020', 'SYS', '100144', 'ENABLED', 1, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100146', 'FEEDBACK_STATUS_REVIEWED', '已查看', 'REVIEWED', '#2080f0', 'SYS', '100144', 'ENABLED', 2, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100147', 'FEEDBACK_STATUS_RESOLVED', '已解决', 'RESOLVED', '#18a058', 'SYS', '100144', 'ENABLED', 3, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100148', 'FEEDBACK_STATUS_CLOSED', '已关闭', 'CLOSED', '#909399', 'SYS', '100144', 'ENABLED', 4, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100149', 'PUBLISH_STATUS', '发布状态', 'PUBLISH_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100150', 'PUBLISH_STATUS_DRAFT', '草稿', 'DRAFT', '#909399', 'SYS', '100149', 'ENABLED', 1, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100151', 'PUBLISH_STATUS_PUBLISHED', '已发布', 'PUBLISHED', '#18a058', 'SYS', '100149', 'ENABLED', 2, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100152', 'PUBLISH_STATUS_REVOKED', '已撤回', 'REVOKED', '#d03050', 'SYS', '100149', 'ENABLED', 3, '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100206', 'TARGET_SCOPE_ACCOUNT_TYPE', '按账号类型', 'ACCOUNT_TYPE', '#722ed1', 'SYS', '100131', 'ENABLED', 2, '2026-08-08 04:14:19.198462', NULL, '2026-08-08 04:14:19.198462', NULL);
INSERT INTO `sys_dict` VALUES ('100210', 'NOTIFICATION_CATEGORY', '通知分类', 'NOTIFICATION_CATEGORY', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100211', 'NOTIFICATION_CATEGORY_ORDER', '订单', 'ORDER', '#2080f0', 'SYS', '100210', 'ENABLED', 1, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100212', 'NOTIFICATION_CATEGORY_APPROVAL', '审批', 'APPROVAL', '#722ed1', 'SYS', '100210', 'ENABLED', 2, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100213', 'NOTIFICATION_CATEGORY_SYSTEM', '系统', 'SYSTEM', '#18a058', 'SYS', '100210', 'ENABLED', 3, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100214', 'NOTIFICATION_CATEGORY_SECURITY', '安全', 'SECURITY', '#d03050', 'SYS', '100210', 'ENABLED', 4, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100215', 'NOTIFICATION_CATEGORY_BIZ', '业务', 'BIZ', '#f0a020', 'SYS', '100210', 'ENABLED', 5, '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100220', 'AUDIT_ACTION_TYPE', '审计操作类型', 'AUDIT_ACTION_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100221', 'AUDIT_ACTION_TYPE_CREATE', '新增', 'CREATE', '#18a058', 'SYS', '100220', 'ENABLED', 1, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100222', 'AUDIT_ACTION_TYPE_UPDATE', '修改', 'UPDATE', '#2080f0', 'SYS', '100220', 'ENABLED', 2, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100223', 'AUDIT_ACTION_TYPE_DELETE', '删除', 'DELETE', '#d03050', 'SYS', '100220', 'ENABLED', 3, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100224', 'AUDIT_ACTION_TYPE_QUERY', '查询', 'QUERY', '#909399', 'SYS', '100220', 'ENABLED', 4, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100225', 'AUDIT_ACTION_TYPE_EXPORT', '导出', 'EXPORT', '#f0a020', 'SYS', '100220', 'ENABLED', 5, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100226', 'AUDIT_ACTION_TYPE_LOGIN', '登录', 'LOGIN', '#18a058', 'SYS', '100220', 'ENABLED', 6, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100227', 'AUDIT_ACTION_TYPE_LOGOUT', '登出', 'LOGOUT', '#722ed1', 'SYS', '100220', 'ENABLED', 7, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('100228', 'AUDIT_ACTION_TYPE_OTHER', '其它', 'OTHER', '#909399', 'SYS', '100220', 'ENABLED', 8, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101500', 'REAL_NAME_BUSINESS_TYPE', '实名业务类型', 'REAL_NAME_BUSINESS_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101501', 'REAL_NAME_BUSINESS_ACCOUNT_VERIFY', '账号实名认证', 'ACCOUNT_VERIFY', '#18a058', 'SYS', '101500', 'ENABLED', 1, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101502', 'REAL_NAME_BUSINESS_ACCOUNT_RECOVERY', '实名找回账号', 'ACCOUNT_RECOVERY', '#909399', 'SYS', '101500', 'ENABLED', 2, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101510', 'IDENTITY_VERIFY_STATUS', '实名认证状态', 'IDENTITY_VERIFY_STATUS', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101511', 'IDENTITY_VERIFY_STATUS_UNVERIFIED', '未认证', 'UNVERIFIED', '#909399', 'SYS', '101510', 'ENABLED', 1, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101512', 'IDENTITY_VERIFY_STATUS_PENDING', '审核中', 'PENDING', '#f0a020', 'SYS', '101510', 'ENABLED', 2, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101513', 'IDENTITY_VERIFY_STATUS_VERIFIED', '已认证', 'VERIFIED', '#18a058', 'SYS', '101510', 'ENABLED', 3, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101514', 'IDENTITY_VERIFY_STATUS_REJECTED', '已驳回', 'REJECTED', '#d03050', 'SYS', '101510', 'ENABLED', 4, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101520', 'IDENTITY_DOCUMENT_TYPE', '证件类型', 'IDENTITY_DOCUMENT_TYPE', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101521', 'IDENTITY_DOCUMENT_ID_CARD', '居民身份证', 'ID_CARD', '#2080f0', 'SYS', '101520', 'ENABLED', 1, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101522', 'IDENTITY_DOCUMENT_PASSPORT', '护照', 'PASSPORT', '#2080f0', 'SYS', '101520', 'ENABLED', 2, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101523', 'IDENTITY_DOCUMENT_EID_CARD', '电子身份证', 'EID_CARD', '#722ed1', 'SYS', '101520', 'ENABLED', 3, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101530', 'IDENTITY_VERIFY_CHANNEL', '认证通道', 'IDENTITY_VERIFY_CHANNEL', '#2080f0', 'SYS', NULL, 'ENABLED', 0, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101531', 'IDENTITY_VERIFY_CHANNEL_THIRD_PARTY', '第三方实人', 'THIRD_PARTY', '#2080f0', 'SYS', '101530', 'ENABLED', 1, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('101532', 'IDENTITY_VERIFY_CHANNEL_MANUAL', '人工审核', 'MANUAL', '#f0a020', 'SYS', '101530', 'ENABLED', 2, '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_dict` VALUES ('7211108434702344271', 'WECHAT_OPEN', '微信开放平台', 'WECHAT_OPEN', NULL, 'SYS', '7656105493975688393', 'ENABLED', 4, '2026-08-12 15:57:59.388141', NULL, '2026-08-12 15:57:59.388141', NULL);
INSERT INTO `sys_dict` VALUES ('7260287663522585806', 'WECHAT_MP', '微信小程序', 'WECHAT_MP', NULL, 'SYS', '7656105493975688393', 'ENABLED', 5, '2026-08-12 15:57:59.392548', NULL, '2026-08-12 15:57:59.392548', NULL);
INSERT INTO `sys_dict` VALUES ('7351536453246655198', 'GITEE', 'Gitee', 'GITEE', NULL, 'SYS', '7656105493975688393', 'ENABLED', 2, '2026-08-12 15:57:59.380323', NULL, '2026-08-12 15:57:59.380323', NULL);
INSERT INTO `sys_dict` VALUES ('7399371307733482375', 'QQ', 'QQ', 'QQ', NULL, 'SYS', '7656105493975688393', 'ENABLED', 3, '2026-08-12 15:57:59.384444', NULL, '2026-08-12 15:57:59.384444', NULL);
INSERT INTO `sys_dict` VALUES ('7587538474363101234', 'GITHUB', 'GitHub', 'GITHUB', NULL, 'SYS', '7656105493975688393', 'ENABLED', 1, '2026-08-12 15:57:59.376072', NULL, '2026-08-12 15:57:59.376072', NULL);
INSERT INTO `sys_dict` VALUES ('7656105493975688393', 'OAUTH_PROVIDER', '三方登录提供商', 'OAUTH_PROVIDER', NULL, 'SYS', NULL, 'ENABLED', 90, '2026-08-12 15:57:59.360135', NULL, '2026-08-12 15:57:59.360135', NULL);

-- ----------------------------
-- Table structure for sys_feedback
-- ----------------------------
DROP TABLE IF EXISTS `sys_feedback`;
CREATE TABLE `sys_feedback`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容',
  `category` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类',
  `contact` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '联系方式',
  `attach_object_names` json NOT NULL COMMENT '用户上传附件 object_name 列表',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '反馈状态：PENDING/REPLIED/CLOSED 等',
  `reply` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '管理员回复内容',
  `replied_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '回复人账户ID',
  `replied_at` datetime(6) NULL DEFAULT NULL COMMENT '管理员回复时间',
  `submitter_account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '提交人账户类型',
  `submitter_account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '提交人账户ID',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '意见反馈' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_feedback
-- ----------------------------
INSERT INTO `sys_feedback` VALUES ('7491849375090675712', '哈哈哈', '撒擦三次', 'SUGGESTION', '擦拭擦拭', '[\"uploads/2026/08/08/be3515e142974cf08b46e348d1c3d8d3.png\"]', 'RESOLVED', 'ok', '1', '2026-08-08 13:40:49.757811', 'PORTAL', '7491847383584804864', '2026-08-08 13:34:42.831539', '7491847383584804864', '2026-08-08 13:40:49.742207', '1');
INSERT INTO `sys_feedback` VALUES ('8300000000000301', '门户登录体验', '希望增加短信验证码登录选项。', 'SUGGESTION', 'bob@demo.local', '[]', 'PENDING', NULL, NULL, NULL, 'PORTAL', '8200000000000211', '2026-08-23 09:00:00.000000', '8200000000000211', '2026-08-23 09:00:00.000000', '8200000000000211');
INSERT INTO `sys_feedback` VALUES ('8300000000000302', '活动页面加载慢', '活动列表在弱网下加载超过 5 秒。', 'BUG', '13800001001', '[]', 'PENDING', NULL, NULL, NULL, 'PORTAL', '8200000000000211', '2026-08-23 09:00:00.000000', '8200000000000211', '2026-08-23 09:00:00.000000', '8200000000000211');
INSERT INTO `sys_feedback` VALUES ('8300000000000303', '个人中心样式问题', '头像上传后预览区域会抖动。', 'BUG', 'alice@demo.local', '[]', 'RESOLVED', '已记录，将在下个版本修复。', '1', '2026-08-23 09:00:00.000000', 'PORTAL', '8200000000000212', '2026-08-23 09:00:00.000000', '8200000000000212', '2026-08-23 09:00:00.000000', '8200000000000212');
INSERT INTO `sys_feedback` VALUES ('8300000000000304', '希望增加深色模式', '门户夜间使用较刺眼，建议支持深色主题。', 'SUGGESTION', 'alice@demo.local', '[]', 'PENDING', NULL, NULL, NULL, 'PORTAL', '8200000000000212', '2026-08-23 09:00:00.000000', '8200000000000212', '2026-08-23 09:00:00.000000', '8200000000000212');
INSERT INTO `sys_feedback` VALUES ('8300000000000305', '历史门户反馈', '早期门户测试账号提交的反馈样例。', 'GENERAL', '13800000000', '[]', 'RESOLVED', 'ok', '1', '2026-08-08 13:40:49.757811', 'PORTAL', '7491847383584804864', '2026-08-23 09:00:00.000000', '7491847383584804864', '2026-08-23 09:00:00.000000', '7491847383584804864');
INSERT INTO `sys_feedback` VALUES ('8300000000000306', '部门树展示优化', '建议部门详情中直接展示成员数量。', 'SUGGESTION', 'iam-admin@demo.local', '[]', 'PENDING', NULL, NULL, NULL, 'ADMIN', '8200000000000201', '2026-08-23 09:00:00.000000', '8200000000000201', '2026-08-23 09:00:00.000000', '8200000000000201');
INSERT INTO `sys_feedback` VALUES ('8300000000000307', '角色授权交互', '授权用户弹窗希望支持按部门筛选。', 'SUGGESTION', 'biz-all@demo.local', '[]', 'PENDING', NULL, NULL, NULL, 'ADMIN', '8200000000000202', '2026-08-23 09:00:00.000000', '8200000000000202', '2026-08-23 09:00:00.000000', '8200000000000202');
INSERT INTO `sys_feedback` VALUES ('8300000000000308', '只读账号权限确认', '只读账号不应看到删除按钮，请确认前端权限控制。', 'BUG', 'readonly@demo.local', '[]', 'RESOLVED', '已确认，仅隐藏无权限按钮。', '1', '2026-08-23 09:00:00.000000', 'ADMIN', '8200000000000206', '2026-08-23 09:00:00.000000', '8200000000000206', '2026-08-23 09:00:00.000000', '8200000000000206');

-- ----------------------------
-- Table structure for sys_file
-- ----------------------------
DROP TABLE IF EXISTS `sys_file`;
CREATE TABLE `sys_file`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `object_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '对象存储中的对象键/路径',
  `original_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户上传时的原始文件名',
  `storage_provider` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '存储服务商：minio/rustfs/oss/s3',
  `bucket` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '对象存储桶名称',
  `content_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'MIME 类型',
  `size` bigint NOT NULL COMMENT '文件大小（字节）',
  `url` varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件访问 URL（可为签名地址）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_file_object_name`(`object_name` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '文件元数据' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_file
-- ----------------------------
INSERT INTO `sys_file` VALUES ('2086404723008208898', 'uploads/2026/08/09/48af08ca31e346d3a75139fa84eb0282.png', 'QR2026080700024_1786146767761.png', 'rustfs', 'defaultbucket', 'image/png', 38636, 'uploads/2026/08/09/48af08ca31e346d3a75139fa84eb0282.png', '2026-08-09 10:50:24.153462', '1', '2026-08-09 10:50:24.153462', '1');
INSERT INTO `sys_file` VALUES ('2086408061170970625', 'uploads/2026/08/09/3dbf7f27b66f4023a0736bdccaa9596b.png', 'QR2026080700024_1786146767761.png', 'rustfs', 'defaultbucket', 'image/png', 38636, 'uploads/2026/08/09/3dbf7f27b66f4023a0736bdccaa9596b.png', '2026-08-09 11:03:40.030025', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2086410328867565570', 'uploads/2026/08/09/02acc3dee5454d34913b07f49fe59cac.png', 'avatar.png', 'rustfs', 'defaultbucket', 'image/png', 193387, 'uploads/2026/08/09/02acc3dee5454d34913b07f49fe59cac.png', '2026-08-09 11:12:40.692242', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2086415187620601857', 'uploads/2026/08/09/85e1b98acfc9465abbbba86ef3b4fec8.jpg', '120153703_touxiang_bobopic (1).jpg', 'rustfs', 'defaultbucket', 'image/jpeg', 65451, 'uploads/2026/08/09/85e1b98acfc9465abbbba86ef3b4fec8.jpg', '2026-08-09 11:31:59.117567', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('2088928636523130882', 'uploads/2026/08/16/947e530849e245d3bcb873354da8f113.txt', 'PyPI-Recovery-Codes-charliebyte-2026-06-14T07_19_21.473015.txt', 'rustfs', 'defaultbucket', 'text/plain', 135, 'uploads/2026/08/16/947e530849e245d3bcb873354da8f113.txt', '2026-08-16 09:59:32.021932', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('7491869364283744256', 'uploads/2026/08/08/93c4b0dba86f483bb1905ba079550521.png', 'QR2026080700024_1786146767761.png', 'minio', 'vms', 'image/png', 38636, 'uploads/2026/08/08/93c4b0dba86f483bb1905ba079550521.png', '2026-08-08 14:54:08.592047', '1', '2026-08-16 10:27:45.912235', '1');
INSERT INTO `sys_file` VALUES ('7491906012023353344', 'uploads/2026/08/08/0e535c4dc69241eab526c5e94d9eb19b.png', 'QR2026080700024_1786146767761.png', 'rustfs', 'defaultbucket', 'image/png', 38636, 'uploads/2026/08/08/0e535c4dc69241eab526c5e94d9eb19b.png', '2026-08-08 17:19:46.108576', '1', '2026-08-16 10:27:45.912235', '1');

-- ----------------------------
-- Table structure for sys_group
-- ----------------------------
DROP TABLE IF EXISTS `sys_group`;
CREATE TABLE `sys_group`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户组名称',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据权限范围）',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '用户组描述',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户组状态：ENABLED/DISABLED',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_group_name`(`name` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户组' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_group
-- ----------------------------
INSERT INTO `sys_group` VALUES ('8200000000000301', '研发组成员', '8200000000000102', '继承 BIZ_DEPT，目录模块本部门数据', 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_group` VALUES ('8200000000000302', '市场运营组', '8200000000000103', '继承 BIZ_ALL，活动模块全量数据', 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_group` VALUES ('8200000000000303', 'IAM协作组', '8200000000000102', '继承 IAM_READONLY，账号只读协作', 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_group` VALUES ('8200000000000304', '前端专项组', '8200000000000104', '前端组成员，用于子部门协作演示', 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_group` VALUES ('8200000000000305', '后端专项组', '8200000000000105', '后端组成员，目录模块本部门数据', 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_group` VALUES ('8200000000000306', '测试专项组', '8200000000000106', '测试组成员，知识分类部门及子部门', 'ENABLED', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_iam_relation
-- ----------------------------
DROP TABLE IF EXISTS `sys_iam_relation`;
CREATE TABLE `sys_iam_relation`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `subject_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主体类型：ACCOUNT/DEPT/ROLE/GROUP/POSITION',
  `subject_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主体记录ID',
  `account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户类型：ADMIN（管理端）/ PORTAL（门户端）',
  `relation_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关系类型：MEMBER/GRANT/OWN 等',
  `target_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标类型：RESOURCE/ROLE/DEPT/DATA_SCOPE 等',
  `target_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标记录ID',
  `target_key` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '目标业务标识（如权限 code）',
  `grant_mode` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '授权模式：DIRECT/INHERIT 等',
  `data_scope` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据范围：ALL/DEPT/DEPT_AND_CHILD/CUSTOM/SELF',
  `custom_scope_dept_ids` json NOT NULL COMMENT '自定义数据范围部门ID列表（JSON 数组）',
  `is_primary` tinyint(1) NOT NULL COMMENT '是否主关系/主岗位：1 是 / 0 否',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关系状态：ACTIVE/INACTIVE',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT 'IAM 关系说明',
  `reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '授权/变更原因',
  `expired_at` datetime(6) NULL DEFAULT NULL COMMENT '关系失效时间（空表示永久）',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_iam_relation_subject_relation_target`(`subject_type` ASC, `subject_id` ASC, `relation_type` ASC, `target_type` ASC, `target_id` ASC, `target_key` ASC, `account_type` ASC) USING BTREE,
  INDEX `ix_sys_iam_relation_account_type_relation`(`account_type` ASC, `relation_type` ASC) USING BTREE,
  INDEX `ix_sys_iam_relation_subject`(`subject_type` ASC, `subject_id` ASC, `relation_type` ASC) USING BTREE,
  INDEX `ix_sys_iam_relation_target`(`target_type` ASC, `target_id` ASC, `target_key` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = 'IAM 关系（权限挂载）' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_iam_relation
-- ----------------------------
INSERT INTO `sys_iam_relation` VALUES ('1', 'ACCOUNT', '1', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '1', '', 'CASCADE', 'SELF', '[]', 0, 99, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 11:56:13.747886', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7100000000000002002', 'RESOURCE', '200002', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'workspace:overview:view', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '工作台总览', NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7109538496802851524', 'RESOURCE', '201060', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7112731414735234196', 'RESOURCE', '203033', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:detail', 'CASCADE', 'SELF', '[]', 0, 30, 'ENABLED', '详情订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7122813088692955083', 'RESOURCE', '202240', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:revoke', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '撤回消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7133323711340623180', 'RESOURCE', '203046', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:list', 'CASCADE', 'DEPT_AND_CHILD', '[]', 0, 90, 'ENABLED', '树列表知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7159667080064467923', 'RESOURCE', '203023', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:detail', 'CASCADE', 'DEPT', '[]', 0, 30, 'ENABLED', '详情目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7174771647983316441', 'RESOURCE', '202205', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7212116468775288981', 'RESOURCE', '204012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:update', 'CASCADE', 'ALL', '[]', 0, 2, 'ENABLED', '编辑任务', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7223792660518235132', 'RESOURCE', '203014', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:update', 'CASCADE', 'ALL', '[]', 0, 40, 'ENABLED', '编辑活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7232506669573029115', 'RESOURCE', '200028', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7284368855990246834', 'RESOURCE', '204001', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分页任务', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7304374671311844075', 'RESOURCE', '203034', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:update', 'CASCADE', 'SELF', '[]', 0, 40, 'ENABLED', '编辑订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7334894311555546200', 'RESOURCE', '203035', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:delete', 'CASCADE', 'SELF', '[]', 0, 50, 'ENABLED', '删除订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7380222627904177407', 'RESOURCE', '202241', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:pin', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '置顶消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7453524161449865528', 'RESOURCE', '203041', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:page', 'CASCADE', 'DEPT_AND_CHILD', '[]', 0, 10, 'ENABLED', '分页知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7518444451967536602', 'RESOURCE', '203032', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:create', 'CASCADE', 'SELF', '[]', 0, 20, 'ENABLED', '新增订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7560408972191285564', 'RESOURCE', '204013', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:delete', 'CASCADE', 'ALL', '[]', 0, 3, 'ENABLED', '删除任务', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7567716937788130247', 'RESOURCE', '203015', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:delete', 'CASCADE', 'ALL', '[]', 0, 50, 'ENABLED', '删除活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7569914124743592951', 'RESOURCE', '202201', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分页消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7574875650561833761', 'RESOURCE', '203043', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:detail', 'CASCADE', 'DEPT_AND_CHILD', '[]', 0, 30, 'ENABLED', '详情知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7575458644059959564', 'RESOURCE', '202202', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7601824633419714671', 'RESOURCE', '203042', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:create', 'CASCADE', 'DEPT_AND_CHILD', '[]', 0, 20, 'ENABLED', '新增知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7606481288132251344', 'RESOURCE', '200029', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7611870955633752502', 'RESOURCE', '203012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:create', 'CASCADE', 'ALL', '[]', 0, 20, 'ENABLED', '新增活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7624067801880049144', 'RESOURCE', '202209', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:publish', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '发布消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7624773106003991812', 'RESOURCE', '203025', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:delete', 'CASCADE', 'DEPT', '[]', 0, 50, 'ENABLED', '删除目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7639169641875772298', 'RESOURCE', '203011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:page', 'CASCADE', 'ALL', '[]', 0, 10, 'ENABLED', '分页活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7654259910312149696', 'RESOURCE', '203013', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestactivity:detail', 'CASCADE', 'ALL', '[]', 0, 30, 'ENABLED', '详情活动', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7660302211516474641', 'RESOURCE', '203024', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:update', 'CASCADE', 'DEPT', '[]', 0, 40, 'ENABLED', '编辑目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7740028803530587951', 'RESOURCE', '204016', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:joblog:page', 'CASCADE', 'ALL', '[]', 0, 6, 'ENABLED', '执行日志', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7742840047784749526', 'RESOURCE', '202204', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7768912188750692632', 'RESOURCE', '203022', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:create', 'CASCADE', 'DEPT', '[]', 0, 20, 'ENABLED', '新增目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7792386902249912041', 'RESOURCE', '204015', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:run', 'CASCADE', 'ALL', '[]', 0, 5, 'ENABLED', '立即执行', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7822968206741092129', 'RESOURCE', '203026', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:list', 'CASCADE', 'DEPT', '[]', 0, 90, 'ENABLED', '树列表目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7837316234393882458', 'RESOURCE', '203031', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestorder:page', 'CASCADE', 'SELF', '[]', 0, 10, 'ENABLED', '分页订单', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7845666016635732956', 'RESOURCE', '203045', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:delete', 'CASCADE', 'DEPT_AND_CHILD', '[]', 0, 50, 'ENABLED', '删除知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7859474578876774469', 'RESOURCE', '201061', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:audit:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', NULL, NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7873407408257995473', 'RESOURCE', '204014', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:detail', 'CASCADE', 'ALL', '[]', 0, 4, 'ENABLED', '任务详情', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7927293661818445174', 'RESOURCE', '202203', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '详情消息', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 13:06:55.992098', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7939060249732762857', 'RESOURCE', '204011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:job:create', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '新增任务', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7986838434768267433', 'RESOURCE', '203021', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestcatalog:page', 'CASCADE', 'DEPT', '[]', 0, 10, 'ENABLED', '分页目录', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('7990168508290078017', 'RESOURCE', '203044', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'biz:cgtestknowledgecategory:update', 'CASCADE', 'DEPT_AND_CHILD', '[]', 0, 40, 'ENABLED', '编辑知识分类', NULL, NULL, '{}', '2026-08-08 13:06:56.015503', NULL, '2026-08-08 13:06:56.015503', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8106481288132251001', 'RESOURCE', '202231', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:realname:verify', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '审核实名认证', NULL, NULL, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000001', 'RESOURCE', '200007', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:page', 'CASCADE', 'CUSTOM', '[\"8200000000000102\", \"8200000000000103\"]', 0, 0, 'ENABLED', '账号管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000002', 'RESOURCE', '200008', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:dept:tree', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '部门管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000003', 'RESOURCE', '200009', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '用户组管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000004', 'RESOURCE', '200010', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:position:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '岗位管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000005', 'RESOURCE', '200011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '角色管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000006', 'RESOURCE', '200012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:list', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '资源管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000007', 'RESOURCE', '200012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:page', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '资源管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000008', 'RESOURCE', '200018', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resourcemodule:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '资源模块管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000009', 'RESOURCE', '200031', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientmodule:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '客户端模块管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000010', 'RESOURCE', '200032', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:list', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '客户端资源管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000011', 'RESOURCE', '200032', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:page', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '客户端资源管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000012', 'RESOURCE', '201101', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增账号', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000013', 'RESOURCE', '201102', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看账号', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000014', 'RESOURCE', '201103', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑账号', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000015', 'RESOURCE', '201104', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除账号', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000016', 'RESOURCE', '201105', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:grantrole', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配角色', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000017', 'RESOURCE', '201106', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:grantgroup', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配用户组', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000018', 'RESOURCE', '201107', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:grantdept', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配部门', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000019', 'RESOURCE', '201108', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:grantresource', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000020', 'RESOURCE', '201109', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:account:grantclientresource', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000021', 'RESOURCE', '201121', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:dept:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增部门', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000022', 'RESOURCE', '201122', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:dept:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看部门', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000023', 'RESOURCE', '201123', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:dept:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑部门', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000024', 'RESOURCE', '201124', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:dept:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除部门', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000025', 'RESOURCE', '201131', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增用户组', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000026', 'RESOURCE', '201132', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看用户组', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000027', 'RESOURCE', '201133', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑用户组', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000028', 'RESOURCE', '201134', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除用户组', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000029', 'RESOURCE', '201135', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:grantuser', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配用户', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000030', 'RESOURCE', '201136', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:grantrole', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配角色', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000031', 'RESOURCE', '201137', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:grantresource', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000032', 'RESOURCE', '201138', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:group:grantclientresource', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000033', 'RESOURCE', '201151', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:position:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增岗位', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000034', 'RESOURCE', '201152', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:position:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看岗位', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000035', 'RESOURCE', '201153', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:position:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑岗位', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000036', 'RESOURCE', '201154', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:position:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除岗位', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000037', 'RESOURCE', '201161', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增角色', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000038', 'RESOURCE', '201162', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看角色', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000039', 'RESOURCE', '201163', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑角色', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000040', 'RESOURCE', '201164', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除角色', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000041', 'RESOURCE', '201165', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:grantresource', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000042', 'RESOURCE', '201167', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:grantuser', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配用户', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000043', 'RESOURCE', '201168', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:role:grantclientresource', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分配客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000044', 'RESOURCE', '201181', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000045', 'RESOURCE', '201182', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000046', 'RESOURCE', '201183', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000047', 'RESOURCE', '201184', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000048', 'RESOURCE', '201186', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:list', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '资源树', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000049', 'RESOURCE', '201185', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resource:grant', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '绑定权限', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000050', 'RESOURCE', '201191', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resourcemodule:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增资源模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000051', 'RESOURCE', '201192', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resourcemodule:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看资源模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000052', 'RESOURCE', '201193', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resourcemodule:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑资源模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000053', 'RESOURCE', '201194', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:resourcemodule:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除资源模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000054', 'RESOURCE', '201311', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientmodule:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增客户端模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000055', 'RESOURCE', '201312', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientmodule:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看客户端模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000056', 'RESOURCE', '201313', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientmodule:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑客户端模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000057', 'RESOURCE', '201314', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientmodule:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除客户端模块', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000058', 'RESOURCE', '201321', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000059', 'RESOURCE', '201322', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000060', 'RESOURCE', '201323', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000061', 'RESOURCE', '201324', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除客户端资源', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000062', 'RESOURCE', '201325', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:list', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '客户端资源树', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107100000000063', 'RESOURCE', '201326', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'iam:clientresource:grant', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '绑定客户端资源权限', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000001', 'RESOURCE', '201011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:dict:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增字典', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000002', 'RESOURCE', '201012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:dict:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看字典', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000003', 'RESOURCE', '201013', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:dict:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑字典', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000004', 'RESOURCE', '201014', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:dict:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除字典', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000005', 'RESOURCE', '201021', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:banner:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增展示图', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000006', 'RESOURCE', '201022', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:banner:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看展示图', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000007', 'RESOURCE', '201023', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:banner:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑展示图', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000008', 'RESOURCE', '201024', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:banner:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除展示图', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000009', 'RESOURCE', '201031', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:file:upload', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '上传文件', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000010', 'RESOURCE', '201032', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:file:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看文件', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000011', 'RESOURCE', '201033', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:file:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑文件', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000012', 'RESOURCE', '201034', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:file:url', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '打开文件', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000013', 'RESOURCE', '201035', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:file:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除文件', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000014', 'RESOURCE', '201041', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'auth:session:tokenlist', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看令牌', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000015', 'RESOURCE', '201042', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'auth:session:exit', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '强退账号', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000016', 'RESOURCE', '201043', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'auth:session:tokenexit', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '强退令牌', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000017', 'RESOURCE', '201051', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增生成方案', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000018', 'RESOURCE', '201052', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看生成方案', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000019', 'RESOURCE', '201053', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑生成方案', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000020', 'RESOURCE', '201054', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除生成方案', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000021', 'RESOURCE', '201055', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:tables', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '读取数据库表', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000022', 'RESOURCE', '201056', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:preview', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '预览代码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000023', 'RESOURCE', '201057', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:download', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '下载代码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000024', 'RESOURCE', '202011', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:config:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增系统配置', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000025', 'RESOURCE', '202012', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:config:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看系统配置', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000026', 'RESOURCE', '202013', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:config:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑系统配置', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000027', 'RESOURCE', '202251', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:weakpassword:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分页弱密码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000028', 'RESOURCE', '202252', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:weakpassword:create', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '新增弱密码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000029', 'RESOURCE', '202253', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:weakpassword:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '编辑弱密码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000030', 'RESOURCE', '202254', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:weakpassword:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除弱密码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000031', 'RESOURCE', '202255', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:weakpassword:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看弱密码', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000032', 'RESOURCE', '202221', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:feedback:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '分页反馈', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000033', 'RESOURCE', '202222', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:feedback:detail', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '查看反馈', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000034', 'RESOURCE', '202223', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:feedback:update', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '处理反馈', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000035', 'RESOURCE', '202224', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:feedback:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除反馈', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000036', 'RESOURCE', '200004', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:dict:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '字典管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000037', 'RESOURCE', '200005', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:banner:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '展示图管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000038', 'RESOURCE', '200023', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:file:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '文件管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000039', 'RESOURCE', '200025', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'auth:session:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '在线会话访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000040', 'RESOURCE', '202230', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:realname:verify', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '实名认证审核访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000041', 'RESOURCE', '202015', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:codegen:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '代码生成访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000042', 'RESOURCE', '202010', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:config:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '系统配置访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000043', 'RESOURCE', '202010', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:weakpassword:page', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '系统配置访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000044', 'RESOURCE', '202220', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:feedback:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '反馈管理访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000045', 'RESOURCE', '202200', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:notice:page', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '通知消息访问', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8107200000000046', 'RESOURCE', '202014', 'ADMIN', 'RESOURCE_PERMISSION', 'PERMISSION', '', 'sys:config:delete', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '删除系统配置', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_iam_relation` VALUES ('8200000000000311', 'GROUP', '8200000000000301', 'ADMIN', 'GROUP_ROLE', 'ROLE', '4', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '研发组绑定业务-本部门角色', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000312', 'GROUP', '8200000000000304', 'ADMIN', 'GROUP_ROLE', 'ROLE', '5', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '前端组绑定业务-仅本人角色', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000313', 'GROUP', '8200000000000302', 'ADMIN', 'GROUP_ROLE', 'ROLE', '3', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '市场组绑定业务-全量角色', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000314', 'GROUP', '8200000000000303', 'ADMIN', 'GROUP_ROLE', 'ROLE', '7', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', 'IAM协作组绑定只读角色', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000315', 'GROUP', '8200000000000305', 'ADMIN', 'GROUP_ROLE', 'ROLE', '4', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '后端组绑定业务-本部门角色', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000316', 'GROUP', '8200000000000306', 'ADMIN', 'GROUP_ROLE', 'ROLE', '6', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '测试组绑定业务-部门及子部门角色', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000401', 'ACCOUNT', '8200000000000201', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000102', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000402', 'ACCOUNT', '8200000000000203', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000104', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000403', 'ACCOUNT', '8200000000000205', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000102', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000404', 'ACCOUNT', '8200000000000207', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000102', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000405', 'ACCOUNT', '8200000000000202', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000103', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000406', 'ACCOUNT', '8200000000000204', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000104', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000407', 'ACCOUNT', '8200000000000206', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000107', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000408', 'ACCOUNT', '8200000000000208', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000105', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000409', 'ACCOUNT', '8200000000000209', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000106', '', 'CASCADE', 'ALL', '[]', 1, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000410', 'ACCOUNT', '8200000000000203', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000102', '', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '前端组员工兼属研发部', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000411', 'ACCOUNT', '8200000000000208', 'ADMIN', 'ACCOUNT_DEPT', 'DEPT', '8200000000000102', '', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '后端组员工兼属研发部', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000501', 'ACCOUNT', '8200000000000201', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '2', '', 'CASCADE', 'ALL', '[]', 0, 1, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000502', 'ACCOUNT', '8200000000000202', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '3', '', 'CASCADE', 'ALL', '[]', 0, 2, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000503', 'ACCOUNT', '8200000000000203', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '4', '', 'CASCADE', 'ALL', '[]', 0, 3, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000504', 'ACCOUNT', '8200000000000204', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '5', '', 'CASCADE', 'ALL', '[]', 0, 4, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000505', 'ACCOUNT', '8200000000000205', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '6', '', 'CASCADE', 'ALL', '[]', 0, 5, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000506', 'ACCOUNT', '8200000000000206', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '7', '', 'CASCADE', 'ALL', '[]', 0, 6, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000507', 'ACCOUNT', '8200000000000208', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '4', '', 'CASCADE', 'ALL', '[]', 0, 7, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000508', 'ACCOUNT', '8200000000000209', 'ADMIN', 'ACCOUNT_ROLE', 'ROLE', '6', '', 'CASCADE', 'ALL', '[]', 0, 8, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000520', 'ACCOUNT', '8200000000000207', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000301', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000521', 'ACCOUNT', '8200000000000203', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000301', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000522', 'ACCOUNT', '8200000000000204', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000304', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000523', 'ACCOUNT', '8200000000000202', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000302', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000524', 'ACCOUNT', '8200000000000201', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000303', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000525', 'ACCOUNT', '8200000000000206', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000303', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000526', 'ACCOUNT', '8200000000000205', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000301', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '研发副主管加入研发组成员', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000527', 'ACCOUNT', '8200000000000208', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000305', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000528', 'ACCOUNT', '8200000000000208', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000301', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '后端组兼属研发大组', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000529', 'ACCOUNT', '8200000000000209', 'ADMIN', 'ACCOUNT_GROUP', 'GROUP', '8200000000000306', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000601', 'ROLE', '2', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200001', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '工作台', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000602', 'ROLE', '2', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200006', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '组织权限', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000603', 'ROLE', '3', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200001', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '工作台', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000604', 'ROLE', '3', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '202004', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '活动', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000605', 'ROLE', '4', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200001', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '工作台', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000606', 'ROLE', '4', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '202005', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '目录', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000607', 'ROLE', '5', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200001', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '工作台', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000608', 'ROLE', '5', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '202006', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '订单', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000609', 'ROLE', '6', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200001', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '工作台', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000610', 'ROLE', '6', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '202007', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '知识分类', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000611', 'ROLE', '7', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '200007', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '账号列表', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_iam_relation` VALUES ('8200000000000612', 'ROLE', '7', 'ADMIN', 'SUBJECT_RESOURCE_GRANT', 'RESOURCE', '201102', '', 'CASCADE', 'ALL', '[]', 0, 0, 'ENABLED', '账号详情', NULL, NULL, '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_job
-- ----------------------------
DROP TABLE IF EXISTS `sys_job`;
CREATE TABLE `sys_job`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `handler` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '任务处理器标识（Boot 为 JobHandler 全限定类名）',
  `trigger_type` varchar(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '触发类型：CRON（表达式）/ FIXED（固定间隔秒）',
  `trigger_config` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '触发配置：Cron 表达式或间隔秒数',
  `params` json NULL COMMENT '执行参数（JSON）',
  `last_run_time` datetime(6) NULL DEFAULT NULL COMMENT '上次调度执行时间',
  `next_run_time` datetime(6) NOT NULL COMMENT '下次计划执行时间',
  `last_result` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '上次执行结果摘要',
  `enabled` tinyint(1) NOT NULL COMMENT '是否启用调度：1 启用 / 0 停用',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '描述说明',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）'
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '内置任务' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_job
-- ----------------------------
INSERT INTO `sys_job` VALUES ('7541000000000000006', '注销账号清理', 'iam_account_purge_cancelled', 'CRON', '0 0 3 * * *', '{\"retentionDays\": 15}', '2026-08-16 09:53:27.108262', '2026-08-23 19:00:00.000000', 'purged=0', 1, '每日清理已取消且超过保留期的账号数据', 6, '2026-08-16 09:53:26.897926', NULL, '2026-08-16 09:53:27.152576', '1');
INSERT INTO `sys_job` VALUES ('7541000000000000007', '审计日志清理', 'sys_audit_log_cleanup', 'CRON', '0 0 4 * * *', '{\"batchSize\": 1000, \"loginRetentionDays\": 180, \"operationRetentionDays\": 365}', '2026-08-16 09:53:27.108262', '2026-08-23 20:00:00.000000', 'deletedLogin=0,deletedOperation=0', 1, '每日按保留天数清理过期登录与操作审计日志', 7, '2026-08-16 09:53:26.897926', NULL, '2026-08-16 09:53:27.152576', '1');
INSERT INTO `sys_job` VALUES ('7541000000000000008', '任务执行日志清理', 'sys_job_log_cleanup', 'CRON', '0 30 3 * * *', '{\"batchSize\": 1000, \"retentionDays\": 30}', '2026-08-16 09:53:27.104743', '2026-08-23 19:30:00.000000', 'deleted=0,retentionDays=30,batchSize=1000', 1, '按保留天数批量清理过期 sys_job_log', 8, '2026-08-16 09:53:26.903237', NULL, '2026-08-16 09:53:27.146058', '1');
INSERT INTO `sys_job` VALUES ('7541000000000000004', '审计告警', 'sys_audit_alert', 'FIXED', '300', '{}', '2026-08-23 09:53:13.513609', '2026-08-23 09:58:13.525501', 'done fired=0', 1, '按配置规则扫描审计日志并发送告警', 4, '2026-08-16 09:53:26.892521', NULL, '2026-08-23 09:53:13.526501', '1');
INSERT INTO `sys_job` VALUES ('7541000000000000001', '示例任务', 'sys_job_sample', 'FIXED', '60', '{}', '2026-08-23 09:53:52.846690', '2026-08-23 09:54:52.846690', 'echo: (empty)', 1, '演示调度链路：回显执行参数', 1, '2026-08-16 09:53:26.876522', NULL, '2026-08-23 09:53:52.848831', '1');
INSERT INTO `sys_job` VALUES ('7541000000000000003', 'Banner 互动计数刷库', 'sys_banner_flush_interactions', 'FIXED', '60', '{}', '2026-08-23 09:53:52.848831', '2026-08-23 09:54:52.850832', 'flushed=0', 1, '将 Redis 互动增量写入 sys_banner.interaction_count', 3, '2026-08-16 09:53:26.887732', NULL, '2026-08-23 09:53:52.852835', '1');
INSERT INTO `sys_job` VALUES ('7541000000000000002', 'Banner 状态同步', 'sys_banner_status_sync', 'FIXED', '60', '{}', '2026-08-23 09:53:52.848831', '2026-08-23 09:54:52.856832', 'expired=0,activated=0', 1, '按 start_at / end_at 激活或过期 Banner', 2, '2026-08-16 09:53:26.882777', NULL, '2026-08-23 09:53:52.858828', '1');

-- ----------------------------
-- Table structure for sys_job_log
-- ----------------------------
DROP TABLE IF EXISTS `sys_job_log`;
CREATE TABLE `sys_job_log`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `job_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '关联任务定义ID（sys_job.id）',
  `params` json NULL COMMENT '执行参数快照（JSON）',
  `started_at` datetime(6) NOT NULL COMMENT '本次执行开始时间',
  `duration_ms` bigint NULL DEFAULT NULL COMMENT '执行用时（毫秒）',
  `success` tinyint(1) NOT NULL COMMENT '执行结果：是否成功',
  `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '执行结果摘要或错误堆栈摘要',
  `executor` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '执行人：人工触发为账户ID，调度为 system',
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '客户端/实例 IP 地址',
  `process_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '执行进程 PID',
  `app_dir` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '执行实例应用目录',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  INDEX `idx_sys_job_log_started_at`(`started_at` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '任务执行日志' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_job_log
-- ----------------------------

-- ----------------------------
-- Table structure for sys_notice
-- ----------------------------
DROP TABLE IF EXISTS `sys_notice`;
CREATE TABLE `sys_notice`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `kind` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息种类：NOTIFICATION（通知）/ ANNOUNCEMENT（公告）',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容',
  `content_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '内容格式：TEXT/HTML/MARKDOWN 等',
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '通知分类编码',
  `severity` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '重要等级：INFO/WARNING/ERROR 等',
  `target_scope` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '投放范围：ALL/ACCOUNT/DEPT/ROLE 等',
  `target_account_types` json NOT NULL COMMENT '目标账户类型列表（JSON 数组）',
  `target_account_ids` json NOT NULL COMMENT '目标账户ID列表（JSON 数组）',
  `target_dept_ids` json NOT NULL COMMENT '目标部门ID列表（JSON 数组）',
  `target_role_ids` json NOT NULL COMMENT '目标角色ID列表（JSON 数组）',
  `publish_locations` json NOT NULL COMMENT '发布位置（公告）',
  `is_pinned` tinyint(1) NOT NULL COMMENT '是否置顶（公告）',
  `pinned_until` datetime(6) NULL DEFAULT NULL COMMENT '置顶截止时间',
  `sender_account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '发送方账户类型',
  `sender_account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '发送方账户ID',
  `source_type` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '来源业务模块标识',
  `source_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '来源业务记录ID',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '发布状态：DRAFT/PUBLISHED/REVOKED 等',
  `publish_at` datetime(6) NULL DEFAULT NULL COMMENT '计划/实际发布时间',
  `revoked_at` datetime(6) NULL DEFAULT NULL COMMENT '撤回时间',
  `expire_at` datetime(6) NULL DEFAULT NULL COMMENT '过期时间（公告有效截止）',
  `view_count` int NOT NULL COMMENT '浏览/查看次数',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_notice_status_kind_publish`(`status` ASC, `kind` ASC, `publish_at` ASC) USING BTREE,
  INDEX `idx_sys_notice_status_pinned_publish`(`status` ASC, `is_pinned` ASC, `publish_at` DESC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '系统公告/通知' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_notice
-- ----------------------------
INSERT INTO `sys_notice` VALUES ('7491842112464527360', 'NOTIFICATION', '急急急', '哈哈哈', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 13:05:46.000000', NULL, NULL, 0, '{}', '2026-08-08 13:05:51.295643', '1', '2026-08-08 13:05:51.295643', '1');
INSERT INTO `sys_notice` VALUES ('7491842211315884032', 'NOTIFICATION', '哈哈哈', '哈哈哈哈哈', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"ADMIN\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 13:06:10.000000', NULL, NULL, 0, '{}', '2026-08-08 13:06:14.871274', '1', '2026-08-08 13:06:14.871274', '1');
INSERT INTO `sys_notice` VALUES ('7491853809015291904', 'ANNOUNCEMENT', '欢迎使用 HEI 门户', '门户账号体系、个人中心与消息中心已就绪。如有问题可通过「我的反馈」提交。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{\"center\": true, \"workspace\": true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 11:52:19.890496', NULL, '2026-11-06 13:52:19.890496', 6, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 14:56:51.531823', '7491847383584804864');
INSERT INTO `sys_notice` VALUES ('7491853809015291905', 'ANNOUNCEMENT', '系统维护预告', '本周日 02:00-04:00 将进行例行维护，期间门户可能短暂不可用，请提前做好安排。', 'markdown', 'SYSTEM', 'WARNING', 'ALL', '[\"PORTAL\", \"ADMIN\"]', '[]', '[]', '[]', '{\"popup\": true, \"center\": true, \"workspace\": true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 12:52:19.890496', NULL, '2026-08-22 13:52:19.890496', 1, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 14:11:23.438139', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809015291906', 'ANNOUNCEMENT', '意见反馈功能上线', '现已支持在线提交意见反馈并查看处理进度。登录后打开用户菜单中的「我的反馈」即可使用。', 'text', 'SYSTEM', 'SUCCESS', 'ACCOUNT_TYPE', '[\"PORTAL\"]', '[]', '[]', '[]', '{\"center\": true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-08 13:22:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809044652032', 'NOTIFICATION', '账号安全提醒', '建议定期修改密码，并确保绑定的手机号与邮箱可用，以便找回账号。', 'text', 'SECURITY', 'WARNING', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, 'SYSTEM', NULL, 'PUBLISHED', '2026-08-08 13:32:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809044652033', 'NOTIFICATION', '新功能提示：消息中心', '右上角铃铛可查看未读通知与公告，支持一键全部已读。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, 'SYSTEM', NULL, 'PUBLISHED', '2026-08-08 13:42:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('7491853809044652034', 'NOTIFICATION', '管理端测试通知', '这是一条仅面向 ADMIN 的通知，用于验证账户类型过滤。', 'text', 'SYSTEM', 'INFO', 'ACCOUNT_TYPE', '[\"ADMIN\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, 'SYSTEM', NULL, 'PUBLISHED', '2026-08-08 13:47:19.890496', NULL, NULL, 0, '{}', '2026-08-08 13:52:19.983927', NULL, '2026-08-08 13:52:19.983927', NULL);
INSERT INTO `sys_notice` VALUES ('8300000000000201', 'NOTIFICATION', '门户新手指引', '完成个人资料填写后可解锁更多功能，建议绑定邮箱与手机号。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 0, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000202', 'NOTIFICATION', 'Bob 的专属通知', '这是面向门户账户的示例通知，用于演示消息列表。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 0, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000203', 'ANNOUNCEMENT', '门户春季活动', '春季主题活动已开启，欢迎参与线上打卡与积分兑换。', 'text', 'SYSTEM', 'SUCCESS', 'ALL', '[\"PORTAL\"]', '[]', '[]', '[]', '{\"center\": true, \"workspace\": true}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 3, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000204', 'NOTIFICATION', '管理端 IAM 更新', '角色、用户组与部门授权演示数据已就绪，可在组织权限模块查看。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"ADMIN\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 0, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000205', 'NOTIFICATION', '审计日志提醒', '关键操作将写入审计日志，请管理员定期查看异常登录与权限变更。', 'text', 'SECURITY', 'WARNING', 'ALL', '[\"ADMIN\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 0, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000206', 'ANNOUNCEMENT', '管理端版本说明', '本次演示环境包含完整 IAM 与内容运营模块，仅供本地开发验证。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"ADMIN\"]', '[]', '[]', '[]', '{\"workspace\": true}', 1, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 2, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000207', 'NOTIFICATION', '双端系统通知', '管理端与门户均可收到本条通知，用于验证目标账户类型多选。', 'text', 'SYSTEM', 'INFO', 'ALL', '[\"ADMIN\", \"PORTAL\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 0, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');
INSERT INTO `sys_notice` VALUES ('8300000000000208', 'NOTIFICATION', '待处理反馈提醒', '有新的用户反馈待处理，请前往反馈管理查看。', 'text', 'SYSTEM', 'WARNING', 'ALL', '[\"ADMIN\"]', '[]', '[]', '[]', '{}', 0, NULL, NULL, NULL, NULL, NULL, 'PUBLISHED', '2026-08-23 09:00:00.000000', NULL, NULL, 0, '{}', '2026-08-23 09:00:00.000000', '1', '2026-08-23 09:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_notice_read
-- ----------------------------
DROP TABLE IF EXISTS `sys_notice_read`;
CREATE TABLE `sys_notice_read`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `notice_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公告/通知ID（sys_notice.id）',
  `account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户类型：ADMIN（管理端）/ PORTAL（门户端）',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账户ID',
  `read_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '用户阅读时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_notice_read_account`(`notice_id` ASC, `account_type` ASC, `account_id` ASC) USING BTREE,
  INDEX `idx_sys_notice_read_account`(`account_type` ASC, `account_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '公告已读记录' ROW_FORMAT = Dynamic;

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
CREATE TABLE `sys_operation_audit_log`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `module` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '业务模块编码（如 sys、iam）',
  `resource_type` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '资源类型编码（如 SysAccount）',
  `resource_id` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '被操作资源主键ID',
  `action` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作动作编码',
  `summary` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作内容可读摘要',
  `before_data` json NULL COMMENT '变更前数据（JSON）',
  `after_data` json NULL COMMENT '变更后数据（JSON）',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作人账户ID',
  `account_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作人账户类型：ADMIN/PORTAL',
  `request_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '请求链路 ID（Trace）',
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '客户端/实例 IP 地址',
  `user_agent` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '客户端 User-Agent',
  `success` tinyint(1) NOT NULL COMMENT '是否成功：1 成功 / 0 失败',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '错误信息',
  `operator_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作人昵称快照（写入时落库）',
  `action_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作名称（前端展示）',
  `action_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作类型：CREATE/UPDATE/DELETE/QUERY/EXPORT/LOGIN/LOGOUT/OTHER',
  `module_label` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '操作模块中文展示名',
  `duration_ms` int NULL DEFAULT NULL COMMENT '耗时（毫秒）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_sys_operation_audit_account_id`(`account_id` ASC) USING BTREE,
  INDEX `idx_sys_operation_audit_created_at`(`created_at` ASC) USING BTREE,
  INDEX `idx_sys_operation_audit_module_action`(`module` ASC, `action` ASC) USING BTREE,
  INDEX `idx_sys_operation_audit_resource`(`resource_type` ASC, `resource_id` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作审计日志' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_operation_audit_log
-- ----------------------------
INSERT INTO `sys_operation_audit_log` VALUES ('2086348203017056258', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', 'aab6d330ea5648ac8609cedf04598930', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 07:05:48.561987');
INSERT INTO `sys_operation_audit_log` VALUES ('2086404723817709569', 'iam', 'sys_file', NULL, 'upload', 'POST /api/v1/admin/sys/file/upload', NULL, NULL, '1', 'admin', '9680f40a4f99443ca7ad392b4a8bcc01', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 10:50:24.191888');
INSERT INTO `sys_operation_audit_log` VALUES ('2086408062643171330', 'iam', 'sys_file', NULL, 'upload', 'POST /api/v1/admin/sys/file/upload', NULL, NULL, '1', 'admin', '0b0147e7fc334c46bd0ce20b4bd3e8e5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:03:40.112575');
INSERT INTO `sys_operation_audit_log` VALUES ('2086413233146265601', 'iam', 'auth', NULL, 'login', 'POST /api/v1/portal/login', NULL, NULL, '7491847383584804864', 'portal', '346ae88b022e4dcca715f45c9d1acd99', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:24:12.949234');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414670496473090', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/portal/logout', NULL, NULL, NULL, NULL, 'f69f316b0d1d456c96267ff44156cc38', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:29:55.529704');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414718772912129', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, 'b57f1bbb39e14d19ab83073979cb9ad9', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:30:07.299008');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414768404111361', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '0c33636455a84a968a602e8ca91efd15', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:30:19.117349');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414793913868289', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, '16ca7bf3218641f199759626faaba64d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:30:25.202419');
INSERT INTO `sys_operation_audit_log` VALUES ('2086414860225814529', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '3f32732033144eeeab4d6ab4fa098853', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:30:41.022353');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415187943563265', 'iam', 'sys_file', NULL, 'upload', 'POST /api/v1/admin/sys/file/upload', NULL, NULL, '1', 'admin', '94e68d53f1c54a39899e57e8777303a4', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:31:59.153041');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415233657282562', 'iam', 'iam_account', NULL, 'update', 'POST /api/v1/admin/sys/accounts/update', NULL, NULL, '1', 'admin', '74e6acb5b526402da0cb5ae3a6bfe261', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:32:10.050871');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415589262958594', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, '61422f3e916c4389b835b42cd7ee89a1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:33:34.827482');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415657378455554', 'iam', 'auth', NULL, 'forgot_password', 'POST /api/v1/admin/forgot-password', NULL, NULL, NULL, NULL, '7d41b61b6e01477490fb699253770a64', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:33:51.061495');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415774965768194', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '964907d0682a4dfea90556655625386e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:34:19.122222');
INSERT INTO `sys_operation_audit_log` VALUES ('2086415885812834306', 'iam', 'iam_account', NULL, 'update', 'POST /api/v1/admin/sys/accounts/update', NULL, NULL, '1', 'admin', '6627a2127a38469ebf9ea00ed3a6a4f9', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:34:45.539416');
INSERT INTO `sys_operation_audit_log` VALUES ('2086416001558847490', 'iam', 'auth', NULL, 'logout', 'POST /api/v1/admin/logout', NULL, NULL, NULL, NULL, 'd022212e771f473aa82cf2f1de58c3b0', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:35:13.139818');
INSERT INTO `sys_operation_audit_log` VALUES ('2086416040901419009', 'iam', 'auth', NULL, 'forgot_password', 'POST /api/v1/admin/forgot-password', NULL, NULL, NULL, NULL, '65956d3abbe443db8be96f482cdf7d60', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:35:22.522867');
INSERT INTO `sys_operation_audit_log` VALUES ('2086418125776633857', 'iam', 'auth', NULL, 'login', 'POST /api/v1/admin/login', NULL, NULL, '1', 'admin', '45ce46c4357e44de944d457b6c5f6bad', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 11:43:39.491727');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340161582051329', 'sys', 'sys_job', '7541000000000000001', 'enabled', '更新了任务 【示例任务】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:28:57.960251Z】修改为【2026-08-23T09:43:03.228986300+08:00】', '{\"id\": \"7541000000000000001\", \"name\": \"示例任务\", \"pkey\": \"7541000000000000001\", \"sort\": \"1\", \"enabled\": \"false\", \"lastResult\": \"echo: (无参数)\", \"description\": \"演示调度链路：回显执行参数\", \"lastRunTime\": \"2026-08-16T10:27:57.95913Z\", \"nextRunTime\": \"2026-08-16T10:28:57.960251Z\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"60\"}', '{\"id\": \"7541000000000000001\", \"name\": \"示例任务\", \"pkey\": \"7541000000000000001\", \"sort\": \"1\", \"enabled\": \"true\", \"lastResult\": \"echo: (无参数)\", \"description\": \"演示调度链路：回显执行参数\", \"lastRunTime\": \"2026-08-16T10:27:57.95913Z\", \"nextRunTime\": \"2026-08-23T09:43:03.2289863+08:00\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"60\"}', '1', 'admin', 'ceb5b229956b4e08a327f05b1940d49f', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 113, '2026-08-23 01:42:03.321620');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340178770309121', 'sys', 'sys_job', '7541000000000000002', 'enabled', '更新了任务 【Banner 状态同步】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:28:57.970142Z】修改为【2026-08-23T09:43:06.637890+08:00】', '{\"id\": \"7541000000000000002\", \"name\": \"Banner 状态同步\", \"pkey\": \"7541000000000000002\", \"sort\": \"2\", \"enabled\": \"false\", \"lastResult\": \"expired=0,activated=0\", \"description\": \"按 start_at / end_at 激活或过期 Banner\", \"lastRunTime\": \"2026-08-16T10:27:57.95913Z\", \"nextRunTime\": \"2026-08-16T10:28:57.970142Z\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"60\"}', '{\"id\": \"7541000000000000002\", \"name\": \"Banner 状态同步\", \"pkey\": \"7541000000000000002\", \"sort\": \"2\", \"enabled\": \"true\", \"lastResult\": \"expired=0,activated=0\", \"description\": \"按 start_at / end_at 激活或过期 Banner\", \"lastRunTime\": \"2026-08-16T10:27:57.95913Z\", \"nextRunTime\": \"2026-08-23T09:43:06.63789+08:00\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"60\"}', '1', 'admin', 'b6552564a29a40529d94f16c7878f149', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 34, '2026-08-23 01:42:06.668025');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340187511238657', 'sys', 'sys_job', '7541000000000000003', 'enabled', '更新了任务 【Banner 互动计数刷库】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:28:57.962143Z】修改为【2026-08-23T09:43:09.774690200+08:00】', '{\"id\": \"7541000000000000003\", \"name\": \"Banner 互动计数刷库\", \"pkey\": \"7541000000000000003\", \"sort\": \"3\", \"enabled\": \"false\", \"lastResult\": \"flushed=0\", \"description\": \"将 Redis 互动增量写入 sys_banner.interaction_count\", \"lastRunTime\": \"2026-08-16T10:27:57.95913Z\", \"nextRunTime\": \"2026-08-16T10:28:57.962143Z\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"60\"}', '{\"id\": \"7541000000000000003\", \"name\": \"Banner 互动计数刷库\", \"pkey\": \"7541000000000000003\", \"sort\": \"3\", \"enabled\": \"true\", \"lastResult\": \"flushed=0\", \"description\": \"将 Redis 互动增量写入 sys_banner.interaction_count\", \"lastRunTime\": \"2026-08-16T10:27:57.95913Z\", \"nextRunTime\": \"2026-08-23T09:43:09.7746902+08:00\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"60\"}', '1', 'admin', 'cbe02fe5d43f4565b03df93a7542b262', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 29, '2026-08-23 01:42:09.804635');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340204632387585', 'sys', 'sys_job', '7541000000000000004', 'enabled', '更新了任务 【审计告警】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:30:56.283976Z】修改为【2026-08-23T09:47:12.752911800+08:00】', '{\"id\": \"7541000000000000004\", \"name\": \"审计告警\", \"pkey\": \"7541000000000000004\", \"sort\": \"4\", \"enabled\": \"false\", \"lastResult\": \"done fired=0\", \"description\": \"按配置规则扫描审计日志并发送告警\", \"lastRunTime\": \"2026-08-16T10:25:56.244268Z\", \"nextRunTime\": \"2026-08-16T10:30:56.283976Z\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"300\"}', '{\"id\": \"7541000000000000004\", \"name\": \"审计告警\", \"pkey\": \"7541000000000000004\", \"sort\": \"4\", \"enabled\": \"true\", \"lastResult\": \"done fired=0\", \"description\": \"按配置规则扫描审计日志并发送告警\", \"lastRunTime\": \"2026-08-16T10:25:56.244268Z\", \"nextRunTime\": \"2026-08-23T09:47:12.7529118+08:00\", \"triggerType\": \"FIXED\", \"triggerConfig\": \"300\"}', '1', 'admin', 'c356ad9788d24a039753c2fddd69df27', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 23, '2026-08-23 01:42:12.777518');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340213335568386', 'sys', 'sys_job', '7541000000000000006', 'enabled', '更新了任务 【注销账号清理】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T19:00Z】修改为【2026-08-24T03:00+08:00】', '{\"id\": \"7541000000000000006\", \"name\": \"注销账号清理\", \"pkey\": \"7541000000000000006\", \"sort\": \"6\", \"enabled\": \"false\", \"lastResult\": \"purged=0\", \"description\": \"每日清理已取消且超过保留期的账号数据\", \"lastRunTime\": \"2026-08-16T09:53:27.108262Z\", \"nextRunTime\": \"2026-08-16T19:00:00Z\", \"triggerType\": \"CRON\", \"triggerConfig\": \"0 0 3 * * *\"}', '{\"id\": \"7541000000000000006\", \"name\": \"注销账号清理\", \"pkey\": \"7541000000000000006\", \"sort\": \"6\", \"enabled\": \"true\", \"lastResult\": \"purged=0\", \"description\": \"每日清理已取消且超过保留期的账号数据\", \"lastRunTime\": \"2026-08-16T09:53:27.108262Z\", \"nextRunTime\": \"2026-08-24T03:00:00+08:00\", \"triggerType\": \"CRON\", \"triggerConfig\": \"0 0 3 * * *\"}', '1', 'admin', 'acfbb59ca08946f2adb296e902127224', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 69, '2026-08-23 01:42:16.059081');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340230611906561', 'sys', 'sys_job', '7541000000000000007', 'enabled', '更新了任务 【审计日志清理】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T19:00Z】修改为【2026-08-24T04:00+08:00】', '{\"id\": \"7541000000000000007\", \"name\": \"审计日志清理\", \"pkey\": \"7541000000000000007\", \"sort\": \"7\", \"enabled\": \"false\", \"lastResult\": \"deletedLogin=0,deletedOperation=0\", \"description\": \"每日按保留天数清理过期登录与操作审计日志\", \"lastRunTime\": \"2026-08-16T09:53:27.108262Z\", \"nextRunTime\": \"2026-08-16T19:00:00Z\", \"triggerType\": \"CRON\", \"triggerConfig\": \"0 0 4 * * *\"}', '{\"id\": \"7541000000000000007\", \"name\": \"审计日志清理\", \"pkey\": \"7541000000000000007\", \"sort\": \"7\", \"enabled\": \"true\", \"lastResult\": \"deletedLogin=0,deletedOperation=0\", \"description\": \"每日按保留天数清理过期登录与操作审计日志\", \"lastRunTime\": \"2026-08-16T09:53:27.108262Z\", \"nextRunTime\": \"2026-08-24T04:00:00+08:00\", \"triggerType\": \"CRON\", \"triggerConfig\": \"0 0 4 * * *\"}', '1', 'admin', '791047e2934d45feb839c83cef5e0774', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 39, '2026-08-23 01:42:19.027121');
INSERT INTO `sys_operation_audit_log` VALUES ('2091340239294115841', 'sys', 'sys_job', '7541000000000000008', 'enabled', '更新了任务 【任务执行日志清理】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T19:30Z】修改为【2026-08-24T03:30+08:00】', '{\"id\": \"7541000000000000008\", \"name\": \"任务执行日志清理\", \"pkey\": \"7541000000000000008\", \"sort\": \"8\", \"enabled\": \"false\", \"lastResult\": \"deleted=0,retentionDays=30,batchSize=1000\", \"description\": \"按保留天数批量清理过期 sys_job_log\", \"lastRunTime\": \"2026-08-16T09:53:27.104743Z\", \"nextRunTime\": \"2026-08-16T19:30:00Z\", \"triggerType\": \"CRON\", \"triggerConfig\": \"0 30 3 * * *\"}', '{\"id\": \"7541000000000000008\", \"name\": \"任务执行日志清理\", \"pkey\": \"7541000000000000008\", \"sort\": \"8\", \"enabled\": \"true\", \"lastResult\": \"deleted=0,retentionDays=30,batchSize=1000\", \"description\": \"按保留天数批量清理过期 sys_job_log\", \"lastRunTime\": \"2026-08-16T09:53:27.104743Z\", \"nextRunTime\": \"2026-08-24T03:30:00+08:00\", \"triggerType\": \"CRON\", \"triggerConfig\": \"0 30 3 * * *\"}', '1', 'admin', 'bac8332d9ef944af873120cb4100111b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '启用任务', 'UPDATE', '系统 - 任务', 26, '2026-08-23 01:42:21.973647');
INSERT INTO `sys_operation_audit_log` VALUES ('2091341545769480194', 'iam', 'iam_resource', '200041', 'update', '更新了资源 【客户端资源授权】：【编码】从【client-resource-auth】修改为【client-resource】；【isAffix】从【false】修改为【true】；【isCache】从【false】修改为【true】；【名称】从【客户端资源授权】修改为【客户端资源】；【path】从【/client-resource-auth】修改为【/client-resource-】', '{\"id\": \"200041\", \"code\": \"client-resource-auth\", \"icon\": \"icon-park-outline:application-one\", \"name\": \"客户端资源授权\", \"path\": \"/client-resource-auth\", \"pkey\": \"200041\", \"sort\": \"16\", \"status\": \"ENABLED\", \"isAffix\": \"false\", \"isCache\": \"false\", \"children\": [], \"moduleId\": \"210001\", \"isVisible\": \"true\", \"description\": \"客户端模块与客户端资源授权配置\", \"resourceType\": \"CATALOG\"}', '{\"id\": \"200041\", \"code\": \"client-resource\", \"icon\": \"icon-park-outline:application-one\", \"name\": \"客户端资源\", \"path\": \"/client-resource-\", \"pkey\": \"200041\", \"sort\": \"16\", \"status\": \"ENABLED\", \"isAffix\": \"true\", \"isCache\": \"true\", \"children\": [], \"moduleId\": \"210001\", \"isVisible\": \"true\", \"description\": \"客户端模块与客户端资源授权配置\", \"resourceType\": \"CATALOG\"}', '1', 'admin', 'e81ab24901fb4d00adbae991cd626f7e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '更新资源', 'UPDATE', '权限 - 资源', 40, '2026-08-23 01:47:32.343592');
INSERT INTO `sys_operation_audit_log` VALUES ('2091446223026950146', 'auth', 'auth', '1', 'login', '账号 【superadmin】登录成功', NULL, '{\"account\": \"superadmin\", \"captchaId\": \"41985019ddb44bfd90cd33fdc7296ff9\", \"loginMode\": \"PASSWORD\", \"rememberMe\": \"true\", \"accountType\": \"ADMIN\", \"captchaValue\": \"WWF4\", \"identityType\": \"ACCOUNT\", \"passwordKeyId\": \"cdc8f476b1304fa49a6f4e47f0a91df2\"}', '1', 'admin', '0a65925e19824081b21936482bc1e0f1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, '超管', '登录', 'LOGIN', '认证 - 账号', 923, '2026-08-23 08:43:29.440247');
INSERT INTO `sys_operation_audit_log` VALUES ('2091452426947825665', 'auth', 'auth', '7491847383584804864', 'logout', '账号 【user】退出成功', NULL, NULL, '7491847383584804864', NULL, '05c1991cfa254e63b4066ac5266d0f82', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, NULL, '退出登录', 'LOGOUT', '认证 - 账号', 17, '2026-08-23 09:08:10.340537');
INSERT INTO `sys_operation_audit_log` VALUES ('2091458990962282498', 'auth', 'auth', '7491847383584804864', 'login', '账号 【user】登录成功', NULL, '{\"account\": \"user\", \"captchaId\": \"da56c65889eb482da13054c7bbb425eb\", \"loginMode\": \"PASSWORD\", \"rememberMe\": \"true\", \"accountType\": \"PORTAL\", \"captchaValue\": \"6DEL\", \"identityType\": \"ACCOUNT\", \"passwordKeyId\": \"cc1d583f9a9f4ba6adc50f04cc27e0da\"}', '7491847383584804864', 'portal', 'e5b0037f5e134acd9c1df647c3b7a2b5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 1, NULL, 'user-171fd244', '登录', 'LOGIN', '认证 - 账号', 368, '2026-08-23 09:34:13.934427');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461077930512386', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', '17d34fe1d62c44709b97058f93de7074', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 72, '2026-08-23 09:42:31.907535');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461086591750145', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', '1b808947b5ad4bc2b10928e044b9e7ab', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 49, '2026-08-23 09:42:33.535284');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461095290736642', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', 'd728393abec64fc79903a77d79d92897', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 45, '2026-08-23 09:42:35.475514');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461095504646146', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', 'c34d6a4a0be7453ca48186a9c3882e6e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 50, '2026-08-23 09:42:36.647415');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461104081997827', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', '432d7038c2524d8c918128e835ef910c', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 35, '2026-08-23 09:42:37.822688');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461104220409857', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', '2a806157e42548c5ae4684669902cf6a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 49, '2026-08-23 09:42:38.305821');
INSERT INTO `sys_operation_audit_log` VALUES ('2091461112780984321', 'real', 'real_name_case', NULL, 'init_third_party', '发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败', NULL, '{\"realName\": \"11\", \"documentNo\": \"11\", \"businessType\": \"ACCOUNT_VERIFY\", \"documentType\": \"ID_CARD\"}', '7491847383584804864', 'portal', 'f5f4d44ead234099af5b7bf062583924', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', 'user-171fd244', '发起第三方实名', 'OTHER', '实名认证 - 工单', 30, '2026-08-23 09:42:39.353753');
INSERT INTO `sys_operation_audit_log` VALUES ('7491824755243462656', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', NULL, NULL, NULL, 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 11:56:52.549136');
INSERT INTO `sys_operation_audit_log` VALUES ('7491826689375354880', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '997f4d2c30b34d2fab35091b27f20dc1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 12:04:33.755576');
INSERT INTO `sys_operation_audit_log` VALUES ('7491842112552607744', 'iam', 'create', NULL, 'post', 'POST /api/v1/admin/sys/notices/create', 'null', 'null', '1', 'ADMIN', '6c8b3cc0d20a4a60866b37c2e8a98baf', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:05:51.329499');
INSERT INTO `sys_operation_audit_log` VALUES ('7491842211382992896', 'iam', 'create', NULL, 'post', 'POST /api/v1/admin/sys/notices/create', 'null', 'null', '1', 'ADMIN', '606a87d34ba3408db97735ddeaabe10a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:06:14.892188');
INSERT INTO `sys_operation_audit_log` VALUES ('7491845090663665664', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/profile/avatar/upload', 'null', 'null', '1', 'ADMIN', 'e779ab8350f7430ca6ab3d3c5689a822', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:17:41.365887');
INSERT INTO `sys_operation_audit_log` VALUES ('7491847385635819520', 'auth', 'account', '7491847383584804864', 'register', 'Portal account registered', 'null', 'null', '7491847383584804864', 'PORTAL', '178b180693724ed4bbdef49d6411609c', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:26:48.530169');
INSERT INTO `sys_operation_audit_log` VALUES ('7491847464409042944', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'ef3e093461f64d63a33ceff7fa281b98', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:27:06.809483');
INSERT INTO `sys_operation_audit_log` VALUES ('7491849365728989184', 'iam', 'upload', NULL, 'post', 'POST /api/v1/portal/sys/file/upload', 'null', 'null', '7491847383584804864', 'PORTAL', 'ef951dc41e29426689d40f63b2650a2c', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:34:40.621796');
INSERT INTO `sys_operation_audit_log` VALUES ('7491849375170367488', 'iam', 'submit', NULL, 'post', 'POST /api/v1/portal/sys/feedbacks/submit', 'null', 'null', '7491847383584804864', 'PORTAL', '45261bcedc5c4900a5daa0e6b51ba112', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:34:42.871686');
INSERT INTO `sys_operation_audit_log` VALUES ('7491850767205326848', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '7491847383584804864', 'PORTAL', '9db3060fba264773ae39314a7a2a3e3d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '403', NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:40:14.759860');
INSERT INTO `sys_operation_audit_log` VALUES ('7491850787824525312', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', 'cbef755d2cf74399bbf78075b1b4a569', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:40:19.210215');
INSERT INTO `sys_operation_audit_log` VALUES ('7491850914173739008', 'iam', 'update', NULL, 'post', 'POST /api/v1/admin/sys/feedbacks/update', 'null', 'null', '1', 'ADMIN', '10c752edb8ee4d45b7a5868a9a5c0bee', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:40:49.799448');
INSERT INTO `sys_operation_audit_log` VALUES ('7491851004623904768', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', '8121b706c0ab43a59cf412d3d4a889d1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:41:10.980972');
INSERT INTO `sys_operation_audit_log` VALUES ('7491852911862013952', 'iam', 'upload', NULL, 'post', 'POST /api/v1/portal/profile/avatar/upload', 'null', 'null', '7491847383584804864', 'PORTAL', 'a340f1b7b3174d1c8b5918818666e40d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:48:46.085676');
INSERT INTO `sys_operation_audit_log` VALUES ('7491853936371097600', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/portal/sys/notices/read-all', 'null', 'null', '7491847383584804864', 'PORTAL', '4a07b49c65ab41aa8f6eb409aeb81f44', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 13:52:50.347166');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856130252136448', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'a38dac89228e43c4809a2d0ccddade60', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:01:33.410073');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856145267744768', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'e4a86183e35b4211b5d81c2d78ac128a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:01:36.989520');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856147004186624', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', '36b62e3e1dc54b438d1c0a694de8b974', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:01:37.403846');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856149260722176', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'f7d8fe609757442c9f42a9557ce49a4b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:01:37.941479');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856211504193536', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', '3b4a2d936c0f44d78d8c9dc24c2cd748', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:01:52.781516');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856220572278784', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'e8b49b5611db4c498bb8fcb2f0ddfe12', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:01:54.944100');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856359823171584', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'e95abd971e9b45f2a0e930e50c233117', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:02:28.143287');
INSERT INTO `sys_operation_audit_log` VALUES ('7491856638173962240', 'iam', 'read-all', NULL, 'post', 'POST /api/v1/admin/sys/notices/read-all', 'null', 'null', '1', 'ADMIN', 'a78eed380e434144a1b1e74729d6e54f', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:03:34.506926');
INSERT INTO `sys_operation_audit_log` VALUES ('7491867956285251584', 'auth', 'account', 'k-oPNnS5p0ls7b7JnuEO8yRR80EZ3FR4CCStRQvGdLo', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', '60488ac9dcec4a59bf287de53e535a92', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:48:32.938124');
INSERT INTO `sys_operation_audit_log` VALUES ('7491867956377526272', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', '60488ac9dcec4a59bf287de53e535a92', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:48:32.976961');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868063009316864', 'auth', 'account', 'P48I4wXyp2T1K6Gws0WwWa-48PnEE1AEQshPdhtbqko', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', '8aee2686706648ac8febf62481191666', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:48:58.392732');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868063068037120', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', '8aee2686706648ac8febf62481191666', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:48:58.413543');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868407806271488', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'd249fc3c14104900b1c5ab59acc4ebdc', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:50:20.145710');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868707275382784', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '3436832124ae49b7bd8885920146ed05', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:51:31.644571');
INSERT INTO `sys_operation_audit_log` VALUES ('7491868946300379136', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', 'd2b5c5b182d844089cde10f6802c702d', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:52:28.992926');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869125392965632', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', 'bfad8619bcb04e929dbe3b8f4ceeb34b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:53:11.692280');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869202706571264', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '54a3eb743686490089f5f64617d88f23', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:53:30.125444');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869288723357696', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', '94a81d39bd6d4179b431c24ba933c36b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:53:50.633310');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869364359241728', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '36a39cc7c1b04ed3bd3f7be81c619a76', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:54:08.665933');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869492054827008', 'auth', 'account', 'Lr0IMLzxN8NytPPOtKb-7axzFWyr4kB2Y1fiu1Ofk7w', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'f7e05c09efd5437c9fa4e07971765789', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:54:39.098545');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869492138713088', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'f7e05c09efd5437c9fa4e07971765789', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:54:39.131162');
INSERT INTO `sys_operation_audit_log` VALUES ('7491869580353314816', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'ed81b9adb9254c798c5389e9304fae3e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:54:59.829300');
INSERT INTO `sys_operation_audit_log` VALUES ('7491870806759415808', 'auth', 'account', '8GXa7mbl8xxMrkS8e50Kcnc6wjnIDWwHCEZ9RZkFi6Y', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', '5d92a1b221184b848bcad88547b77ec2', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:59:52.539321');
INSERT INTO `sys_operation_audit_log` VALUES ('7491870806889439232', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', '5d92a1b221184b848bcad88547b77ec2', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:59:52.592863');
INSERT INTO `sys_operation_audit_log` VALUES ('7491870836610277376', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '9f84348d2ea048569a529eb0db336177', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 14:59:59.115155');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871812532523008', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', '6505ddfea04149a39f336a1a62857bde', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:03:52.356699');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871872020336640', 'auth', 'account', '-1hh-Gk-I2IVhH5QvuGxD6aKn71kIqmWf2LHo6gW7tA', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'ae74d9a084604aa28c8f3c54af94de95', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:04:06.528253');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871872095834112', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'ae74d9a084604aa28c8f3c54af94de95', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:04:06.557240');
INSERT INTO `sys_operation_audit_log` VALUES ('7491871928467279872', 'auth', 'account', '7491847383584804864', 'forgot_password', 'PORTAL password reset requested', 'null', 'null', '7491847383584804864', 'PORTAL', '57b11ecf39894f2e8e40bcdd2f636fe6', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:04:19.399391');
INSERT INTO `sys_operation_audit_log` VALUES ('7491872608439390208', 'auth', 'account', '7491847383584804864', 'forgot_password', 'PORTAL password reset requested', 'null', 'null', '7491847383584804864', 'PORTAL', '6e791c69c0cc4d858d75af582256e3f4', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:07:01.437344');
INSERT INTO `sys_operation_audit_log` VALUES ('7491872894360899584', 'auth', 'account', '7491872891940786176', 'register', 'Portal account registered', 'null', 'null', '7491872891940786176', 'PORTAL', 'bd310a77940d4c1d8e3ad519c3cc09c0', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:08:10.284656');
INSERT INTO `sys_operation_audit_log` VALUES ('7491875322690949120', 'iam', 'send-login-code', NULL, 'post', 'POST /api/v1/portal/send-login-code', 'null', 'null', NULL, NULL, 'a270b9c191dd4879bbfbb36ba9ea9cd4', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '400', NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:17:49.244173');
INSERT INTO `sys_operation_audit_log` VALUES ('7491875361068830720', 'iam', 'send-login-code', NULL, 'post', 'POST /api/v1/portal/send-login-code', 'null', 'null', NULL, NULL, '2eb366121b4a49989fa62fa361f605f2', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:17:58.394406');
INSERT INTO `sys_operation_audit_log` VALUES ('7491875444590006272', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'c2b3d7256b6e4c54a489859978feb288', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:18:18.208300');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876150378123264', 'auth', 'account', 'cKRn3a7CVntuJV6gMUKhYxVdj66PajD-FmVnQE8WHGQ', 'logout', 'Logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'c81c429632dd4d6cb6c5e71ea025465b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:21:06.560599');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876150457815040', 'iam', 'logout', NULL, 'post', 'POST /api/v1/portal/logout', 'null', 'null', '7491847383584804864', 'PORTAL', 'c81c429632dd4d6cb6c5e71ea025465b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:21:06.597915');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876284402913280', 'auth', 'account', '218B1W9aw378shi2KJwdvkHL6dy4KOxhy0RThJIwBkU', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', '4c93ad1119ae46fda36dfcbce10e964a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:21:38.515325');
INSERT INTO `sys_operation_audit_log` VALUES ('7491876284495187968', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', '4c93ad1119ae46fda36dfcbce10e964a', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 15:21:38.555219');
INSERT INTO `sys_operation_audit_log` VALUES ('7491886508237017088', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '06ffe05bd308431a85c804da1deff708', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 16:02:15.699453');
INSERT INTO `sys_operation_audit_log` VALUES ('7491890281995022336', 'iam', 'interaction', NULL, 'post', 'POST /api/v1/portal/sys/banners/interaction', 'null', 'null', NULL, NULL, '557c0d9c0e7e4a0186aff815a7655afb', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 16:17:15.820139');
INSERT INTO `sys_operation_audit_log` VALUES ('7491905950954287104', 'iam', 'batch-save', NULL, 'post', 'POST /api/v1/admin/sys/config/batch-save', 'null', 'null', '1', 'ADMIN', 'fa3e17a02e6b4df3b99c9d7e4ebbf6ee', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 17:19:31.591092');
INSERT INTO `sys_operation_audit_log` VALUES ('7491906012094656512', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', 'ad21a931262a45928ee617d6a1425ece', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 17:19:46.166700');
INSERT INTO `sys_operation_audit_log` VALUES ('7491978456146837504', 'iam', 'interaction', NULL, 'post', 'POST /api/v1/portal/sys/banners/interaction', 'null', 'null', NULL, NULL, '605ec219cc7c4aa9988ea4a65a790f1e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:07:38.175092');
INSERT INTO `sys_operation_audit_log` VALUES ('7491978848712720384', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '2353489a36fc450ba6fc61e02c043450', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:09:10.963818');
INSERT INTO `sys_operation_audit_log` VALUES ('7491978930061246464', 'iam', 'delete', NULL, 'post', 'POST /api/v1/admin/sys/banners/delete', 'null', 'null', '1', 'ADMIN', '2f55e8ff734c4795bb34ee2d53d21f54', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:09:31.165058');
INSERT INTO `sys_operation_audit_log` VALUES ('7491984102703431680', 'auth', 'account', 'isGNdSOJtqTMMLd5tE89-zWE9E-gZTI2e9QtAN9O5IY', 'logout', 'Logout', 'null', 'null', '1', 'ADMIN', 'e91dd3e589f04d1fbd2847d159404e59', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:30:04.376178');
INSERT INTO `sys_operation_audit_log` VALUES ('7491984102850232320', 'iam', 'logout', NULL, 'post', 'POST /api/v1/admin/logout', 'null', 'null', '1', 'ADMIN', 'e91dd3e589f04d1fbd2847d159404e59', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:30:04.454205');
INSERT INTO `sys_operation_audit_log` VALUES ('7491984343368400896', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '565285e8a92d4db183eb55f22bf5ebd5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:31:01.163330');
INSERT INTO `sys_operation_audit_log` VALUES ('7491985237854060544', 'auth', 'account', '7491847383584804864', 'forgot_password', 'PORTAL password reset requested', 'null', 'null', '7491847383584804864', 'PORTAL', 'e2553b09785a4da580824d763b076a86', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:34:34.641992');
INSERT INTO `sys_operation_audit_log` VALUES ('7491985393911529472', 'auth', 'account', '7491847383584804864', 'reset_password', 'PORTAL password reset', 'null', 'null', '7491847383584804864', 'PORTAL', 'ef552be360104587bf1e3160d613aeb1', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:35:11.342559');
INSERT INTO `sys_operation_audit_log` VALUES ('7491985465474744320', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', 'b2cf5db7a39c42d792ee141212e7892e', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-08 22:35:29.009350');
INSERT INTO `sys_operation_audit_log` VALUES ('7492049918010503168', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', 'bb7ce0ff89594460ac78a96c5f680e69', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 02:51:35.362724');
INSERT INTO `sys_operation_audit_log` VALUES ('7492050216368123904', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '45af2d5c4eb3416e953db56f294eb386', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 02:52:47.145857');
INSERT INTO `sys_operation_audit_log` VALUES ('7492050342729920512', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '3104ce6d274243619c356adc1f04279b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 02:53:17.271997');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070224125145088', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', '187c6bf3c132445e87e101755a4a236b', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:12:16.604286');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070372163104768', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', 'bc76b2dc3b9e45108466763edffa9ba0', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:12:52.661485');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070399157645312', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '8f9d09659bae4815806de838cf2699f5', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:12:59.097782');
INSERT INTO `sys_operation_audit_log` VALUES ('7492070449363464192', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '9ce941c98b414b3cbc84e53cf7b25849', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:13:11.067709');
INSERT INTO `sys_operation_audit_log` VALUES ('7492074203068411904', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '69519ea4cf5542ccabebc8c6182125f6', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:28:06.020572');
INSERT INTO `sys_operation_audit_log` VALUES ('7492074211410882560', 'iam', 'upload', NULL, 'post', 'POST /api/v1/admin/sys/file/upload', 'null', 'null', '1', 'ADMIN', '4228001a67784ec7a6201e518de4fd92', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, '500', NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:28:08.009335');
INSERT INTO `sys_operation_audit_log` VALUES ('7492078890945560576', 'auth', 'account', '1', 'login', 'ADMIN login succeeded', 'null', 'null', '1', 'ADMIN', 'b4da294ce5f746989d1ae6bebad11c65', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:46:43.055912');
INSERT INTO `sys_operation_audit_log` VALUES ('7492079018624368640', 'auth', 'account', '7491847383584804864', 'login', 'PORTAL login succeeded', 'null', 'null', '7491847383584804864', 'PORTAL', '11a3c39a52e84f90bdd30204d1f90a90', '127.0.0.1', 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36', 0, NULL, NULL, NULL, NULL, NULL, NULL, '2026-08-09 04:47:13.496031');

-- ----------------------------
-- Table structure for sys_operation_audit_outbox
-- ----------------------------
DROP TABLE IF EXISTS `sys_operation_audit_outbox`;
CREATE TABLE `sys_operation_audit_outbox`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `payload` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '事件载荷（JSON）',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '投递状态：PENDING/CLAIMED/DONE/DEAD',
  `attempts` int NOT NULL COMMENT '重试次数',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `claimed_at` datetime(6) NULL DEFAULT NULL COMMENT '消费者认领时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '操作审计 Outbox' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_operation_audit_outbox
-- ----------------------------
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340157119311874', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000001\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【示例任务】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:28:57.960251Z】修改为【2026-08-23T09:43:03.228986300+08:00】\",\"before_data\":{\"description\":\"演示调度链路：回显执行参数\",\"enabled\":\"false\",\"id\":\"7541000000000000001\",\"lastResult\":\"echo: (无参数)\",\"lastRunTime\":\"2026-08-16T10:27:57.95913Z\",\"name\":\"示例任务\",\"nextRunTime\":\"2026-08-16T10:28:57.960251Z\",\"pkey\":\"7541000000000000001\",\"sort\":\"1\",\"triggerConfig\":\"60\",\"triggerType\":\"FIXED\"},\"after_data\":{\"description\":\"演示调度链路：回显执行参数\",\"enabled\":\"true\",\"id\":\"7541000000000000001\",\"lastResult\":\"echo: (无参数)\",\"lastRunTime\":\"2026-08-16T10:27:57.95913Z\",\"name\":\"示例任务\",\"nextRunTime\":\"2026-08-23T09:43:03.2289863+08:00\",\"pkey\":\"7541000000000000001\",\"sort\":\"1\",\"triggerConfig\":\"60\",\"triggerType\":\"FIXED\"},\"duration_ms\":\"113\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"ceb5b229956b4e08a327f05b1940d49f\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:03.321620Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:03.329829', '2026-08-23 01:42:04.187633');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340171061178370', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000002\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【Banner 状态同步】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:28:57.970142Z】修改为【2026-08-23T09:43:06.637890+08:00】\",\"before_data\":{\"description\":\"按 start_at / end_at 激活或过期 Banner\",\"enabled\":\"false\",\"id\":\"7541000000000000002\",\"lastResult\":\"expired=0,activated=0\",\"lastRunTime\":\"2026-08-16T10:27:57.95913Z\",\"name\":\"Banner 状态同步\",\"nextRunTime\":\"2026-08-16T10:28:57.970142Z\",\"pkey\":\"7541000000000000002\",\"sort\":\"2\",\"triggerConfig\":\"60\",\"triggerType\":\"FIXED\"},\"after_data\":{\"description\":\"按 start_at / end_at 激活或过期 Banner\",\"enabled\":\"true\",\"id\":\"7541000000000000002\",\"lastResult\":\"expired=0,activated=0\",\"lastRunTime\":\"2026-08-16T10:27:57.95913Z\",\"name\":\"Banner 状态同步\",\"nextRunTime\":\"2026-08-23T09:43:06.63789+08:00\",\"pkey\":\"7541000000000000002\",\"sort\":\"2\",\"triggerConfig\":\"60\",\"triggerType\":\"FIXED\"},\"duration_ms\":\"34\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"b6552564a29a40529d94f16c7878f149\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:06.668024600Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:06.670033', '2026-08-23 01:42:08.459353');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340184235487234', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000003\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【Banner 互动计数刷库】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:28:57.962143Z】修改为【2026-08-23T09:43:09.774690200+08:00】\",\"before_data\":{\"description\":\"将 Redis 互动增量写入 sys_banner.interaction_count\",\"enabled\":\"false\",\"id\":\"7541000000000000003\",\"lastResult\":\"flushed=0\",\"lastRunTime\":\"2026-08-16T10:27:57.95913Z\",\"name\":\"Banner 互动计数刷库\",\"nextRunTime\":\"2026-08-16T10:28:57.962143Z\",\"pkey\":\"7541000000000000003\",\"sort\":\"3\",\"triggerConfig\":\"60\",\"triggerType\":\"FIXED\"},\"after_data\":{\"description\":\"将 Redis 互动增量写入 sys_banner.interaction_count\",\"enabled\":\"true\",\"id\":\"7541000000000000003\",\"lastResult\":\"flushed=0\",\"lastRunTime\":\"2026-08-16T10:27:57.95913Z\",\"name\":\"Banner 互动计数刷库\",\"nextRunTime\":\"2026-08-23T09:43:09.7746902+08:00\",\"pkey\":\"7541000000000000003\",\"sort\":\"3\",\"triggerConfig\":\"60\",\"triggerType\":\"FIXED\"},\"duration_ms\":\"29\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"cbe02fe5d43f4565b03df93a7542b262\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:09.804635400Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:09.806935', '2026-08-23 01:42:10.547294');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340196726124546', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000004\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【审计告警】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T10:30:56.283976Z】修改为【2026-08-23T09:47:12.752911800+08:00】\",\"before_data\":{\"description\":\"按配置规则扫描审计日志并发送告警\",\"enabled\":\"false\",\"id\":\"7541000000000000004\",\"lastResult\":\"done fired=0\",\"lastRunTime\":\"2026-08-16T10:25:56.244268Z\",\"name\":\"审计告警\",\"nextRunTime\":\"2026-08-16T10:30:56.283976Z\",\"pkey\":\"7541000000000000004\",\"sort\":\"4\",\"triggerConfig\":\"300\",\"triggerType\":\"FIXED\"},\"after_data\":{\"description\":\"按配置规则扫描审计日志并发送告警\",\"enabled\":\"true\",\"id\":\"7541000000000000004\",\"lastResult\":\"done fired=0\",\"lastRunTime\":\"2026-08-16T10:25:56.244268Z\",\"name\":\"审计告警\",\"nextRunTime\":\"2026-08-23T09:47:12.7529118+08:00\",\"pkey\":\"7541000000000000004\",\"sort\":\"4\",\"triggerConfig\":\"300\",\"triggerType\":\"FIXED\"},\"duration_ms\":\"23\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"c356ad9788d24a039753c2fddd69df27\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:12.777518100Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:12.779528', '2026-08-23 01:42:14.628082');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340210462470145', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000006\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【注销账号清理】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T19:00Z】修改为【2026-08-24T03:00+08:00】\",\"before_data\":{\"description\":\"每日清理已取消且超过保留期的账号数据\",\"enabled\":\"false\",\"id\":\"7541000000000000006\",\"lastResult\":\"purged=0\",\"lastRunTime\":\"2026-08-16T09:53:27.108262Z\",\"name\":\"注销账号清理\",\"nextRunTime\":\"2026-08-16T19:00:00Z\",\"pkey\":\"7541000000000000006\",\"sort\":\"6\",\"triggerConfig\":\"0 0 3 * * *\",\"triggerType\":\"CRON\"},\"after_data\":{\"description\":\"每日清理已取消且超过保留期的账号数据\",\"enabled\":\"true\",\"id\":\"7541000000000000006\",\"lastResult\":\"purged=0\",\"lastRunTime\":\"2026-08-16T09:53:27.108262Z\",\"name\":\"注销账号清理\",\"nextRunTime\":\"2026-08-24T03:00:00+08:00\",\"pkey\":\"7541000000000000006\",\"sort\":\"6\",\"triggerConfig\":\"0 0 3 * * *\",\"triggerType\":\"CRON\"},\"duration_ms\":\"69\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"acfbb59ca08946f2adb296e902127224\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:16.059081400Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:16.061974', '2026-08-23 01:42:16.706826');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340222944718849', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000007\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【审计日志清理】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T19:00Z】修改为【2026-08-24T04:00+08:00】\",\"before_data\":{\"description\":\"每日按保留天数清理过期登录与操作审计日志\",\"enabled\":\"false\",\"id\":\"7541000000000000007\",\"lastResult\":\"deletedLogin=0,deletedOperation=0\",\"lastRunTime\":\"2026-08-16T09:53:27.108262Z\",\"name\":\"审计日志清理\",\"nextRunTime\":\"2026-08-16T19:00:00Z\",\"pkey\":\"7541000000000000007\",\"sort\":\"7\",\"triggerConfig\":\"0 0 4 * * *\",\"triggerType\":\"CRON\"},\"after_data\":{\"description\":\"每日按保留天数清理过期登录与操作审计日志\",\"enabled\":\"true\",\"id\":\"7541000000000000007\",\"lastResult\":\"deletedLogin=0,deletedOperation=0\",\"lastRunTime\":\"2026-08-16T09:53:27.108262Z\",\"name\":\"审计日志清理\",\"nextRunTime\":\"2026-08-24T04:00:00+08:00\",\"pkey\":\"7541000000000000007\",\"sort\":\"7\",\"triggerConfig\":\"0 0 4 * * *\",\"triggerType\":\"CRON\"},\"duration_ms\":\"39\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"791047e2934d45feb839c83cef5e0774\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:19.027120700Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:19.031692', '2026-08-23 01:42:20.807088');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091340235250806786', '{\"resource_type\":\"sys_job\",\"resource_id\":\"7541000000000000008\",\"action\":\"enabled\",\"action_name\":\"启用任务\",\"action_type\":\"UPDATE\",\"module_label\":\"系统 - 任务\",\"operator_name\":\"superadmin\",\"summary\":\"更新了任务 【任务执行日志清理】：【启用】从【false】修改为【true】；【nextRunTime】从【2026-08-16T19:30Z】修改为【2026-08-24T03:30+08:00】\",\"before_data\":{\"description\":\"按保留天数批量清理过期 sys_job_log\",\"enabled\":\"false\",\"id\":\"7541000000000000008\",\"lastResult\":\"deleted=0,retentionDays=30,batchSize=1000\",\"lastRunTime\":\"2026-08-16T09:53:27.104743Z\",\"name\":\"任务执行日志清理\",\"nextRunTime\":\"2026-08-16T19:30:00Z\",\"pkey\":\"7541000000000000008\",\"sort\":\"8\",\"triggerConfig\":\"0 30 3 * * *\",\"triggerType\":\"CRON\"},\"after_data\":{\"description\":\"按保留天数批量清理过期 sys_job_log\",\"enabled\":\"true\",\"id\":\"7541000000000000008\",\"lastResult\":\"deleted=0,retentionDays=30,batchSize=1000\",\"lastRunTime\":\"2026-08-16T09:53:27.104743Z\",\"name\":\"任务执行日志清理\",\"nextRunTime\":\"2026-08-24T03:30:00+08:00\",\"pkey\":\"7541000000000000008\",\"sort\":\"8\",\"triggerConfig\":\"0 30 3 * * *\",\"triggerType\":\"CRON\"},\"duration_ms\":\"26\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/jobs/enabled\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"bac8332d9ef944af873120cb4100111b\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:42:21.973646700Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:42:21.977119', '2026-08-23 01:42:22.895351');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091341537091465217', '{\"resource_type\":\"iam_resource\",\"resource_id\":\"200041\",\"action\":\"update\",\"action_name\":\"更新资源\",\"action_type\":\"UPDATE\",\"module_label\":\"权限 - 资源\",\"operator_name\":\"superadmin\",\"summary\":\"更新了资源 【客户端资源授权】：【编码】从【client-resource-auth】修改为【client-resource】；【isAffix】从【false】修改为【true】；【isCache】从【false】修改为【true】；【名称】从【客户端资源授权】修改为【客户端资源】；【path】从【/client-resource-auth】修改为【/client-resource-】\",\"before_data\":{\"children\":[],\"code\":\"client-resource-auth\",\"description\":\"客户端模块与客户端资源授权配置\",\"icon\":\"icon-park-outline:application-one\",\"id\":\"200041\",\"isAffix\":\"false\",\"isCache\":\"false\",\"isVisible\":\"true\",\"moduleId\":\"210001\",\"name\":\"客户端资源授权\",\"path\":\"/client-resource-auth\",\"pkey\":\"200041\",\"resourceType\":\"CATALOG\",\"sort\":\"16\",\"status\":\"ENABLED\"},\"after_data\":{\"children\":[],\"code\":\"client-resource\",\"description\":\"客户端模块与客户端资源授权配置\",\"icon\":\"icon-park-outline:application-one\",\"id\":\"200041\",\"isAffix\":\"true\",\"isCache\":\"true\",\"isVisible\":\"true\",\"moduleId\":\"210001\",\"name\":\"客户端资源\",\"path\":\"/client-resource-\",\"pkey\":\"200041\",\"resourceType\":\"CATALOG\",\"sort\":\"16\",\"status\":\"ENABLED\"},\"duration_ms\":\"40\",\"method\":\"POST\",\"path\":\"/api/v1/admin/sys/resources/update\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"e81ab24901fb4d00adbae991cd626f7e\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T01:47:32.343591800Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 01:47:32.347052', '2026-08-23 01:47:34.360242');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091446214722228225', '{\"resource_type\":\"auth\",\"resource_id\":\"1\",\"action\":\"login\",\"action_name\":\"登录\",\"action_type\":\"LOGIN\",\"module_label\":\"认证 - 账号\",\"operator_name\":\"superadmin\",\"summary\":\"账号 【superadmin】登录成功\",\"before_data\":null,\"after_data\":{\"account\":\"superadmin\",\"accountType\":\"ADMIN\",\"captchaId\":\"41985019ddb44bfd90cd33fdc7296ff9\",\"captchaValue\":\"WWF4\",\"identityType\":\"ACCOUNT\",\"loginMode\":\"PASSWORD\",\"passwordKeyId\":\"cdc8f476b1304fa49a6f4e47f0a91df2\",\"rememberMe\":\"true\"},\"duration_ms\":\"923\",\"method\":\"POST\",\"path\":\"/api/v1/admin/login\",\"status_code\":\"200\",\"account_id\":\"1\",\"account_type\":\"admin\",\"request_id\":\"0a65925e19824081b21936482bc1e0f1\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T08:43:29.440247Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 08:43:29.445342', '2026-08-23 08:43:31.322108');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091452426104770562', '{\"resource_type\":\"auth\",\"resource_id\":\"7491847383584804864\",\"action\":\"logout\",\"action_name\":\"退出登录\",\"action_type\":\"LOGOUT\",\"module_label\":\"认证 - 账号\",\"operator_name\":null,\"summary\":\"账号 【user】退出成功\",\"before_data\":null,\"after_data\":null,\"duration_ms\":\"17\",\"method\":\"POST\",\"path\":\"/api/v1/portal/logout\",\"status_code\":\"200\",\"account_id\":\"7491847383584804864\",\"account_type\":null,\"request_id\":\"05c1991cfa254e63b4066ac5266d0f82\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:08:10.340536600Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:08:10.349472', '2026-08-23 09:08:10.509557');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091458984247201793', '{\"resource_type\":\"auth\",\"resource_id\":\"7491847383584804864\",\"action\":\"login\",\"action_name\":\"登录\",\"action_type\":\"LOGIN\",\"module_label\":\"认证 - 账号\",\"operator_name\":\"user\",\"summary\":\"账号 【user】登录成功\",\"before_data\":null,\"after_data\":{\"account\":\"user\",\"accountType\":\"PORTAL\",\"captchaId\":\"da56c65889eb482da13054c7bbb425eb\",\"captchaValue\":\"6DEL\",\"identityType\":\"ACCOUNT\",\"loginMode\":\"PASSWORD\",\"passwordKeyId\":\"cc1d583f9a9f4ba6adc50f04cc27e0da\",\"rememberMe\":\"true\"},\"duration_ms\":\"368\",\"method\":\"POST\",\"path\":\"/api/v1/portal/login\",\"status_code\":\"200\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"e5b0037f5e134acd9c1df647c3b7a2b5\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:34:13.934426500Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:34:13.934427', '2026-08-23 09:34:15.500046');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461072914124801', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"72\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"17d34fe1d62c44709b97058f93de7074\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:31.907535300Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:31.907535', '2026-08-23 09:42:33.082755');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461079742451713', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"49\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"1b808947b5ad4bc2b10928e044b9e7ab\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:33.535284400Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:33.535284', '2026-08-23 09:42:35.144068');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461087837458434', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"45\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"d728393abec64fc79903a77d79d92897\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:35.475514300Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:35.475514', '2026-08-23 09:42:37.209070');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461092795125761', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"50\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"c34d6a4a0be7453ca48186a9c3882e6e\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:36.647415Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:36.647415', '2026-08-23 09:42:37.209070');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461097727627266', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"35\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"432d7038c2524d8c918128e835ef910c\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:37.822687600Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:37.824691', '2026-08-23 09:42:39.314852');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461099749281793', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"49\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"2a806157e42548c5ae4684669902cf6a\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:38.305820700Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:38.305821', '2026-08-23 09:42:39.314852');
INSERT INTO `sys_operation_audit_outbox` VALUES ('2091461104144912386', '{\"resource_type\":\"real_name_case\",\"resource_id\":null,\"action\":\"init_third_party\",\"action_name\":\"发起第三方实名\",\"action_type\":\"OTHER\",\"module_label\":\"实名认证 - 工单\",\"operator_name\":\"user\",\"summary\":\"发起第三方实名认证 【user】（认证方式：人工审核，证件类型：身份证）失败\",\"before_data\":null,\"after_data\":{\"businessType\":\"ACCOUNT_VERIFY\",\"documentNo\":\"11\",\"documentType\":\"ID_CARD\",\"realName\":\"11\"},\"duration_ms\":\"30\",\"method\":\"POST\",\"path\":\"/api/v1/portal/real-name/case/init-third-party\",\"status_code\":\"500\",\"account_id\":\"7491847383584804864\",\"account_type\":\"portal\",\"request_id\":\"f5f4d44ead234099af5b7bf062583924\",\"ip\":\"127.0.0.1\",\"user_agent\":\"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/151.0.0.0 Safari/537.36\",\"occurred_at\":\"2026-08-23T09:42:39.353752800Z\",\"message_id\":null}', 'DONE', 1, '2026-08-23 09:42:39.353753', '2026-08-23 09:42:41.392377');

-- ----------------------------
-- Table structure for sys_position
-- ----------------------------
DROP TABLE IF EXISTS `sys_position`;
CREATE TABLE `sys_position`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '职位名称',
  `category` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '职位类别',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据权限范围）',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_virtual` tinyint(1) NOT NULL COMMENT '是否虚拟组织：1 虚拟 / 0 实体',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '职位状态：ENABLED/DISABLED',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '职位描述',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '职位' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_position
-- ----------------------------
INSERT INTO `sys_position` VALUES ('8200000000000801', '研发总监', 'MANAGEMENT', '8200000000000102', 1, 0, 'ENABLED', '研发部管理岗', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000802', '高级工程师', 'TECHNICAL', '8200000000000102', 2, 0, 'ENABLED', '研发部技术骨干', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000803', '前端工程师', 'TECHNICAL', '8200000000000104', 1, 0, 'ENABLED', '前端组默认岗位', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000804', '后端工程师', 'TECHNICAL', '8200000000000105', 1, 0, 'ENABLED', '后端组默认岗位', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000805', '测试工程师', 'TECHNICAL', '8200000000000106', 1, 0, 'ENABLED', '测试组默认岗位', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000806', '市场经理', 'MANAGEMENT', '8200000000000103', 1, 0, 'ENABLED', '市场部管理岗', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000807', '销售专员', 'OPERATION', '8200000000000103', 2, 0, 'ENABLED', '市场部业务岗', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_position` VALUES ('8200000000000808', '人事专员', 'SUPPORT', '8200000000000107', 1, 0, 'ENABLED', '人事行政支持岗', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_resource
-- ----------------------------
DROP TABLE IF EXISTS `sys_resource`;
CREATE TABLE `sys_resource`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `parent_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '父级资源ID（菜单树）',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `resource_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源类型：MENU/BUTTON/API 等',
  `module_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属资源模块ID',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '路径',
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '前端路由组件路径',
  `redirect` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '路由重定向地址',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标标识',
  `color` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '颜色值',
  `href` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '外链地址',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `is_visible` tinyint(1) NOT NULL COMMENT '是否可见：1 可见 / 0 隐藏',
  `is_cache` tinyint(1) NOT NULL COMMENT '是否缓存路由：1 缓存 / 0 不缓存',
  `is_affix` tinyint(1) NOT NULL COMMENT '是否固定标签页：1 固定 / 0 不固定',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '资源状态：ENABLED/DISABLED',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '资源描述说明',
  `layout` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '页面布局类型',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_resource_module_id_code`(`module_id` ASC, `code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '菜单/按钮/API 资源' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_resource
-- ----------------------------
INSERT INTO `sys_resource` VALUES ('200001', NULL, 'workspace', '工作台', 'MENU', '210001', '/workspace', '/workspace/index.vue', NULL, 'icon-park-outline:analysis', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200002', '200001', 'workspace-overview', '查看工作台总览', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200003', NULL, 'ops', '系统运维', 'CATALOG', '210001', '/sys', NULL, NULL, 'icon-park-outline:setting-two', NULL, NULL, 25, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200004', '200003', 'sys-dict', '字典管理', 'MENU', '210001', '/sys/dict', '/sys/dict/index.vue', NULL, 'icon-park-outline:file-search', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200005', '200019', 'content-banner', '展示图管理', 'MENU', '210001', '/sys/banner', '/sys/banner/index.vue', NULL, 'icon-park-outline:ad-product', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200006', NULL, 'org', '组织权限', 'CATALOG', '210001', '/iam', NULL, NULL, 'icon-park-outline:people', NULL, NULL, 10, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200007', '200006', 'iam-account', '账号管理', 'MENU', '210001', '/iam/account', '/iam/account/index.vue', NULL, 'icon-park-outline:people', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200008', '200006', 'iam-dept', '部门管理', 'MENU', '210001', '/iam/dept', '/iam/dept/index.vue', NULL, 'icon-park-outline:tree-diagram', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200009', '200006', 'iam-group', '用户组管理', 'MENU', '210001', '/iam/group', '/iam/group/index.vue', NULL, 'icon-park-outline:group', NULL, NULL, 3, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200010', '200006', 'iam-position', '岗位管理', 'MENU', '210001', '/iam/position', '/iam/position/index.vue', NULL, 'icon-park-outline:people-bottom', NULL, NULL, 4, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200011', '200006', 'iam-role', '角色管理', 'MENU', '210001', '/iam/role', '/iam/role/index.vue', NULL, 'icon-park-outline:peoples', NULL, NULL, 5, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200012', '200040', 'iam-resource', '资源管理', 'MENU', '210001', '/iam/resource', '/iam/resource/index.vue', NULL, 'icon-park-outline:all-application', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200018', '200040', 'iam-resourcemodule', '资源模块管理', 'MENU', '210001', '/iam/resource_module', '/iam/resource_module/index.vue', NULL, 'icon-park-outline:blocks-and-arrows', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200019', NULL, 'content', '内容运营', 'CATALOG', '210001', '/content', NULL, '/sys/notice', 'icon-park-outline:picture-album', NULL, NULL, 20, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200023', '200003', 'sys-file', '文件管理', 'MENU', '210001', '/sys/file', '/sys/file/index.vue', NULL, 'icon-park-outline:file-code', NULL, NULL, 3, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200025', '200003', 'sys-session', '在线会话', 'MENU', '210001', '/sys/session', '/auth/session/index.vue', NULL, 'icon-park-outline:connection', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-06-30 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200027', '200003', 'sys-audit-api', '操作审计接口', 'API_GROUP', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', '操作审计后端权限组', NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200028', '200003', 'sys-login-log', '登录日志', 'MENU', '210001', '/sys/login-log', '/sys/login-log/index.vue', NULL, 'icon-park-outline:log', NULL, NULL, 5, 1, 0, 0, 'ENABLED', '登录成功/失败历史记录', NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200029', '200003', 'sys-audit', '操作审计', 'MENU', '210001', '/sys/audit', '/sys/audit/index.vue', NULL, 'icon-park-outline:audit', NULL, NULL, 7, 1, 0, 0, 'ENABLED', '系统操作审计日志', NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200031', '200041', 'iam-clientmodule', '客户端模块管理', 'MENU', '210001', '/iam/client_module', '/iam/client_module/index.vue', NULL, 'icon-park-outline:application-one', NULL, NULL, 1, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200032', '200041', 'iam-clientresource', '客户端资源管理', 'MENU', '210001', '/iam/client_resource', '/iam/client_resource/index.vue', NULL, 'icon-park-outline:page-template', NULL, NULL, 2, 1, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200040', NULL, 'resource-auth', '资源授权', 'CATALOG', '210001', '/resource-auth', NULL, NULL, 'icon-park-outline:all-application', NULL, NULL, 15, 1, 0, 0, 'ENABLED', '菜单资源与资源模块授权配置', NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('200041', NULL, 'client-resource', '客户端资源', 'CATALOG', '210001', '/client-resource-', NULL, NULL, 'icon-park-outline:application-one', NULL, NULL, 16, 1, 1, 1, 'ENABLED', '客户端模块与客户端资源授权配置', NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_resource` VALUES ('201011', '200004', 'sys-dict-create', '新增字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201012', '200004', 'sys-dict-detail', '查看字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201013', '200004', 'sys-dict-update', '编辑字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201014', '200004', 'sys-dict-delete', '删除字典', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201021', '200005', 'sys-banner-create', '新增展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201022', '200005', 'sys-banner-detail', '查看展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201023', '200005', 'sys-banner-update', '编辑展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201024', '200005', 'sys-banner-delete', '删除展示图', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201025', '200005', 'sys-banner-create-page', '新增展示图页', 'PAGE', '210001', '/sys/banner/create', '/sys/banner/form.vue', NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201026', '200005', 'sys-banner-edit-page', '编辑展示图页', 'PAGE', '210001', '/sys/banner/edit', '/sys/banner/form.vue', NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201027', '200005', 'sys-banner-detail-page', '展示图详情页', 'PAGE', '210001', '/sys/banner/detail', '/sys/banner/detail.vue', NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201031', '200023', 'sys-file-upload', '上传文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201032', '200023', 'sys-file-detail', '查看文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201033', '200023', 'sys-file-update', '编辑文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201034', '200023', 'sys-file-url', '打开文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201035', '200023', 'sys-file-delete', '删除文件', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201041', '200025', 'sys-session-tokenlist', '查看令牌', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201042', '200025', 'sys-session-exit', '强退账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201043', '200025', 'sys-session-tokenexit', '强退令牌', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201051', '202015', 'sys-codegen-create', '新增生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201052', '202015', 'sys-codegen-detail', '查看生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201053', '202015', 'sys-codegen-update', '编辑生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201054', '202015', 'sys-codegen-delete', '删除生成方案', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201055', '202015', 'sys-codegen-tables', '读取数据库表', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201056', '202015', 'sys-codegen-preview', '预览代码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 60, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201057', '202015', 'sys-codegen-download', '下载代码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 70, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-07-18 16:10:45.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201060', '200028', 'sys-login-log-detail', '查看登录日志', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201061', '200029', 'sys-audit-detail', '查看审计详情', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201101', '200007', 'iam-account-create', '新增账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201102', '200007', 'iam-account-detail', '查看账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201103', '200007', 'iam-account-update', '编辑账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201104', '200007', 'iam-account-delete', '删除账号', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201105', '200007', 'iam-account-grant-role', '分配角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201106', '200007', 'iam-account-grant-group', '分配用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201107', '200007', 'iam-account-grant-dept', '分配部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201108', '200007', 'iam-account-grant-resource', '分配资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201109', '200007', 'iam-account-grant-client-resource', '分配客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 9, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201121', '200008', 'iam-dept-create', '新增部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201122', '200008', 'iam-dept-detail', '查看部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201123', '200008', 'iam-dept-update', '编辑部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201124', '200008', 'iam-dept-delete', '删除部门', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201131', '200009', 'iam-group-create', '新增用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201132', '200009', 'iam-group-detail', '查看用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201133', '200009', 'iam-group-update', '编辑用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201134', '200009', 'iam-group-delete', '删除用户组', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201135', '200009', 'iam-group-grant-user', '分配用户', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201136', '200009', 'iam-group-grant-role', '分配角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201137', '200009', 'iam-group-grant-resource', '分配资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201138', '200009', 'iam-group-grant-client-resource', '分配客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201151', '200010', 'iam-position-create', '新增岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201152', '200010', 'iam-position-detail', '查看岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201153', '200010', 'iam-position-update', '编辑岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201154', '200010', 'iam-position-delete', '删除岗位', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201161', '200011', 'iam-role-create', '新增角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201162', '200011', 'iam-role-detail', '查看角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201163', '200011', 'iam-role-update', '编辑角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201164', '200011', 'iam-role-delete', '删除角色', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201165', '200011', 'iam-role-grant-resource', '分配资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201167', '200011', 'iam-role-grant-user', '分配用户', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201168', '200011', 'iam-role-grant-client-resource', '分配客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201181', '200012', 'iam-resource-create', '新增资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201182', '200012', 'iam-resource-detail', '查看资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201183', '200012', 'iam-resource-update', '编辑资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201184', '200012', 'iam-resource-delete', '删除资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201185', '200012', 'iam-resource-grant', '绑定权限', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201186', '200012', 'iam-resource-list', '资源树', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201191', '200018', 'iam-resourcemodule-create', '新增资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201192', '200018', 'iam-resourcemodule-detail', '查看资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201193', '200018', 'iam-resourcemodule-update', '编辑资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201194', '200018', 'iam-resourcemodule-delete', '删除资源模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-03 00:00:00.000000', NULL, '2026-07-03 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201311', '200031', 'iam-clientmodule-create', '新增客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201312', '200031', 'iam-clientmodule-detail', '查看客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201313', '200031', 'iam-clientmodule-update', '编辑客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201314', '200031', 'iam-clientmodule-delete', '删除客户端模块', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201321', '200032', 'iam-clientresource-create', '新增客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201322', '200032', 'iam-clientresource-detail', '查看客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201323', '200032', 'iam-clientresource-update', '编辑客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201324', '200032', 'iam-clientresource-delete', '删除客户端资源', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201325', '200032', 'iam-clientresource-list', '客户端资源树', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('201326', '200032', 'iam-clientresource-grant', '绑定客户端资源权限', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202001', NULL, 'devtools', '开发工具', 'CATALOG', '210001', '/test', NULL, '/sys/codegen', 'icon-park-outline:code', NULL, NULL, 90, 1, 0, 0, 'ENABLED', '系统模块测试页面目录', NULL, '{}', '2026-07-18 12:39:16.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202002', '202001', 'system-test-editor', '编辑器测试', 'MENU', '210001', '/test/editor', '/test/editor/index.vue', NULL, 'icon-park-outline:edit', NULL, NULL, 2, 1, 0, 0, 'ENABLED', 'Markdown、富文本和代码编辑器组件测试页面', NULL, '{}', '2026-07-18 12:39:16.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202003', '202001', 'system-test-icon', '图标选择器测试', 'MENU', '210001', '/test/icon', '/test/icon/index.vue', NULL, 'icon-park-outline:all-application', NULL, NULL, 3, 1, 0, 0, 'ENABLED', 'Iconify 离线图标选择器测试页面', NULL, '{}', '2026-07-18 12:49:42.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202004', '202030', 'biz-cgtestactivity', '代码生成测试-活动', 'MENU', '210001', '/biz/cg-test-activity', '/biz/cg-test-activity/index.vue', NULL, 'icon-park-outline:calendar', NULL, NULL, 1, 1, 0, 0, 'ENABLED', '代码生成 CRUD 样例', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202005', '202030', 'biz-cgtestcatalog', '代码生成测试-目录树', 'MENU', '210001', '/biz/cg-test-catalog', '/biz/cg-test-catalog/index.vue', NULL, 'icon-park-outline:tree-list', NULL, NULL, 2, 1, 0, 0, 'ENABLED', '代码生成树表样例', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202006', '202030', 'biz-cgtestorder', '代码生成测试-订单', 'MENU', '210001', '/biz/cg-test-order', '/biz/cg-test-order/index.vue', NULL, 'icon-park-outline:transaction-order', NULL, NULL, 3, 1, 0, 0, 'ENABLED', '代码生成主子表样例', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202007', '202030', 'biz-cgtestknowledgecategory', '代码生成测试-知识分类', 'MENU', '210001', '/biz/cg-test-knowledge-category', '/biz/cg-test-knowledge-category/index.vue', NULL, 'icon-park-outline:book-open', NULL, NULL, 4, 1, 0, 0, 'ENABLED', '代码生成左树右表样例', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202010', '200003', 'system-config', '系统配置', 'MENU', '210001', '/sys/config', '/sys/config/index.vue', NULL, 'icon-park-outline:setting-config', NULL, NULL, 4, 1, 0, 0, 'ENABLED', '系统配置管理页面', NULL, '{}', '2026-07-18 14:07:48.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202011', '202010', 'sys:config:create', '新增系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48.000000', NULL, '2026-07-18 14:07:48.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202012', '202010', 'sys:config:detail', '查看系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48.000000', NULL, '2026-07-18 14:07:48.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202013', '202010', 'sys:config:update', '编辑系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48.000000', NULL, '2026-07-18 14:07:48.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202014', '202010', 'sys:config:delete', '删除系统配置', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-18 14:07:48.000000', NULL, '2026-07-18 14:07:48.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202015', '202001', 'sys-codegen', '代码生成', 'MENU', '210001', '/sys/codegen', '/sys/codegen/index.vue', NULL, 'icon-park-outline:code', NULL, NULL, 1, 1, 0, 0, 'ENABLED', '代码生成管理', NULL, '{}', '2026-07-18 16:10:45.000000', NULL, '2026-08-09 00:00:00.000000', '1');
INSERT INTO `sys_resource` VALUES ('202030', NULL, 'biz-demo', '业务示例', 'CATALOG', '210001', '/biz', NULL, NULL, 'icon-park-outline:application-one', NULL, NULL, 40, 1, 0, 0, 'ENABLED', '代码生成业务示例页面', NULL, '{}', '2026-08-09 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202200', '200019', 'content-notice', '通知消息', 'MENU', '210001', '/sys/notice', '/sys/notice/index.vue', NULL, 'icon-park-outline:message', NULL, NULL, 2, 1, 0, 0, 'ENABLED', '消息管理', NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202201', '202200', 'sys-notice-page', '分页消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202202', '202200', 'sys-notice-create', '新增消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202203', '202200', 'sys-notice-detail', '详情消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202204', '202200', 'sys-notice-update', '编辑消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202205', '202200', 'sys-notice-delete', '删除消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202206', '202200', 'sys-notice-create-page', '新增消息页', 'PAGE', '210001', '/sys/notice/create', '/sys/notice/form.vue', NULL, NULL, NULL, NULL, 60, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202207', '202200', 'sys-notice-edit-page', '编辑消息页', 'PAGE', '210001', '/sys/notice/edit', '/sys/notice/form.vue', NULL, NULL, NULL, NULL, 70, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202208', '202200', 'sys-notice-detail-page', '消息详情页', 'PAGE', '210001', '/sys/notice/detail', '/sys/notice/detail.vue', NULL, NULL, NULL, NULL, 80, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202209', '202200', 'sys-notice-publish', '发布消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 55, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202220', '200019', 'content-feedback', '反馈管理', 'MENU', '210001', '/sys/feedback', '/sys/feedback/index.vue', NULL, 'icon-park-outline:write', NULL, NULL, 3, 1, 0, 0, 'ENABLED', '意见反馈管理', NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-08-09 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202221', '202220', 'sys-feedback-page', '分页反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202222', '202220', 'sys-feedback-detail', '查看反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202223', '202220', 'sys-feedback-update', '处理反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202224', '202220', 'sys-feedback-delete', '删除反馈', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202225', '202220', 'sys-feedback-edit-page', '处理反馈页', 'PAGE', '210001', '/sys/feedback/edit', '/sys/feedback/form.vue', NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202226', '202220', 'sys-feedback-detail-page', '反馈详情页', 'PAGE', '210001', '/sys/feedback/detail', '/sys/feedback/detail.vue', NULL, NULL, NULL, NULL, 60, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-07-24 00:00:00.000000', NULL, '2026-07-24 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202230', '200003', 'sys-real-name', '实名认证审核', 'MENU', '210001', '/sys/real-name', '/sys/real-name/index.vue', NULL, 'icon-park-outline:id-card', NULL, NULL, 8, 1, 0, 0, 'ENABLED', '实名认证待审队列与审核', NULL, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202231', '202230', 'sys-real-name-review', '审核实名认证', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-22 00:00:00.000000', NULL, '2026-08-22 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202240', '202200', 'sys-notice-revoke', '撤回消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 56, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202241', '202200', 'sys-notice-pin', '置顶消息', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 57, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202251', '202010', 'sys-weakpassword-page', '分页弱密码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202252', '202010', 'sys-weakpassword-create', '新增弱密码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 11, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202253', '202010', 'sys-weakpassword-update', '编辑弱密码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 12, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202254', '202010', 'sys-weakpassword-delete', '删除弱密码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 13, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('202255', '202010', 'sys-weakpassword-detail', '查看弱密码', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 14, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-23 00:00:00.000000', NULL, '2026-08-23 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203011', '202004', 'biz-cgtestactivity-page', '分页活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203012', '202004', 'biz-cgtestactivity-create', '新增活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203013', '202004', 'biz-cgtestactivity-detail', '详情活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203014', '202004', 'biz-cgtestactivity-update', '编辑活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203015', '202004', 'biz-cgtestactivity-delete', '删除活动', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203021', '202005', 'biz-cgtestcatalog-page', '分页目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203022', '202005', 'biz-cgtestcatalog-create', '新增目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203023', '202005', 'biz-cgtestcatalog-detail', '详情目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203024', '202005', 'biz-cgtestcatalog-update', '编辑目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203025', '202005', 'biz-cgtestcatalog-delete', '删除目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203026', '202005', 'biz-cgtestcatalog-list', '树列表目录', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 90, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203031', '202006', 'biz-cgtestorder-page', '分页订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203032', '202006', 'biz-cgtestorder-create', '新增订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203033', '202006', 'biz-cgtestorder-detail', '详情订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203034', '202006', 'biz-cgtestorder-update', '编辑订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203035', '202006', 'biz-cgtestorder-delete', '删除订单', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203041', '202007', 'biz-cgtestknowledgecategory-page', '分页知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203042', '202007', 'biz-cgtestknowledgecategory-create', '新增知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 20, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203043', '202007', 'biz-cgtestknowledgecategory-detail', '详情知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 30, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203044', '202007', 'biz-cgtestknowledgecategory-update', '编辑知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 40, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203045', '202007', 'biz-cgtestknowledgecategory-delete', '删除知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 50, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('203046', '202007', 'biz-cgtestknowledgecategory-list', '树列表知识分类', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 90, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204001', '200003', 'sys-job', '任务管理', 'MENU', '210001', '/sys/job', '/sys/job/index.vue', NULL, 'icon-park-outline:timer', NULL, NULL, 4, 1, 0, 0, 'ENABLED', '任务调度管理（CRON / 固定间隔，Redis 锁防多实例重复执行）', NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204011', '204001', 'sys-job-create', '新增任务', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 1, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204012', '204001', 'sys-job-update', '编辑任务', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 2, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204013', '204001', 'sys-job-delete', '删除任务', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 3, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204014', '204001', 'sys-job-detail', '任务详情', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 4, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204015', '204001', 'sys-job-run', '立即执行', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 5, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204016', '204001', 'sys-job-log', '执行日志', 'BUTTON', '210001', NULL, NULL, NULL, NULL, NULL, NULL, 6, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204021', '204001', 'sys-job-create-page', '新增任务页', 'PAGE', '210001', '/sys/job/create', '/sys/job/form.vue', NULL, NULL, NULL, NULL, 7, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204022', '204001', 'sys-job-edit-page', '编辑任务页', 'PAGE', '210001', '/sys/job/edit', '/sys/job/form.vue', NULL, NULL, NULL, NULL, 8, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204023', '204001', 'sys-job-detail-page', '任务详情页', 'PAGE', '210001', '/sys/job/detail', '/sys/job/detail.vue', NULL, NULL, NULL, NULL, 9, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);
INSERT INTO `sys_resource` VALUES ('204024', '204001', 'sys-job-log-page', '任务执行记录页', 'PAGE', '210001', '/sys/job/log', '/sys/job/log.vue', NULL, NULL, NULL, NULL, 10, 0, 0, 0, 'ENABLED', NULL, NULL, '{}', '2026-08-16 00:00:00.000000', NULL, '2026-08-16 00:00:00.000000', NULL);

-- ----------------------------
-- Table structure for sys_resource_module
-- ----------------------------
DROP TABLE IF EXISTS `sys_resource_module`;
CREATE TABLE `sys_resource_module`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `client` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属客户端：admin/portal 等',
  `icon` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '图标标识',
  `color` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '颜色值',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '模块状态：ENABLED/DISABLED',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '资源模块描述',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_resource_module_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '资源模块' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_resource_module
-- ----------------------------
INSERT INTO `sys_resource_module` VALUES ('210001', '管理端', 'admin', 'ADMIN', 'icon-park-outline:all-application', '#2080f0', 1, 'ENABLED', '管理端菜单与权限资源模块', '{}', '2026-06-30 00:00:00.000000', NULL, '2026-06-30 00:00:00.000000', NULL);

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `code` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '编码',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `category` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色分类',
  `scope_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色作用域：GLOBAL/DEPT 等',
  `owner_dept_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '所属部门ID（数据权限范围）',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `status` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色状态：ENABLED/DISABLED',
  `is_builtin` tinyint(1) NOT NULL COMMENT '是否内置角色：1 内置 / 0 自定义',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL COMMENT '角色描述',
  `extra` json NOT NULL COMMENT '扩展信息（JSON）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_role_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '角色' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES ('1', 'SUPER_ADMIN', '超级管理员', 'SYS', 'PLATFORM', NULL, 1, 'ENABLED', 0, '系统内置超级管理员角色', '{}', '2026-08-08 11:56:13.747886', NULL, '2026-08-08 11:56:13.747886', NULL);
INSERT INTO `sys_role` VALUES ('2', 'IAM_ADMIN', 'IAM 管理员', 'SYS', 'PLATFORM', '8200000000000102', 10, 'ENABLED', 0, '组织权限管理（非超管）', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_role` VALUES ('3', 'BIZ_ALL', '业务-全部数据', 'SYS', 'PLATFORM', '8200000000000103', 20, 'ENABLED', 0, '活动模块 ALL', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_role` VALUES ('4', 'BIZ_DEPT', '业务-本部门', 'SYS', 'PLATFORM', '8200000000000102', 21, 'ENABLED', 0, '目录模块 DEPT', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_role` VALUES ('5', 'BIZ_SELF', '业务-仅本人', 'SYS', 'PLATFORM', '8200000000000104', 22, 'ENABLED', 0, '订单模块 SELF', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_role` VALUES ('6', 'BIZ_CHILD', '业务-部门及子部门', 'SYS', 'PLATFORM', '8200000000000102', 23, 'ENABLED', 0, '知识分类 DEPT_AND_CHILD', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');
INSERT INTO `sys_role` VALUES ('7', 'IAM_READONLY', 'IAM 只读', 'SYS', 'PLATFORM', '8200000000000107', 30, 'ENABLED', 0, '账号管理只读', '{}', '2026-08-23 08:00:00.000000', '1', '2026-08-23 08:00:00.000000', '1');

-- ----------------------------
-- Table structure for sys_weak_password
-- ----------------------------
DROP TABLE IF EXISTS `sys_weak_password`;
CREATE TABLE `sys_weak_password`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '弱口令明文（用于注册/改密校验）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `idx_sys_weak_password_password`(`password` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '弱口令库' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_weak_password
-- ----------------------------
INSERT INTO `sys_weak_password` VALUES ('7142597855705676121', 'qwerty', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_weak_password` VALUES ('7267772085910261393', '123456', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_weak_password` VALUES ('7404805181363764417', 'admin123', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_weak_password` VALUES ('7411256677926569870', '111111', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);
INSERT INTO `sys_weak_password` VALUES ('7661438436788304682', 'password', '2026-08-08 00:00:00.000000', NULL, '2026-08-08 00:00:00.000000', NULL);

-- ----------------------------
-- Table structure for sys_workspace_shortcut
-- ----------------------------
DROP TABLE IF EXISTS `sys_workspace_shortcut`;
CREATE TABLE `sys_workspace_shortcut`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主键ID',
  `account_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '所属账号ID',
  `resource_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '快捷菜单资源ID（sys_resource）',
  `sort` int NOT NULL COMMENT '排序号（越小越靠前）',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `created_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '创建人（账户ID）',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  `updated_by` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL COMMENT '更新人（账户ID）',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `uq_sys_workspace_shortcut_account_resource`(`account_id` ASC, `resource_id` ASC) USING BTREE,
  INDEX `ix_sys_workspace_shortcut_account_sort`(`account_id` ASC, `sort` ASC) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '工作台个人快捷应用' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of sys_workspace_shortcut
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
