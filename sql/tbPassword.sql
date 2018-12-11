CREATE TABLE `tbPassword` (
  `iPasswordID`  bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT                                        COMMENT '主键',
  `iUserID`      bigint(20)          NOT NULL DEFAULT '0'                                           COMMENT 'tbUser表主键',
  `eType`        enum('1','2')       NOT NULL DEFAULT '1'                                           COMMENT 'ID类型',
  `vLoginName`   varchar(32)         NOT NULL DEFAULT ''                                            COMMENT '登录名',
  `vPassword`    varchar(64)         NOT NULL DEFAULT ''                                            COMMENT '密码',
  `dtCreateTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP                             COMMENT '创建时间',
  `dtModifyTime` datetime            NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '修改时间',
  PRIMARY KEY (`iPasswordID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
