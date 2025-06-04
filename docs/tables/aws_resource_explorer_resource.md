---
title: "Steampipe Table: aws_resource_explorer_resource - Query AWS Resource Explorer Resources using SQL"
description: "Allows users to query AWS Resource Explorer Resources, providing comprehensive information about AWS resources across regions in your account."
folder: "Resource Explorer"
---

# Table: aws_resource_explorer_resource - Query AWS Resource Explorer Resources using SQL

AWS Resource Explorer is a resource search and discovery service that helps you find and explore resources across your AWS account. It provides an internet search engine-like experience for discovering and exploring AWS resources like EC2 instances, S3 buckets, Lambda functions, and more. Resource Explorer maintains an index of your resources and their metadata, enabling quick and efficient searches.

## Table Usage Guide

The `aws_resource_explorer_resource` table provides insights into resources indexed by AWS Resource Explorer. As a cloud administrator or DevOps engineer, you can use this table to explore and discover resources across your AWS account, helping with tasks such as:

- Identifying resources by type, region, or tags
- Auditing resource metadata and properties
- Finding resources with specific configurations or attributes
- Cross-region resource discovery and management

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
  properties,
  tags
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
  properties,
  tags
from
  aws_resource_explorer_resource
where
  resource_type = 'AWS::EC2::Instance'
order by
  region;
```

### Find resources by tag

Locate resources that have specific tags.

```sql+postgres
select
  arn,
  resource_type,
  region,
  title,
  tags
from
  aws_resource_explorer_resource
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  arn,
  resource_type,
  region,
  title,
  tags
from
  aws_resource_explorer_resource
where
  json_extract(tags, '$.Environment') = 'Production';
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
