package handler

import "database/sql"

// EnsurePlusResourceSchema creates the application-scoped resource tables.
// Every resource is tied to apps.id so later client delivery cannot leak data
// between applications. Pricing fields are included from day one for paid plugins.
func EnsurePlusResourceSchema(db *sql.DB) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS app_plugin_sources (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			app_id BIGINT UNSIGNED NOT NULL,
			name VARCHAR(100) NOT NULL,
			url VARCHAR(500) NOT NULL,
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			sort INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_app_plugin_source_url (app_id, url(191)),
			KEY idx_app_plugin_source (app_id, enabled, sort)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='按应用隔离的插件源'`,
		`CREATE TABLE IF NOT EXISTS app_plugins (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			app_id BIGINT UNSIGNED NOT NULL,
			plugin_key VARCHAR(100) NOT NULL,
			name VARCHAR(150) NOT NULL,
			description VARCHAR(1000) NOT NULL DEFAULT '',
			category VARCHAR(50) NOT NULL DEFAULT 'other',
			version VARCHAR(40) NOT NULL DEFAULT '0.0.0',
			package_url VARCHAR(1000) NOT NULL DEFAULT '',
			sha256 CHAR(64) NOT NULL DEFAULT '',
			pricing_type ENUM('free','one_time','subscription','plan','trial') NOT NULL DEFAULT 'free',
			price DECIMAL(12,2) NOT NULL DEFAULT 0,
			currency CHAR(3) NOT NULL DEFAULT 'CNY',
			trial_days INT UNSIGNED NOT NULL DEFAULT 0,
			requires_license TINYINT(1) NOT NULL DEFAULT 0,
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			published TINYINT(1) NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_app_plugin_key (app_id, plugin_key),
			KEY idx_app_plugin_published (app_id, published, enabled)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='按应用隔离的插件目录和商业属性'`,
		`CREATE TABLE IF NOT EXISTS app_home_templates (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			app_id BIGINT UNSIGNED NOT NULL,
			template_key VARCHAR(100) NOT NULL,
			name VARCHAR(150) NOT NULL,
			description VARCHAR(1000) NOT NULL DEFAULT '',
			version VARCHAR(40) NOT NULL DEFAULT '0.0.0',
			schema_version INT NOT NULL DEFAULT 1,
			content_path VARCHAR(1000) NOT NULL DEFAULT '',
			preview_path VARCHAR(1000) NOT NULL DEFAULT '',
			sha256 CHAR(64) NOT NULL DEFAULT '',
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			published TINYINT(1) NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_app_template_key (app_id, template_key),
			KEY idx_app_template_published (app_id, published, enabled)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='按应用隔离的首页模板'`,
		`CREATE TABLE IF NOT EXISTS app_advertisements (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			app_id BIGINT UNSIGNED NOT NULL,
			position VARCHAR(40) NOT NULL,
			title VARCHAR(200) NOT NULL,
			image_url VARCHAR(1000) NOT NULL DEFAULT '',
			destination_url VARCHAR(1000) NOT NULL DEFAULT '',
			description VARCHAR(1000) NOT NULL DEFAULT '',
			weight INT NOT NULL DEFAULT 0,
			start_at DATETIME DEFAULT NULL,
			end_at DATETIME DEFAULT NULL,
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			KEY idx_app_ad_position_time (app_id, position, enabled, start_at, end_at)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='按应用隔离的广告投放'`,
		`CREATE TABLE IF NOT EXISTS plugin_entitlements (
			id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
			app_id BIGINT UNSIGNED NOT NULL,
			plugin_id BIGINT UNSIGNED NOT NULL,
			user_id BIGINT UNSIGNED DEFAULT NULL,
			agent_id BIGINT UNSIGNED DEFAULT NULL,
			license_id BIGINT UNSIGNED DEFAULT NULL,
			order_id BIGINT UNSIGNED DEFAULT NULL,
			status ENUM('trialing','active','expired','revoked') NOT NULL DEFAULT 'active',
			granted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME DEFAULT NULL,
			max_devices INT UNSIGNED NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			KEY idx_entitlement_lookup (app_id, plugin_id, user_id, status),
			KEY idx_entitlement_agent (app_id, plugin_id, agent_id, status)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='按应用隔离的付费插件权益'`,
	}
	for _, statement := range statements {
		if _, err := db.Exec(statement); err != nil {
			return err
		}
	}
	return nil
}
