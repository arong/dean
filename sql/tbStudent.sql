CREATE TABLE `tbStudent`  (
  `iUserID`       bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `iClassID`      tinyint(1) UNSIGNED NOT NULL DEFAULT '0'                                              COMMENT '年级编号',
  `vRegistNumber` varchar(16)         NOT NULL DEFAULT ''                                               COMMENT '学号',
  `vUserName`     varchar(32)         NOT NULL DEFAULT ''                                               COMMENT '学生姓名',
  `eGender`       enum('1','2','3')   NOT NULL DEFAULT '3'                                              COMMENT '性别： 1男 2女 3未知',
  `eStatus`       tinyint(4)          NOT NULL DEFAULT 1                                                COMMENT '逻辑状态',
  `dtCreateTime`  datetime(0)         NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime`  datetime(0)         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iUserID`) USING BTREE
) ENGINE = InnoDB  CHARACTER SET = utf8;
