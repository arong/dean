CREATE TABLE tbStudentScore (
  -- `iStudentScoreID` INT(11) UNSIGNED NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `iStudentID`      BIGINT(20)       NOT NULL DEFAULT '0'                                              COMMENT '学生ID',
  `iTermID`         INT(10)          NOT NULL DEFAULT '0'                                              COMMENT '学期ID',
  `eExam`           TINYINT(1)       NOT NULL DEFAULT '0'                                              COMMENT '考试编号',
  `iSubjectID`      INT(10)          NOT NULL DEFAULT '0'                                              COMMENT '课程ID',
  `iScore`          SMALLINT(5)      NOT NULL DEFAULT '0'                                              COMMENT '分数',
  `dtCreateTime`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime`    DATETIME         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iStudentID`,`iTermID`,`eExam`,`iSubjectID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
