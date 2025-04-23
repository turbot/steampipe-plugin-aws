---
title: "Steampipe Table: aws_cloudfront_origin_request_policy - Query AWS CloudFront Origin Request Policies using SQL"
description: "Allows users to query AWS CloudFront Origin Request Policies, providing details about each policy such as ID, name, comment, cookies configuration, headers configuration, query strings configuration, and more."
folder: "CloudFront"
---

# Table: aws_cloudfront_origin_request_policy - Query AWS CloudFront Origin Request Policies using SQL

The AWS CloudFront Origin Request Policy is a feature of Amazon CloudFront, a content delivery network service. It allows you to control how much information about the viewer's request is forwarded to the origin. This includes headers, cookies, and URL query strings, enabling you to customize the content returned by your origin based on the values in the request.

## Table Usage Guide

The `aws_cloudfront_origin_request_policy` table in Steampipe provides you with information about Origin Request Policies within AWS CloudFront. This table allows you, as a DevOps engineer, to query policy-specific details, including ID, name, comment, cookies configuration, headers configuration, query strings configuration, and more. You can utilize this table to gather insights on policies, such as policy configurations and associated metadata. The schema outlines the various attributes of the Origin Request Policy for you, including the policy ID, creation date, last modified date, and associated tags.

## Examples

### Basic info
Explore which AWS Cloudfront origin request policies have been modified recently, gaining insights into potential changes and updates. This can be useful for maintaining security compliance and ensuring correct configuration.

```sql+postgres
select
  name,
  id,
  comment,
  etag,
  last_modified_time
from
  aws_cloudfront_origin_request_policy;
```

```sql+sqlite
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
Determine the characteristics of HTTP headers related to each origin request policy. This can be useful to understand how your CloudFront distributions are configured, which can help in optimizing your web content delivery and troubleshooting issues.

```sql+postgres
select
  name,
  id,
  headers_config ->> 'HeaderBehavior' as header_behavior,
  headers_config ->> 'Headers' as headers
from
  aws_cloudfront_origin_request_policy;
```

```sql+sqlite
select
  name,
  id,
  json_extract(headers_config, '$.HeaderBehavior') as header_behavior,
  json_extract(headers_config, '$.Headers') as headers
from
  aws_cloudfront_origin_request_policy;
```