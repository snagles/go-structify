CREATE DATABASE IF NOT EXISTS test;
DROP TABLE IF EXISTS test.`all_data_types`;
CREATE TABLE test.`all_data_types` (
`varchar` VARCHAR( 20 ) NOT NULL ,
`tinyint` TINYINT NOT NULL ,
`text` TEXT NOT NULL ,
`date` DATE NOT NULL ,
`smallint` SMALLINT NOT NULL ,
`mediumint` MEDIUMINT NOT NULL ,
`int` INT NOT NULL ,
`bigint` BIGINT NOT NULL ,
`float` FLOAT( 10, 2 ) NOT NULL ,
`double` DOUBLE NOT NULL ,
`decimal` DECIMAL( 10, 2 ) NOT NULL ,
`datetime` DATETIME NOT NULL ,
`timestamp` TIMESTAMP NOT NULL ,
`time` TIME NOT NULL ,
`year` YEAR NOT NULL ,
`char` CHAR( 10 ) NOT NULL ,
`tinyblob` TINYBLOB NOT NULL ,
`tinytext` TINYTEXT NOT NULL ,
`blob` BLOB NOT NULL ,
`mediumblob` MEDIUMBLOB NOT NULL ,
`mediumtext` MEDIUMTEXT NOT NULL ,
`longblob` LONGBLOB NOT NULL ,
`longtext` LONGTEXT NOT NULL ,
`enum` ENUM( '1', '2', '3' ) NOT NULL ,
`set` SET( '1', '2', '3' ) NOT NULL ,
`bool` BOOL NOT NULL ,
`binary` BINARY( 20 ) NOT NULL ,
`varbinary` VARBINARY( 20 ) NOT NULL,
`varcharnull` VARCHAR( 20 ) NULL ,
`tinyintnull` TINYINT NULL ,
`textnull` TEXT NULL ,
`datenull` DATE NULL ,
`smallintnull` SMALLINT NULL ,
`mediumintnull` MEDIUMINT NULL ,
`intnull` INT NULL ,
`bigintnull` BIGINT NULL ,
`floatnul` FLOAT( 10, 2 ) NULL ,
`doublenull` DOUBLE NULL ,
`decimalnull` DECIMAL( 10, 2 ) NULL ,
`datetimenull` DATETIME NULL ,
`timestampnull` TIMESTAMP NULL ,
`timenull` TIME NULL ,
`yearnull` YEAR NULL ,
`charnull` CHAR( 10 ) NULL ,
`tinyblobnull` TINYBLOB NULL ,
`tinytextnull` TINYTEXT NULL ,
`blobnull` BLOB NULL ,
`mediumblobnull` MEDIUMBLOB NULL ,
`mediumtextnull` MEDIUMTEXT NULL ,
`longblobnull` LONGBLOB NULL ,
`longtextnull` LONGTEXT NULL ,
`enumnull` ENUM( '1', '2', '3' ) NULL ,
`setnull` SET( '1', '2', '3' ) NULL ,
`boolnull` BOOL NULL ,
`binarynull` BINARY( 20 ) NULL ,
`varbinarynull` VARBINARY( 20 ) NOT NULL
);