
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `authorizations` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `token` mediumblob,
  `identifierValue` varchar(255) NOT NULL,
  `identifierType` int NOT NULL,
  `registrationID` bigint(20) NOT NULL,
  `status` int NOT NULL,
  `expires` datetime NOT NULL,
  `allowedChallenges` tinyint,
  `validated` datetime NOT NULL,
  `validationRecord` mediumblob,
  `error` mediumblob,
  `thumbprint` mediumblob,

  PRIMARY KEY (`id`),

  KEY `regId_authorizations_idx` (`registrationID`) COMMENT 'Common lookup',
  CONSTRAINT `regId_authorizations`  FOREIGN KEY (`registrationID`) REFERENCES `registrations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE `authorizations`;
