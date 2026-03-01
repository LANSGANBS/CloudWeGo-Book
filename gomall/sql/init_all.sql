-- Gomall 数据库初始化脚本
-- 在 MacBook M2 上首次启动时执行

-- ==================== 创建数据库 ====================

CREATE DATABASE IF NOT EXISTS `cart` DEFAULT CHARACTER SET = 'utf8mb4';
CREATE DATABASE IF NOT EXISTS `checkout` DEFAULT CHARACTER SET = 'utf8mb4';
CREATE DATABASE IF NOT EXISTS `order` DEFAULT CHARACTER SET = 'utf8mb4';
CREATE DATABASE IF NOT EXISTS `payment` DEFAULT CHARACTER SET = 'utf8mb4';
CREATE DATABASE IF NOT EXISTS `product` DEFAULT CHARACTER SET = 'utf8mb4';
CREATE DATABASE IF NOT EXISTS `user` DEFAULT CHARACTER SET = 'utf8mb4';
CREATE DATABASE IF NOT EXISTS `email` DEFAULT CHARACTER SET = 'utf8mb4';

-- ==================== user 数据库 ====================

USE `user`;

CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `password_hashed` longtext,
  `is_admin` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `user` (id, created_at, updated_at, email, password_hashed, is_admin) VALUES 
(1, '2023-12-26 09:46:19.852', '2023-12-26 09:46:19.852', '123@admin.com', '$2a$10$jTvUFh7Z8Kw0hLV8WrAws.PRQTeuH4gopJ7ZMoiFvwhhz5Vw.bj7C', 0),
(2, '2026-02-09 14:05:55.001', '2026-02-09 14:05:55.001', 'lansganbs@qq.com', '$2a$10$gbXnOKhwO3c/qz.uaD9w6.Del1zBbjTbYwO6Y0mpSmKVs4FOGHsDS', 1);

-- ==================== product 数据库 ====================

USE `product`;

CREATE TABLE IF NOT EXISTS `category` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `description` longtext,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `category` (id, created_at, updated_at, name, description, deleted_at) VALUES 
(1, '2023-12-06 15:05:06.000', '2023-12-06 15:05:06.000', 'T-Shirt', 'T-Shirt', null),
(2, '2023-12-06 15:05:06.000', '2023-12-06 15:05:06.000', 'Sticker', 'Sticker', null),
(5, null, null, 'Test', 'Test', null);

