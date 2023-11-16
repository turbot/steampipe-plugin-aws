---
title: "Table: aws_ram_resource_association - Query AWS RAM Resource Associations using SQL"
description: "Allows users to query AWS RAM Resource Associations to retrieve information about the associations between resources and resource shares."
---

# Table: aws_ram_resource_association - Query AWS RAM Resource Associations using SQL

The `aws_ram_resource_association` table in Steampipe provides information about resource associations within AWS Resource Access Manager (RAM). This table allows DevOps engineers to query association-specific details, including the associated resource ARN, resource share ARN, association type, and status. Users can utilize this table to gather insights on resource associations, such as resources associated with a specific resource share, status of the association, and more. The schema outlines the various attributes of the resource association, including the resource ARN, resource share ARN, association type, status, and creation time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ram_resource_association` table, you can use the `.inspect aws_ram_resource_association` command in Steampipe.

**Key columns**:

- `resource_arn`: The Amazon Resource Name (ARN) of the resource. It can be used to join with other tables to get more information about the resource.
- `resource_share_arn`: The ARN of the resource share. This can be used to join with the `aws_ram_resource_share` table to get more details about the resource share.
- `association_type`: The type of association, either 'PRINCIPAL' or 'RESOURCE'. This column is useful to filter the associations based on their type.

## Examples

### Basic info

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_resource_association;
```

### List permissions attached with each shared resource associated

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  p ->> 'Arn' as resource_share_permission_arn,
  p ->> 'Status' as resource_share_permission_status
from
  aws_ram_resource_association,
  jsonb_array_elements(resource_share_permission) p;
```

### Get resources that failed association

```sql
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_resource_association
where
  status = 'FAILED';
```