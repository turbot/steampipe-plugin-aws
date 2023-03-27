# Table: aws_athena_workgroup

A workgroup is an Athena configuration containing information like engine version, output location, data encryption, etc.

## Examples

### List all workgroups with basic information
```sql
select 
  name, 
  description, 
  effective_engine_version, 
  output_location, 
  creation_time 
from 
  aws_athena_workgroup 
order by 
  creation_time
```

### List all workgroups using engine 3

```sql
select 
  name, 
  description 
from 
  aws_athena_workgroup 
where 
  effective_engine_version = 'Athena engine version 3'
```

### Count workgroups in each region
```sql
select 
  region, 
  count(*) 
from 
  aws_athena_workgroup 
group by 
  region
```

### List disabled workgroups
```sql
select 
  name, 
  description, 
  creation_time
from 
  aws_athena_workgroup 
where
  state = 'DISABLED'
```
