# Table: aws_cloudfront_response_headers_policy

An origin request policy controls the values that are include in requests that CloudFront sends to your origin.

## Examples

### Basic info

```sql
select
  name,
  id,
  response_headers_policy_config ->> 'Comment' as description,
  type,
  last_modified_time
from
  aws_cloudfront_response_headers_policy;
```

### List user created response header policies only

```sql
select
  name,
  id,
  response_headers_policy_config ->> 'Comment' as description,
  type,
  last_modified_time
from
  aws_cloudfront_response_headers_policy
where
  type = "custom";
```

### List response header policies that were modified in the last hour

```sql
select
  name,
  id,
  last_modified_time
from
  aws_cloudfront_response_headers_policy
where
  last_modified_time >= (now() - interval '1' hour)
order by
  last_modified_time DESC;
```
