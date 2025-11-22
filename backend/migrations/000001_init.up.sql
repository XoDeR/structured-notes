CREATE TABLE IF NOT EXISTS `users` (
  `id` BIGINT(20) UNSIGNED NOT NULL,
  `username` varchar(25) NOT NULL,
  `firstname` varchar(25) DEFAULT NULL,
  `lastname` varchar(25) DEFAULT NULL,
  `role` int NOT NULL DEFAULT '1',
  `avatar` varchar(75) DEFAULT NULL,
  `email` varchar(50) NOT NULL,
  `password` varchar(255) NOT NULL,
  `password_reset_token` varchar(255) DEFAULT NULL,
  `created_timestamp` bigint NOT NULL,
  `updated_timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `nodes` (
    `id` BIGINT UNSIGNED PRIMARY KEY,
    `user_id` BIGINT UNSIGNED NOT NULL,
    `parent_id` BIGINT UNSIGNED NULL,
    `name` VARCHAR(50) NOT NULL,
    `description` VARCHAR(255) NULL,
    `tags` VARCHAR(200) NULL,
    `role` TINYINT NOT NULL COMMENT '1=workspace, 2=category, 3=document, 4=media',
    `color` INT NULL,
    `icon` TEXT NULL,
    `thumbnail` TEXT NULL,
    `theme` VARCHAR(30) NULL,
    `accessibility` tinyint(1) NOT NULL DEFAULT 1 COMMENT '1=private, 2=public',
    `access` INT NOT NULL DEFAULT 0,
    `display` TINYINT NULL,
    `order` INT NULL,
    `content` LONGTEXT NULL,
    `content_compiled` LONGTEXT NULL,
    `size` INT NULL,
    `metadata` JSON NULL,
    `created_timestamp` BIGINT NOT NULL,
    `updated_timestamp` BIGINT NOT NULL,
    CONSTRAINT `nodes_parent_fk` FOREIGN KEY (`parent_id`) REFERENCES `nodes` (`id`) ON DELETE CASCADE,
    CONSTRAINT `nodes_users_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `chk_parent_not_self` CHECK (`parent_id` IS NULL OR `parent_id` != `id`)
);

CREATE TABLE IF NOT EXISTS `sessions` (
  `id` BIGINT(20) UNSIGNED NOT NULL,
  `user_id` BIGINT(20) UNSIGNED NOT NULL,
  `refresh_token` varchar(255) DEFAULT NULL,
  `expire_token` bigint DEFAULT NULL,
  `last_refresh_timestamp` bigint DEFAULT NULL,
  `active` int DEFAULT NULL,
  `login_timestamp` bigint DEFAULT NULL,
  `logout_timestamp` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `users_auth_users_id_fk` (`user_id`),
  CONSTRAINT `users_auth_users_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `connections_logs` (
  `id` BIGINT(20) UNSIGNED NOT NULL,
  `user_id` BIGINT(20) UNSIGNED NOT NULL,
  `ip_adress` varchar(50) DEFAULT NULL,
  `user_agent` varchar(200) DEFAULT NULL,
  `location` varchar(100) DEFAULT NULL,
  `type` varchar(10) NOT NULL,
  `timestamp` bigint NOT NULL,
  PRIMARY KEY (`id`),
  KEY `connections_logs_users_id_fk` (`user_id`),
  CONSTRAINT `connections_logs_users_id_fk` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS `permissions` (
    `id` BIGINT UNSIGNED PRIMARY KEY,
    node_id BIGINT UNSIGNED NOT NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    permission TINYINT NOT NULL DEFAULT 0,
    `created_timestamp` BIGINT NOT NULL,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_nodes_user_id ON nodes(user_id);
CREATE INDEX idx_nodes_parent_id ON nodes(parent_id);
CREATE INDEX idx_nodes_accessibility ON nodes(accessibility);
CREATE INDEX idx_nodes_created ON nodes(created_timestamp);
CREATE INDEX idx_nodes_updated ON nodes(updated_timestamp);

CREATE INDEX idx_permissions_node_id ON permissions(node_id);
CREATE INDEX idx_permissions_user_id ON permissions(user_id);
CREATE INDEX idx_permissions_node_user ON permissions(node_id, user_id);