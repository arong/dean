CREATE TABLE `tbClassTeacherRelation` (
  `iClassTeacherRelationID` int(11) unsigned NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `iClassID`                int(11)          NOT NULL DEFAULT '0'                                              COMMENT '班级表主键',
  `iSubjectID`              int(10)          NOT NULL DEFAULT '0'                                              COMMENT '科目',
  `iTeacherID`              int(11)          NOT NULL DEFAULT '0'                                              COMMENT '教师表主键',
  `eStatus`                 tinyint(1)       NOT NULL DEFAULT '1'                                              COMMENT '逻辑状态',
  `dtCreateTime`            datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime`            datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '最后修改时间',
  PRIMARY KEY (`iClassTeacherRelationID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
