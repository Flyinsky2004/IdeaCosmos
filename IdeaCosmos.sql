/*
 Navicat Premium Dump SQL

 Source Server         : mysql8
 Source Server Type    : MySQL
 Source Server Version : 80040 (8.0.40)
 Source Host           : localhost:3306
 Source Schema         : IdeaCosmos

 Target Server Type    : MySQL
 Target Server Version : 80040 (8.0.40)
 File Encoding         : 65001

 Date: 25/12/2024 14:35:10
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `username` varchar(50) DEFAULT NULL,
  `password` longtext,
  `email` longtext,
  `tokens` bigint DEFAULT NULL,
  `permission` tinyint unsigned DEFAULT NULL,
  `group` tinyint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` (`id`, `created_at`, `updated_at`, `deleted_at`, `username`, `password`, `email`, `tokens`, `permission`, `group`) VALUES (2, '2024-12-24 15:14:56.196', '2024-12-24 15:14:56.196', NULL, 'flyinsky', '$2a$10$2juZ9iQcHNP3Ft98F9K5y.MELqutQML81IhJx46wzbDk.SbBl/HsK', 'w2084151024@gmail.com', 0, 0, 0);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
