CREATE TABLE `tbSubject`  (
  `iSubjectID`    int(20) UNSIGNED    NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `vSubjectKey`   varchar(32)         NOT NULL DEFAULT ''                                               COMMENT '课程key',
  `vSubjectName`  varchar(32)         NOT NULL DEFAULT ''                                               COMMENT '课程名称',
  `eStatus`       tinyint(4)          NOT NULL DEFAULT 1                                                COMMENT '逻辑状态',
  `dtCreateTime`  datetime(0)         NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime`  datetime(0)         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iSubjectID`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;
