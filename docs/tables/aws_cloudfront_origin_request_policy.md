---
title: "Table: aws_cloudfront_origin_request_policy - Query AWS CloudFront Origin Request Policies using SQL"
description: "Allows users to query AWS CloudFront Origin Request Policies, providing details about each policy such as ID, name, comment, cookies configuration, headers configuration, query strings configuration, and more."
---

# Table: aws_cloudfront_origin_request_policy - Query AWS CloudFront Origin Request Policies using SQL

The `aws_cloudfront_origin_request_policy` table in Steampipe provides information about Origin Request Policies within AWS CloudFront. This table allows DevOps engineers to query policy-specific details, including ID, name, comment, cookies configuration, headers configuration, query strings configuration, and more. Users can utilize this table to gather insights on policies, such as policy configurations and associated metadata. The schema outlines the various attributes of the Origin Request Policy, including the policy ID, creation date, last modified date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudfront_origin_request_policy` table, you can use the `.inspect aws_cloudfront_origin_request_policy` command in Steampipe.

### Key columns:

- `id`: The ID of the origin request policy. This can be used to join this table with other tables that reference CloudFront origin request policies.
- `name`: The name of the origin request policy. This can provide a more human-readable reference when joining with other tables.
- `last_modified_time`: The time when the origin request policy was last modified. This can be useful for tracking changes and updates to policies.

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
