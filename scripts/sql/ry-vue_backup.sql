-- MySQL dump 10.13  Distrib 9.5.0, for macos15.7 (arm64)
--
-- Host: localhost    Database: ry-vue
-- ------------------------------------------------------
-- Server version	9.5.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ 'abcbc794-aef0-11ef-a2c6-56f7d33db7dd:1-1004';

--
-- Table structure for table `gen_table`
--

DROP TABLE IF EXISTS `gen_table`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gen_table` (
  `table_id` bigint NOT NULL COMMENT '编号',
  `data_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '数据源名称',
  `table_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '表名称',
  `table_comment` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '表描述',
  `sub_table_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联子表的表名',
  `sub_table_fk_name` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '子表关联的外键名',
  `class_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '实体类名称',
  `tpl_category` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT 'crud' COMMENT '使用的模板（crud单表操作 tree树表操作）',
  `package_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '生成包路径',
  `module_name` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '生成模块名',
  `business_name` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '生成业务名',
  `function_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '生成功能名',
  `function_author` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '生成功能作者',
  `gen_type` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '生成代码方式（0zip压缩包 1自定义路径）',
  `gen_path` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '/' COMMENT '生成路径（不填默认项目路径）',
  `options` varchar(1000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '其它生成选项',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`table_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='代码生成业务表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gen_table`
--

