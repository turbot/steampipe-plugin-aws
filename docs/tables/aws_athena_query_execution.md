# Table: aws_athena_query_execution

A query execution is all the information about a single instance of an Athena query execution. 

## Examples

### List all queries that still use older engine

```sql
select
  id,
  workgroup,
  query
from
  aws_athena_query_execution
where
  effective_engine_version != 'Athena engine version 3';
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
