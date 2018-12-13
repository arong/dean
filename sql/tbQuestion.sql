create table tbQuestion (
 `iQuestionID` int(10) not null AUTO_INCREMENT comment '问卷id',
 `iQuestionaireID` int(10) not null default '0' comment '问卷id',
 `vQuestion` varchar(64) not null default '' comment '问卷标题',
 `iIndex` tinyint(2) not null default '0' comment '问题编号',
 `eType` tinyint(1) not null default '0' comment '问题类型, 1: 单选, 2: 多选, 3: 文本',
 `bRequired` tinyint(1)  not null default '0' comment '是否为必填项',
 PRIMARY KEY (`iQuestionID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
