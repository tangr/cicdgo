-- Create DATABASE cicd CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci;
-- ./gf gen dao --gf.gcfg.file=config/config-dev.toml

CREATE TABLE IF NOT EXISTS `cicd_pipeline` (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `pipeline_name` varchar(255) NOT NULL,
  `group_id` int(11) NOT NULL,
  `agent_id` int(11) NOT NULL,
  `concurrency` int(11) DEFAULT 1,
  `body` JSON NOT NULL,
  -- `realpipelineid` int(11) DEFAULT NULL,
  -- `version` int(11) DEFAULT NULL,
  `author` varchar(255) NOT NULL,
  `updated_at` bigint(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `pipeline_name` (`pipeline_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cicd_job` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pipeline_id` int(11) NOT NULL,
  `agent_id` int(11) NOT NULL,
  `concurrency` int(11) DEFAULT 1,
  `job_type` varchar(255) NOT NULL,
  `job_status` varchar(255) NOT NULL,
  `script` JSON NOT NULL,
  `comment` varchar(255) NOT NULL,
  `author` varchar(255) NOT NULL,
  `created_at` bigint(10) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cicd_log` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pipeline_id` int(11) NOT NULL,
  `agent_id` int(11) NOT NULL,
  `job_type` varchar(255) NOT NULL,
  `job_id` int(11) NOT NULL,
  `task_status` varchar(255) NOT NULL,
  `ipaddr` varchar(255) NOT NULL,
  `updated_at` bigint(10) NOT NULL,
  `output` longtext DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `log_id` (`job_id`, `ipaddr`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cicd_package` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `pipeline_id` int(11) NOT NULL,
  `job_id` int(11) NOT NULL,
  `job_status` varchar(255) NOT NULL,
  `package_name` varchar(255) NOT NULL,
  `created_at` bigint(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `job_id` (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=10000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cicd_script` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `script_name` varchar(255) NOT NULL,
  `script_body` longtext NOT NULL,
  `author` varchar(255) NOT NULL,
  `updated_at` bigint(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `script_name` (`script_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cicd_agent` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `agent_name` varchar(255) NOT NULL,
  `ipaddr` varchar(255) NOT NULL,
  `updated_at` bigint(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `agent_name` (`agent_name`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `cicd_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_name` varchar(255) NOT NULL,
  `parent_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `cicd_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `group_id` JSON NOT NULL,
  `updated_at` bigint(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
