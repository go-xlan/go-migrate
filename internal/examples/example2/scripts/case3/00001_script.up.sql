CREATE TABLE `users` (`id` integer PRIMARY KEY AUTOINCREMENT,`created_at` datetime,`updated_at` datetime,`deleted_at` datetime,`username` text);

CREATE INDEX `idx_users_deleted_at` ON `users`(`deleted_at`);
