CREATE TABLE `tb2`
(
    `id`         bigint unsigned AUTO_INCREMENT,
    `created_at` datetime(3) NULL,
    `updated_at` datetime(3) NULL,
    `a`          varchar(255) COMMENT 'A',
    `b`          varchar(255) COMMENT 'B',
    `c`          varchar(255) COMMENT 'C',
    `d`          varchar(255) COMMENT 'D',
    PRIMARY KEY (`id`)
);
