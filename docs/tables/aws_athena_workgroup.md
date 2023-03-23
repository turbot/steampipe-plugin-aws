# Table: aws_athena_workgroup

A workgroup is an Athena configuration containing information like engine version, output location, data encryption, etc.

## Examples

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
