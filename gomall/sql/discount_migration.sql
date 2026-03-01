-- 添加折扣相关字段到 product 表
ALTER TABLE `product` 
ADD COLUMN `discount_type` TINYINT DEFAULT 0 COMMENT '折扣类型: 0=无折扣, 1=折扣率(几折), 2=固定降价(减多少元)' AFTER `price`,
ADD COLUMN `discount_value` DECIMAL(10,2) DEFAULT 0 COMMENT '折扣值: 折扣率如0.8表示8折, 固定降价如10表示减10元' AFTER `discount_type`,
ADD COLUMN `discount_start_time` DATETIME DEFAULT NULL COMMENT '限时特惠开始时间, NULL表示普通折扣' AFTER `discount_value`,
ADD COLUMN `discount_end_time` DATETIME DEFAULT NULL COMMENT '限时特惠结束时间, NULL表示无限时长' AFTER `discount_start_time`,
ADD COLUMN `original_price` DECIMAL(10,2) DEFAULT NULL COMMENT '原价(折扣前的价格)' AFTER `discount_end_time`;

-- 添加索引以支持折扣筛选
ALTER TABLE `product` 
ADD INDEX `idx_discount_type` (`discount_type`),
ADD INDEX `idx_discount_time` (`discount_start_time`, `discount_end_time`);

-- 创建价格变动历史记录表
CREATE TABLE IF NOT EXISTS `product_price_history` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  `product_id` INT UNSIGNED NOT NULL COMMENT '商品ID',
  `change_type` TINYINT NOT NULL COMMENT '变动类型: 1=普通价格调整, 2=设置折扣, 3=设置限时特惠, 4=取消折扣, 5=限时特惠到期恢复',
  `old_price` DECIMAL(10,2) DEFAULT NULL COMMENT '变动前价格',
  `new_price` DECIMAL(10,2) DEFAULT NULL COMMENT '变动后价格',
  `old_discount_type` TINYINT DEFAULT NULL COMMENT '变动前折扣类型',
  `new_discount_type` TINYINT DEFAULT NULL COMMENT '变动后折扣类型',
  `old_discount_value` DECIMAL(10,2) DEFAULT NULL COMMENT '变动前折扣值',
  `new_discount_value` DECIMAL(10,2) DEFAULT NULL COMMENT '变动后折扣值',
  `discount_start_time` DATETIME DEFAULT NULL COMMENT '限时特惠开始时间',
  `discount_end_time` DATETIME DEFAULT NULL COMMENT '限时特惠结束时间',
  `operator_id` INT UNSIGNED DEFAULT 0 COMMENT '操作人ID',
  `operator_name` VARCHAR(50) DEFAULT '' COMMENT '操作人姓名',
  `remark` VARCHAR(255) DEFAULT '' COMMENT '备注',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_product_id` (`product_id`),
  KEY `idx_change_type` (`change_type`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='商品价格变动历史表';
