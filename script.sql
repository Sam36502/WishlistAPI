-- MySQL dump 10.15  Distrib 10.0.28-MariaDB, for debian-linux-gnueabihf (armv7l)
--
-- Host: wishlist    Database: wishlist
-- ------------------------------------------------------
-- Server version	10.0.28-MariaDB-2+b1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `tbl_item`
--

DROP TABLE IF EXISTS `tbl_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
  CONSTRAINT `tbl_item_ibfk_1` FOREIGN KEY (`status_id`) REFERENCES `tbl_status` (`id_status`) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT `tbl_item_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `tbl_user` (`id_user`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_item`
--

LOCK TABLES `tbl_item` WRITE;
/*!40000 ALTER TABLE `tbl_item` DISABLE KEYS */;
INSERT INTO `tbl_item` VALUES (1,'The Hobbit (De Luxe Edition)','The nice edition of The Hobbit by J.R.R. Tolkien',1,4,7178),(10,'The Lord of The Rings (De Luxe Edition)','The nice edition of The Lord of The Rings by J.R.R. Tolkien',1,4,8420),(12,'Ku: Toki Pona Dictionary','Finally BookDepository has a toki pona book in English.',1,4,2727);
/*!40000 ALTER TABLE `tbl_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tbl_link`
--

DROP TABLE IF EXISTS `tbl_link`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
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
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_link`
--

LOCK TABLES `tbl_link` WRITE;
/*!40000 ALTER TABLE `tbl_link` DISABLE KEYS */;
INSERT INTO `tbl_link` VALUES (1,'Book Depository','https://www.bookdepository.com/Hobbit-J-R-R-Tolkien/9780007118359?ref=pd_gw_1_pd_gateway_1_1',1),(10,'Book Depository','https://www.bookdepository.com/Lord-Rings-J-R-R-Tolkien/9780007182367?ref=grid-view&qid=1632527627167&sr=1-7',10),(11,'Alternative','https://www.abebooks.co.uk/servlet/BookDetailsPL?bi=30026172939&searchurl=x%3D0%26fe%3Don%26y%3D0%26bi%3D0%26ds%3D30%26bx%3Doff%26sortby%3D1%26tn%3Dthe%2Bhobbit%26an%3Dtolkien&cm_sp=snippet-_-srp1-_-image1#&gid=1&pid=5',10),(14,'Book Depository','https://www.bookdepository.com/Toki-Pona-Dictionary-Sonja-Lang/9780978292362?ref=grid-view&qid=1632530625281&sr=1-2',12);
/*!40000 ALTER TABLE `tbl_link` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tbl_status`
--

DROP TABLE IF EXISTS `tbl_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_status` (
  `id_status` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(32) NOT NULL,
  `desc` text,
  PRIMARY KEY (`id_status`),
  UNIQUE KEY `id_status` (`id_status`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_status`
--

LOCK TABLES `tbl_status` WRITE;
/*!40000 ALTER TABLE `tbl_status` DISABLE KEYS */;
INSERT INTO `tbl_status` VALUES (1,'Available','No one is planning to get this, yet.'),(2,'Reserved','Someone is planning to get this already.'),(3,'Received','This was on the wishlist, but has since been received.');
/*!40000 ALTER TABLE `tbl_status` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tbl_user`
--

DROP TABLE IF EXISTS `tbl_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `tbl_user` (
  `id_user` int(11) NOT NULL AUTO_INCREMENT,
  `email` varchar(128) NOT NULL,
  `password` varchar(128) NOT NULL,
  `domain` varchar(255) NOT NULL DEFAULT 'http://localhost',
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`id_user`),
  UNIQUE KEY `id_user` (`id_user`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tbl_user`
--

LOCK TABLES `tbl_user` WRITE;
/*!40000 ALTER TABLE `tbl_user` DISABLE KEYS */;
INSERT INTO `tbl_user` VALUES (2,'eyamver.istouped','*65F2FE1BC9981E1B387835AE5E59E3FDF25AE752','http://localhost','Bogos Binted'),(4,'michael.oxlong','*65F2FE1BC9981E1B387835AE5E59E3FDF25AE752','http://localhost','longcokc28');
/*!40000 ALTER TABLE `tbl_user` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-10-03  0:11:52
