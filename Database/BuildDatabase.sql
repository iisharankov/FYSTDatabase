drop table if exists `Log`; 
drop table if exists `ObjectFile`;
drop table if exists `Rule`; 
drop table if exists `BackupLocation`; 
drop table if exists `Instrument`;

use mydb;

select * from Instrument;
select * from BackupLocation;
select * from ObjectFile;
select * from Log;
select * from Rule;

DELETE FROM ObjectFile WHERE FileID=1;


insert into Instrument values(1,'CHAI','expanded CHAI name','Multi-pixel heterodyne receiver for 2 frequency bands observable simultaneously', 64,450,495,99);
insert into Instrument values(2,'PCAM','Primary Camera','It is very big',1024,25,900,99);


-- Locations
insert into BackupLocation values(1, 'Observatory', 'fyst', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);




insert into BackupLocation values(1,'fyst','FYST Servers','192.168.1.1');
insert into BackupLocation values(2,'germany', 'Max Plack Bonn','192.168.1.1');
insert into BackupLocation values(3,'cornell','USA','192.168.1.1');
insert into BackupLocation values(4,'toronto','CITA Location','192.168.1.1');


-- Rules
insert into Rule values(1, 'copy all files onto FYST server', 1, 1, 1);
insert into Rule values(2, 'copy all files onto FYST server', 2, 1, 1);
insert into Rule values(3, 'copy all CHAI files to Germany', 1, 2, 1);
insert into Rule values(4, 'copy all CHAI files to Cornell', 1, 3, 0);
insert into Rule values(5, 'copy all CHAI files to Toronto', 1, 4, 1);
insert into Rule values(6, 'copy all PCAM files to Germany', 2, 2, 0);
insert into Rule values(7, 'copy all PCAM files to Cornell', 2, 3, 1);
insert into Rule values(8, 'copy all PCAM files to Toronto', 2, 4, 1);

-- Data! 
insert into ObjectFile values(1, '2020-10-11 12:50:00', 1, 1024, '658b939dc9896f7', 'home/');
insert into ObjectFile values(2, '2020-10-11 12:51:18', 1, 4096, 'bd5b9e4ba981812', 'home/');
insert into ObjectFile values(3, '2020-10-11 12:50:28', 1, 256, '8a969c45f75b6f0', 'home/');
insert into ObjectFile values(4, '2020-10-11 15:50:26', 1, 54254, '80f904907a51cde', 'home/');
insert into ObjectFile values(5, '2020-10-11 18:25:41', 2, 3245, '59797d35de5c74b', 'home/', '??');
insert into ObjectFile values(6, '2020-10-11 18:50:14', 2, 15432, '2b86224636d8352', 'home/', '??');
insert into ObjectFile values(7, '2020-10-11 20:46:28', 2, 1534, '634eafcb989dd94', 'home/', '??');
insert into ObjectFile values(8, '2020-10-11 20:50:00', 2, 107356, 'a8a993f2f918979', 'home/', '??');
insert into ObjectFile values(9, '2020-10-12 07:14:53', 2, 0001, 'f1fb435dcabb51e', 'home/', '??');
insert into ObjectFile values(10, '2020-10-12 08:25:45', 2, 43423, 'ad8ba2641e8f65f', 'home/', '??');
insert into ObjectFile values(11, '2020-10-12 08:55:40', 2, 23423, '658b939dc9896f7', 'home/', '??');
insert into ObjectFile values(12, '2020-10-12 10:20:23', 2, 34224, 'bd5b9e4ba981812', 'home/', '??');
insert into ObjectFile values(13, '2020-10-12 10:42:32', 2, 4332, '8a969c45f75b6f0', 'home/', '??');
insert into ObjectFile values(14, '2020-10-12 10:50:48', 2, 12523, '80f904907a51cde', 'home/', '??');
insert into ObjectFile values(15, '2020-10-12 15:01:00', 1, 3591, '59797d35de5c74b', 'home/', '??');
insert into ObjectFile values(16, '2020-10-12 15:05:09', 1, 4593, '2b86224636d8352', 'home/', '??');
insert into ObjectFile values(17, '2020-10-12 15:05:10', 1, 345345, '634eafcb989dd94', 'home/', '??');
insert into ObjectFile values(18, '2020-10-12 15:05:22', 1, 6753, 'a8a993f2f918979', 'home/', '??');
insert into ObjectFile values(19, '2020-10-12 16:50:00', 1, 23571, 'f1fb435dcabb51e', 'home/', '??');
insert into ObjectFile values(20, '2020-10-13 12:50:00', 2, 2358, 'ad8ba2641e8f65f', 'home/', '??');
insert into ObjectFile values(21, '2020-10-17 12:50:00', 2, 25637, 'edhjwedewdww', 'home/', '??');
insert into ObjectFile values(22, '2020-10-18 12:50:00', 2, 23234, 'wqewqeqwe', 'home/', '??');
insert into ObjectFile values(23, '2020-10-18 12:50:00', 1, 23234, 'wqewqeqwe', 'home/', '??');
insert into ObjectFile values(24, '2020-10-18 12:50:00', 3, 23234, '472383214', 'home/', '??');
insert into ObjectFile values(25, '2020-10-6 12:50:00', 3, 23234, '472383214', 'home/', '??');


