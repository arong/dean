CREATE TABLE `tbVote` (
  `iVoteID`      int(11) unsigned NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `vVoteCode`    varchar(16)      NOT NULL DEFAULT ''                                               COMMENT '投票码',
  `vVoteDetail`  varchar(256)     NOT NULL DEFAULT ''                                               COMMENT '投票详情',
  `dtCreateTime` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iVoteID`,`vVoteCode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
