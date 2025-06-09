ALTER TABLE `users` ADD `nickname` text;

CREATE UNIQUE INDEX `idx_users_username` ON `users`(`username`);
