CREATE DATABASE `user` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci */ /*!80016 DEFAULT ENCRYPTION='N' */

CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `email` varchar(191) DEFAULT NULL,
  `password_hashed` longtext,
  `is_admin` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci

INSERT INTO user.user (id, created_at, updated_at, email, password_hashed, is_admin) VALUES (1, '2023-12-26 09:46:19.852', '2023-12-26 09:46:19.852', '123@admin.com', '$2a$10$jTvUFh7Z8Kw0hLV8WrAws.PRQTeuH4gopJ7ZMoiFvwhhz5Vw.bj7C', 0);
INSERT INTO user.user (id, created_at, updated_at, email, password_hashed, is_admin) VALUES (2, '2026-02-09 14:05:55.001', '2026-02-09 14:05:55.001', 'lansganbs@qq.com', '$2a$10$gbXnOKhwO3c/qz.uaD9w6.Del1zBbjTbYwO6Y0mpSmKVs4FOGHsDS', 1);
