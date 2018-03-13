/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 100129
 Source Host           : localhost:3306
 Source Schema         : myadmin

 Target Server Type    : MySQL
 Target Server Version : 100129
 File Encoding         : 65001

 Date: 13/03/2018 20:30:35
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for myadmin_menu
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_menu`;
CREATE TABLE `myadmin_menu`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '菜单id',
  `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '标识',
  `title` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '标题',
  `parent_id` int(11) NULL DEFAULT NULL COMMENT '父菜单id',
  `seq` int(11) NOT NULL DEFAULT 0 COMMENT '序号',
  `icon` varchar(32) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '图标',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 28 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_menu
-- ----------------------------
INSERT INTO `myadmin_menu` VALUES (1, 'set', '设置', NULL, 1000, 'set');
INSERT INTO `myadmin_menu` VALUES (11, 'user', '用户管理', 1, 100, 'user');
INSERT INTO `myadmin_menu` VALUES (12, 'resource', '资源管理', 1, 0, 'resource');
INSERT INTO `myadmin_menu` VALUES (13, 'role', '角色管理', 1, 9, 'role');
INSERT INTO `myadmin_menu` VALUES (14, 'menu', '菜单管理', 1, 0, 'menu');
INSERT INTO `myadmin_menu` VALUES (16, 'server', '服务器管理', NULL, 990, 'server');
INSERT INTO `myadmin_menu` VALUES (17, 'game_server', '游戏服管理', 16, 0, 'menu_item');
INSERT INTO `myadmin_menu` VALUES (18, 'server_node', '节点管理', 16, 0, 'menu_item');
INSERT INTO `myadmin_menu` VALUES (19, 'tool', '开发工具', NULL, 0, 'tool');
INSERT INTO `myadmin_menu` VALUES (20, 'debug', '调试工具', 19, 0, 'menu_item');
INSERT INTO `myadmin_menu` VALUES (21, 'operate_tool', '运营用具', NULL, 0, 'menu_level_1');
INSERT INTO `myadmin_menu` VALUES (22, 'statistics', '数据和统计', NULL, 0, 'statistics');
INSERT INTO `myadmin_menu` VALUES (23, 'notice', '公告', 21, 0, 'menu_item');
INSERT INTO `myadmin_menu` VALUES (24, 'log', '日志查询', NULL, 0, 'log');
INSERT INTO `myadmin_menu` VALUES (25, 'log_login', '登录日志', 24, 0, 'menu_item');
INSERT INTO `myadmin_menu` VALUES (26, 'statistics_online', '在线统计', 22, 0, 'menu_item');
INSERT INTO `myadmin_menu` VALUES (27, 'operate_chat', '聊天监控', 21, 0, 'menu_item');

-- ----------------------------
-- Table structure for myadmin_resource
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_resource`;
CREATE TABLE `myadmin_resource`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '资源id',
  `name` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '资源名称',
  `parent_id` int(11) NULL DEFAULT NULL COMMENT '父资源id',
  `url_for` varchar(256) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '地址',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 37 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_resource
-- ----------------------------
INSERT INTO `myadmin_resource` VALUES (3, '资源控制器', NULL, 'ResourceController.*');
INSERT INTO `myadmin_resource` VALUES (14, '用户控制器', NULL, 'UserController.*');
INSERT INTO `myadmin_resource` VALUES (15, '角色控制器', NULL, 'RoleController.*');
INSERT INTO `myadmin_resource` VALUES (17, '菜单控制器', NULL, 'MenuController.*');

-- ----------------------------
-- Table structure for myadmin_role
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role`;
CREATE TABLE `myadmin_role`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '角色名称',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 42 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_role
-- ----------------------------
INSERT INTO `myadmin_role` VALUES (22, '超级管理员');
INSERT INTO `myadmin_role` VALUES (41, 'test');

