CREATE TABLE `tbTeacherScore`  (
  `iTeacherScoreID` int(11) unsigned NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `iTeacherID`      int(11)          NOT NULL DEFAULT '0'                                              COMMENT '教师表主键',
  `iScore`          int(11)          NOT NULL DEFAULT '0'                                              COMMENT '分数',
  `eStatus`         tinyint(1)       NOT NULL DEFAULT '1'                                              COMMENT '逻辑状态',
  `dtCreateTime`    datetime(0)      NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime`    datetime(0)      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iTeacherScoreID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
