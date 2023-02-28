# Table: aws_cloudfront_origin_request_policy

An origin request policy controls the values that are include in requests that CloudFront sends to your origin.

## Examples

### Basic info

```sql
select
  name,
  id,
  comment,
  etag,
  last_modified_time
from
  aws_cloudfront_origin_request_policy;
```

### Get details of HTTP headers associated with each origin request policy

```sql
select
  name,
  id,
  headers_config ->> 'HeaderBehavior' as header_behavior,
  headers_config ->> 'Headers' as headers
from
  aws_cloudfront_origin_request_policy;
```
