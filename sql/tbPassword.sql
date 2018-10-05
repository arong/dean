CREATE TABLE `tbPassword`  (
  `iUserID`      bigint(20)  NOT NULL DEFAULT '0'                                              COMMENT 'tbUser表主键',
  `vPassword`    varchar(64) NOT NULL DEFAULT ''                                               COMMENT '密码',
  `dtCreateTime` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iUserID`) USING BTREE
) ENGINE = InnoDB DEFAULT CHARSET=utf8;
