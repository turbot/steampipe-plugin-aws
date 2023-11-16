---
title: "Table: aws_resource_explorer_index - Query AWS Resource Explorer Index using SQL"
description: "Allows users to query AWS Resource Explorer Index, providing a comprehensive view of all resources across different AWS services in a single table."
---

# Table: aws_resource_explorer_index - Query AWS Resource Explorer Index using SQL

The `aws_resource_explorer_index` table in Steampipe provides information about all resources across different AWS services. This table allows DevOps engineers to query resource-specific details, including the resource type, ARN, region, and associated metadata. Users can utilize this table to gather insights on resources, such as resource distribution across services, regions, and resource types. The schema outlines the various attributes of the resource, including the resource id, service, type, region, and account.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_resource_explorer_index` table, you can use the `.inspect aws_resource_explorer_index` command in Steampipe.

**Key columns**:

- `id`: This is the unique identifier of the resource. It can be used to join this table with other tables to fetch detailed information about specific resources.
- `service`: This column represents the AWS service that the resource belongs to. It can be used to filter resources based on the service.
- `type`: This column indicates the type of the resource. It can be used to filter resources based on their type.

## Examples

### Basic info

```sql
select
  arn,
  region,
  type
from
  aws_resource_explorer_index;
```

### Get the details for the aggregator index

```sql
select
  arn,
  region,
  type
from
  aws_resource_explorer_index
where
  type = 'AGGREGATOR';
```
