create table tbQuestion (
 `iQuestionID`      int(10)    NOT NULL AUTO_INCREMENT COMMENT '问卷id',
 `iQuestionnaireID` int(10)    NOT NULL DEFAULT '0' COMMENT '问卷id',
 `iIndex`           tinyint(2) NOT NULL DEFAULT '0' COMMENT '问题编号',
 `eType`            tinyint(1) NOT NULL DEFAULT '0' COMMENT '问题类型, 1: 单选, 2: 多选, 3: 文本',
 `bRequired`        tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否为必填项',
 `vQuestion`        text       NOT NULL COMMENT '问卷标题, base64编码的文本',
 `vContent`         text       NOT NULL COMMENT '题目内容, base64编码的json字符串',
 PRIMARY KEY (`iQuestionID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
