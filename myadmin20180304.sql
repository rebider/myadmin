/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50547
Source Host           : localhost:3306
Source Database       : myadmin

Target Server Type    : MYSQL
Target Server Version : 50547
File Encoding         : 65001

Date: 2018-03-04 22:27:29
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `myadmin_resource`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_resource`;
CREATE TABLE `myadmin_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '资源id',
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '资源类型[0:菜单 1:接口]',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '资源名称',
  `parent_id` int(11) DEFAULT NULL COMMENT '父资源id',
  `seq` int(11) NOT NULL DEFAULT '0',
  `icon` varchar(32) NOT NULL DEFAULT '' COMMENT '图标',
  `url_for` varchar(256) NOT NULL DEFAULT '' COMMENT '地址',
  `title` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_resource
-- ----------------------------
INSERT INTO `myadmin_resource` VALUES ('1', '0', 'set', null, '100', 'example', '/set', '设置');
INSERT INTO `myadmin_resource` VALUES ('11', '1', 'user', '1', '100', 'example', '', '用户管理');
INSERT INTO `myadmin_resource` VALUES ('12', '0', 'resource', '1', '0', 'example', '', '资源管理');
INSERT INTO `myadmin_resource` VALUES ('13', '0', 'role', '1', '0', 'example', '', '角色管理');

-- ----------------------------
-- Table structure for `myadmin_role`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role`;
CREATE TABLE `myadmin_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '角色名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role
-- ----------------------------
INSERT INTO `myadmin_role` VALUES ('22', '超级管理员');
INSERT INTO `myadmin_role` VALUES ('24', '研发');
INSERT INTO `myadmin_role` VALUES ('26', '运营');
INSERT INTO `myadmin_role` VALUES ('27', '客服');
INSERT INTO `myadmin_role` VALUES ('28', '151651');
INSERT INTO `myadmin_role` VALUES ('29', 'wwwwweqwe');

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
) ENGINE=InnoDB AUTO_INCREMENT=508 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_resource_rel
-- ----------------------------
INSERT INTO `myadmin_role_resource_rel` VALUES ('498', '24', '11', '2018-03-04 04:17:33');
INSERT INTO `myadmin_role_resource_rel` VALUES ('499', '24', '13', '2018-03-04 04:17:33');
INSERT INTO `myadmin_role_resource_rel` VALUES ('500', '28', '13', '2018-03-04 04:18:52');
INSERT INTO `myadmin_role_resource_rel` VALUES ('503', '22', '11', '2018-03-04 08:15:22');
INSERT INTO `myadmin_role_resource_rel` VALUES ('504', '22', '12', '2018-03-04 08:15:22');
INSERT INTO `myadmin_role_resource_rel` VALUES ('505', '22', '13', '2018-03-04 08:15:22');

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
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_user_rel
-- ----------------------------
INSERT INTO `myadmin_role_user_rel` VALUES ('2', '24', '3', '2018-02-27 00:29:49');
INSERT INTO `myadmin_role_user_rel` VALUES ('3', '27', '3', '2018-02-27 00:29:49');
INSERT INTO `myadmin_role_user_rel` VALUES ('4', '22', '1', '2018-02-27 23:36:34');
INSERT INTO `myadmin_role_user_rel` VALUES ('5', '24', '1', '2018-02-27 23:36:34');

-- ----------------------------
-- Table structure for `myadmin_user`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_user`;
CREATE TABLE `myadmin_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '帐号id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '名称',
  `account` varchar(255) NOT NULL DEFAULT '' COMMENT '登录帐号',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '登录密码',
  `status` int(11) NOT NULL DEFAULT '0' COMMENT '状态',
  `mobile` varchar(16) NOT NULL DEFAULT '',
  `login_times` int(11) NOT NULL DEFAULT '0' COMMENT '登录次数',
  `last_login_time` int(11) NOT NULL DEFAULT '0' COMMENT '最近登录时间',
  `last_login_ip` varchar(64) NOT NULL DEFAULT '0' COMMENT '最近登录',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_user
-- ----------------------------
INSERT INTO `myadmin_user` VALUES ('1', 'admin', 'admin', '21232f297a57a5a743894a0e4a801fc3', '1', '', '0', '1520173421', '127.0.0.1');
INSERT INTO `myadmin_user` VALUES ('3', '张三2222', 'zhangsan', 'e10adc3949ba59abbe56e057f20f883e', '0', '', '0', '0', '');
INSERT INTO `myadmin_user` VALUES ('5', '李四', 'lisi', 'e10adc3949ba59abbe56e057f20f883e', '0', '', '0', '0', '127.0.0.1');

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
