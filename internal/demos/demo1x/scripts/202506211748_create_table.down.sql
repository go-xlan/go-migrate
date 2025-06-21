-- reverse -- CREATE TABLE `infos` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`name` varchar(191),`cate` tinyint,PRIMARY KEY (`id`),INDEX `idx_infos_deleted_at` (`deleted_at`),CONSTRAINT `uni_infos_name` UNIQUE (`name`));
SELECT 1 / 0; -- DROP TABLE; -- TODO

-- reverse -- CREATE TABLE `users` (`id` bigint unsigned AUTO_INCREMENT,`created_at` datetime(3) NULL,`updated_at` datetime(3) NULL,`deleted_at` datetime(3) NULL,`username` varchar(191),`nickname` longtext,`rank` bigint unsigned,`score` longtext,PRIMARY KEY (`id`),INDEX `idx_users_deleted_at` (`deleted_at`),CONSTRAINT `uni_users_username` UNIQUE (`username`));
SELECT 1 / 0; -- DROP TABLE; -- TODO
