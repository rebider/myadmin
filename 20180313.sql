/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 50547
Source Host           : localhost:3306
Source Database       : myadmin

Target Server Type    : MYSQL
Target Server Version : 50547
File Encoding         : 65001

Date: 2018-03-13 08:26:01
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `c_game_server`
-- ----------------------------
DROP TABLE IF EXISTS `c_game_server`;
CREATE TABLE `c_game_server` (
  `platform_id` int(11) NOT NULL DEFAULT '0',
  `sid` varchar(255) NOT NULL,
  `desc` varchar(255) NOT NULL DEFAULT '',
  `node` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`sid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of c_game_server
-- ----------------------------

-- ----------------------------
-- Table structure for `c_server_node`
-- ----------------------------
DROP TABLE IF EXISTS `c_server_node`;
CREATE TABLE `c_server_node` (
  `node` varchar(255) NOT NULL,
  `ip` varchar(255) NOT NULL DEFAULT '',
  `port` int(11) NOT NULL DEFAULT '0',
  `type` int(11) NOT NULL DEFAULT '0',
  `zone_node` varchar(255) NOT NULL DEFAULT '',
  `server_version` varchar(255) NOT NULL DEFAULT '',
  `client_version` varchar(255) NOT NULL DEFAULT '',
  `open_time` int(11) NOT NULL DEFAULT '0',
  `is_test` int(11) NOT NULL DEFAULT '0',
  `platform_id` int(11) NOT NULL DEFAULT '0',
  `state` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`node`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of c_server_node
-- ----------------------------

-- ----------------------------
-- Table structure for `myadmin_menu`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_menu`;
CREATE TABLE `myadmin_menu` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '菜单id',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '标识',
  `title` varchar(64) NOT NULL DEFAULT '' COMMENT '标题',
  `parent_id` int(11) DEFAULT NULL COMMENT '父菜单id',
  `seq` int(11) NOT NULL DEFAULT '0' COMMENT '序号',
  `icon` varchar(32) NOT NULL DEFAULT '' COMMENT '图标',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_menu
-- ----------------------------
INSERT INTO `myadmin_menu` VALUES ('1', 'set', '设置1', null, '0', 'example');
INSERT INTO `myadmin_menu` VALUES ('11', 'user', '用户管理', '1', '0', 'example');
INSERT INTO `myadmin_menu` VALUES ('12', 'resource', '资源管理', '1', '0', 'example');
INSERT INTO `myadmin_menu` VALUES ('13', 'role', '角色管理', '1', '0', 'example');
INSERT INTO `myadmin_menu` VALUES ('14', 'menu', '菜单管理', '1', '0', 'example');
INSERT INTO `myadmin_menu` VALUES ('15', 'test', 'test', '1', '0', '');

-- ----------------------------
-- Table structure for `myadmin_resource`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_resource`;
CREATE TABLE `myadmin_resource` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '资源id',
  `name` varchar(64) NOT NULL DEFAULT '' COMMENT '资源名称',
  `parent_id` int(11) DEFAULT NULL COMMENT '父资源id',
  `url_for` varchar(256) NOT NULL DEFAULT '' COMMENT '地址',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_resource
-- ----------------------------
INSERT INTO `myadmin_resource` VALUES ('14', '用户控制器', '0', 'UserController.*');
INSERT INTO `myadmin_resource` VALUES ('15', '角色控制器', '0', 'RoleController.*');

-- ----------------------------
-- Table structure for `myadmin_role`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role`;
CREATE TABLE `myadmin_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '角色id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '角色名称',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role
-- ----------------------------
INSERT INTO `myadmin_role` VALUES ('22', '超级管理员');
INSERT INTO `myadmin_role` VALUES ('24', '研发');
INSERT INTO `myadmin_role` VALUES ('26', '运营');
INSERT INTO `myadmin_role` VALUES ('27', '客服');
INSERT INTO `myadmin_role` VALUES ('28', '151651');

-- ----------------------------
-- Table structure for `myadmin_role_menu_rel`
-- ----------------------------
DROP TABLE IF EXISTS `myadmin_role_menu_rel`;
CREATE TABLE `myadmin_role_menu_rel` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `role_id` int(11) NOT NULL,
  `menu_id` int(11) NOT NULL,
  `created` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_menu_rel
-- ----------------------------

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
) ENGINE=InnoDB AUTO_INCREMENT=521 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_resource_rel
-- ----------------------------
INSERT INTO `myadmin_role_resource_rel` VALUES ('498', '24', '11', '2018-03-04 04:17:33');
INSERT INTO `myadmin_role_resource_rel` VALUES ('499', '24', '13', '2018-03-04 04:17:33');
INSERT INTO `myadmin_role_resource_rel` VALUES ('500', '28', '13', '2018-03-04 04:18:52');
INSERT INTO `myadmin_role_resource_rel` VALUES ('515', '22', '1', '2018-03-07 00:13:55');
INSERT INTO `myadmin_role_resource_rel` VALUES ('516', '22', '11', '2018-03-07 00:13:55');
INSERT INTO `myadmin_role_resource_rel` VALUES ('517', '22', '12', '2018-03-07 00:13:55');
INSERT INTO `myadmin_role_resource_rel` VALUES ('518', '22', '13', '2018-03-07 00:13:55');
INSERT INTO `myadmin_role_resource_rel` VALUES ('519', '22', '14', '2018-03-07 00:13:55');
INSERT INTO `myadmin_role_resource_rel` VALUES ('520', '22', '15', '2018-03-07 00:13:55');

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
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_role_user_rel
-- ----------------------------
INSERT INTO `myadmin_role_user_rel` VALUES ('10', '24', '3', '2018-03-05 14:42:13');
INSERT INTO `myadmin_role_user_rel` VALUES ('11', '27', '3', '2018-03-05 14:42:13');
INSERT INTO `myadmin_role_user_rel` VALUES ('14', '24', '5', '2018-03-05 14:44:31');
INSERT INTO `myadmin_role_user_rel` VALUES ('24', '22', '1', '2018-03-11 14:04:47');

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
  `is_super` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of myadmin_user
-- ----------------------------
INSERT INTO `myadmin_user` VALUES ('1', '管理员', 'admin', '21232f297a57a5a743894a0e4a801fc3', '1', '', '0', '1520900608', '127.0.0.1', '1');
INSERT INTO `myadmin_user` VALUES ('3', '张三2222', 'zhangsan', 'e10adc3949ba59abbe56e057f20f883e', '0', '', '0', '0', '', null);
INSERT INTO `myadmin_user` VALUES ('5', '李四', 'lisi', 'e10adc3949ba59abbe56e057f20f883e', '0', '', '0', '0', '', null);

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