-- insert into Log values(1, 1, '2020-10-12 12:50:00');
-- insert into Log values(1, 3, '2020-10-12 12:55:00');
-- insert into Log values(2, 1, '2020-10-13 12:50:00');
-- insert into Log values(2, 3, '2020-10-13 12:55:00');
-- insert into Log values(3, 1, '2020-10-13 12:50:00');
-- insert into Log values(3, 3, '2020-10-14 12:55:00');



insert into Log values (1, 1, date_add(sysdate(), interval 1 day), 0, '');
insert into Log values (2, 1, date_add(sysdate(), interval 2 day), 0, '');
insert into Log values (3, 1, date_add(sysdate(), interval 3 day), 0, '');
insert into Log values (4, 1, date_add(sysdate(), interval 4 day), 0, '');
insert into Log values (15, 1, date_add(sysdate(), interval 15 day), 0, '');
insert into Log values (16, 1, date_add(sysdate(), interval 16 day), 0, '');
insert into Log values (17, 1, date_add(sysdate(), interval 17 day), 0, '');
insert into Log values (18, 1, date_add(sysdate(), interval 18 day), 0, '');
insert into Log values (19, 1, date_add(sysdate(), interval 19 day), 0, '');
insert into Log values (1, 3, date_add(sysdate(), interval 1 day), 0, '');
insert into Log values (2, 3, date_add(sysdate(), interval 2 day), 0, '');
insert into Log values (3, 3, date_add(sysdate(), interval 3 day), 0, '');
insert into Log values (4, 3, date_add(sysdate(), interval 4 day), 0, '');
insert into Log values (15, 3, date_add(sysdate(), interval 15 day), 0, '');
insert into Log values (16, 3, date_add(sysdate(), interval 16 day), 0, '');
insert into Log values (17, 3, date_add(sysdate(), interval 17 day), 0, '');
insert into Log values (18, 3, date_add(sysdate(), interval 18 day), 0, '');
insert into Log values (19, 3, date_add(sysdate(), interval 19 day), 0, '');
insert into Log values (5, 5, date_add(sysdate(), interval 5 day), 0, '');
insert into Log values (6, 5, date_add(sysdate(), interval 6 day), 0, '');
insert into Log values (7, 5, date_add(sysdate(), interval 7 day), 0, '');
insert into Log values (8, 5, date_add(sysdate(), interval 8 day), 0, '');
insert into Log values (9, 5, date_add(sysdate(), interval 9 day), 0, '');
insert into Log values (10, 5, date_add(sysdate(), interval 10 day), 0, '');
insert into Log values (11, 5, date_add(sysdate(), interval 11 day), 0, '');
insert into Log values (12, 5, date_add(sysdate(), interval 12 day), 0, '');
insert into Log values (13, 5, date_add(sysdate(), interval 13 day), 0, '');
insert into Log values (14, 5, date_add(sysdate(), interval 14 day), 0, '');
insert into Log values (20, 5, date_add(sysdate(), interval 20 day), 0, '');
insert into Log values (5, 6, date_add(sysdate(), interval 5 day), 0, '');
insert into Log values (6, 6, date_add(sysdate(), interval 6 day), 0, '');
insert into Log values (7, 6, date_add(sysdate(), interval 7 day), 0, '');
insert into Log values (8, 6, date_add(sysdate(), interval 8 day), 0, '');
insert into Log values (9, 6, date_add(sysdate(), interval 9 day), 0, '');
insert into Log values (10, 6, date_add(sysdate(), interval 10 day), 0, '');
insert into Log values (11, 6, date_add(sysdate(), interval 11 day), 0, '');
insert into Log values (12, 6, date_add(sysdate(), interval 12 day), 0, '');
insert into Log values (13, 6, date_add(sysdate(), interval 13 day), 0, '');
insert into Log values (14, 6, date_add(sysdate(), interval 14 day), 0, '');
insert into Log values (20, 6, date_add(sysdate(), interval 20 day), 0, '');
insert into Log values (21, 5, date_add(sysdate(), interval 23 day), 0, '');



