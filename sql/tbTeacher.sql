CREATE TABLE `tbTeacher` (
  `iTeacherID`   int(11) unsigned   NOT NULL AUTO_INCREMENT                                           COMMENT '主键',
  `iPrimarySubjectID`   int(11) unsigned   NOT NULL  DEFAULT 0                                        COMMENT '主授课程',
  `eGender`      enum('1','2','3')  NOT NULL                                                          COMMENT '性别, 1: 男, 2: 女, 3: 未知',
  `vName`        varchar(32)        NOT NULL DEFAULT ''                                               COMMENT '姓名',
  `vMobile`      varchar(16)        NOT NULL DEFAULT ''                                               COMMENT '手机号',
  `vAddress`     varchar(128)       NOT NULL DEFAULT ''                                               COMMENT '家庭地址',
  `eStatus`      tinyint(1)         NOT NULL DEFAULT '1'                                              COMMENT '逻辑状态',
  `dtBirthday`   date               NOT NULL DEFAULT '2008-01-01 00:00:00'                            COMMENT '生日',
  `dtCreateTime` datetime           NOT NULL DEFAULT CURRENT_TIMESTAMP                                COMMENT '创建时间',
  `dtModifyTime` datetime           NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '修改时间',
  PRIMARY KEY (`iTeacherID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
