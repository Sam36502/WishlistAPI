-- DATABASE STRUCTURE SCRIPT

DROP DATABASE IF EXISTS `wishlist`;
CREATE DATABASE `wishlist` CHARACTER SET utf8;

GRANT INSERT, SELECT, UPDATE, DELETE ON `wishlist`.* TO `wishlist_user`;
USE `wishlist`;

DROP TABLE IF EXISTS `tbl_user`;
CREATE TABLE `tbl_user` (
  `id_user` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL,
  `password` varchar(128) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id_user`),
  UNIQUE KEY `id_user` (`id_user`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tbl_status`;
CREATE TABLE `tbl_status` (
  `id_status` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `desc` text,
  PRIMARY KEY (`id_status`),
  UNIQUE KEY `id_status` (`id_status`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

INSERT INTO `tbl_status` VALUES (1,'Available','No one is planning to get this, yet.'),(2,'Reserved','Someone is planning to get this already.'),(3,'Received','This was on the wishlist, but has since been received.');

DROP TABLE IF EXISTS `tbl_item`;
CREATE TABLE `tbl_item` (
  `id_item` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(64) NOT NULL,
  `desc` text NOT NULL,
  `status_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `price` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_item`),
  UNIQUE KEY `id_item` (`id_item`),
  KEY `status_id` (`status_id`),
  KEY `user_id` (`user_id`),
  KEY `reserved_by_user_id` (`reserved_by_user_id`),
  CONSTRAINT `tbl_item_ibfk_1` FOREIGN KEY (`status_id`) REFERENCES `tbl_status` (`id_status`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `tbl_item_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `tbl_user` (`id_user`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `tbl_item_ibfk_3` FOREIGN KEY (`reserved_by_user_id`) REFERENCES `tbl_user` (`id_user`) ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `tbl_link`;
CREATE TABLE `tbl_link` (
  `id_link` int(11) NOT NULL AUTO_INCREMENT,
  `text` varchar(128) NOT NULL,
  `hyperlink` varchar(256) NOT NULL,
  `item_id` int(11) NOT NULL,
  PRIMARY KEY (`id_link`),
  UNIQUE KEY `id_link` (`id_link`),
  KEY `item_id` (`item_id`),
  CONSTRAINT `tbl_link_ibfk_1` FOREIGN KEY (`item_id`) REFERENCES `tbl_item` (`id_item`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8;
