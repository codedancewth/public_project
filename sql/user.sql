CREATE TABLE `user` (
                        `id` int NOT NULL AUTO_INCREMENT,
                        `user_name` varchar(128) NOT NULL DEFAULT '' COMMENT '用户的姓名',
                        `user_account`  varchar(128) NOT NULL DEFAULT '' COMMENT '用户的账号',
                        `user_password` varchar(1024) NOT NULL DEFAULT '' COMMENT '用户的密码',
                        `status` tinyint NOT NULL DEFAULT '0' COMMENT '是否活跃',
                        `created_time` int NOT NULL COMMENT '创建时间',
                        `updated_time` int DEFAULT '0' COMMENT '更新时间',
                        `is_deleted` tinyint NOT NULL DEFAULT '0' COMMENT '是否已经删除,0否 1是',
                        PRIMARY KEY (`id`),
                        CONSTRAINT idx_user_name UNIQUE (user_name),
                        CONSTRAINT idx_user_account UNIQUE (user_account)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COMMENT='用户表';


INSERT INTO `user` ( `user_name`, `user_account`, `user_password`, `status`, `created_time`, `updated_time`, `is_deleted` ) VALUES ( 'dance', 'ww',  '123456',  1,  UNIX_TIMESTAMP(),  0,  0 );