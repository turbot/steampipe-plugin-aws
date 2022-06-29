# Table: aws_cloudfront_function

Contains configuration information and metadata about a CloudFront function.

## Examples

### Basic info

```sql
select
  name,
  status,
  arn,
  e_tag,
  function_config
from
  aws_cloudfront_function;
```

### List details of all functions deployed to the live stage

```sql
select
  name,
  function_config ->> 'Comment' as comment,
  arn,
  status,
  e_tag
from
  aws_cloudfront_function
where
  function_metadata ->> 'Stage' = 'LIVE';
```

### List functions ordered by its creation time starting with latest first

```sql
select
  name,
  arn,
  function_metadata ->> 'Stage' as stage,
  status,
  function_metadata ->> 'CreatedTime' as created_time,
  function_metadata ->> 'LastModifiedTime' as last_modified_time
 from
  aws_cloudfront_function
order by
  function_metadata ->> 'CreatedTime' DESC;
```

### List functions updated in the last hour with latest first

```sql
select
  name,
  arn,
  function_metadata ->> 'Stage' as stage,
  status,
  function_metadata ->> 'LastModifiedTime' as last_modified_time
from
  aws_cloudfront_function
where
  (function_metadata ->> 'LastModifiedTime')::timestamp >= (now() - interval '1' hour)
order by
  function_metadata ->> 'LastModifiedTime' DESC;
```
