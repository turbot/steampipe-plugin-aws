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

### Return user created response header policies only

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

### Return response header policies only modified in the last hour
