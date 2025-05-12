---
title: "Steampipe Table: aws_cloudfront_response_headers_policy - Query AWS CloudFront Response Headers Policy using SQL"
description: "Allows users to query AWS CloudFront Response Headers Policies, providing information about the policy configurations that determine the headers CloudFront includes in HTTP responses."
folder: "CloudFront"
---

# Table: aws_cloudfront_response_headers_policy - Query AWS CloudFront Response Headers Policy using SQL

The AWS CloudFront Response Headers Policy is a feature within AWS CloudFront that allows you to manage and customize the HTTP headers returned in the response from your CloudFront distributions. This can be used to enhance the security of your application, improve the caching efficiency, or to provide additional information to the clients. With this policy, you can add, remove, or modify the values of HTTP header fields, providing you with greater control over your content delivery.

## Table Usage Guide

The `aws_cloudfront_response_headers_policy` table in Steampipe provides you with information about the Response Headers Policies within AWS CloudFront. This table allows you, as a DevOps engineer, to query policy-specific details, including policy ID, name, header behavior, and associated custom headers. You can utilize this table to gather insights on policies, such as custom header configurations, header behavior settings, and more. The schema outlines the various attributes of the Response Headers Policy for you, including the policy ARN, creation time, last modified time, and associated tags.

**Important Notes**
- This table supports the optional quals `type`.
- Queries with optional quals are optimised to use additional filtering provided by the AWS API function.

## Examples

### Basic info
Discover the segments that have been recently modified in your AWS Cloudfront response headers policy. This can be useful for assessing the elements within the policy including their names, IDs, and descriptions, and understanding any changes or updates that have been made.

```sql+postgres
select
  name,
  id,
  response_headers_policy_config ->> 'Comment' as description,
  type,
  last_modified_time
from
  aws_cloudfront_response_headers_policy;
```

```sql+sqlite
select
  name,
  id,
  json_extract(response_headers_policy_config, '$.Comment') as description,
  type,
  last_modified_time
from
  aws_cloudfront_response_headers_policy;
```

### List user created response header policies only
Determine the areas in which user-created response header policies exist within the AWS Cloudfront service. This query is beneficial for understanding the custom configurations that have been implemented, along with their last modification time.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  json_extract(response_headers_policy_config, '$.Comment') as description,
  type,
  last_modified_time
from
  aws_cloudfront_response_headers_policy
where
  type = 'custom';
```

### List response header policies that were modified in the last hour
Determine the areas in which response header policies have been recently updated within the last hour. This is useful to track changes and maintain the security and efficiency of your AWS Cloudfront configurations.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  last_modified_time
from
  aws_cloudfront_response_headers_policy
where
  last_modified_time >= (datetime('now','-1 hours'))
order by
  last_modified_time DESC;
```