---
title: "Table: aws_cloudfront_response_headers_policy - Query AWS CloudFront Response Headers Policy using SQL"
description: "Allows users to query AWS CloudFront Response Headers Policies, providing information about the policy configurations that determine the headers CloudFront includes in HTTP responses."
---

# Table: aws_cloudfront_response_headers_policy - Query AWS CloudFront Response Headers Policy using SQL

The `aws_cloudfront_response_headers_policy` table in Steampipe provides information about the Response Headers Policies within AWS CloudFront. This table allows DevOps engineers to query policy-specific details, including policy ID, name, header behavior, and associated custom headers. Users can utilize this table to gather insights on policies, such as custom header configurations, header behavior settings, and more. The schema outlines the various attributes of the Response Headers Policy, including the policy ARN, creation time, last modified time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudfront_response_headers_policy` table, you can use the `.inspect aws_cloudfront_response_headers_policy` command in Steampipe.

Key columns:

- `id`: This is the unique identifier of the Response Headers Policy. It can be used to join this table with other tables that also contain policy IDs.
- `arn`: This is the Amazon Resource Name (ARN) of the Response Headers Policy. It can be used to join with other tables that contain ARNs.
- `name`: This is the name of the Response Headers Policy. It can be used to join with other tables that contain policy names, allowing for more human-readable queries and results.

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
