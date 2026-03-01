CREATE DATABASE `order` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */

CREATE TABLE `order` (
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
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (1, '2026-02-09 14:14:07.027', '2026-02-09 14:14:07.027', '8a55de5b-057e-11f1-b8a3-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'placed');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (2, '2026-02-09 14:16:09.594', '2026-02-09 14:16:09.793', 'd366ce87-057e-11f1-b8a3-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (3, '2026-02-23 16:23:23.438', '2026-02-23 16:23:23.587', 'eb4d96eb-1090-11f1-bd28-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (4, '2026-02-24 19:09:20.861', '2026-02-24 19:09:21.078', '44cce522-1171-11f1-8a9f-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (5, '2026-02-24 19:11:37.182', '2026-02-24 19:11:37.267', '9610390a-1171-11f1-8a9f-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (6, '2026-02-24 19:15:11.319', '2026-02-24 19:15:11.402', '15b2f73b-1172-11f1-8a9f-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (7, '2026-02-24 19:25:44.518', '2026-02-24 19:25:44.618', '8f1be37f-1173-11f1-8dd0-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (8, '2026-02-24 19:44:16.386', '2026-02-24 19:44:16.553', '25d71901-1176-11f1-8dd0-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (9, '2026-02-24 19:55:01.462', '2026-02-24 19:55:01.547', 'a65350af-1177-11f1-be19-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');
INSERT INTO `order`.`order` (id, created_at, updated_at, order_id, user_id, user_currency, email, street_address, city, state, country, zip_code, order_state) VALUES (10, '2026-02-24 20:13:51.024', '2026-02-24 20:13:51.296', '4799fd26-117a-11f1-b5a1-0250508de514', 2, 'USD', 'abc@example.com', '7th street', 'hangzhou', 'zhejiang', 'china', 0, 'paid');


CREATE TABLE `order_item` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `product_id` int unsigned DEFAULT NULL,
  `order_id_refer` varchar(256) DEFAULT NULL,
  `quantity` int DEFAULT NULL,
  `cost` float DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_order_item_order_id_refer` (`order_id_refer`),
  CONSTRAINT `fk_order_order_items` FOREIGN KEY (`order_id_refer`) REFERENCES `order` (`order_id`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (1, '2026-02-09 14:14:07.104', '2026-02-09 14:14:07.104', 3, '8a55de5b-057e-11f1-b8a3-0250508de514', 1, 6.6);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (2, '2026-02-09 14:14:07.104', '2026-02-09 14:14:07.104', 7, '8a55de5b-057e-11f1-b8a3-0250508de514', 1, 4.8);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (3, '2026-02-09 14:16:09.602', '2026-02-09 14:16:09.602', 3, 'd366ce87-057e-11f1-b8a3-0250508de514', 1, 6.6);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (4, '2026-02-23 16:23:23.449', '2026-02-23 16:23:23.449', 1, 'eb4d96eb-1090-11f1-bd28-0250508de514', 1, 9.9);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (5, '2026-02-23 16:23:23.449', '2026-02-23 16:23:23.449', 12, 'eb4d96eb-1090-11f1-bd28-0250508de514', 1, 114514);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (6, '2026-02-24 19:09:20.900', '2026-02-24 19:09:20.900', 5, '44cce522-1171-11f1-8a9f-0250508de514', 10, 11);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (7, '2026-02-24 19:09:20.900', '2026-02-24 19:09:20.900', 1, '44cce522-1171-11f1-8a9f-0250508de514', 50, 495);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (8, '2026-02-24 19:11:37.195', '2026-02-24 19:11:37.195', 5, '9610390a-1171-11f1-8a9f-0250508de514', 10, 11);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (9, '2026-02-24 19:15:11.328', '2026-02-24 19:15:11.328', 5, '15b2f73b-1172-11f1-8a9f-0250508de514', 1, 1.1);
INSERT INTO `order`.order_item (id, created_at, updated_at, product_id, order_id_refer, quantity, cost) VALUES (10, '2026-02-24 19:25:44.533', '2026-02-24 19:25:44.533', 5, '8f1be37f-1173-11f1-8dd0-0250508de514', 1, 1.1);
