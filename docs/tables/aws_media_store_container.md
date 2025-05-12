---
title: "Steampipe Table: aws_media_store_container - Query AWS MediaStore Container using SQL"
description: "Allows users to query AWS MediaStore Container information, including ARN, creation time, status, and access logging details."
folder: "MediaStore"
---

# Table: aws_media_store_container - Query AWS MediaStore Container using SQL

The AWS MediaStore Container is a high-performance storage service for media data. It offers the performance, consistency, and low latency required to deliver live streaming video content. MediaStore Container is designed to enable you to deliver video content to consumers more quickly and reliably.

## Table Usage Guide

The `aws_media_store_container` table in Steampipe provides you with information about containers within AWS Elemental MediaStore. This table allows you, as a DevOps engineer, to query container-specific details, such as ARN, creation time, status, and access logging details. You can utilize this table to gather insights on containers, such as the container's lifecycle policy, CORS policy, and more. The schema outlines the various attributes of the MediaStore container for you, including the container ARN, creation time, status, access logging details, and associated tags.

## Examples

### Basic info
Explore the status and settings of your AWS Media Store containers to assess their accessibility and operational status. This is beneficial in understanding the overall health and configuration of your media storage infrastructure.

```sql+postgres
select
  name,
  arn,
  status,
  access_logging_enabled,
  creation_time,
  endpoint
from
  aws_media_store_container;
```

```sql+sqlite
select
  name,
  arn,
  status,
  access_logging_enabled,
  creation_time,
  endpoint
from
  aws_media_store_container;
```

### List containers which are in 'CREATING' state
Identify instances where AWS Media Store containers are still in the process of being created. This can be useful in managing resources and troubleshooting potential issues with container creation.

```sql+postgres
select
  name,
  arn,
  status,
  access_logging_enabled,
  creation_time,
  endpoint
from
  aws_media_store_container
where
  status = 'CREATING';
```

```sql+sqlite
select
  name,
  arn,
  status,
  access_logging_enabled,
  creation_time,
  endpoint
from
  aws_media_store_container
where
  status = 'CREATING';
```

### List policy details for the containers
Gain insights into the policy details associated with your containers to better understand their configuration and ensure adherence to your organization's security standards. This can be particularly useful for auditing purposes or to identify any potential areas of risk.

```sql+postgres
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_media_store_container;
```

```sql+sqlite
select
  name,
  policy,
  policy_std
from
  aws_media_store_container;
```

### List containers with access logging enabled
Identify instances where access logging is enabled for AWS Media Store containers. This could be beneficial for auditing purposes or to ensure compliance with data access policies.

```sql+postgres
select
  name,
  arn,
  access_logging_enabled
from
  aws_media_store_container
where
  access_logging_enabled;
```

```sql+sqlite
select
  name,
  arn,
  access_logging_enabled
from
  aws_media_store_container
where
  access_logging_enabled = 1;
```