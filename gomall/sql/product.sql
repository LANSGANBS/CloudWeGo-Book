CREATE DATABASE `product` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */

CREATE TABLE `cart` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL,
  `product_id` int unsigned DEFAULT NULL,
  `qty` int unsigned DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

CREATE TABLE `category` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `name` longtext,
  `description` longtext,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_category_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO product.category (id, created_at, updated_at, name, description, deleted_at) VALUES (1, '2023-12-06 15:05:06.000', '2023-12-06 15:05:06.000', 'T-Shirt', 'T-Shirt', null);
INSERT INTO product.category (id, created_at, updated_at, name, description, deleted_at) VALUES (2, '2023-12-06 15:05:06.000', '2023-12-06 15:05:06.000', 'Sticker', 'Sticker', null);
INSERT INTO product.category (id, created_at, updated_at, name, description, deleted_at) VALUES (3, null, null, 'test', 'test', '2026-02-10 17:07:20.578');
INSERT INTO product.category (id, created_at, updated_at, name, description, deleted_at) VALUES (4, null, null, 'test', 'test', '2026-02-23 13:51:30.163');
INSERT INTO product.category (id, created_at, updated_at, name, description, deleted_at) VALUES (5, null, null, 'Test', 'Test', null);

CREATE TABLE `product` (
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
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (1, '2024-12-06 15:26:19.000', '2026-02-24 23:26:44.601', 'Notebook', 'The cloudwego notebook is a highly efficient and feature-rich notebook designed to meet all your note-taking needs. ', '/static/image/notebook.jpeg', 9.9, 0, null, 20);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (2, '2023-12-06 15:26:19.000', '2026-02-24 22:54:32.827', 'Mouse-Pad', 'The cloudwego mouse pad is a premium-grade accessory designed to enhance your computer usage experience. ', '/static/image/mouse-pad.jpeg', 8.8, 0, null, 1021);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (3, '2023-12-06 15:26:19.000', '2026-02-24 22:54:39.751', 'T-Shirt', 'The cloudwego t-shirt is a stylish and comfortable clothing item that allows you to showcase your fashion sense while enjoying maximum comfort.', '/static/image/t-shirt.jpeg', 6.6, 0, null, 1190);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (4, '2023-12-06 15:26:19.000', '2026-02-24 22:54:46.926', 'T-Shirt', 'The cloudwego t-shirt is a stylish and comfortable clothing item that allows you to showcase your fashion sense while enjoying maximum comfort.', '/static/image/t-shirt-1.jpeg', 2.2, 0, null, 2);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (5, '2025-12-06 15:26:19.000', '2026-02-24 22:54:53.220', 'Sweatshirt', 'The cloudwego Sweatshirt is a cozy and fashionable garment that provides warmth and style during colder weather.', '/static/image/sweatshirt.jpeg', 1.1, 0, null, 5);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (6, '2023-12-06 15:26:19.000', '2026-02-24 22:55:00.233', 'T-Shirt', 'The cloudwego t-shirt is a stylish and comfortable clothing item that allows you to showcase your fashion sense while enjoying maximum comfort.', '/static/image/t-shirt-2.jpeg', 1.8, 0, null, 100);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (7, '2024-12-06 15:26:19.000', '2026-02-24 22:55:10.390', 'mascot', 'The cloudwego mascot is a charming and captivating representation of the brand, designed to bring joy and a playful spirit to any environment.', '/static/image/logo.jpg', 4.8, 0, null, 15);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (8, null, null, 'test', 'test', 'https://www.lansganbs.cn/images/friends/lansganbs.png', 0.01, 0, '2026-02-10 16:43:26.982', 0);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (9, null, null, 'test', 'test', 'https://www.lansganbs.cn/images/friends/lansganbs.png', 0.02, 0, '2026-02-10 17:05:43.441', 0);
INSERT INTO product.product (id, created_at, updated_at, name, description, picture, price, stock, deleted_at, sales) VALUES (10, null, null, 'test', 'test', 'https://www.lansganbs.cn/images/friends/lansganbs.png', 0.01, 0, '2026-02-10 17:23:05.018', 0);

CREATE TABLE `product_category` (
  `category_id` bigint NOT NULL,
  `product_id` bigint NOT NULL,
  PRIMARY KEY (`category_id`,`product_id`),
  KEY `fk_product_category_product` (`product_id`),
  CONSTRAINT `fk_product_category_category` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`),
  CONSTRAINT `fk_product_category_product` FOREIGN KEY (`product_id`) REFERENCES `product` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO product.product_category (category_id, product_id) VALUES (2, 1);
INSERT INTO product.product_category (category_id, product_id) VALUES (2, 2);
INSERT INTO product.product_category (category_id, product_id) VALUES (1, 3);
INSERT INTO product.product_category (category_id, product_id) VALUES (1, 4);
INSERT INTO product.product_category (category_id, product_id) VALUES (1, 5);
INSERT INTO product.product_category (category_id, product_id) VALUES (1, 6);
INSERT INTO product.product_category (category_id, product_id) VALUES (2, 7);
INSERT INTO product.product_category (category_id, product_id) VALUES (1, 8);
INSERT INTO product.product_category (category_id, product_id) VALUES (1, 9);
INSERT INTO product.product_category (category_id, product_id) VALUES (4, 10);

CREATE TABLE `stock` (
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
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存主表'

INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (1, 1, 1000, 0, 1000, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 15:26:44', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (2, 2, 960, 0, 960, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 23:21:04', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (3, 3, 989, 0, 989, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-25 14:03:23', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (4, 4, 1000, 0, 1000, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:54:46', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (5, 5, 999, 0, 999, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-25 00:25:38', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (6, 6, 1000, 0, 1000, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:54:59', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (7, 7, 985, 0, 985, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-25 00:23:53', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (8, 8, 0, 0, 0, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:26:18', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (9, 9, 0, 0, 0, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:26:18', null);
INSERT INTO product.stock (id, product_id, quantity, reserved, available, min_stock, max_stock, safety_stock, unit, warehouse_id, location, batch_no, expired_at, status, created_at, updated_at, deleted_at) VALUES (10, 10, 0, 0, 0, 10, 1000, 20, '件', 1, '', '', null, 1, '2026-02-24 14:26:18', '2026-02-24 14:26:18', null);

CREATE TABLE `stock_alert` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `alert_type` tinyint NOT NULL COMMENT '预警类型: 1=低库存, 2=超储, 3=即将过期, 4=已过期',
  `alert_level` tinyint DEFAULT '1' COMMENT '预警级别: 1=提示, 2=警告, 3=严重',
  `threshold` bigint NOT NULL COMMENT '阈值',
  `current_value` bigint NOT NULL COMMENT '当前值',
  `status` tinyint DEFAULT '0' COMMENT '状态: 0=待处理, 1=已处理, 2=已忽略',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存预警表'

CREATE TABLE `stock_check` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `check_no` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '盘点单号',
  `warehouse_id` int unsigned DEFAULT '1' COMMENT '仓库ID',
  `status` tinyint DEFAULT '0' COMMENT '状态: 0=待盘点, 1=盘点中, 2=已完成, 3=已取消',
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存盘点主表'

CREATE TABLE `stock_check_item` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存盘点明细表'

CREATE TABLE `stock_dlq` (
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
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存死信队列表'

CREATE TABLE `stock_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `product_id` int unsigned NOT NULL COMMENT '商品ID',
  `order_no` varchar(64) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '关联订单号',
  `change_type` tinyint NOT NULL COMMENT '变动类型: 1=采购入库, 2=销售出库, 3=退货入库, 4=调整, 5=盘点, 6=预留, 7=释放, 8=损耗, 9=调拨',
  `change_qty` bigint NOT NULL COMMENT '变动数量(正数增加,负数减少)',
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
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='库存变动日志表'

INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (1, 2, '', 2, -6, 1000, 994, 0, 'system', '', 1, '2026-02-24 23:11:33', '2026-02-24 23:11:33', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (2, 2, '', 2, -6, 994, 988, 0, 'system', '', 1, '2026-02-24 23:11:38', '2026-02-24 23:11:38', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (3, 2, '', 2, -6, 988, 982, 0, 'system', '', 1, '2026-02-24 23:11:51', '2026-02-24 23:11:51', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (4, 2, '', 2, -6, 982, 976, 0, 'system', '', 1, '2026-02-24 23:20:43', '2026-02-24 23:20:43', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (5, 2, '', 2, -16, 976, 960, 0, 'system', '', 1, '2026-02-24 23:21:04', '2026-02-24 23:21:04', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (6, 3, '', 2, -6, 1000, 994, 0, 'system', 'RocketMQ deduct, msgId=stock_3_1771949928462965800', 1, '2026-02-25 00:18:50', '2026-02-25 00:18:50', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (7, 7, '', 2, -15, 1000, 985, 0, 'system', 'RocketMQ deduct, msgId=stock_7_1771950202381031600', 1, '2026-02-25 00:23:53', '2026-02-25 00:23:53', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (8, 5, '', 2, -1, 1000, 999, 0, 'system', 'RocketMQ deduct, msgId=stock_5_1771950315622131700', 1, '2026-02-25 00:25:38', '2026-02-25 00:25:38', null);
INSERT INTO product.stock_log (id, product_id, order_no, change_type, change_qty, before_qty, after_qty, operator_id, operator_name, remark, warehouse_id, created_at, updated_at, deleted_at) VALUES (9, 3, '', 2, -5, 994, 989, 0, 'system', 'RocketMQ deduct, msgId=stock_3_1771999402650342200', 1, '2026-02-25 14:03:23', '2026-02-25 14:03:23', null);

CREATE TABLE `stock_message_log` (
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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO product.stock_message_log (id, message_id, product_id, operation, quantity, status, processed_at, created_at) VALUES (1, 'stock_3_1771999402650342200', 3, 'deduct', 5, 'processed', '2026-02-25 14:03:23', '2026-02-25 14:03:23');
