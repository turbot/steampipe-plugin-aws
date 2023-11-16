---
title: "Table: aws_ram_principal_association - Query AWS RAM Principal Associations using SQL"
description: "Allows users to query AWS RAM Principal Associations. The `aws_ram_principal_association` table in Steampipe provides information about principal associations within AWS Resource Access Manager (RAM). This table allows DevOps engineers to query principal-specific details, including resource share ARN, principal ARN, creation time, and associated tags. Users can utilize this table to gather insights on principal associations, such as their status, external status, and more. The schema outlines the various attributes of the principal association, including the resource share ARN, principal ARN, creation time, and associated tags."
---

# Table: aws_ram_principal_association - Query AWS RAM Principal Associations using SQL

The `aws_ram_principal_association` table in Steampipe provides information about principal associations within AWS Resource Access Manager (RAM). This table allows DevOps engineers to query principal-specific details, including resource share ARN, principal ARN, creation time, and associated tags. Users can utilize this table to gather insights on principal associations, such as their status, external status, and more. The schema outlines the various attributes of the principal association, including the resource share ARN, principal ARN, creation time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ram_principal_association` table, you can use the `.inspect aws_ram_principal_association` command in Steampipe.

**Key columns**:

- `resource_share_arn`: The Amazon Resource Name (ARN) of the resource share. This can be used to join with other tables that contain resource share ARNs.
- `principal_arn`: The ARN of the principal. This can be used to join with other tables that contain principal ARNs.
- `creation_time`: The time when the association was created. This can be used for sorting and filtering based on the time of creation.

## Examples

### Basic info

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association;
```

### List permissions attached with each principal associated

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  p ->> 'Arn' as resource_share_permission_arn,
  p ->> 'Status' as resource_share_permission_status
from
  aws_ram_principal_association,
  jsonb_array_elements(resource_share_permission) p;
```

### Get principals that failed association

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association
where
  status = 'FAILED';
```