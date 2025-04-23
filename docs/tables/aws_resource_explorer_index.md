---
title: "Steampipe Table: aws_resource_explorer_index - Query AWS Resource Explorer Index using SQL"
description: "Allows users to query AWS Resource Explorer Index, providing a comprehensive view of all resources across different AWS services in a single table."
folder: "Resource Explorer"
---

# Table: aws_resource_explorer_index - Query AWS Resource Explorer Index using SQL

The AWS Resource Explorer Index is a tool within the AWS Management Console that allows you to view and navigate through all of your AWS resources. It provides a unified, searchable interface to find and manage all of your resources. This service enhances visibility and control over your AWS environment, ensuring you can easily locate, monitor, and manage your resources.

## Table Usage Guide

The `aws_resource_explorer_index` table in Steampipe provides you with information about all resources across different AWS services. This table allows you, as a DevOps engineer, to query resource-specific details, including the resource type, ARN, region, and associated metadata. You can utilize this table to gather insights on resources, such as resource distribution across services, regions, and resource types. The schema outlines the various attributes of the resource, including the resource id, service, type, region, and account for you.

## Examples

### Basic info
Explore which resources are available in your AWS environment and where they are located. This can help in managing resources effectively and optimizing regional distribution.

```sql+postgres
select
  arn,
  region,
  type
from
  aws_resource_explorer_index;
```

```sql+sqlite
select
  arn,
  region,
  type
from
  aws_resource_explorer_index;
```

### Get the details for the aggregator index
Identify instances where the type of resource index in AWS is an aggregator. This can be useful to understand where and how these specific resource indexes are being utilized across different regions.

```sql+postgres
select
  arn,
  region,
  type
from
  aws_resource_explorer_index
where
  type = 'AGGREGATOR';
```

```sql+sqlite
select
  arn,
  region,
  type
from
  aws_resource_explorer_index
where
  type = 'AGGREGATOR';
```