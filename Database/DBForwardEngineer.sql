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
-- Table `mydb`.`Locations`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Locations` (
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
-- Table `mydb`.`Instruments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Instruments` (
  `InstrumentID` INT NOT NULL AUTO_INCREMENT,
  `InstrumentName` VARCHAR(10) NOT NULL COMMENT 'Mianly CHAI and P-CAM for now',
  `FullName` VARCHAR(45) NOT NULL COMMENT 'Unabbreviated name',
  `Description` VARCHAR(2000) NOT NULL,
  `NumberOfPixels` INT NOT NULL,
  `FrequencyMin` INT NOT NULL,
  `FrequencyMax` INT NOT NULL,
  PRIMARY KEY (`InstrumentID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Rules`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Rules` (
  `RuleID` INT NOT NULL AUTO_INCREMENT,
  `RuleDescription` VARCHAR(2000) NULL DEFAULT NULL,
  `InstrumentID` INT NOT NULL,
  `LocationID` INT NOT NULL,
  `Active` TINYINT NOT NULL,
  PRIMARY KEY (`RuleID`),
  INDEX `FK_Rules_LocationID` (`LocationID` ASC) VISIBLE,
  INDEX `FK_Rules_InstrumentID` (`InstrumentID` ASC) VISIBLE,
  CONSTRAINT `FK_Rules_InstrumentID`
    FOREIGN KEY (`InstrumentID`)
    REFERENCES `mydb`.`Instruments` (`InstrumentID`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT,
  CONSTRAINT `FK_Rules_LocationID`
    FOREIGN KEY (`LocationID`)
    REFERENCES `mydb`.`Locations` (`LocationID`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Files`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Files` (
  `FileID` INT NOT NULL AUTO_INCREMENT,
  `DateCreated` DATETIME NOT NULL,
  `InstrumentID` INT NOT NULL,
  `Size` INT UNSIGNED NOT NULL COMMENT 'In bytes',
  `HashOfBytes` VARCHAR(500) NOT NULL,
  PRIMARY KEY (`FileID`),
  INDEX `FK_Files_InstrumentID` (`InstrumentID` ASC) VISIBLE,
  CONSTRAINT `FK_Files_InstrumentID`
    FOREIGN KEY (`InstrumentID`)
    REFERENCES `mydb`.`Instruments` (`InstrumentID`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Records`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Records` (
  `FileID` INT NOT NULL,
  `RuleID` INT NOT NULL,
  `BackupDate` DATETIME NOT NULL,
  INDEX `FK_Records_RuleID` (`RuleID` ASC) VISIBLE,
  INDEX `FK_Records_FileID` (`FileID` ASC) VISIBLE,
  CONSTRAINT `FK_Logs_RuleID`
    FOREIGN KEY (`RuleID`)
    REFERENCES `mydb`.`Rules` (`RuleID`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT,
  CONSTRAINT `FK_Logs_FileID`
    FOREIGN KEY (`FileID`)
    REFERENCES `mydb`.`Files` (`FileID`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `mydb`.`Copies`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Copies` (
  `FileID` INT NULL,
  `LocationID` INT NULL,
  `URL` VARCHAR(1000) NULL,
  INDEX `FK_Copies_FileID` (`FileID` ASC) VISIBLE,
  INDEX `FK_Copies_LocationID` (`LocationID` ASC) VISIBLE,
  CONSTRAINT `FK_Copies_FileID`
    FOREIGN KEY (`FileID`)
    REFERENCES `mydb`.`Files` (`FileID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `FK_Copies_LocationID`
    FOREIGN KEY (`LocationID`)
    REFERENCES `mydb`.`Locations` (`LocationID`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
