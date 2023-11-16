---
title: "Table: aws_media_store_container - Query AWS MediaStore Container using SQL"
description: "Allows users to query AWS MediaStore Container information, including ARN, creation time, status, and access logging details."
---

# Table: aws_media_store_container - Query AWS MediaStore Container using SQL

The `aws_media_store_container` table in Steampipe provides information about containers within AWS Elemental MediaStore. This table allows DevOps engineers to query container-specific details, such as ARN, creation time, status, and access logging details. Users can utilize this table to gather insights on containers, such as the container's lifecycle policy, CORS policy, and more. The schema outlines the various attributes of the MediaStore container, including the container ARN, creation time, status, access logging details, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_media_store_container` table, you can use the `.inspect aws_media_store_container` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Name (ARN) of the container. This can be used to join with other tables that contain MediaStore container ARNs.
- `name`: The name of the container. This can be used to join with other tables that contain MediaStore container names.
- `status`: The status of the container. This is useful for filtering containers based on their status.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_media_store_container;
```

### List containers with access logging enabled

```sql
select
  name,
  arn,
  access_logging_enabled
from
  aws_media_store_container
where
  access_logging_enabled;
```
