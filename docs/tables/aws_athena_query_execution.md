# Table: aws_athena_query_execution

A query execution is all the information about a single instance of an Athena query execution. 

**Important notes:**

- You **_must_** specify `workgroup` in a `where` clause in order to use this table.
- It is possible to join data iwth `aws_athena_workgroup` table to get all query executions from all workgroups (see example below)

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
  workgroup = 'primary'
and
  error_message is not null;
```

### Estimate data read by each workgroup

```sql
select 
  workgroup, 
  sum(data_scanned_in_bytes) 
from 
  aws_athena_query_execution q
join
  aws_athena_workgroup w
on 
  q.workgroup = w.name
group by 
  workgroup;
```
