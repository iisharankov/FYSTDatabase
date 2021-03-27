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
  `LocationID` INT NOT NULL AUTO_INCREMENT,
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
  `InstrumentID` INT NOT NULL AUTO_INCREMENT,
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
  `RuleID` INT NOT NULL AUTO_INCREMENT,
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
  `FileID` INT NOT NULL AUTO_INCREMENT,
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