CREATE TABLE IF NOT EXISTS `product` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `description` longtext,
  `picture` longtext,
  `price` float DEFAULT NULL,
  `stock` int DEFAULT '0' COMMENT '库存数量',
  `deleted_at` datetime(3) DEFAULT NULL,
  `sales` bigint DEFAULT '0',
  PRIMARY KEY (`id`),
  KEY `idx_product_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `product` (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES 
(1, '2024-12-06 15:26:19.000', '2026-02-24 23:26:44.601', 'Notebook', 'The cloudwego notebook is a highly efficient and feature-rich notebook designed to meet all your note-taking needs.', '/static/image/notebook.jpeg', 9.9, 1000, null, 20),
(2, '2023-12-06 15:26:19.000', '2026-02-24 22:54:32.827', 'Mouse-Pad', 'The cloudwego mouse pad is a premium-grade accessory designed to enhance your computer usage experience.', '/static/image/mouse-pad.jpeg', 8.8, 960, null, 1021),
(3, '2023-12-06 15:26:19.000', '2026-02-24 22:54:39.751', 'T-Shirt', 'The cloudwego t-shirt is a stylish and comfortable clothing item that allows you to showcase your fashion sense while enjoying maximum comfort.', '/static/image/t-shirt.jpeg', 6.6, 989, null, 1190),
(4, '2023-12-06 15:26:19.000', '2026-02-24 22:54:46.926', 'T-Shirt', 'The cloudwego t-shirt is a stylish and comfortable clothing item that allows you to showcase your fashion sense while enjoying maximum comfort.', '/static/image/t-shirt-1.jpeg', 2.2, 1000, null, 2),
(5, '2025-12-06 15:26:19.000', '2026-02-24 22:54:53.220', 'Sweatshirt', 'The cloudwego Sweatshirt is a cozy and fashionable garment that provides warmth and style during colder weather.', '/static/image/sweatshirt.jpeg', 1.1, 999, null, 5),
(6, '2023-12-06 15:26:19.000', '2026-02-24 22:55:00.233', 'T-Shirt', 'The cloudwego t-shirt is a stylish and comfortable clothing item that allows you to showcase your fashion sense while enjoying maximum comfort.', '/static/image/t-shirt-2.jpeg', 1.8, 1000, null, 100),
(7, '2024-12-06 15:26:19.000', '2026-02-24 22:55:10.390', 'mascot', 'The cloudwego mascot is a charming and captivating representation of the brand, designed to bring joy and a playful spirit to any environment.', '/static/image/logo.jpg', 4.8, 985, null, 15);

CREATE TABLE IF NOT EXISTS `product_category` (
  `category_id` bigint NOT NULL,
  `product_id` bigint NOT NULL,
  PRIMARY KEY (`category_id`,`product_id`),
  KEY `fk_product_category_product` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO `product_category` (category_id, product_id) VALUES 
(2, 1), (2, 2), (1, 3), (1, 4), (1, 5), (1, 6), (2, 7);

CREATE TABLE IF NOT EXISTS `stock` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `quantity` bigint NOT NULL DEFAULT '0' COMMENT '库存数量',
  `reserved` bigint NOT NULL DEFAULT '0' COMMENT '预留数量',
  `available` bigint NOT NULL DEFAULT '0' COMMENT '可用数量',
  `min_stock` bigint NOT NULL DEFAULT '10' COMMENT '最低库存预警值',
  `max_stock` bigint NOT NULL DEFAULT '1000' COMMENT '最高库存预警值',
  `safety_stock` bigint NOT NULL DEFAULT '20' COMMENT '安全库存',
  `unit` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '件' COMMENT '库存单位',
  `warehouse_id` int unsigned DEFAULT '1' COMMENT '仓库ID',
  `location` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '库位',
  `batch_no` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '批次号',
  `expired_at` datetime DEFAULT NULL COMMENT '过期时间',
  `status` tinyint DEFAULT '1' COMMENT '状态: 0=禁用, 1=正常, 2=锁定',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_product_id` (`product_id`),
  KEY `idx_warehouse_id` (`warehouse_id`),
  KEY `idx_status` (`status`),
  KEY `idx_stock_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存主表';

INSERT INTO `stock` (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES 
(1, 1, 1000, 0, 1000, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 15:26:44', null),
(2, 2, 960, 0, 960, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 23:21:04', null),
(3, 3, 989, 0, 989, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-25 14:03:23', null),
(4, 4, 1000, 0, 1000, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:54:46', null),
(5, 5, 999, 0, 999, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-25 00:25:38', null),
(6, 6, 1000, 0, 1000, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:54:59', null),
(7, 7, 985, 0, 985, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-25 00:23:53', null);

CREATE TABLE IF NOT EXISTS `stock_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '关联订单号',
  `change_type` tinyint NOT NULL COMMENT '变动类型',
  `change_qty` bigint NOT NULL COMMENT '变动数量',
  `before_qty` bigint NOT NULL COMMENT '变动前数量',
  `after_qty` bigint NOT NULL COMMENT '变动后数量',
  `operator_id` int unsigned DEFAULT '0' COMMENT '操作人ID',
  `operator_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作人姓名',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
  `warehouse_id` int unsigned DEFAULT '1' COMMENT '仓库ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_order_no` (`order_no`),
  KEY `idx_change_type` (`change_type`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_stock_log_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存变动日志表';

CREATE TABLE IF NOT EXISTS `stock_alert` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `alert_type` tinyint NOT NULL COMMENT '预警类型',
  `alert_level` tinyint DEFAULT '1' COMMENT '预警级别',
  `threshold` bigint NOT NULL COMMENT '阈值',
  `current_value` bigint NOT NULL COMMENT '当前值',
  `status` tinyint DEFAULT '0' COMMENT '状态',
  `handled_at` datetime DEFAULT NULL COMMENT '处理时间',
  `handler_id` int unsigned DEFAULT '0' COMMENT '处理人ID',
  `handler_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '处理人姓名',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '处理备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_status` (`status`),
  KEY `idx_alert_type` (`alert_type`),
  KEY `idx_stock_alert_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存预警表';

CREATE TABLE IF NOT EXISTS `stock_check` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `check_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '盘点单号',
  `warehouse_id` int unsigned DEFAULT '1' COMMENT '仓库ID',
  `status` tinyint DEFAULT '0' COMMENT '状态',
  `total_items` int DEFAULT '0' COMMENT '总商品数',
  `diff_items` int DEFAULT '0' COMMENT '差异商品数',
  `operator_id` int unsigned DEFAULT '0' COMMENT '操作人ID',
  `operator_name` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作人姓名',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
  `finished_at` datetime DEFAULT NULL COMMENT '完成时间',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_check_no` (`check_no`),
  KEY `idx_status` (`status`),
  KEY `idx_warehouse_id` (`warehouse_id`),
  KEY `idx_stock_check_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存盘点主表';

CREATE TABLE IF NOT EXISTS `stock_check_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `check_id` bigint unsigned NOT NULL COMMENT '盘点主表ID',
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `system_qty` bigint NOT NULL COMMENT '系统库存',
  `actual_qty` bigint NOT NULL COMMENT '实际库存',
  `diff_qty` bigint NOT NULL COMMENT '差异数量',
  `remark` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '备注',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_check_id` (`check_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_stock_check_item_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存盘点明细表';

CREATE TABLE IF NOT EXISTS `stock_dlq` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `message_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '消息ID',
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `quantity` bigint NOT NULL COMMENT '数量',
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '订单号',
  `user_id` int unsigned DEFAULT '0' COMMENT '用户ID',
  `operation` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作类型',
  `retry_count` int DEFAULT '0' COMMENT '重试次数',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT '错误信息',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_id` (`message_id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存死信队列表';

CREATE TABLE IF NOT EXISTS `stock_message_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `message_id` varchar(128) NOT NULL,
  `product_id` int unsigned NOT NULL,
  `operation` varchar(32) NOT NULL,
  `quantity` bigint NOT NULL,
  `status` varchar(32) DEFAULT 'processed',
  `processed_at` datetime NOT NULL,
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `message_id` (`message_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ==================== cart 数据库 ====================

USE `cart`;

CREATE TABLE IF NOT EXISTS `cart` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL,
  `product_id` int unsigned DEFAULT NULL,
  `qty` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ==================== order 数据库 ====================

USE `order`;

CREATE TABLE IF NOT EXISTS `order` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `order_id` varchar(256) DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL,
  `user_currency` longtext,
  `email` longtext,
  `street_address` longtext,
  `city` longtext,
  `state` longtext,
  `country` longtext,
  `zip_code` int DEFAULT NULL,
  `order_state` longtext,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_order_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE IF NOT EXISTS `order_item` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `product_id` int unsigned DEFAULT NULL,
  `order_id_refer` varchar(256) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  `cost` float DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_item_order_id_refer` (`order_id_refer`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ==================== payment 数据库 ====================

USE `payment`;

CREATE TABLE IF NOT EXISTS `payment` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL,
  `order_id` longtext,
  `transaction_id` longtext,
  `amount` float DEFAULT NULL,
  `pay_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ==================== checkout 数据库 ====================

USE `checkout`;

CREATE TABLE IF NOT EXISTS `local_message` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `message_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
  `business_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Business operation ID',
  `message_type` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Message type',
  `target_service` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Target service name',
  `target_method` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Target method name',
  `payload` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Message payload in JSON format',
  `status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT 'pending' COMMENT 'Message status',
  `retry_count` int DEFAULT '0' COMMENT 'Number of retry attempts',
  `max_retry` int DEFAULT '5' COMMENT 'Maximum retry attempts',
  `next_retry_at` datetime DEFAULT NULL COMMENT 'Next scheduled retry time',
  `last_retry_at` datetime DEFAULT NULL COMMENT 'Last retry attempt time',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT 'Last error message',
  `confirm_at` datetime DEFAULT NULL COMMENT 'Confirmation time',
  `priority` int DEFAULT '0' COMMENT 'Message priority',
  `correlation_id` varchar(128) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Correlation ID for tracing',
  `callback_url` varchar(256) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Callback URL',
  `callback_status` varchar(32) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT 'Callback status',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_message_id` (`message_id`),
  KEY `idx_business_id` (`business_id`),
  KEY `idx_status` (`status`),
  KEY `idx_status_next_retry` (`status`,`next_retry_at`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_priority` (`priority`),
  KEY `idx_correlation_id` (`correlation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Local message table for outbox pattern';

CREATE TABLE IF NOT EXISTS `message_retry_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  `message_id` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'Message ID',
  `retry_count` int NOT NULL COMMENT 'Retry attempt number',
  `attempt_at` datetime NOT NULL COMMENT 'Attempt timestamp',
  `success` tinyint(1) NOT NULL COMMENT 'Whether attempt succeeded',
  `error_message` text COLLATE utf8mb4_unicode_ci COMMENT 'Error message if failed',
  `duration` bigint DEFAULT '0' COMMENT 'Attempt duration in milliseconds',
  PRIMARY KEY (`id`),
  KEY `idx_message_id` (`message_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Message retry log for auditing';

-- ==================== email 数据库 ====================

USE `email`;

CREATE TABLE IF NOT EXISTS `email_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `to_email` varchar(255) DEFAULT NULL,
  `subject` varchar(500) DEFAULT NULL,
  `status` varchar(50) DEFAULT 'pending',
  `sent_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
