# Table: aws_cloudfront_function

Contains configuration information and metadata about a CloudFront function.

## Examples

### Basic info

```sql
select
  account_id,
  akas,
  arn,
  e_tag,
  function_config,
  function_metadata,
  name,
  status,
  title
 from
  aws_cloudfront_function;
```

### List details about published functions

```sql
select
  name,
  arn,
  status,
  e_tag,
  function_config,
  function_metadata
 from
  aws_cloudfront_function
where
  function_metadata ->> 'Stage' = 'LIVE';
```

### List functions ordered by its creation time

```sql
select
  name,
  arn,
  status,
  function_metadata ->> 'CreatedTime' as created_time,
  function_metadata ->> 'LastModifiedTime' as last_modified_time,
  function_metadata ->> 'Stage' as stage
 from
  aws_cloudfront_function
order by
  function_metadata ->> 'CreatedTime' DESC;
```

### List functions updated in the last hour

```sql
select
  name,
  arn,
  status,
  function_metadata ->> 'CreatedTime' as created_time,
  function_metadata ->> 'LastModifiedTime' as last_modified_time,
  function_metadata ->> 'Stage' as stage
from
  aws_cloudfront_function
where
  to_timestamp(function_metadata ->> 'LastModifiedTime', 'yyyy-MM-ddTHH24:MI:SS.MSZ')  >= (now() - interval '1' hour)
order by
  function_metadata ->> 'LastModifiedTime' DESC;
```
