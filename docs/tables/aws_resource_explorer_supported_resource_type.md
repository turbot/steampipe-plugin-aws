---
title: "Steampipe Table: aws_resource_explorer_supported_resource_type - Query AWS Resource Explorer Supported Resource Types using SQL"
description: "Allows users to query AWS Resource Explorer Supported Resource Types to obtain details about supported resource types for AWS Resource Groups."
folder: "Resource Explorer"
---

# Table: aws_resource_explorer_supported_resource_type - Query AWS Resource Explorer Supported Resource Types using SQL

The AWS Resource Explorer Supported Resource Types is a feature of AWS Resource Groups that allows you to view and manage your AWS resources. This service organizes your resources based on their types and the AWS services that they belong to. By using SQL queries, you can easily locate and manage your resources, allowing for efficient resource management.

## Table Usage Guide

The `aws_resource_explorer_supported_resource_type` table in Steampipe provides you with information about supported resource types for AWS Resource Groups. This table allows you, as a DevOps engineer or other technical professional, to query details about supported resource types, including their names and whether they can be included in a resource group. You can utilize this table to gather insights on resource types, such as which resources can be grouped together for easy management and configuration. The schema outlines the various attributes of the supported resource type for you, including the resource type and whether it is groupable.

## Examples

### Basic info
Determine the areas in which specific services and resources types are supported within AWS. This can help in understanding the availability and compatibility of various services across different resource types.

```sql+postgres
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type;
```

```sql+sqlite
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type;
```

### List supported IAM resource types
Determine the areas in which IAM resources are supported within a specific service, allowing for more efficient management and allocation of resources.

```sql+postgres
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type
where
  service = 'iam';
```

```sql+sqlite
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type
where
  service = 'iam';
```