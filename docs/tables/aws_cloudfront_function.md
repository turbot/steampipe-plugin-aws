# Table: aws_cloudfront_function

CloudFront Functions is ideal for lightweight, short-running functions for use cases like the following:

- Cache key normalization – You can transform HTTP request attributes (headers, query strings, cookies, even the URL path) to create an optimal cache key, which can improve your cache hit ratio.
- Header manipulation – You can insert, modify, or delete HTTP headers in the request or response. For example, you can add a True-Client-IP header to every request.
- URL redirects or rewrites – You can redirect viewers to other pages based on information in the request, or rewrite all requests from one path to another.
- Request authorization – You can validate hashed authorization tokens, such as JSON web tokens (JWT), by inspecting authorization headers or other request metadata.

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
