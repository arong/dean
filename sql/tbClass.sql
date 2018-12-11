CREATE TABLE `tbClass` (
  `iClassID`     int(10)     unsigned NOT NULL AUTO_INCREMENT                                           COMMENT '班级表主键',
  `iGrade`       tinyint(1)  unsigned NOT NULL DEFAULT '0'                                              COMMENT '年级编号',
  `iIndex`       tinyint(2)  unsigned NOT NULL DEFAULT '0'                                              COMMENT '班级编号',
  `vName`        varchar(16)          NOT NULL DEFAULT ''                                               COMMENT '班级名称',
  `iMasterID`    int(10)              NOT NULL DEFAULT '0'                                              COMMENT '班主任老师ID',
  `iStartYear`   int(10)              NOT NULL DEFAULT '0'                                              COMMENT '开学年份',
  `eTerm`        tinyint(1)           NOT NULL DEFAULT '0'                                              COMMENT '学期',
  `eStatus`      tinyint(1)           NOT NULL DEFAULT '1'                                              COMMENT '逻辑状态',
  `dtCreateTime` datetime             NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime` datetime             NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iClassID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
