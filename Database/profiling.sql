set profiling=1;
show profiles;

select f.FileID, r.RuleID, f.Size, i.InstrumentID,
i.InstrumentName, f.DateCreated, f.HashOfBytes, BL.S3Bucket
from Files f
join Copies c on c.FileID=f.FileID
join Instruments i on i.InstrumentID=f.InstrumentID
join Rules r on r.InstrumentID=i.InstrumentID
join Locations BL on BL.LocationID=r.LocationID
left join Records l on l.FileID=f.FileID and r.RuleID=l.RuleID
where l.FileID is null AND r.Active=1 order by f.DateCreated ;


CREATE TABLE profileResults select * from profiles;
set profiling=0;

SELECT EVENT_ID, TRUNCATE(TIMER_WAIT/1000000000000,6) as Duration, SQL_TEXT FROM performance_schema.events_statements_history_long;
select * FROM performance_schema.events_statements_history_long;


SELECT * FROM INFORMATION_SCHEMA.PROFILING WHERE QUERY_ID=1063;
SELECT Query_ID, sum(Duration), sum(CPU_USER), sum(CPU_SYSTEM) FROM INFORMATION_SCHEMA.PROFILING WHERE QUERY_ID=1065;


INSERT INTO Duration (Query_ID, Status, DURATION)
SELECT QUERY_ID, STATE, t.DURATION
FROM   INFORMATION_SCHEMA.PROFILING JOIN (
  SELECT   QUERY_ID, MAX(SEQ) AS SEQ, SUM(DURATION) AS DURATION
  FROM     INFORMATION_SCHEMA.PROFILING
  WHERE    QUERY_ID = 1065
  GROUP BY QUERY_ID  -- superfluous in the presence of the above filter
) t USING (QUERY_ID, SEQ);


select * from performance_schema.threads

 SELECT THREAD_ID, EVENT_ID, END_EVENT_ID, SQL_TEXT, NESTING_EVENT_ID
 FROM events_statements_history_long
 WHERE THREAD_ID = 1063
 AND EVENT_NAME = 'statement/sql/select'
 ORDER BY EVENT_ID DESC LIMIT 3