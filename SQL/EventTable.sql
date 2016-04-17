CREATE DATABASE `ScoreBot` /*!40100 DEFAULT CHARACTER SET utf8 */;

CREATE TABLE `Events` (
  `EventId` int(11) NOT NULL,
  `Description` varchar(45) DEFAULT NULL,
  `I` varchar(45) DEFAULT NULL,
  `Games_GameId` int(11) NOT NULL,
  `Games_Teams_TeamId` int(11) NOT NULL,
  `Games_Teams_Sports_SportsId` int(11) NOT NULL,
  PRIMARY KEY (`EventId`,`Games_GameId`,`Games_Teams_TeamId`,`Games_Teams_Sports_SportsId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Games` (
  `GameId` int(11) NOT NULL AUTO_INCREMENT,
  `AwayId` int(11) DEFAULT NULL,
  `Start` datetime DEFAULT NULL,
  `Finish` datetime DEFAULT NULL,
  `HomeScore` int(11) DEFAULT NULL,
  `AwayScore` int(11) DEFAULT NULL,
  `Status` varchar(45) DEFAULT NULL,
  `homeId` int(11) DEFAULT NULL,
  `url` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`GameId`)
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8;

CREATE TABLE `Sports` (
  `SportsId` int(11) NOT NULL,
  `Name` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`SportsId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Subscription` (
  `SubscriptionId` int(11) NOT NULL AUTO_INCREMENT,
  `Users_UserId` int(11) NOT NULL,
  `Teams_TeamId` int(11) NOT NULL,
  PRIMARY KEY (`SubscriptionId`,`Users_UserId`,`Teams_TeamId`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

CREATE TABLE `Teams` (
  `TeamId` int(11) NOT NULL,
  `Name` varchar(45) DEFAULT NULL,
  `Sports_SportsId` int(11) NOT NULL,
  PRIMARY KEY (`TeamId`,`Sports_SportsId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `Users` (
  `UserId` int(11) NOT NULL AUTO_INCREMENT,
  `UserName` varchar(45) DEFAULT NULL,
  `Platform` varchar(45) DEFAULT NULL,
  `Phone` varchar(45) DEFAULT NULL,
  `Country` varchar(45) DEFAULT NULL,
  `Joined` datetime DEFAULT NULL,
  PRIMARY KEY (`UserId`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8;
