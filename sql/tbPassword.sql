CREATE TABLE `tbPassword` (
  `iPasswordID` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `iUserID` bigint(20) NOT NULL DEFAULT '0' COMMENT 'tbUser表主键',
  `eType` enum('2','1') NOT NULL DEFAULT '1',
  `vLoginName` varchar(32) not null default '' comment '登录名',
  `vPassword` varchar(64) NOT NULL DEFAULT '' COMMENT '密码',
  `dtCreateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `dtModifyTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`iPasswordID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
