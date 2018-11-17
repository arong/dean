CREATE TABLE `tbTerm` (
  `iTermID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `dtRegisteDay` date NOT NULL DEFAULT '' comment '开学日期',
  PRIMARY KEY (`iTermID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;