LOCK TABLES `gen_table` WRITE;
/*!40000 ALTER TABLE `gen_table` DISABLE KEYS */;
/*!40000 ALTER TABLE `gen_table` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `gen_table_column`
--

DROP TABLE IF EXISTS `gen_table_column`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gen_table_column` (
  `column_id` bigint NOT NULL COMMENT '编号',
  `table_id` bigint DEFAULT NULL COMMENT '归属表编号',
  `column_name` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '列名称',
  `column_comment` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '列描述',
  `column_type` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '列类型',
  `java_type` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'JAVA类型',
  `java_field` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'JAVA字段名',
  `is_pk` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否主键（1是）',
  `is_increment` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否自增（1是）',
  `is_required` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否必填（1是）',
  `is_insert` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否为插入字段（1是）',
  `is_edit` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否编辑字段（1是）',
  `is_list` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否列表字段（1是）',
  `is_query` char(1) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '是否查询字段（1是）',
  `query_type` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT 'EQ' COMMENT '查询方式（等于、不等于、大于、小于、范围）',
  `html_type` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '显示类型（文本框、文本域、下拉框、复选框、单选框、日期控件）',
  `dict_type` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '字典类型',
  `sort` int DEFAULT NULL COMMENT '排序',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`column_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='代码生成业务表字段';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `gen_table_column`
--

LOCK TABLES `gen_table_column` WRITE;
/*!40000 ALTER TABLE `gen_table_column` DISABLE KEYS */;
/*!40000 ALTER TABLE `gen_table_column` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_client`
--

DROP TABLE IF EXISTS `sys_client`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_client` (
  `id` bigint NOT NULL COMMENT 'id',
  `client_id` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户端id',
  `client_key` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户端key',
  `client_secret` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '客户端秘钥',
  `grant_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '授权类型',
  `device_type` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '设备类型',
  `active_timeout` int DEFAULT '1800' COMMENT 'token活跃超时时间',
  `timeout` int DEFAULT '604800' COMMENT 'token固定超时',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统授权表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_client`
--

LOCK TABLES `sys_client` WRITE;
/*!40000 ALTER TABLE `sys_client` DISABLE KEYS */;
INSERT INTO `sys_client` VALUES (1,'e5cd7e4891bf95d1d19206ce24a7b32e','pc','pc123','password,social','pc',1800,604800,'0','0',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04'),(2,'428a8310cd442757ae699df5d894f051','app','app123','password,sms,social','android',1800,604800,'0','0',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04');
/*!40000 ALTER TABLE `sys_client` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `config_id` bigint NOT NULL COMMENT '参数主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `config_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '参数名称',
  `config_key` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '参数键名',
  `config_value` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '参数键值',
  `config_type` char(1) COLLATE utf8mb4_unicode_ci DEFAULT 'N' COMMENT '系统内置（Y是 N否）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`config_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='参数配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
INSERT INTO `sys_config` VALUES (1,'000000','主框架页-默认皮肤样式名称','sys.index.skinName','skin-blue','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'蓝色 skin-blue、绿色 skin-green、紫色 skin-purple、红色 skin-red、黄色 skin-yellow'),(2,'000000','用户管理-账号初始密码','sys.user.initPassword','123456','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'初始化密码 123456'),(3,'000000','主框架页-侧边栏主题','sys.index.sideTheme','theme-dark','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'深色主题theme-dark，浅色主题theme-light'),(5,'000000','账号自助-是否开启用户注册功能','sys.account.registerUser','false','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'是否开启注册用户功能（true开启，false关闭）'),(11,'000000','OSS预览列表资源开关','sys.oss.previewListResource','true','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'true:开启, false:关闭');
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dept`
--

DROP TABLE IF EXISTS `sys_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dept` (
  `dept_id` bigint NOT NULL COMMENT '部门id',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `parent_id` bigint DEFAULT '0' COMMENT '父部门id',
  `ancestors` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '祖级列表',
  `dept_name` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '部门名称',
  `dept_category` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '部门类别编码',
  `order_num` int DEFAULT '0' COMMENT '显示顺序',
  `leader` bigint DEFAULT NULL COMMENT '负责人',
  `phone` varchar(11) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系电话',
  `email` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '邮箱',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '部门状态（0正常 1停用）',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='部门表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dept`
--

LOCK TABLES `sys_dept` WRITE;
/*!40000 ALTER TABLE `sys_dept` DISABLE KEYS */;
INSERT INTO `sys_dept` VALUES (100,'000000',0,'0','XXX科技',NULL,0,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(101,'000000',100,'0,100','深圳总公司',NULL,1,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(102,'000000',100,'0,100','长沙分公司',NULL,2,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(103,'000000',101,'0,100,101','研发部门',NULL,1,1,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(104,'000000',101,'0,100,101','市场部门',NULL,2,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(105,'000000',101,'0,100,101','测试部门',NULL,3,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(106,'000000',101,'0,100,101','财务部门',NULL,4,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(107,'000000',101,'0,100,101','运维部门',NULL,5,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(108,'000000',102,'0,100,102','市场部门',NULL,1,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(109,'000000',102,'0,100,102','财务部门',NULL,2,NULL,'15888888888','xxx@qq.com','0','0',103,1,'2026-01-02 21:32:02',NULL,NULL),(2009971775458099201,'000000',101,'0,100,101','aa','aa',0,NULL,NULL,NULL,'0','1',103,1,'2026-01-10 20:53:09',1,'2026-01-10 20:53:14');
/*!40000 ALTER TABLE `sys_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict_data`
--

DROP TABLE IF EXISTS `sys_dict_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_data` (
  `dict_code` bigint NOT NULL COMMENT '字典编码',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `dict_sort` int DEFAULT '0' COMMENT '字典排序',
  `dict_label` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '字典标签',
  `dict_value` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '字典键值',
  `dict_type` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '字典类型',
  `css_class` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '样式属性（其他样式扩展）',
  `list_class` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '表格回显样式',
  `is_default` char(1) COLLATE utf8mb4_unicode_ci DEFAULT 'N' COMMENT '是否默认（Y是 N否）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`dict_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典数据表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_data`
--

LOCK TABLES `sys_dict_data` WRITE;
/*!40000 ALTER TABLE `sys_dict_data` DISABLE KEYS */;
INSERT INTO `sys_dict_data` VALUES (1,'000000',1,'男','0','sys_user_sex','','','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'性别男'),(2,'000000',2,'女','1','sys_user_sex','','','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'性别女'),(3,'000000',3,'未知','2','sys_user_sex','','','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'性别未知'),(4,'000000',1,'显示','0','sys_show_hide','','primary','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'显示菜单'),(5,'000000',2,'隐藏','1','sys_show_hide','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'隐藏菜单'),(6,'000000',1,'正常','0','sys_normal_disable','','primary','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'正常状态'),(7,'000000',2,'停用','1','sys_normal_disable','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'停用状态'),(12,'000000',1,'是','Y','sys_yes_no','','primary','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'系统默认是'),(13,'000000',2,'否','N','sys_yes_no','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'系统默认否'),(14,'000000',1,'通知','1','sys_notice_type','','warning','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'通知'),(15,'000000',2,'公告','2','sys_notice_type','','success','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'公告'),(16,'000000',1,'正常','0','sys_notice_status','','primary','Y',103,1,'2026-01-02 21:32:04',NULL,NULL,'正常状态'),(17,'000000',2,'关闭','1','sys_notice_status','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'关闭状态'),(18,'000000',1,'新增','1','sys_oper_type','','info','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'新增操作'),(19,'000000',2,'修改','2','sys_oper_type','','info','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'修改操作'),(20,'000000',3,'删除','3','sys_oper_type','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'删除操作'),(21,'000000',4,'授权','4','sys_oper_type','','primary','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'授权操作'),(22,'000000',5,'导出','5','sys_oper_type','','warning','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'导出操作'),(23,'000000',6,'导入','6','sys_oper_type','','warning','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'导入操作'),(24,'000000',7,'强退','7','sys_oper_type','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'强退操作'),(25,'000000',8,'生成代码','8','sys_oper_type','','warning','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'生成操作'),(26,'000000',9,'清空数据','9','sys_oper_type','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'清空操作'),(27,'000000',1,'成功','0','sys_common_status','','primary','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'正常状态'),(28,'000000',2,'失败','1','sys_common_status','','danger','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'停用状态'),(29,'000000',99,'其他','0','sys_oper_type','','info','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'其他操作'),(30,'000000',0,'密码认证','password','sys_grant_type','el-check-tag','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'密码认证'),(31,'000000',0,'短信认证','sms','sys_grant_type','el-check-tag','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'短信认证'),(32,'000000',0,'邮件认证','email','sys_grant_type','el-check-tag','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'邮件认证'),(33,'000000',0,'小程序认证','xcx','sys_grant_type','el-check-tag','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'小程序认证'),(34,'000000',0,'三方登录认证','social','sys_grant_type','el-check-tag','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'三方登录认证'),(35,'000000',0,'PC','pc','sys_device_type','','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'PC'),(36,'000000',0,'安卓','android','sys_device_type','','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'安卓'),(37,'000000',0,'iOS','ios','sys_device_type','','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'iOS'),(38,'000000',0,'小程序','xcx','sys_device_type','','default','N',103,1,'2026-01-02 21:32:04',NULL,NULL,'小程序');
/*!40000 ALTER TABLE `sys_dict_data` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dict_type`
--

DROP TABLE IF EXISTS `sys_dict_type`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dict_type` (
  `dict_id` bigint NOT NULL COMMENT '字典主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `dict_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '字典名称',
  `dict_type` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '字典类型',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`dict_id`),
  UNIQUE KEY `tenant_id` (`tenant_id`,`dict_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='字典类型表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dict_type`
--

LOCK TABLES `sys_dict_type` WRITE;
/*!40000 ALTER TABLE `sys_dict_type` DISABLE KEYS */;
INSERT INTO `sys_dict_type` VALUES (1,'000000','用户性别','sys_user_sex',103,1,'2026-01-02 21:32:04',NULL,NULL,'用户性别列表'),(2,'000000','菜单状态','sys_show_hide',103,1,'2026-01-02 21:32:04',NULL,NULL,'菜单状态列表'),(3,'000000','系统开关','sys_normal_disable',103,1,'2026-01-02 21:32:04',NULL,NULL,'系统开关列表'),(6,'000000','系统是否','sys_yes_no',103,1,'2026-01-02 21:32:04',NULL,NULL,'系统是否列表'),(7,'000000','通知类型','sys_notice_type',103,1,'2026-01-02 21:32:04',NULL,NULL,'通知类型列表'),(8,'000000','通知状态','sys_notice_status',103,1,'2026-01-02 21:32:04',NULL,NULL,'通知状态列表'),(9,'000000','操作类型','sys_oper_type',103,1,'2026-01-02 21:32:04',NULL,NULL,'操作类型列表'),(10,'000000','系统状态','sys_common_status',103,1,'2026-01-02 21:32:04',NULL,NULL,'登录状态列表'),(11,'000000','授权类型','sys_grant_type',103,1,'2026-01-02 21:32:04',NULL,NULL,'认证授权类型'),(12,'000000','设备类型','sys_device_type',103,1,'2026-01-02 21:32:04',NULL,NULL,'客户端设备类型'),(268231441890934784,'000000','aaaa','aaaa',103,1,NULL,NULL,NULL,'aaaaa');
/*!40000 ALTER TABLE `sys_dict_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_logininfor`
--

DROP TABLE IF EXISTS `sys_logininfor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_logininfor` (
  `info_id` bigint NOT NULL COMMENT '访问ID',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `user_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户账号',
  `client_key` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '客户端',
  `device_type` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '设备类型',
  `ipaddr` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '登录IP地址',
  `login_location` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '登录地点',
  `browser` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '浏览器类型',
  `os` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作系统',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '登录状态（0成功 1失败）',
  `msg` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '提示消息',
  `login_time` datetime DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`info_id`),
  KEY `idx_sys_logininfor_s` (`status`),
  KEY `idx_sys_logininfor_lt` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统访问记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_logininfor`
--

LOCK TABLES `sys_logininfor` WRITE;
/*!40000 ALTER TABLE `sys_logininfor` DISABLE KEYS */;
INSERT INTO `sys_logininfor` VALUES (268593931539709952,'000000','admin','pc','pc','192.168.3.216','','Chrome','macOS','0','登录成功','2026-01-11 12:16:24'),(268604438325755904,'000000','admin','pc','pc','192.168.3.216','','Chrome','macOS','0','登录成功','2026-01-11 12:58:09'),(268693116523905024,'000000','admin','pc','pc','192.168.3.216','','Chrome','macOS','0','登录成功','2026-01-11 18:50:31'),(268735266175320064,'000000','admin','pc','pc','192.168.3.216','','Chrome','macOS','0','登录成功','2026-01-11 21:38:00'),(268742762763911168,'000000','admin','pc','pc','192.168.3.216','','Chrome','macOS','0','登录成功','2026-01-11 22:07:48'),(2009478729205395458,'000000','admin','pc','pc','0:0:0:0:0:0:0:1','内网IP','Chrome','OSX','0','登录成功','2026-01-09 12:13:58'),(2009480166878588930,'000000','admin','pc','pc','0:0:0:0:0:0:0:1','内网IP','Chrome','OSX','0','登录成功','2026-01-09 12:19:41'),(2009482126079930369,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','退出成功','2026-01-09 12:27:28'),(2009482150327201794,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-09 12:27:34'),(2009504487181234177,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-09 13:56:19'),(2009513116571545602,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','退出成功','2026-01-09 14:30:37'),(2009513264844386305,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-09 14:31:12'),(2009585304528089090,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-09 19:17:28'),(2009615587969155073,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-09 21:17:48'),(2009823910102384641,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-10 11:05:36'),(2009863201402732545,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-10 13:41:43'),(2009964074493657090,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-10 20:22:33'),(2009964498617483266,'000000','admin','pc','pc','192.168.3.216','内网IP','Unknown','Unknown','0','Login successful','2026-01-10 20:24:14'),(2009973057879130113,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','退出成功','2026-01-10 20:58:15'),(2010205943156957186,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-11 12:23:39'),(2010299865212379137,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-11 18:36:52'),(2010300130334334978,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','Login successful','2026-01-11 18:37:55'),(2010344960569630722,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-11 21:36:04'),(2010350056707710977,'000000','admin','pc','pc','192.168.3.216','内网IP','Chrome','OSX','0','登录成功','2026-01-11 21:56:19');
/*!40000 ALTER TABLE `sys_logininfor` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menu`
--

DROP TABLE IF EXISTS `sys_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menu` (
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  `menu_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '菜单名称',
  `parent_id` bigint DEFAULT '0' COMMENT '父菜单ID',
  `order_num` int DEFAULT '0' COMMENT '显示顺序',
  `path` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '路由地址',
  `component` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '组件路径',
  `query_param` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '路由参数',
  `is_frame` int DEFAULT '1' COMMENT '是否为外链（0是 1否）',
  `is_cache` int DEFAULT '0' COMMENT '是否缓存（0缓存 1不缓存）',
  `menu_type` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '菜单类型（M目录 C菜单 F按钮）',
  `visible` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '显示状态（0显示 1隐藏）',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '菜单状态（0正常 1停用）',
  `perms` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '权限标识',
  `icon` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '#' COMMENT '菜单图标',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
  PRIMARY KEY (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单权限表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menu`
--

LOCK TABLES `sys_menu` WRITE;
/*!40000 ALTER TABLE `sys_menu` DISABLE KEYS */;
INSERT INTO `sys_menu` VALUES (1,'系统管理',0,1,'system',NULL,'',1,0,'M','0','0','','system',103,1,'2026-01-02 21:32:02',NULL,NULL,'系统管理目录'),(2,'系统监控',0,3,'monitor',NULL,'',1,0,'M','0','0','','monitor',103,1,'2026-01-02 21:32:02',NULL,NULL,'系统监控目录'),(100,'用户管理',1,1,'user','system/user/index','',1,0,'C','0','0','system:user:list','user',103,1,'2026-01-02 21:32:02',NULL,NULL,'用户管理菜单'),(101,'角色管理',1,2,'role','system/role/index','',1,0,'C','0','0','system:role:list','peoples',103,1,'2026-01-02 21:32:02',NULL,NULL,'角色管理菜单'),(102,'菜单管理',1,3,'menu','system/menu/index','',1,0,'C','0','0','system:menu:list','tree-table',103,1,'2026-01-02 21:32:02',NULL,NULL,'菜单管理菜单'),(103,'部门管理',1,4,'dept','system/dept/index','',1,0,'C','0','0','system:dept:list','tree',103,1,'2026-01-02 21:32:02',NULL,NULL,'部门管理菜单'),(104,'岗位管理',1,5,'post','system/post/index','',1,0,'C','0','0','system:post:list','post',103,1,'2026-01-02 21:32:02',NULL,NULL,'岗位管理菜单'),(105,'字典管理',1,6,'dict','system/dict/index','',1,0,'C','0','0','system:dict:list','dict',103,1,'2026-01-02 21:32:02',NULL,NULL,'字典管理菜单'),(106,'参数设置',1,7,'config','system/config/index','',1,0,'C','0','0','system:config:list','edit',103,1,'2026-01-02 21:32:02',NULL,NULL,'参数设置菜单'),(107,'通知公告',1,8,'notice','system/notice/index','',1,0,'C','0','0','system:notice:list','message',103,1,'2026-01-02 21:32:02',NULL,NULL,'通知公告菜单'),(108,'日志管理',1,9,'log','','',1,0,'M','0','0','','log',103,1,'2026-01-02 21:32:02',NULL,NULL,'日志管理菜单'),(109,'在线用户',2,1,'online','monitor/online/index','',1,0,'C','0','0','monitor:online:list','online',103,1,'2026-01-02 21:32:02',NULL,NULL,'在线用户菜单'),(113,'缓存监控',2,5,'cache','monitor/cache/index','',1,0,'C','0','0','monitor:cache:list','redis',103,1,'2026-01-02 21:32:02',NULL,NULL,'缓存监控菜单'),(118,'文件管理',1,10,'oss','system/oss/index','',1,0,'C','0','0','system:oss:list','upload',103,1,'2026-01-02 21:32:02',NULL,NULL,'文件管理菜单'),(123,'客户端管理',1,11,'client','system/client/index','',1,0,'C','0','0','system:client:list','international',103,1,'2026-01-02 21:32:02',NULL,NULL,'客户端管理菜单'),(130,'分配用户',1,2,'role-auth/user/:roleId','system/role/authUser','',1,1,'C','1','0','system:role:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,'/system/role'),(131,'分配角色',1,1,'user-auth/role/:userId','system/user/authRole','',1,1,'C','1','0','system:user:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,'/system/user'),(132,'字典数据',1,6,'dict-data/index/:dictId','system/dict/data','',1,1,'C','1','0','system:dict:list','#',103,1,'2026-01-02 21:32:02',NULL,NULL,'/system/dict'),(133,'文件配置管理',1,10,'oss-config/index','system/oss/config','',1,1,'C','1','0','system:ossConfig:list','#',103,1,'2026-01-02 21:32:02',NULL,NULL,'/system/oss'),(500,'操作日志',108,1,'operlog','monitor/operlog/index','',1,0,'C','0','0','monitor:operlog:list','form',103,1,'2026-01-02 21:32:02',NULL,NULL,'操作日志菜单'),(501,'登录日志',108,2,'logininfor','monitor/logininfor/index','',1,0,'C','0','0','monitor:logininfor:list','logininfor',103,1,'2026-01-02 21:32:02',NULL,NULL,'登录日志菜单'),(1001,'用户查询',100,1,'','','',1,0,'F','0','0','system:user:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1002,'用户新增',100,2,'','','',1,0,'F','0','0','system:user:add','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1003,'用户修改',100,3,'','','',1,0,'F','0','0','system:user:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1004,'用户删除',100,4,'','','',1,0,'F','0','0','system:user:remove','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1005,'用户导出',100,5,'','','',1,0,'F','0','0','system:user:export','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1006,'用户导入',100,6,'','','',1,0,'F','0','0','system:user:import','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1007,'重置密码',100,7,'','','',1,0,'F','0','0','system:user:resetPwd','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1008,'角色查询',101,1,'','','',1,0,'F','0','0','system:role:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1009,'角色新增',101,2,'','','',1,0,'F','0','0','system:role:add','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1010,'角色修改',101,3,'','','',1,0,'F','0','0','system:role:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1011,'角色删除',101,4,'','','',1,0,'F','0','0','system:role:remove','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1012,'角色导出',101,5,'','','',1,0,'F','0','0','system:role:export','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1013,'菜单查询',102,1,'','','',1,0,'F','0','0','system:menu:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1014,'菜单新增',102,2,'','','',1,0,'F','0','0','system:menu:add','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1015,'菜单修改',102,3,'','','',1,0,'F','0','0','system:menu:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1016,'菜单删除',102,4,'','','',1,0,'F','0','0','system:menu:remove','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1017,'部门查询',103,1,'','','',1,0,'F','0','0','system:dept:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1018,'部门新增',103,2,'','','',1,0,'F','0','0','system:dept:add','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1019,'部门修改',103,3,'','','',1,0,'F','0','0','system:dept:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1020,'部门删除',103,4,'','','',1,0,'F','0','0','system:dept:remove','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1021,'岗位查询',104,1,'','','',1,0,'F','0','0','system:post:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1022,'岗位新增',104,2,'','','',1,0,'F','0','0','system:post:add','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1023,'岗位修改',104,3,'','','',1,0,'F','0','0','system:post:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1024,'岗位删除',104,4,'','','',1,0,'F','0','0','system:post:remove','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1025,'岗位导出',104,5,'','','',1,0,'F','0','0','system:post:export','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1026,'字典查询',105,1,'#','','',1,0,'F','0','0','system:dict:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1027,'字典新增',105,2,'#','','',1,0,'F','0','0','system:dict:add','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1028,'字典修改',105,3,'#','','',1,0,'F','0','0','system:dict:edit','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1029,'字典删除',105,4,'#','','',1,0,'F','0','0','system:dict:remove','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1030,'字典导出',105,5,'#','','',1,0,'F','0','0','system:dict:export','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1031,'参数查询',106,1,'#','','',1,0,'F','0','0','system:config:query','#',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(1032,'参数新增',106,2,'#','','',1,0,'F','0','0','system:config:add','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1033,'参数修改',106,3,'#','','',1,0,'F','0','0','system:config:edit','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1034,'参数删除',106,4,'#','','',1,0,'F','0','0','system:config:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1035,'参数导出',106,5,'#','','',1,0,'F','0','0','system:config:export','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1036,'公告查询',107,1,'#','','',1,0,'F','0','0','system:notice:query','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1037,'公告新增',107,2,'#','','',1,0,'F','0','0','system:notice:add','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1038,'公告修改',107,3,'#','','',1,0,'F','0','0','system:notice:edit','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1039,'公告删除',107,4,'#','','',1,0,'F','0','0','system:notice:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1040,'操作查询',500,1,'#','','',1,0,'F','0','0','monitor:operlog:query','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1041,'操作删除',500,2,'#','','',1,0,'F','0','0','monitor:operlog:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1042,'日志导出',500,4,'#','','',1,0,'F','0','0','monitor:operlog:export','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1043,'登录查询',501,1,'#','','',1,0,'F','0','0','monitor:logininfor:query','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1044,'登录删除',501,2,'#','','',1,0,'F','0','0','monitor:logininfor:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1045,'日志导出',501,3,'#','','',1,0,'F','0','0','monitor:logininfor:export','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1046,'在线查询',109,1,'#','','',1,0,'F','0','0','monitor:online:query','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1047,'批量强退',109,2,'#','','',1,0,'F','0','0','monitor:online:batchLogout','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1048,'单条强退',109,3,'#','','',1,0,'F','0','0','monitor:online:forceLogout','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1050,'账户解锁',501,4,'#','','',1,0,'F','0','0','monitor:logininfor:unlock','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1061,'客户端管理查询',123,1,'#','','',1,0,'F','0','0','system:client:query','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1062,'客户端管理新增',123,2,'#','','',1,0,'F','0','0','system:client:add','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1063,'客户端管理修改',123,3,'#','','',1,0,'F','0','0','system:client:edit','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1064,'客户端管理删除',123,4,'#','','',1,0,'F','0','0','system:client:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1065,'客户端管理导出',123,5,'#','','',1,0,'F','0','0','system:client:export','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1600,'文件查询',118,1,'#','','',1,0,'F','0','0','system:oss:query','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1601,'文件上传',118,2,'#','','',1,0,'F','0','0','system:oss:upload','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1602,'文件下载',118,3,'#','','',1,0,'F','0','0','system:oss:download','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1603,'文件删除',118,4,'#','','',1,0,'F','0','0','system:oss:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1620,'配置列表',118,5,'#','','',1,0,'F','0','0','system:ossConfig:list','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1621,'配置添加',118,6,'#','','',1,0,'F','0','0','system:ossConfig:add','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1622,'配置编辑',118,6,'#','','',1,0,'F','0','0','system:ossConfig:edit','#',103,1,'2026-01-02 21:32:03',NULL,NULL,''),(1623,'配置删除',118,6,'#','','',1,0,'F','0','0','system:ossConfig:remove','#',103,1,'2026-01-02 21:32:03',NULL,NULL,'');
/*!40000 ALTER TABLE `sys_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_notice`
--

DROP TABLE IF EXISTS `sys_notice`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_notice` (
  `notice_id` bigint NOT NULL COMMENT '公告ID',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `notice_title` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公告标题',
  `notice_type` char(1) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公告类型（1通知 2公告）',
  `notice_content` longblob COMMENT '公告内容',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '公告状态（0正常 1关闭）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`notice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='通知公告表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_notice`
--

LOCK TABLES `sys_notice` WRITE;
/*!40000 ALTER TABLE `sys_notice` DISABLE KEYS */;
INSERT INTO `sys_notice` VALUES (1,'000000','温馨提醒：2018-07-01 新版本发布啦','2',_binary '<p>新版本内容111</p>','0',103,1,'2026-01-02 21:32:04',1,NULL,'管理员'),(2,'000000','维护通知：2018-07-01 系统凌晨维护','1',_binary '维护内容','0',103,1,'2026-01-02 21:32:04',NULL,NULL,'管理员');
/*!40000 ALTER TABLE `sys_notice` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_oper_log`
--

DROP TABLE IF EXISTS `sys_oper_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_oper_log` (
  `oper_id` bigint NOT NULL COMMENT '日志主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `title` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '模块标题',
  `business_type` int DEFAULT '0' COMMENT '业务类型（0其它 1新增 2修改 3删除）',
  `method` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '方法名称',
  `request_method` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求方式',
  `operator_type` int DEFAULT '0' COMMENT '操作类别（0其它 1后台用户 2手机端用户）',
  `oper_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作人员',
  `dept_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '部门名称',
  `oper_url` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求URL',
  `oper_ip` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '主机地址',
  `oper_location` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作地点',
  `oper_param` varchar(4000) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '请求参数',
  `json_result` varchar(4000) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '返回参数',
  `status` int DEFAULT '0' COMMENT '操作状态（0正常 1异常）',
  `error_msg` varchar(4000) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '错误消息',
  `oper_time` datetime DEFAULT NULL COMMENT '操作时间',
  `cost_time` bigint DEFAULT '0' COMMENT '消耗时间',
  PRIMARY KEY (`oper_id`),
  KEY `idx_sys_oper_log_bt` (`business_type`),
  KEY `idx_sys_oper_log_s` (`status`),
  KEY `idx_sys_oper_log_ot` (`oper_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_oper_log`
--

LOCK TABLES `sys_oper_log` WRITE;
/*!40000 ALTER TABLE `sys_oper_log` DISABLE KEYS */;
INSERT INTO `sys_oper_log` VALUES (2009482281256595457,'000000','用户管理',1,'org.dromara.system.controller.system.SysUserController.add()','POST',1,'admin','研发部门','/system/user','192.168.3.216','内网IP','{\"createDept\":null,\"createBy\":null,\"createTime\":null,\"updateBy\":null,\"updateTime\":null,\"userId\":\"2009482281181097986\",\"deptId\":100,\"userName\":\"21321323\",\"nickName\":\"11\",\"userType\":null,\"email\":null,\"phonenumber\":null,\"sex\":null,\"status\":\"0\",\"remark\":\"\",\"roleIds\":[4],\"postIds\":[],\"roleId\":null,\"userIds\":null,\"excludeUserIds\":null,\"superAdmin\":false}','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:28:05',106),(2009482351238557698,'000000','用户管理',3,'org.dromara.system.controller.system.SysUserController.remove()','DELETE',1,'admin','研发部门','/system/user/3','192.168.3.216','内网IP','[3]','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:28:22',19),(2009482746887254017,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/5','192.168.3.216','内网IP','5','{\"code\":601,\"msg\":\"存在子菜单,不允许删除\",\"data\":null}',0,'','2026-01-09 12:29:56',8),(2009482764692074498,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/5','192.168.3.216','内网IP','5','{\"code\":601,\"msg\":\"存在子菜单,不允许删除\",\"data\":null}',0,'','2026-01-09 12:30:00',7),(2009482793901207554,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1506','192.168.3.216','内网IP','1506','{\"code\":601,\"msg\":\"存在子菜单,不允许删除\",\"data\":null}',0,'','2026-01-09 12:30:07',3),(2009482806739972098,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1506','192.168.3.216','内网IP','1506','{\"code\":601,\"msg\":\"存在子菜单,不允许删除\",\"data\":null}',0,'','2026-01-09 12:30:10',8),(2009482840596393986,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1507','192.168.3.216','内网IP','1507','{\"code\":601,\"msg\":\"菜单已分配,不允许删除\",\"data\":null}',0,'','2026-01-09 12:30:18',13),(2009482855783968769,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1511','192.168.3.216','内网IP','1511','{\"code\":601,\"msg\":\"菜单已分配,不允许删除\",\"data\":null}',0,'','2026-01-09 12:30:22',9),(2009482866689159170,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1511','192.168.3.216','内网IP','1511','{\"code\":601,\"msg\":\"菜单已分配,不允许删除\",\"data\":null}',0,'','2026-01-09 12:30:24',7),(2009483324971397122,'000000','角色管理',3,'org.dromara.system.controller.system.SysRoleController.remove()','DELETE',1,'admin','研发部门','/system/role/4','192.168.3.216','内网IP','[4]','',1,'仅本人已分配，不能删除!','2026-01-09 12:32:14',16),(2009483341501149186,'000000','角色管理',3,'org.dromara.system.controller.system.SysRoleController.remove()','DELETE',1,'admin','研发部门','/system/role/3','192.168.3.216','内网IP','[3]','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:18',20),(2009483379665121282,'000000','角色管理',2,'org.dromara.system.controller.system.SysRoleController.edit()','PUT',1,'admin','研发部门','/system/role','192.168.3.216','内网IP','{\"createDept\":null,\"createBy\":null,\"createTime\":\"2026-01-02 21:32:02\",\"updateBy\":null,\"updateTime\":null,\"roleId\":4,\"roleName\":\"仅本人\",\"roleKey\":\"test2\",\"roleSort\":4,\"dataScope\":\"5\",\"menuCheckStrictly\":true,\"deptCheckStrictly\":true,\"status\":\"0\",\"remark\":\"\",\"menuIds\":[],\"deptIds\":[],\"superAdmin\":false}','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:27',21),(2009483429640253442,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1507','192.168.3.216','内网IP','1507','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:39',10),(2009483443045249026,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1508','192.168.3.216','内网IP','1508','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:42',14),(2009483453371625473,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1509','192.168.3.216','内网IP','1509','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:44',10),(2009483465002430465,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1510','192.168.3.216','内网IP','1510','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:47',15),(2009483475274280961,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1511','192.168.3.216','内网IP','1511','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:50',11),(2009483487098023937,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1506','192.168.3.216','内网IP','1506','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:32:52',12),(2009483498460393473,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1500','192.168.3.216','内网IP','1500','{\"code\":601,\"msg\":\"存在子菜单,不允许删除\",\"data\":null}',0,'','2026-01-09 12:32:55',11),(2009483527656943617,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1502','192.168.3.216','内网IP','1502','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:02',13),(2009483537698107393,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1504','192.168.3.216','内网IP','1504','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:04',9),(2009483548309696513,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1505','192.168.3.216','内网IP','1505','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:07',11),(2009483561320427522,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1503','192.168.3.216','内网IP','1503','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:10',12),(2009483573450354690,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1501','192.168.3.216','内网IP','1501','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:13',10),(2009483585810968578,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1500','192.168.3.216','内网IP','1500','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:16',12),(2009483596468695041,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/5','192.168.3.216','内网IP','5','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:18',17),(2009483610926460930,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/4','192.168.3.216','内网IP','4','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:22',10),(2009483652395544577,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1060','192.168.3.216','内网IP','1060','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:32',9),(2009483663116185602,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1059','192.168.3.216','内网IP','1059','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:34',12),(2009483676231774209,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1057','192.168.3.216','内网IP','1057','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:37',8),(2009483686126137346,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1058','192.168.3.216','内网IP','1058','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:40',10),(2009483697534644226,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1056','192.168.3.216','内网IP','1056','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:43',10),(2009483709274501122,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1055','192.168.3.216','内网IP','1055','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:45',9),(2009483721207296001,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/116','192.168.3.216','内网IP','116','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:48',12),(2009483732674523138,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/115','192.168.3.216','内网IP','115','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:51',9),(2009483745370681345,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/3','192.168.3.216','内网IP','3','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:33:54',6),(2009483891068219393,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/120','192.168.3.216','内网IP','120','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:34:29',8),(2009483902057295874,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/117','192.168.3.216','内网IP','117','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 12:34:31',13),(2009504911787405314,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1615','192.168.3.216','内网IP','1615','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:00',10),(2009504937599152129,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1614','192.168.3.216','内网IP','1614','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:07',6),(2009504953772388353,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1613','192.168.3.216','内网IP','1613','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:10',11),(2009504967613591553,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1612','192.168.3.216','内网IP','1612','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:14',14),(2009504980070674434,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1611','192.168.3.216','内网IP','1611','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:17',9),(2009504994125787137,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/122','192.168.3.216','内网IP','122','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:20',11),(2009505020218552322,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1610','192.168.3.216','内网IP','1610','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:26',10),(2009505032293953537,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1609','192.168.3.216','内网IP','1609','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:29',15),(2009505044520349697,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1608','192.168.3.216','内网IP','1608','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:32',13),(2009505055035469825,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1607','192.168.3.216','内网IP','1607','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:35',9),(2009505067396083714,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/1606','192.168.3.216','内网IP','1606','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:38',17),(2009505080352288769,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/121','192.168.3.216','内网IP','121','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:41',10),(2009505094344486913,'000000','菜单管理',3,'org.dromara.system.controller.system.SysMenuController.remove()','DELETE',1,'admin','研发部门','/system/menu/6','192.168.3.216','内网IP','6','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-09 13:58:44',14),(2009971775521013762,'000000','部门管理',1,'org.dromara.system.controller.system.SysDeptController.add()','POST',1,'admin','研发部门','/system/dept','192.168.3.216','内网IP','{\"createDept\":null,\"createBy\":null,\"createTime\":null,\"updateBy\":null,\"updateTime\":null,\"deptId\":null,\"parentId\":101,\"deptName\":\"aa\",\"deptCategory\":\"aa\",\"orderNum\":0,\"leader\":null,\"phone\":null,\"email\":null,\"status\":\"0\",\"belongDeptId\":null}','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-10 20:53:09',18),(2009971794336661505,'000000','部门管理',3,'org.dromara.system.controller.system.SysDeptController.remove()','DELETE',1,'admin','研发部门','/system/dept/2009971775458099201','192.168.3.216','内网IP','\"2009971775458099201\"','{\"code\":200,\"msg\":\"操作成功\",\"data\":null}',0,'','2026-01-10 20:53:14',33);
/*!40000 ALTER TABLE `sys_oper_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_oss`
--

DROP TABLE IF EXISTS `sys_oss`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_oss` (
  `oss_id` bigint NOT NULL COMMENT '对象存储主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `file_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件名',
  `original_name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '原名',
  `file_suffix` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '文件后缀名',
  `url` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'URL地址',
  `ext1` text COLLATE utf8mb4_unicode_ci COMMENT '扩展字段',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `create_by` bigint DEFAULT NULL COMMENT '上传人',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新人',
  `service` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'minio' COMMENT '服务商',
  PRIMARY KEY (`oss_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='OSS对象存储表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_oss`
--

LOCK TABLES `sys_oss` WRITE;
/*!40000 ALTER TABLE `sys_oss` DISABLE KEYS */;
INSERT INTO `sys_oss` VALUES (268730381954449408,'000000','2026/01/11/1345cc240d6048168410ef27bcea8203.jpeg','get.jpeg','.jpeg','http://127.0.0.1:9000/ruoyi/2026/01/11/1345cc240d6048168410ef27bcea8203.jpeg','{\"fileSize\":51736,\"contentType\":\"image/jpeg\"}',103,NULL,1,NULL,NULL,'minio'),(268738273663254528,'000000','2026/01/11/83e3563733454008bea667ae14e78dc1.jpeg','get.jpeg','.jpeg','http://127.0.0.1:9000/ruoyi/2026/01/11/83e3563733454008bea667ae14e78dc1.jpeg','{\"fileSize\":251757,\"contentType\":\"image/png\"}',103,NULL,1,NULL,NULL,'minio'),(268739580851650560,'000000','2026/01/11/617999357ff94a2783db9834709777cd.jpeg','get.jpeg','.jpeg','http://127.0.0.1:9000/ruoyi/2026/01/11/617999357ff94a2783db9834709777cd.jpeg','{\"fileSize\":251757,\"contentType\":\"image/png\"}',103,NULL,1,NULL,NULL,'minio'),(268742981341675520,'000000','2026/01/11/49cd87c32b254f7e8a706a85f21e4fd2.jpg','jvm_gc.jpg','.jpg','http://127.0.0.1:9000/ruoyi/2026/01/11/49cd87c32b254f7e8a706a85f21e4fd2.jpg','{\"fileSize\":235532,\"contentType\":\"image/png\"}',103,NULL,1,NULL,NULL,'minio');
/*!40000 ALTER TABLE `sys_oss` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_oss_config`
--

DROP TABLE IF EXISTS `sys_oss_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_oss_config` (
  `oss_config_id` bigint NOT NULL COMMENT '主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `config_key` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '配置key',
  `access_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT 'accessKey',
  `secret_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '秘钥',
  `bucket_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '桶名称',
  `prefix` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '前缀',
  `endpoint` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '访问站点',
  `domain` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '自定义域名',
  `is_https` char(1) COLLATE utf8mb4_unicode_ci DEFAULT 'N' COMMENT '是否https（Y=是,N=否）',
  `region` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '域',
  `access_policy` char(1) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '1' COMMENT '桶权限类型(0=private 1=public 2=custom)',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '1' COMMENT '是否默认（0=是,1=否）',
  `ext1` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '扩展字段',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`oss_config_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对象存储配置表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_oss_config`
--

LOCK TABLES `sys_oss_config` WRITE;
/*!40000 ALTER TABLE `sys_oss_config` DISABLE KEYS */;
INSERT INTO `sys_oss_config` VALUES (1,'000000','minio','minioadmin','minioadmin','ruoyi','','127.0.0.1:9000','','N','','1','0','',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04',NULL),(2,'000000','qiniu','XXXXXXXXXXXXXXX','XXXXXXXXXXXXXXX','ruoyi','','s3-cn-north-1.qiniucs.com','','N','','1','1','',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04',NULL),(3,'000000','aliyun','XXXXXXXXXXXXXXX','XXXXXXXXXXXXXXX','ruoyi','','oss-cn-beijing.aliyuncs.com','','N','','1','1','',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04',NULL),(4,'000000','qcloud','XXXXXXXXXXXXXXX','XXXXXXXXXXXXXXX','ruoyi-1240000000','','cos.ap-beijing.myqcloud.com','','N','ap-beijing','1','1','',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04',NULL),(5,'000000','image','ruoyi','ruoyi123','ruoyi','image','127.0.0.1:9000','','N','','1','1','',103,1,'2026-01-02 21:32:04',1,'2026-01-02 21:32:04',NULL);
/*!40000 ALTER TABLE `sys_oss_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_post`
--

DROP TABLE IF EXISTS `sys_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_post` (
  `post_id` bigint NOT NULL COMMENT '岗位ID',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `dept_id` bigint NOT NULL COMMENT '部门id',
  `post_code` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '岗位编码',
  `post_category` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '岗位类别编码',
  `post_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '岗位名称',
  `post_sort` int NOT NULL COMMENT '显示顺序',
  `status` char(1) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '状态（0正常 1停用）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='岗位信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_post`
--

LOCK TABLES `sys_post` WRITE;
/*!40000 ALTER TABLE `sys_post` DISABLE KEYS */;
INSERT INTO `sys_post` VALUES (1,'000000',103,'ceo',NULL,'董事长',1,'0',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(2,'000000',100,'se',NULL,'项目经理',2,'0',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(3,'000000',100,'hr',NULL,'人力资源',3,'0',103,1,'2026-01-02 21:32:02',NULL,NULL,''),(4,'000000',100,'user',NULL,'普通员工',4,'0',103,1,'2026-01-02 21:32:02',NULL,NULL,'');
/*!40000 ALTER TABLE `sys_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role`
--

DROP TABLE IF EXISTS `sys_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `role_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名称',
  `role_key` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色权限字符串',
  `role_sort` int NOT NULL COMMENT '显示顺序',
  `data_scope` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '1' COMMENT '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限 6：部门及以下或本人数据权限）',
  `menu_check_strictly` tinyint(1) DEFAULT '1' COMMENT '菜单树选择项是否关联显示',
  `dept_check_strictly` tinyint(1) DEFAULT '1' COMMENT '部门树选择项是否关联显示',
  `status` char(1) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色状态（0正常 1停用）',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role`
--

LOCK TABLES `sys_role` WRITE;
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` VALUES (1,'000000','超级管理员','superadmin',1,'1',1,1,'0','0',103,1,'2026-01-02 21:32:02',NULL,NULL,'超级管理员'),(3,'000000','本部门及以下','test1',3,'4',1,1,'0','1',103,1,'2026-01-02 21:32:02',1,'2026-01-09 12:32:18',''),(4,'000000','仅本人','test2',4,'5',1,1,'0','0',103,1,'2026-01-02 21:32:02',1,'2026-01-09 12:32:27',''),(268698441482240000,'000000','11','11',1,'1',0,0,'0','0',103,1,NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_dept`
--

DROP TABLE IF EXISTS `sys_role_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_dept` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `dept_id` bigint NOT NULL COMMENT '部门ID',
  PRIMARY KEY (`role_id`,`dept_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色和部门关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_dept`
--

LOCK TABLES `sys_role_dept` WRITE;
/*!40000 ALTER TABLE `sys_role_dept` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_role_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_menu`
--

DROP TABLE IF EXISTS `sys_role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_menu` (
  `role_id` bigint NOT NULL COMMENT '角色ID',
  `menu_id` bigint NOT NULL COMMENT '菜单ID',
  PRIMARY KEY (`role_id`,`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色和菜单关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_menu`
--

LOCK TABLES `sys_role_menu` WRITE;
/*!40000 ALTER TABLE `sys_role_menu` DISABLE KEYS */;
INSERT INTO `sys_role_menu` VALUES (268698441482240000,1),(268698441482240000,2),(268698441482240000,100),(268698441482240000,101),(268698441482240000,102),(268698441482240000,103),(268698441482240000,104),(268698441482240000,105),(268698441482240000,106),(268698441482240000,107),(268698441482240000,108),(268698441482240000,109),(268698441482240000,113),(268698441482240000,118),(268698441482240000,123),(268698441482240000,130),(268698441482240000,131),(268698441482240000,132),(268698441482240000,133),(268698441482240000,500),(268698441482240000,501),(268698441482240000,1001),(268698441482240000,1002),(268698441482240000,1003),(268698441482240000,1004),(268698441482240000,1005),(268698441482240000,1006),(268698441482240000,1007),(268698441482240000,1008),(268698441482240000,1009),(268698441482240000,1010),(268698441482240000,1011),(268698441482240000,1012),(268698441482240000,1013),(268698441482240000,1014),(268698441482240000,1015),(268698441482240000,1016),(268698441482240000,1017),(268698441482240000,1018),(268698441482240000,1019),(268698441482240000,1020),(268698441482240000,1021),(268698441482240000,1022),(268698441482240000,1023),(268698441482240000,1024),(268698441482240000,1025),(268698441482240000,1026),(268698441482240000,1027),(268698441482240000,1028),(268698441482240000,1029),(268698441482240000,1030),(268698441482240000,1031),(268698441482240000,1032),(268698441482240000,1033),(268698441482240000,1034),(268698441482240000,1035),(268698441482240000,1036),(268698441482240000,1037),(268698441482240000,1038),(268698441482240000,1039),(268698441482240000,1040),(268698441482240000,1041),(268698441482240000,1042),(268698441482240000,1043),(268698441482240000,1044),(268698441482240000,1045),(268698441482240000,1046),(268698441482240000,1047),(268698441482240000,1048),(268698441482240000,1050),(268698441482240000,1061),(268698441482240000,1062),(268698441482240000,1063),(268698441482240000,1064),(268698441482240000,1065),(268698441482240000,1600),(268698441482240000,1601),(268698441482240000,1602),(268698441482240000,1603),(268698441482240000,1620),(268698441482240000,1621),(268698441482240000,1622),(268698441482240000,1623);
/*!40000 ALTER TABLE `sys_role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_social`
--

DROP TABLE IF EXISTS `sys_social`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_social` (
  `id` bigint NOT NULL COMMENT '主键',
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户id',
  `auth_id` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台+平台唯一id',
  `source` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户来源',
  `open_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '平台编号唯一id',
  `user_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '登录账号',
  `nick_name` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户昵称',
  `email` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户邮箱',
  `avatar` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '头像地址',
  `access_token` varchar(2000) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户的授权令牌',
  `expire_in` int DEFAULT NULL COMMENT '用户的授权令牌的有效期，部分平台可能没有',
  `refresh_token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '刷新令牌，部分平台可能没有',
  `access_code` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '平台的授权信息，部分平台可能没有',
  `union_id` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户的 unionid',
  `scope` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '授予的权限，部分平台可能没有',
  `token_type` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '个别平台的授权信息，部分平台可能没有',
  `id_token` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'id token，部分平台可能没有',
  `mac_algorithm` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '小米平台用户的附带属性，部分平台可能没有',
  `mac_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '小米平台用户的附带属性，部分平台可能没有',
  `code` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '用户的授权code，部分平台可能没有',
  `oauth_token` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Twitter平台用户的附带属性，部分平台可能没有',
  `oauth_token_secret` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Twitter平台用户的附带属性，部分平台可能没有',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='社会化关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_social`
--

LOCK TABLES `sys_social` WRITE;
/*!40000 ALTER TABLE `sys_social` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_social` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_tenant`
--

DROP TABLE IF EXISTS `sys_tenant`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_tenant` (
  `id` bigint NOT NULL COMMENT 'id',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租户编号',
  `contact_user_name` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系人',
  `contact_phone` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '联系电话',
  `company_name` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '企业名称',
  `license_number` varchar(30) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '统一社会信用代码',
  `address` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '地址',
  `intro` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '企业简介',
  `domain` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '域名',
  `remark` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `package_id` bigint DEFAULT NULL COMMENT '租户套餐编号',
  `expire_time` datetime DEFAULT NULL COMMENT '过期时间',
  `account_count` int DEFAULT '-1' COMMENT '用户数量（-1不限制）',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '租户状态（0正常 1停用）',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_tenant`
--

LOCK TABLES `sys_tenant` WRITE;
/*!40000 ALTER TABLE `sys_tenant` DISABLE KEYS */;
INSERT INTO `sys_tenant` VALUES (1,'000000','管理组','15888888888','XXX有限公司',NULL,NULL,'多租户通用后台管理管理系统',NULL,NULL,NULL,NULL,-1,'0','0',103,1,'2026-01-02 21:32:02',NULL,NULL);
/*!40000 ALTER TABLE `sys_tenant` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_tenant_package`
--

DROP TABLE IF EXISTS `sys_tenant_package`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_tenant_package` (
  `package_id` bigint NOT NULL COMMENT '租户套餐id',
  `package_name` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '套餐名称',
  `menu_ids` varchar(3000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '关联菜单id',
  `remark` varchar(200) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  `menu_check_strictly` tinyint(1) DEFAULT '1' COMMENT '菜单树选择项是否关联显示',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '状态（0正常 1停用）',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`package_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='租户套餐表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_tenant_package`
--

LOCK TABLES `sys_tenant_package` WRITE;
/*!40000 ALTER TABLE `sys_tenant_package` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_tenant_package` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user`
--

DROP TABLE IF EXISTS `sys_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `dept_id` bigint DEFAULT NULL COMMENT '部门ID',
  `user_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户账号',
  `nick_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户昵称',
  `user_type` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT 'sys_user' COMMENT '用户类型（sys_user系统用户）',
  `email` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '用户邮箱',
  `phonenumber` varchar(11) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '手机号码',
  `sex` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '用户性别（0男 1女 2未知）',
  `avatar` bigint DEFAULT NULL COMMENT '头像地址',
  `password` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '密码',
  `status` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '帐号状态（0正常 1停用）',
  `del_flag` char(1) COLLATE utf8mb4_unicode_ci DEFAULT '0' COMMENT '删除标志（0代表存在 1代表删除）',
  `login_ip` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '最后登录IP',
  `login_date` datetime DEFAULT NULL COMMENT '最后登录时间',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_by` bigint DEFAULT NULL COMMENT '创建者',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新者',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `remark` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user`
--

LOCK TABLES `sys_user` WRITE;
/*!40000 ALTER TABLE `sys_user` DISABLE KEYS */;
INSERT INTO `sys_user` VALUES (1,'000000',103,'admin','疯狂的狮子Li','sys_user','crazyLionLi@163.com','15888888888','1',268742981341675520,'$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2','0','0','192.168.3.216','2026-01-11 21:56:19',103,1,'2026-01-02 21:32:02',-1,'2026-01-11 21:56:19','管理员'),(3,'000000',108,'test','本部门及以下 密码666666','sys_user','','','0',NULL,'$2a$10$b8yUzN0C71sbz.PhNOCgJe.Tu1yWC3RNrTyjSQ8p1W0.aaUXUJ.Ne','0','1','127.0.0.1','2026-01-02 21:32:02',103,1,'2026-01-02 21:32:02',1,'2026-01-09 12:28:22',NULL),(4,'000000',102,'test1','仅本人 密码666666','sys_user','','','0',NULL,'$2a$10$b8yUzN0C71sbz.PhNOCgJe.Tu1yWC3RNrTyjSQ8p1W0.aaUXUJ.Ne','0','0','127.0.0.1','2026-01-02 21:32:02',103,1,'2026-01-02 21:32:02',4,'2026-01-02 21:32:02',NULL),(2009482281181097986,'000000',100,'21321323','11','sys_user','','','0',NULL,'$2a$10$ZdTk.GuQ3lCt0hgN51iqdepRaJJCjk7HDmFK0tgGeU/xlaO36k8cW','0','0','',NULL,103,1,'2026-01-09 12:28:05',1,'2026-01-09 12:28:05','');
/*!40000 ALTER TABLE `sys_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_post`
--

DROP TABLE IF EXISTS `sys_user_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_post` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `post_id` bigint NOT NULL COMMENT '岗位ID',
  PRIMARY KEY (`user_id`,`post_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户与岗位关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_post`
--

LOCK TABLES `sys_user_post` WRITE;
/*!40000 ALTER TABLE `sys_user_post` DISABLE KEYS */;
INSERT INTO `sys_user_post` VALUES (1,1);
/*!40000 ALTER TABLE `sys_user_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_user_role`
--

DROP TABLE IF EXISTS `sys_user_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_user_role` (
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `role_id` bigint NOT NULL COMMENT '角色ID',
  PRIMARY KEY (`user_id`,`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户和角色关联表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_user_role`
--

LOCK TABLES `sys_user_role` WRITE;
/*!40000 ALTER TABLE `sys_user_role` DISABLE KEYS */;
INSERT INTO `sys_user_role` VALUES (1,1),(4,4),(2009482281181097986,4);
/*!40000 ALTER TABLE `sys_user_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `test_demo`
--

DROP TABLE IF EXISTS `test_demo`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `test_demo` (
  `id` bigint NOT NULL COMMENT '主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `dept_id` bigint DEFAULT NULL COMMENT '部门id',
  `user_id` bigint DEFAULT NULL COMMENT '用户id',
  `order_num` int DEFAULT '0' COMMENT '排序号',
  `test_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'key键',
  `value` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '值',
  `version` int DEFAULT '0' COMMENT '版本',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新人',
  `del_flag` int DEFAULT '0' COMMENT '删除标志',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test_demo`
--

LOCK TABLES `test_demo` WRITE;
/*!40000 ALTER TABLE `test_demo` DISABLE KEYS */;
INSERT INTO `test_demo` VALUES (1,'000000',102,4,1,'测试数据权限','测试',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(2,'000000',102,3,2,'子节点1','111',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(3,'000000',102,3,3,'子节点2','222',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(4,'000000',108,4,4,'测试数据','demo',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(5,'000000',108,3,13,'子节点11','1111',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(6,'000000',108,3,12,'子节点22','2222',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(7,'000000',108,3,11,'子节点33','3333',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(8,'000000',108,3,10,'子节点44','4444',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(9,'000000',108,3,9,'子节点55','5555',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(10,'000000',108,3,8,'子节点66','6666',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(11,'000000',108,3,7,'子节点77','7777',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(12,'000000',108,3,6,'子节点88','8888',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(13,'000000',108,3,5,'子节点99','9999',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0);
/*!40000 ALTER TABLE `test_demo` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `test_tree`
--

DROP TABLE IF EXISTS `test_tree`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `test_tree` (
  `id` bigint NOT NULL COMMENT '主键',
  `tenant_id` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '000000' COMMENT '租户编号',
  `parent_id` bigint DEFAULT '0' COMMENT '父id',
  `dept_id` bigint DEFAULT NULL COMMENT '部门id',
  `user_id` bigint DEFAULT NULL COMMENT '用户id',
  `tree_name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '值',
  `version` int DEFAULT '0' COMMENT '版本',
  `create_dept` bigint DEFAULT NULL COMMENT '创建部门',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `create_by` bigint DEFAULT NULL COMMENT '创建人',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  `update_by` bigint DEFAULT NULL COMMENT '更新人',
  `del_flag` int DEFAULT '0' COMMENT '删除标志',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='测试树表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `test_tree`
--

LOCK TABLES `test_tree` WRITE;
/*!40000 ALTER TABLE `test_tree` DISABLE KEYS */;
INSERT INTO `test_tree` VALUES (1,'000000',0,102,4,'测试数据权限',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(2,'000000',1,102,3,'子节点1',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(3,'000000',2,102,3,'子节点2',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(4,'000000',0,108,4,'测试树1',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(5,'000000',4,108,3,'子节点11',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(6,'000000',4,108,3,'子节点22',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(7,'000000',4,108,3,'子节点33',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(8,'000000',5,108,3,'子节点44',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(9,'000000',6,108,3,'子节点55',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(10,'000000',7,108,3,'子节点66',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(11,'000000',7,108,3,'子节点77',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(12,'000000',10,108,3,'子节点88',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0),(13,'000000',10,108,3,'子节点99',0,103,'2026-01-02 21:32:04',1,NULL,NULL,0);
/*!40000 ALTER TABLE `test_tree` ENABLE KEYS */;
UNLOCK TABLES;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-01-11 22:43:38