select concat('insert into Log values (', a.fileID, ', ', r.RuleID, ', ', 'date_add(sysdate(), interval ', a.fileId, ' day)', ');') from ObjectFile a
join  Rule r on a.InstrumentId=r.InstrumentID
where r.active = 1;




-- Tells you which RuleID corresponds to a FileID and Location string (locationname in BackupLocation)
select r.RuleID from Rule r
join ObjectFile o on o.InstrumentID=r.InstrumentID
join BackupLocation b on b.LocationID=r.LocationID
where o.FileId = 2 and b.LocationName = "fyst"


-- Tells you which Instruments have active Rules to be copied to Germany
select b.InsturmentName, c.LocationName from Rule a
join Instrument b on a.InstrumentID=b.InstrumentID
join BackupLocation c on c.LocationID=a.LocationID
where c.LocationName='Germany' and a.active=1; 

-- All files within certian date range
select * from Log
where backupdate between '2020-10-12' and '2020-11-28'
order by backupDate;


-- How many files per location and their average size
select lo.LocationName, count(*), avg(f.Size)
from Log l 
join Rule r on l.RuleID=r.RuleID
join BackupLocation lo on lo.locationId=r.locationID
join ObjectFile f on f.FileID=l.FileID
group by lo.LocationName;

-- How many file by each location by day
select lo.LocationName, backupdate, count(*)
from Log l 
join Rule r on l.RuleID=r.RuleID
join BackupLocation lo on lo.locationId=r.locationID
group by lo.LocationName, l.backupDate 
order by l.backupDate desc, lo.LocationName;


-- Find each location that has more than 1 file/day
select lo.locationName, backupdate, count(*)
from Log l 
join Rule r on l.RuleID=r.RuleID
join BackupLocation lo on lo.locationId=r.locationID
group by lo.LocationName, l.backupDate having count(*)>1
order by l.backupDate desc, lo.LocationName;


-- Find all Rules where no files are copied
select  r.RuleDescription, l.FileId from Rule r
left join Log l on l.RuleID=r.RuleID 
where l.RuleId is null;

-- Alter table Instrument  rename column InsturmentName to InstrumentName;
-- What files have not been copied (most 'at risk' data):
select distinct f.FileID, f.Size, i.InstrumentName
select distinct f.*
from ObjectFile f 
join Instrument i on i.InstrumentID=f.InstrumentID
left join Log l on l.FileID=f.FileID
where l.FileID is null;

-- How many files have not ben copied w active rule
select distinct f.FileID, f.Size, i.InstrumentName, f.DateCreated, f.ObjectStorage
from ObjectFile f 
join Instrument i on i.InstrumentID=f.InstrumentID
left join Log l on l.FileID=f.FileID
join Rule r on r.InstrumentID=i.InstrumentID
where l.FileID is null AND r.Active=1 order by f.DateCreated;


-- How many files have not ben copied w active rule for each rule
select f.FileID, r.RuleID, f.Size, i.InstrumentName, f.DateCreated, f.ObjectStorage, f.HashOfBytes, BL.LocationName
from ObjectFile f 
join Instrument i on i.InstrumentID=f.InstrumentID
join Rule r on r.InstrumentID=i.InstrumentID
join BackupLocation BL on BL.LocationID=r.LocationID
left join Log l on l.FileID=f.FileID and r.RuleId=l.RuleID
where l.FileID is null AND r.Active=1 order by f.DateCreated;


-- What files have <2 copies
select f.FileID, count(l.FileID)
from ObjectFile f 
left join Log l on l.FileID=f.FileID
group by f.FileID having count(l.FileID)<2;

-- All PCAM files using nested query 
select * from ObjectFile where InstrumentID in 
(select InstrumentId from Instrument where InstrumentName = 'PCAM');

-- Get last record of Log
SELECT * FROM Log ORDER BY FileID DESC LIMIT 1;
SELECT * FROM ObjectFile ORDER BY FileID DESC LIMIT 1;

-- select FileID from ObjectFile where FileID = (select max(FileID) from Log)
-- Get last row of ObjectFile table that has been added to Log
select * from ObjectFile where FileID = (select max(FileID) from Log)
-- (SELECT * FROM l Log ORDER BY FileID DESC LIMIT 1);
-- PCAM will have cornell, chai goes to germany, housekeeping to all three (tentative/general idea)