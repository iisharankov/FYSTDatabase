-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mydb` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `mydb` ;

-- -----------------------------------------------------
-- Table `mydb`.`BackupLocation`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`BackupLocation` (
  `LocationID` INT NOT NULL,
  `LocationName` VARCHAR(45) NOT NULL,
  `S3Bucket` VARCHAR(45) NOT NULL,
  `IPAddress` VARCHAR(100) NULL DEFAULT NULL,
  `Port` INT NULL,
  `Username` VARCHAR(100) NULL,
  `Password` VARCHAR(100) NULL,
  `SSL` TINYINT NULL,
  PRIMARY KEY (`LocationID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Instrument`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Instrument` (
  `InstrumentID` INT NOT NULL,
  `InstrumentName` VARCHAR(10) NOT NULL COMMENT 'Mianly CHAI and P-CAM for now',
  `FullName` VARCHAR(45) NOT NULL COMMENT 'Unabbreviated name',
  `Description` VARCHAR(2000) NOT NULL,
  `NumberOfPixels` INT NOT NULL,
  `FrequencyMin` INT NOT NULL,
  `FrequencyMax` INT NOT NULL,
  `TempRange` INT NOT NULL,
  PRIMARY KEY (`InstrumentID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Rule`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Rule` (
  `RuleID` INT NOT NULL,
  `RuleDescription` VARCHAR(2000) NULL DEFAULT NULL,
  `InstrumentID` INT NOT NULL,
  `LocationID` INT NOT NULL,
  `Active` TINYINT NOT NULL,
  PRIMARY KEY (`RuleID`),
  INDEX `FK_BackupLocation_idx` (`LocationID` ASC) VISIBLE,
  INDEX `FK_InstrumentID_idx` (`InstrumentID` ASC) VISIBLE,
  CONSTRAINT `FK_BackupRules_InstrumentID`
    FOREIGN KEY (`InstrumentID`)
    REFERENCES `mydb`.`Instrument` (`InstrumentID`),
  CONSTRAINT `FK_BackupRules_LocationID`
    FOREIGN KEY (`LocationID`)
    REFERENCES `mydb`.`BackupLocation` (`LocationID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`ObjectFile`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`ObjectFile` (
  `FileID` INT NOT NULL,
  `DateCreated` DATETIME NOT NULL,
  `InstrumentID` INT NOT NULL,
  `Size` INT UNSIGNED NOT NULL COMMENT 'In bytes',
  `HashOfBytes` VARCHAR(500) NOT NULL,
  `ObjectStorage` VARCHAR(1000) NOT NULL COMMENT 'Not sure how to map ObjectStorage, never confronted with it before\\n',
  PRIMARY KEY (`FileID`),
  INDEX `FK_InstrumentID_idx` (`InstrumentID` ASC) VISIBLE,
  CONSTRAINT `FK_Files_InstrumentID`
    FOREIGN KEY (`InstrumentID`)
    REFERENCES `mydb`.`Instrument` (`InstrumentID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Log`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Log` (
  `FileID` INT NOT NULL,
  `RuleID` INT NOT NULL,
  `BackupDate` DATETIME NOT NULL,
  `IsCopying` TINYINT(3) UNSIGNED ZEROFILL NULL DEFAULT NULL,
  `URL` VARCHAR(1000) NOT NULL,
  INDEX `FK_BackupRule_idx` (`RuleID` ASC) VISIBLE,
  INDEX `FK_BackupLog_FileID` (`FileID` ASC) VISIBLE,
  CONSTRAINT `FK_BackupLog_BackupRuleID`
    FOREIGN KEY (`RuleID`)
    REFERENCES `mydb`.`Rule` (`RuleID`),
  CONSTRAINT `FK_BackupLog_FileID`
    FOREIGN KEY (`FileID`)
    REFERENCES `mydb`.`ObjectFile` (`FileID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;



-- Populate Database with mock metadata
use mydb 

-- Add Instruments
insert into Instrument values(1,'CHAI','expanded CHAI name','Multi-pixel heterodyne receiver for 2 frequency bands observable simultaneously', 64,450,495,99);
insert into Instrument values(2,'PCAM','Primary Camera','It is very big',1024,25,900,99);


-- Add Locations
insert into BackupLocation values(1, 'Observatory', 'fyst', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);
insert into BackupLocation values(2, 'Max Plack Bonn', 'germany', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);
insert into BackupLocation values(3, 'Cornell USA', 'cornell', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);
insert into BackupLocation values(4, 'CITA Canada', 'toronto', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);


-- Add Rules
insert into Rule values(1, 'copy all files onto FYST server', 1, 1, 1);
insert into Rule values(2, 'copy all files onto FYST server', 2, 1, 1);
insert into Rule values(3, 'copy all CHAI files to Germany', 1, 2, 1);
insert into Rule values(4, 'copy all CHAI files to Cornell', 1, 3, 0);
insert into Rule values(5, 'copy all CHAI files to Toronto', 1, 4, 1);
insert into Rule values(6, 'copy all PCAM files to Germany', 2, 2, 0);
insert into Rule values(7, 'copy all PCAM files to Cornell', 2, 3, 1);
insert into Rule values(8, 'copy all PCAM files to Toronto', 2, 4, 1);