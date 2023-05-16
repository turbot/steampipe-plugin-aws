# Table: aws_athena_query_execution

A query execution is all the information about a single instance of an Athena query execution. 

## Examples

### List all queries in error

```sql
select
  id,
  query,
  error_message,
  error_type
from
  aws_athena_query_execution
where
  error_message is not null;
```

### Estimate data read by each workgroup

```sql
select 
  workgroup, 
  sum(data_scanned_in_bytes) 
from 
  aws_athena_query_execution
group by 
  workgroup;
```

### Find queries with biggest execution time

```sql
select
  id,
  query,
  workgroup,
  engine_execution_time_in_millis 
from
  aws_athena_query_execution 
order by
  engine_execution_time_in_millis limit 5;
```

### Find most used databases

```sql
select
  database,
  count(id) as nb_query 
from
  aws_athena_query_execution 
group by
  database 
order by
  nb_query limit 5;
```
