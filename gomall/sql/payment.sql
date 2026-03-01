CREATE DATABASE `payment` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */

CREATE TABLE `payment` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `user_id` int unsigned DEFAULT NULL,
  `order_id` longtext,
  `transaction_id` longtext,
  `amount` float DEFAULT NULL,
  `pay_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (1, '2026-02-09 14:16:09.650', '2026-02-09 14:16:09.650', 2, 'd366ce87-057e-11f1-b8a3-0250508de514', '37b3537a-e9c9-466c-ac8d-4e1a8854185a', 6.6, '2026-02-09 14:16:09.648');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (2, '2026-02-23 16:23:23.553', '2026-02-23 16:23:23.553', 2, 'eb4d96eb-1090-11f1-bd28-0250508de514', '0af73f3e-deb6-4b2e-bda2-1fdb23b16936', 114524, '2026-02-23 16:23:23.550');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (3, '2026-02-24 19:09:21.013', '2026-02-24 19:09:21.013', 2, '44cce522-1171-11f1-8a9f-0250508de514', '7a5a5870-3352-4006-856e-b80bdacd03c3', 506, '2026-02-24 19:09:21.010');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (4, '2026-02-24 19:11:37.245', '2026-02-24 19:11:37.245', 2, '9610390a-1171-11f1-8a9f-0250508de514', '8dbcc6c1-8a6f-4f74-8245-9ca934b2a0e1', 11, '2026-02-24 19:11:37.245');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (5, '2026-02-24 19:15:11.378', '2026-02-24 19:15:11.378', 2, '15b2f73b-1172-11f1-8a9f-0250508de514', '7bc1ddcd-bee5-4ad9-8ce4-86ab24ebf4be', 1.1, '2026-02-24 19:15:11.378');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (6, '2026-02-24 19:25:44.586', '2026-02-24 19:25:44.586', 2, '8f1be37f-1173-11f1-8dd0-0250508de514', '72efc5c9-c356-4ad3-a86f-ac3916a875e7', 1.1, '2026-02-24 19:25:44.586');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (7, '2026-02-24 19:44:16.505', '2026-02-24 19:44:16.505', 2, '25d71901-1176-11f1-8dd0-0250508de514', '5e0b0bde-e22d-40e2-894b-d277616132e6', 9.9, '2026-02-24 19:44:16.505');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (8, '2026-02-24 19:55:01.523', '2026-02-24 19:55:01.523', 2, 'a65350af-1177-11f1-be19-0250508de514', '97eb773f-c6b4-4897-a768-61310c974692', 198, '2026-02-24 19:55:01.523');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (9, '2026-02-24 20:13:51.156', '2026-02-24 20:13:51.156', 2, '4799fd26-117a-11f1-b5a1-0250508de514', '6aa2040e-1967-40e4-afe2-01dc9bee2aee', 257.4, '2026-02-24 20:13:51.155');
INSERT INTO payment.payment (id, created_at, updated_at, user_id, order_id, transaction_id, amount, pay_at) VALUES (10, '2026-02-24 20:29:51.004', '2026-02-24 20:29:51.004', 2, '83bef4dd-117c-11f1-b93c-0250508de514', '73039fe0-5e91-4efe-a72d-79f7123d603f', 132, '2026-02-24 20:29:51.004');
