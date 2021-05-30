use mydb;
DROP PROCEDURE IF EXISTS testwhile;
delimiter $$
create procedure testwhile(size int)
begin
 declare n int;
 declare curID int default 0;
 
 CREATE temporary table tmp(A int);
 insert into tmp values (1);
 insert into tmp values (3);
 insert into tmp values (5);


 select ifnull(MAX(FileID), 0) into curID FROM Files; 
 --
 set n:= curID;
 while n < curID + size do
    set n := n+1;
    
    insert into Files values(n, '2020-10-11 12:50:00', 1, 1024, MD5(n));
    insert into Copies values (n, 1, "test");
	
    -- insert into Records values (n, 1, sysdate());
    -- insert into Records values (n, 3, sysdate());
    -- insert into Records values (n, 5, sysdate());
 end while;
 insert into Records select f.FileID, t.A, sysdate() from Files f join tmp t;
 drop table tmp;
end;
$$

call testwhile(5);

select * from tmp
select * from Records