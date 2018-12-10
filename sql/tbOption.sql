create table tbOption (
 `iOptionID` int(10) not AUTO_INCREMENT '0' comment '问卷id',
 `iQuestionID` int(10) not null default '0' comment '问卷id',
 `vOption` varchar(64) not null default '' comment '问卷标题',
 `iIndex` tinyint(2) not null default '0' comment '问题编号',
 PRIMARY KEY (`iOptionID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
