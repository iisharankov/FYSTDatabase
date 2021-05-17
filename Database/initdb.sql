
-- Populate Database with mock metadata
use mydb 

-- Add Instruments
insert into Instruments values(1,'CHAI','expanded CHAI name','Multi-pixel heterodyne receiver for 2 frequency bands observable simultaneously', 64,450,495,99);
insert into Instruments values(2,'PCAM','Primary Camera','It is very big',1024,25,900,99);


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
insert into Rules values(5, 'copy all CHAI files to Toronto', 1, 4, 1);
insert into Rules values(6, 'copy all PCAM files to Germany', 2, 2, 0);
insert into Rules values(7, 'copy all PCAM files to Cornell', 2, 3, 1);
insert into Rules values(8, 'copy all PCAM files to Toronto', 2, 4, 1); 