-- ----------------------------
-- Table structure for myadmin_role_menu_rel
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role_menu_rel`;
CREATE TABLE `myadmin_role_menu_rel`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `menu_id` int(11) NOT NULL,
  `created` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 60 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_role_menu_rel
-- ----------------------------
INSERT INTO `myadmin_role_menu_rel` VALUES (44, 22, 1, '2018-03-13 11:43:56');
INSERT INTO `myadmin_role_menu_rel` VALUES (45, 22, 11, '2018-03-13 11:43:56');
INSERT INTO `myadmin_role_menu_rel` VALUES (46, 22, 12, '2018-03-13 11:43:56');
INSERT INTO `myadmin_role_menu_rel` VALUES (47, 22, 13, '2018-03-13 11:43:56');
INSERT INTO `myadmin_role_menu_rel` VALUES (48, 22, 14, '2018-03-13 11:43:56');
INSERT INTO `myadmin_role_menu_rel` VALUES (54, 41, 1, '2018-03-13 11:51:37');
INSERT INTO `myadmin_role_menu_rel` VALUES (55, 41, 11, '2018-03-13 11:51:37');
INSERT INTO `myadmin_role_menu_rel` VALUES (56, 41, 12, '2018-03-13 11:51:37');
INSERT INTO `myadmin_role_menu_rel` VALUES (57, 41, 13, '2018-03-13 11:51:37');
INSERT INTO `myadmin_role_menu_rel` VALUES (58, 41, 14, '2018-03-13 11:51:37');
INSERT INTO `myadmin_role_menu_rel` VALUES (59, 41, 16, '2018-03-13 11:51:37');

-- ----------------------------
-- Table structure for myadmin_role_resource_rel
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role_resource_rel`;
CREATE TABLE `myadmin_role_resource_rel`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `resource_id` int(11) NOT NULL,
  `created` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 620 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_role_resource_rel
-- ----------------------------
INSERT INTO `myadmin_role_resource_rel` VALUES (612, 22, 3, '2018-03-13 11:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES (613, 22, 14, '2018-03-13 11:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES (614, 22, 15, '2018-03-13 11:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES (615, 22, 17, '2018-03-13 11:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES (616, 41, 3, '2018-03-13 11:45:34');
INSERT INTO `myadmin_role_resource_rel` VALUES (617, 41, 14, '2018-03-13 11:45:34');
INSERT INTO `myadmin_role_resource_rel` VALUES (618, 41, 15, '2018-03-13 11:45:34');
INSERT INTO `myadmin_role_resource_rel` VALUES (619, 41, 17, '2018-03-13 11:45:34');

-- ----------------------------
-- Table structure for myadmin_role_user_rel
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role_user_rel`;
CREATE TABLE `myadmin_role_user_rel`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `created` datetime(0) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 41 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_role_user_rel
-- ----------------------------
INSERT INTO `myadmin_role_user_rel` VALUES (36, 22, 1, '2018-03-13 11:04:41');
INSERT INTO `myadmin_role_user_rel` VALUES (40, 41, 14, '2018-03-13 11:44:16');

-- ----------------------------
-- Table structure for myadmin_user
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_user`;
CREATE TABLE `myadmin_user`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '帐号id',
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '名称',
  `account` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录帐号',
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '登录密码',
  `status` int(11) NOT NULL DEFAULT 0 COMMENT '状态',
  `mobile` varchar(16) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '',
  `login_times` int(11) NOT NULL DEFAULT 0 COMMENT '登录次数',
  `last_login_time` int(11) NOT NULL DEFAULT 0 COMMENT '最近登录时间',
  `last_login_ip` varchar(64) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '0' COMMENT '最近登录',
  `is_super` int(11) NULL DEFAULT NULL COMMENT '是否超级管理员',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 15 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of myadmin_user
-- ----------------------------
INSERT INTO `myadmin_user` VALUES (1, 'admin', 'admin', '21232f297a57a5a743894a0e4a801fc3', 1, '', 0, 1520944207, '127.0.0.1', 1);
INSERT INTO `myadmin_user` VALUES (14, 'test1', 'test', '098f6bcd4621d373cade4e832627b4f6', 1, '', 0, 1520941540, '127.0.0.1', 0);

-- ----------------------------
-- Table structure for myadmin_user_myadmin_roles
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_user_myadmin_roles`;
CREATE TABLE `myadmin_user_myadmin_roles`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `myadmin_user_id` int(11) NOT NULL,
  `myadmin_role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Compact;

SET FOREIGN_KEY_CHECKS = 1;
