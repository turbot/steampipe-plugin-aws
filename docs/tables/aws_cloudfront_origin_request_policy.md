# Table: aws_cloudfront_origin_request_policy

An Origin Request Policy controls the values that are include in requests that CloudFront sends to your origin.

## Examples

### Basic info

```sql
select
  id,
  name,
  comment,
  e_tag,
  last_modified_time,
  type
from
  aws_cloudfront_origin_request_policy;
```

### Get details of http headers associated with origin request policy

```sql
select
  name,
  id,
  headers_config ->> 'HeaderBehavior' as header_behavior,
  headers_config ->> 'Headers' as headers
from
  aws_cloudfront_origin_request_policy;
```
