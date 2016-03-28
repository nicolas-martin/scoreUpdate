CREATE TABLE `Event` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Type` varchar(45) NOT NULL,
  `Media` varchar(45) DEFAULT NULL,
  `MatchId` int(11) DEFAULT NULL,
  `Score` varchar(45) DEFAULT NULL,
  `IsSent` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Id_UNIQUE` (`Id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

