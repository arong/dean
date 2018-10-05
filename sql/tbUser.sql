CREATE TABLE `tbuser` (
  `iUserID` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `vRegistNumber` varchar(16) NOT NULL DEFAULT '' COMMENT '学号',
  `vUserName` varchar(32) NOT NULL DEFAULT '' COMMENT '学生姓名',
  `eGender` enum('3','2','1') NOT NULL DEFAULT '3' COMMENT '性别： 1男 2女 3未知',
  `eStatus` tinyint(4) NOT NULL DEFAULT '1' COMMENT '逻辑状态',
  `dtCreateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `dtModifyTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`iUserID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
