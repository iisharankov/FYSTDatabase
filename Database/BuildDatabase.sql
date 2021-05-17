drop table if exists `Records`; 
drop table if exists `Files`;
drop table if exists `Rules`; 
drop table if exists `Locations`; 
drop table if exists `Instrument`;

use mydb;
DELETE FROM Files WHERE FileID between 0 and 100;

select * from Instruments;
select * from Locations;
select * from Files;
select * from Records;
select * from Rules;

UPDATE Records SET IsCopying=001 WHERE FileID=1 AND RuleID = 1;

select f.FileID, r.RuleID, f.Size, i.InstrumentID,
i.InstrumentName, f.DateCreated, f.ObjectStorage, f.HashOfBytes, BL.S3Bucket, l.IsCopying
from Files f
join Instruments i on i.InstrumentID=f.InstrumentID
join Rules r on r.InstrumentID=i.InstrumentID
join Locations BL on BL.LocationID=r.LocationID
left join Records l on l.FileID=f.FileID and r.RuleID=l.RuleID
where l.FileID is null AND r.Active=1 order by f.DateCreated;



select l.URL from Records l join Rules r on r.RuleID=l.RuleID 
join Locations BL on BL.LocationID=r.LocationID 
where l.FileID=1 and BL.S3Bucket="fyst";


select * from Locations b where b.LocationID = 1;
    
-- All lgos that have more than 3 entries per fileID
select * from Records where FileID in (
select FileID from Records group by FileID having count(*)>3);



-- Add Instruments
insert into Instruments values(1,'CHAI','expanded CHAI name','Multi-pixel heterodyne receiver for 2 frequency bands observable simultaneously', 64,450,495);
insert into Instruments values(2,'PCAM','Primary Camera','It is very big',1024,25,900);


-- Add Locations
insert into Locations values(1, 'Observatory', 'fyst', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);
insert into Locations values(2, 'Max Plack Bonn', 'germany', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);
insert into Locations values(3, 'Cornell USA', 'cornell', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);
insert into Locations values(4, 'CITA Canada', 'toronto', '0.0.0.0', '9000', 'iisharankov', 'iisharankov', false);


-- Add Rules
insert into Rules values(1, 'copy all files onto FYST server', 1, 1, 1);
insert into Rules values(2, 'copy all files onto FYST server', 2, 1, 1);
insert into Rules values(3, 'copy all CHAI files to Germany', 1, 2, 1);
insert into Rules values(4, 'copy all CHAI files to Cornell', 1, 3, 0);
insert into Ruless values(5, 'copy all CHAI files to Toronto', 1, 4, 1);
insert into Ruless values(6, 'copy all PCAM files to Germany', 2, 2, 0);
insert into Ruless values(7, 'copy all PCAM files to Cornell', 2, 3, 1);
insert into Ruless values(8, 'copy all PCAM files to Toronto', 2, 4, 1); 


-- Data! 
insert into Files values(7, '2020-10-11 12:50:00', 1, 1024, '658b939dc9896f7');

insert into Records values (112, 1, date_add(sysdate(), interval 1 day), 1, '');



select concat('insert into Records values (', a.fileID, ', ', r.RuleID, ', ', 'date_add(sysdate(), interval ', a.fileId, ' day)', ');') from Files a
join  Rules r on a.InstrumentId=r.InstrumentID
where r.active = 1;




-- Tells you which RuleID corresponds to a FileID and Location string (locationname in Locations)
select r.RuleID from Rules r
join Files o on o.InstrumentID=r.InstrumentID
join Locations b on b.LocationID=r.LocationID
where o.FileId = 2 and b.LocationName = "fyst"


-- Tells you which Instruments have active Ruless to be copied to Germany
select b.InsturmentName, c.LocationName from Rules a
join Instrumentsb on a.InstrumentID=b.InstrumentID
join Locations c on c.LocationID=a.LocationID
where c.LocationName='Germany' and a.active=1; 

-- All files within certian date range
select * from Records
where backupdate between '2020-10-12' and '2020-11-28'
order by backupDate;


-- How many files per location and their average size
select lo.LocationName, count(*), avg(f.Size)
from Records l 
join Rules r on l.RuleID=r.RuleID
join Locations lo on lo.locationId=r.locationID
join Files f on f.FileID=l.FileID
group by lo.LocationName;

-- How many file by each location by day
select lo.LocationName, backupdate, count(*)
from Records l 
join Rules r on l.RuleID=r.RuleID
join Locations lo on lo.locationId=r.locationID
group by lo.LocationName, l.backupDate 
order by l.backupDate desc, lo.LocationName;


-- Find each location that has more than 1 file/day
select lo.locationName, backupdate, count(*)
from Records l 
join Rules r on l.RuleID=r.RuleID
join Locations lo on lo.locationId=r.locationID
group by lo.LocationName, l.backupDate having count(*)>1
order by l.backupDate desc, lo.LocationName;


-- Find all Ruless where no files are copied
select  r.RulesDescription, l.FileId from Rules r
left join Records l on l.RuleID=r.RuleID 
where l.RuleID is null;

-- Alter table Instruments rename column InsturmentName to InstrumentName;
-- What files have not been copied (most 'at risk' data):
select distinct f.FileID, f.Size, i.InstrumentName
select distinct f.*
from Files f 
join Instrumentsi on i.InstrumentID=f.InstrumentID
left join Records l on l.FileID=f.FileID
where l.FileID is null;

-- How many files have not ben copied w active Rules
select distinct f.FileID, r.RuleID, f.Size, i.InstrumentName, f.DateCreated, f.ObjectStorage
from Files f 
join Instrumentsi on i.InstrumentID=f.InstrumentID
left join Records l on l.FileID=f.FileID
join Rules r on r.InstrumentID=i.InstrumentID
where l.FileID is null AND r.Active=1 order by f.DateCreated;


-- How many files have not ben copied w active Rules for each Rules
select f.FileID, r.RuleID, f.Size, i.InstrumentName, f.DateCreated, f.ObjectStorage, f.HashOfBytes, BL.LocationName
from Files f 
join Instrumentsi on i.InstrumentID=f.InstrumentID
join Rules r on r.InstrumentID=i.InstrumentID
join Locations BL on BL.LocationID=r.LocationID
left join Records l on l.FileID=f.FileID and r.RuleID=l.RuleID
where l.FileID is null AND r.Active=1 order by f.DateCreated;


-- What files have <2 copies
select f.FileID, count(l.FileID)
from Files f 
left join Records l on l.FileID=f.FileID
group by f.FileID having count(l.FileID)<2;

-- All PCAM files using nested query 
-- select * from Files where InstrumentID in 
-- (select InstrumentId from Instrumentswhere InstrumentName = 'PCAM');

-- Get last record of Records
SELECT * FROM Records ORDER BY FileID DESC LIMIT 1;
SELECT * FROM Files ORDER BY FileID DESC LIMIT 1;

-- select FileID from Files where FileID = (select max(FileID) from Records)
-- Get last row of Files table that has been added to Records
-- select * from Files where FileID = (select max(FileID) from Records)
-- (SELECT * FROM l Records ORDER BY FileID DESC LIMIT 1);
-- PCAM will have cornell, chai goes to germany, housekeeping to all three (tentative/general idea)