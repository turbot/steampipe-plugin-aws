---
title: "Steampipe Table: aws_cloudfront_function - Query AWS CloudFront Functions using SQL"
description: "Allows users to query AWS CloudFront Functions to retrieve detailed information about each function, including its ARN, stage, status, and more."
folder: "CloudFront"
---

# Table: aws_cloudfront_function - Query AWS CloudFront Functions using SQL

The AWS CloudFront Function is a feature of Amazon CloudFront that allows you to write lightweight functions in JavaScript for high-scale, latency-sensitive CDN customizations. These functions execute at the edge locations, closer to the viewer, allowing you to manipulate HTTP request and response headers, URL, and methods. This feature helps in delivering a highly personalized content with low latency to your viewers.

## Table Usage Guide

The `aws_cloudfront_function` table in Steampipe provides you with information about functions within AWS CloudFront. This table allows you, as a DevOps engineer, to query function-specific details, including the function's ARN, stage, status, and associated metadata. You can utilize this table to gather insights on functions, such as their status, the events they are associated with, and more. The schema outlines the various attributes of the CloudFront function for you, including the function ARN, creation timestamp, last modified timestamp, and associated tags.

## Examples

### Basic info

```sql+postgres
select
  name,
  status,
  arn,
  e_tag,
  function_config
from
  aws_cloudfront_function;
```

```sql+sqlite
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

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(function_config, '$.Comment') as comment,
  arn,
  status,
  e_tag
from
  aws_cloudfront_function
where
  json_extract(function_metadata, '$.Stage') = 'LIVE';
```

### List functions ordered by its creation time starting with latest first

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  json_extract(function_metadata, '$.Stage') as stage,
  status,
  json_extract(function_metadata, '$.CreatedTime') as created_time,
  json_extract(function_metadata, '$.LastModifiedTime') as last_modified_time
from
  aws_cloudfront_function
order by
  json_extract(function_metadata, '$.CreatedTime') DESC;
```

### List functions updated in the last hour with latest first

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  json_extract(function_metadata, '$.Stage') as stage,
  status,
  json_extract(function_metadata, '$.LastModifiedTime') as last_modified_time
from
  aws_cloudfront_function
where
  datetime(json_extract(function_metadata, '$.LastModifiedTime')) >= datetime('now', '-1 hour')
order by
  json_extract(function_metadata, '$.LastModifiedTime') DESC;
```