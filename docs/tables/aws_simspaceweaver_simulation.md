# Table: aws_simspaceweaver_simulation

AWS SimSpace Weaver is a service that you can use to build and run large-scale spatial simulations in the AWS Cloud.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  status,
  execution_id,
  schema_error
from
  aws_simspaceweaver_simulation;
```

### List simulations older than 30 days

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_simspaceweaver_simulation
where
  creation_time >= now() - interval '30' day;
```

### List failed simulations

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_simspaceweaver_simulation
where
  status = 'FAILED';
```

### Get logging configurations for simulations

```sql
select
  name,
  arn,
  jsonb_pretty(d)
from
  aws_simspaceweaver_simulation,
  jsonb_array_elements(logging_configuration -> 'Destinations') as d;
```

### Get bucket details for simulations

```sql
select
  s.name,
  s.arn,
  s.schema_s3_location ->> 'BucketName' as bucket_name,
  s.schema_s3_location ->> 'ObjectKey' as object_key,
  b.versioning_enabled,
  b.block_public_acls,
  b.acl
from
  aws_simspaceweaver_simulation as s,
  aws_s3_bucket as b
where
  s.schema_s3_location ->> 'BucketName' = b.name;
```
