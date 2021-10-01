-- MySQL dump 10.13  Distrib 8.0.26, for Linux (x86_64)
--
-- Host: localhost    Database: chatroom
-- ------------------------------------------------------
-- Server version	8.0.26

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `member`
--

DROP TABLE IF EXISTS `member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `member` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` char(20) DEFAULT NULL,
  `password` char(20) DEFAULT NULL,
  `email` char(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `member`
--

LOCK TABLES `member` WRITE;
/*!40000 ALTER TABLE `member` DISABLE KEYS */;
INSERT INTO `member` VALUES (1,'Miles123','Miles123','Miles@123'),(2,'Aaaa','Aaa123456','aaa@123'),(3,'AAAA','Aaa123456','AAA@123'),(4,'BBBB','Bbb123456','bbb@123'),(5,'CCCC','Ccc123456','ccc@123'),(6,'miles','Miles123','miles@123'),(8,'DDDD','Ddd123456','ddd@123'),(9,'dasdas','Ddd123456','dasd@ads'),(10,'Eric','Eric123456','eric@123');
/*!40000 ALTER TABLE `member` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `my_chatroom`
--

DROP TABLE IF EXISTS `my_chatroom`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `my_chatroom` (
  `room_name` char(20) NOT NULL,
  `room_owner` char(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `my_chatroom`
--

LOCK TABLES `my_chatroom` WRITE;
/*!40000 ALTER TABLE `my_chatroom` DISABLE KEYS */;
INSERT INTO `my_chatroom` VALUES ('aaaa','Aaaa'),('bbbb','Aaaa'),('cccc','Aaaa'),('dddd','Aaaa'),('sdaa','Aaaa'),('aassq','DDDD'),('aqqqq','Aaaa'),('fsdfsdf','Aaaa'),('cvcxv','Aaaa'),('qqqq','DDDD'),('uuuuu','DDDD'),('1123456','DDDD'),('fdsfsd','DDDD'),('qwqwq','DDDD'),('rrryry','Miles123'),('dasdadsa','DDDD'),('wwwww','DDDD'),('hthth','DDDD'),('5656565','DDDD'),('7777777','DDDD'),('4545454','DDDD'),('efewfewf','miles'),('dqqqq','DDDD'),('4546464','Aaaa'),('5767676','Aaaa'),('546546546','Aaaa'),('23123123','Aaaa'),('5757575','Aaaa'),('fdsfdsfs','Aaaa'),('fdsfsdfs','Aaaa'),('wdwdw','Aaaa'),('hthtrhr','Aaaa'),('miles_room','miles'),('newroommmm','Aaaa'),('qqqqqqqqqqqq','Aaaa'),('eric_room','Eric'),('fefefe','Aaaa'),('gdfgdfg','miles'),('hgfhgf','miles'),('1111111111111','miles'),('qqsqqq','DDDD');
/*!40000 ALTER TABLE `my_chatroom` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2021-10-01  6:11:05
