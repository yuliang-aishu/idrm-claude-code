-- 数据源配置表
-- 创建时间: 2026-01-04

CREATE TABLE `datasource` (
    `id` CHAR(36) NOT NULL COMMENT 'ID (UUID v7)',
    `name` varchar(100) NOT NULL COMMENT '数据源名称',
    `type` varchar(50) NOT NULL COMMENT '数据源类型：mysql/postgresql/redis/mongodb/sqlserver',
    `host` varchar(200) NOT NULL COMMENT '连接地址',
    `port` int NOT NULL COMMENT '连接端口',
    `database` varchar(100) DEFAULT NULL COMMENT '数据库名',
    `username` varchar(100) NOT NULL COMMENT '连接用户名',
    `password` varchar(500) NOT NULL COMMENT '连接密码（AES-256-GCM 加密）',
    `description` varchar(500) DEFAULT NULL COMMENT '描述信息',
    `status` varchar(20) NOT NULL DEFAULT 'enabled' COMMENT '状态：enabled/disabled',
    `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序权重',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据源配置表';
