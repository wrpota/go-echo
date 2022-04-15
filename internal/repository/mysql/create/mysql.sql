-- 管理员
CREATE TABLE `admin` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary Key',
    `username` varchar(32) NOT NULL DEFAULT '' COMMENT '用户名',
    `password` varchar(100) NOT NULL DEFAULT '' COMMENT '密码',
    `nickname` varchar(60) NOT NULL DEFAULT '' COMMENT '昵称',
    `mobile` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
    `is_used` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用 1:是  -1:否',
    `is_deleted` tinyint(1) NOT NULL DEFAULT '-1' COMMENT '是否删除 1:是  -1:否',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `unique_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理员表';


-- 定时任务 
CREATE TABLE `cron_task` (
    `id` int(11) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name` varchar(64) NOT NULL DEFAULT '' COMMENT '任务名称',
    `spec` varchar(64) NOT NULL DEFAULT '' COMMENT 'crontab 表达式',
    `command` varchar(255) NOT NULL DEFAULT '' COMMENT '执行命令',
    `protocol` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '执行方式 1:shell 2:http',
    `http_method` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT 'http 请求方式 1:get 2:post',
    `timeout` int(11) unsigned NOT NULL DEFAULT '60' COMMENT '超时时间(单位:秒)',
    `retry_times` tinyint(1) NOT NULL DEFAULT '3' COMMENT '重试次数',
    `retry_interval` int(11) NOT NULL DEFAULT '60' COMMENT '重试间隔(单位:秒)',
    `notify_status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '执行结束是否通知 1:不通知 2:失败通知 3:结束通知 4:结果关键字匹配通知',
    `notify_type` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '通知类型 1:邮件 2:webhook',
    `notify_receiver_email` varchar(255) NOT NULL DEFAULT '' COMMENT '通知者邮箱地址(多个用,分割)',
    `notify_keyword` varchar(255) NOT NULL DEFAULT '' COMMENT '通知匹配关键字(多个用,分割)',
    `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
    `is_used` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用 1:是  -1:否',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `created_user` varchar(60) NOT NULL DEFAULT '' COMMENT '创建人',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `updated_user` varchar(60) NOT NULL DEFAULT '' COMMENT '更新人',
PRIMARY KEY (`id`),
KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='后台任务表';