SET NAMES 'utf8';
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

USE shop;

DROP TABLE IF EXISTS products;
CREATE TABLE products (
    `id` int NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `category` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `login` varchar(255) NOT NULL,
    `password` varchar(255) NOT NULL,
    `confirm` TINYINT NOT NULL,
    PRIMARY KEY (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `confirmations`;
CREATE TABLE `confirmations`(
    `login` varchar(255) NOT NULL,
    `token` varchar(255) NOT NULL,
    `expire` BIGINT NOT NULL,
    PRIMARY KEY (`token`),
    INDEX (`login`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `sessions`;
CREATE TABLE `sessions` (
    `token` varchar(255) NOT NULL,
    `login` varchar(255) NOT NULL,
    `expire` varchar(255) NOT NULL,
    PRIMARY KEY (`login`),
    INDEX (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;