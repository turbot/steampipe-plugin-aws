# Table: aws_cloudfront_response_headers_policy

AWS CloudFront is a globally-distributed network offered by AWS used to securely transfers digital content to its clients using a high transfer speed.

Response headers policies are provided to simplify the process of HTTP header response manipulation.
They allow the user to define CORS, security, and custom response headers as a configuration setting in CloudFront.

This table details the contents of the response header policies.

**Important notes:**

This table supports the optional quals `type`.
Queries with optional quals are optimised to use additional filtering provided by the AWS API function.

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
  type = 'custom';
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
