---
title: "Steampipe Table: aws_resource_explorer_resource - Query AWS Resource Explorer Resources using SQL"
description: "Allows users to query AWS Resource Explorer Resources, providing comprehensive information about AWS resources across regions in your account."
folder: "Resource Explorer"
---

# Table: aws_resource_explorer_resource - Query AWS Resource Explorer Resources using SQL

AWS Resource Explorer is a resource search and discovery service that helps you find and explore resources across your AWS account. It provides an internet search engine-like experience for discovering and exploring AWS resources like EC2 instances, S3 buckets, Lambda functions, and more. Resource Explorer maintains an index of your resources and their metadata, enabling quick and efficient searches.

## Table Usage Guide

The `aws_resource_explorer_resource` table provides insights into resources indexed by AWS Resource Explorer. As a cloud administrator or DevOps engineer, you can use this table to explore and discover resources across your AWS account, helping with tasks such as:

**Important Notes**
- For improved performance, it is advised that you use the optional qual `filter` to limit the result set to a specific time period. For information about the supported syntax, see [Search query reference for Resource Explorer](https://docs.aws.amazon.com/resource-explorer/latest/userguide/using-search-query-syntax.html) in the AWS Resource Explorer User Guide.
- This table supports optional quals. Queries with optional quals are optimised. Optional quals are supported for the following columns:
  - `filter`
  - `view_arn`

## Examples

### Basic info
Get a simple overview of resources with essential fields.

```sql+postgres
select
  arn,
  resource_type,
  service,
  region,
  title
from
  aws_resource_explorer_resource
limit 10;
```

```sql+sqlite
select
  arn,
  resource_type,
  service,
  region,
  title
from
  aws_resource_explorer_resource
limit 10;
```

### List all EC2 instances
Find all EC2 instances across regions in your account.

```sql+postgres
select
  arn,
  region,
  title,
  properties
from
  aws_resource_explorer_resource
where
  resource_type = 'AWS::EC2::Instance'
order by
  region;
```

```sql+sqlite
select
  arn,
  region,
  title,
  properties
from
  aws_resource_explorer_resource
where
  resource_type = 'AWS::EC2::Instance'
order by
  region;
```

### Count resources by type and region
Get a summary of resource distribution across types and regions.

```sql+postgres
select
  resource_type,
  region,
  count(*) as resource_count
from
  aws_resource_explorer_resource
group by
  resource_type,
  region
order by
  resource_count desc;
```

```sql+sqlite
select
  resource_type,
  region,
  count(*) as resource_count
from
  aws_resource_explorer_resource
group by
  resource_type,
  region
order by
  resource_count desc;
```

### List all EC2 resources
Find all EC2 resources across regions in your account.

```sql+postgres
select
  arn,
  region,
  title,
  properties
from
  aws_resource_explorer_resource
where
  filter = 'service:EC2'
order by
  region;
```

```sql+sqlite
select
  arn,
  region,
  title,
  properties
from
  aws_resource_explorer_resource
where
  filter = 'service:EC2'
order by
  region;
```