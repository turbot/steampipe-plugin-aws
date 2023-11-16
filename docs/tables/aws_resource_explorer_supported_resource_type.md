---
title: "Table: aws_resource_explorer_supported_resource_type - Query AWS Resource Explorer Supported Resource Types using SQL"
description: "Allows users to query AWS Resource Explorer Supported Resource Types to obtain details about supported resource types for AWS Resource Groups."
---

# Table: aws_resource_explorer_supported_resource_type - Query AWS Resource Explorer Supported Resource Types using SQL

The `aws_resource_explorer_supported_resource_type` table in Steampipe provides information about supported resource types for AWS Resource Groups. This table allows DevOps engineers, and other technical professionals to query details about supported resource types, including their names and whether they can be included in a resource group. Users can utilize this table to gather insights on resource types, such as which resources can be grouped together for easy management and configuration. The schema outlines the various attributes of the supported resource type, including the resource type and whether it is groupable.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_resource_explorer_supported_resource_type` table, you can use the `.inspect aws_resource_explorer_supported_resource_type` command in Steampipe.

Key columns:

- `name`: The name of the supported resource type. This is a key column because it can be used to join this table with other tables that contain information about specific resource types in AWS.
- `is_groupable`: Indicates whether the resource type can be included in a resource group. This is a useful column for understanding which resource types can be grouped together for management and configuration purposes.

## Examples

### Basic info

```sql
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type;
```

### List supported IAM resource types

```sql
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type
where
  service = 'iam';
```
