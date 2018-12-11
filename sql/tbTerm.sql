CREATE TABLE `tbTerm` (
  `iTermID`     int(10) unsigned    NOT NULL AUTO_INCREMENT       COMMENT '主键',
  `iSchoolYear` int(10) unsigned    NOT NULL DEFAULT '0'          COMMENT '学年',
  `eTerm`       tinyint(1) unsigned NOT NULL DEFAULT '0'          COMMENT '学期',
  `dtBegin`     date                NOT NULL DEFAULT '0000-00-00' COMMENT '学期开始日期',
  `dtEnd`       date                NOT NULL DEFAULT '0000-00-00' COMMENT '学期结束日期',
  PRIMARY KEY (`iTermID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
