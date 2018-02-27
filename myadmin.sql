/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50547
Source Host           : localhost:3306
Source Database       : myadmin

Target Server Type    : MYSQL
Target Server Version : 50547
File Encoding         : 65001

Date: 2018-01-18 08:23:59
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `myadmin_resource`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_resource`;
CREATE TABLE `myadmin_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `rtype` int(11) NOT NULL DEFAULT '0',
  `name` varchar(64) NOT NULL DEFAULT '',
  `parent_id` int(11) DEFAULT NULL,
  `seq` int(11) NOT NULL DEFAULT '0',
  `icon` varchar(32) NOT NULL DEFAULT '',
  `url_for` varchar(256) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=63 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_resource
-- ----------------------------
INSERT INTO `myadmin_resource` VALUES ('8', '0', '系统菜单', null, '10000', '', '');
INSERT INTO `myadmin_resource` VALUES ('9', '1', '资源管理', '14', '100', '', 'ResourceController.Index');
INSERT INTO `myadmin_resource` VALUES ('12', '1', '角色管理', '14', '100', '', 'RoleController.Index');
INSERT INTO `myadmin_resource` VALUES ('13', '1', '用户管理', '14', '100', '', 'UserController.Index');
INSERT INTO `myadmin_resource` VALUES ('14', '1', '系统管理', '8', '90', 'fa fa-gears', '');
INSERT INTO `myadmin_resource` VALUES ('23', '1', '系统日志', '14', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('25', '2', '编辑', '9', '100', 'fa fa-pencil', 'ResourceController.Edit');
INSERT INTO `myadmin_resource` VALUES ('26', '2', '编辑', '13', '100', 'fa fa-pencil', 'UserController.Edit');
INSERT INTO `myadmin_resource` VALUES ('27', '2', '删除', '9', '100', 'fa fa-trash', 'ResourceController.Delete');
INSERT INTO `myadmin_resource` VALUES ('29', '2', '删除', '13', '100', 'fa fa-trash', 'UserController.Delete');
INSERT INTO `myadmin_resource` VALUES ('30', '2', '编辑', '12', '100', 'fa fa-pencil', 'RoleController.Edit');
INSERT INTO `myadmin_resource` VALUES ('31', '2', '删除', '12', '100', 'fa fa-trash', 'RoleController.Delete');
INSERT INTO `myadmin_resource` VALUES ('32', '2', '分配资源', '12', '100', 'fa fa-th', 'RoleController.Allocate');
INSERT INTO `myadmin_resource` VALUES ('35', '1', ' 首页', null, '1', 'fa fa-dashboard', 'HomeController.Index');
INSERT INTO `myadmin_resource` VALUES ('41', '1', '数据和统计', null, '200', 'fa fa-book', '');
INSERT INTO `myadmin_resource` VALUES ('42', '1', '运营工具', null, '300', 'fa fa-book', '');
INSERT INTO `myadmin_resource` VALUES ('43', '1', '在线数据', '41', '100', 'fa-plus', '');
INSERT INTO `myadmin_resource` VALUES ('44', '1', '日志查询', '42', '100', 'fa fa-book', '');
INSERT INTO `myadmin_resource` VALUES ('45', '1', '充值数据', '41', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('46', '1', '充值列表', '45', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('47', '1', '数据分析', '41', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('48', '1', '留存数据', '47', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('49', '1', '服务器管理', null, '5000', 'fa fa-book', '');
INSERT INTO `myadmin_resource` VALUES ('50', '1', '游戏服管理', '49', '100', '', 'GameServerController.List');
INSERT INTO `myadmin_resource` VALUES ('51', '1', '节点管理', '49', '100', '', 'ServerNodeController.List');
INSERT INTO `myadmin_resource` VALUES ('52', '1', '开发工具', null, '100', 'fa fa-book', '');
INSERT INTO `myadmin_resource` VALUES ('53', '1', '调试工具', '52', '100', '', 'ToolController.Build');
INSERT INTO `myadmin_resource` VALUES ('54', '1', '活跃分析', '47', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('55', '1', '角色列表', '41', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('56', '1', '当前在线', '43', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('57', '1', '在线统计', '43', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('58', '1', '公告', '42', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('59', '1', '玩家查询', '41', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('60', '1', '聊天监控', '42', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('61', '1', '登录日志', '44', '100', '', '');
INSERT INTO `myadmin_resource` VALUES ('62', '0', '业务菜单', null, '99', '', '');

-- ----------------------------
-- Table structure for `myadmin_role`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role`;
CREATE TABLE `myadmin_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL DEFAULT '',
  `seq` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role
-- ----------------------------
INSERT INTO `myadmin_role` VALUES ('22', '超级管理员', '10');
INSERT INTO `myadmin_role` VALUES ('24', '研发', '20');
INSERT INTO `myadmin_role` VALUES ('26', '运营', '30');
INSERT INTO `myadmin_role` VALUES ('27', '客服', '40');

-- ----------------------------
-- Table structure for `myadmin_role_resource_rel`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role_resource_rel`;
CREATE TABLE `myadmin_role_resource_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `resource_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=472 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_resource_rel
-- ----------------------------
INSERT INTO `myadmin_role_resource_rel` VALUES ('448', '24', '8', '2017-12-19 06:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES ('449', '24', '14', '2017-12-19 06:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES ('450', '24', '23', '2017-12-19 06:40:16');
INSERT INTO `myadmin_role_resource_rel` VALUES ('451', '22', '35', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('458', '22', '8', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('459', '22', '14', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('460', '22', '23', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('462', '22', '9', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('463', '22', '25', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('464', '22', '27', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('465', '22', '12', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('466', '22', '30', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('467', '22', '31', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('468', '22', '32', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('469', '22', '13', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('470', '22', '26', '2018-01-16 06:46:10');
INSERT INTO `myadmin_role_resource_rel` VALUES ('471', '22', '29', '2018-01-16 06:46:10');

-- ----------------------------
-- Table structure for `myadmin_role_user_rel`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role_user_rel`;
CREATE TABLE `myadmin_role_user_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=70 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_user_rel
-- ----------------------------
INSERT INTO `myadmin_role_user_rel` VALUES ('69', '22', '1', '2018-01-16 06:43:57');

-- ----------------------------
-- Table structure for `myadmin_user`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_user`;
CREATE TABLE `myadmin_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '帐号id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  `account` varchar(255) NOT NULL DEFAULT '' COMMENT '登录帐号',
  `password` varchar(255) NOT NULL DEFAULT ''  COMMENT '登录密码',
--   `is_super` tinyint(1) NOT NULL DEFAULT '0',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态',
   `mobile` varchar(16) NOT NULL DEFAULT '',
  `login_times` int(11) NOT NULL DEFAULT '0' COMMENT '登录次数',
  `last_login_time` int(11) NOT NULL DEFAULT '0' COMMENT '最近登录时间',
  `last_login_ip` varchar(64) NOT NULL DEFAULT '0' COMMENT '最近登录',

--   `email` varchar(256) NOT NULL DEFAULT '',
--   `avatar` varchar(256) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_user
-- ----------------------------
INSERT INTO `myadmin_user` VALUES ('1', 'admin', 'admin', 'e10adc3949ba59abbe56e057f20f883e', '1', '', '0', '0', '127.0.0.1');
INSERT INTO `myadmin_user` VALUES ('3', '张三', 'zhangsan', 'e10adc3949ba59abbe56e057f20f883e', '0', '', '0', '0', '127.0.0.1');
INSERT INTO `myadmin_user` VALUES ('5', '李四', 'lisi', 'e10adc3949ba59abbe56e057f20f883e', '0', '', '0', '0', '127.0.0.1');
INSERT INTO `myadmin_user` VALUES ('7', 'test', 'test', 'd41d8cd98f00b204e9800998ecf8427e', '1', '', '0', '0', '127.0.0.1');

-- ----------------------------
-- Table structure for `myadmin_user_myadmin_roles`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_user_myadmin_roles`;
CREATE TABLE `myadmin_user_myadmin_roles` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `myadmin_user_id` int(11) NOT NULL,
  `myadmin_role_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_user_myadmin_roles
-- ----------------------------
