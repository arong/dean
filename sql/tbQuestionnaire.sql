CREATE TABLE `tbQuestionnaire` (
   `iQuestionnaireID` int(10) NOT NULL AUTO_INCREMENT COMMENT '问卷id',
   `vTitle` varchar(64) NOT NULL DEFAULT '' COMMENT '问卷标题',
   `eDraftStatus` tinyint(1) NOT NULL DEFAULT '0' COMMENT '文稿状态',
   `dtStartTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '开发日期',
   `dtStopTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT '截至日期',
   `vEditorName` varchar(32) NOT NULL DEFAULT '' COMMENT '编辑',
   PRIMARY KEY (`iQuestionnaireID`)
 ) ENGINE=InnoDB DEFAULT CHARSET=utf8
