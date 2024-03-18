CREATE TABLE IF NOT EXISTS permissions (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	uuid CHAR(36) UNIQUE,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP NULL DEFAULT NULL
);