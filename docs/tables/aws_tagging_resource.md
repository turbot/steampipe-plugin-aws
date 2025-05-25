---
title: "Steampipe Table: aws_tagging_resource - Query AWS Resource Tagging API using SQL"
description: "Allows users to query AWS Resource Tagging API to get details about resources and their associated tags."
folder: "Resource Tagging"
---

# Table: aws_tagging_resource - Query AWS Resource Tagging API using SQL

The AWS Resource Tagging API allows you to manage tags for AWS resources. It provides a uniform way to categorize resources by purpose, owner, environment, or other criteria. Using this API, you can apply tags to your AWS resources, work with the resource groups you create, and more.

## Table Usage Guide

The `aws_tagging_resource` table in Steampipe provides you with information about resources and their associated tags in AWS. This table allows you, as a DevOps engineer, to query resource-specific details, including resource ARN, resource type, and associated tags. You can utilize this table to gather insights on resources, such as resources with specific tags, resources of a certain type, and more. The schema outlines the various attributes of the AWS resource for you, including the resource ARN, resource type, and associated tags.

## Examples

### Basic info
Explore which resources are compliant and their associated tags within a specific region. This will help in managing and organizing resources effectively, ensuring compliance and efficient resource allocation.

```sql+postgres
select
  name,
  arn,
  compliance_status,
  tags,
  region
from
  aws_tagging_resource;
```

```sql+sqlite
select
  name,
  arn,
  compliance_status,
  tags,
  region
from
  aws_tagging_resource;
```

### List resources which are compliant with effective tag policy
Explore which resources are adhering to the effective tag policy. This is useful in maintaining consistency in resource tagging, aiding in cost tracking, resource organization, and access control.

```sql+postgres
select
  name,
  arn,
  tags,
  compliance_status
from
  aws_tagging_resource
where
  compliance_status;
```

```sql+sqlite
select
  name,
  arn,
  tags,
  compliance_status
from
  aws_tagging_resource
where
  compliance_status is not null;
```

### Filter Resources by Resource Types

Filter results to retrieve only resources from specific AWS services or resource types. The `resource_types` column accepts a JSON array of strings in two formats:

- `service` — All resources from a service (e.g., `"ec2"`)
- `service:resourceType` — Specific resource type (e.g., `"ec2:instance"`)

#### Examples

**Get tags for EC2 instances only:**
```sql
select
  name,
  arn,
  tags,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:instance"]';
```

**Get tags for multiple resource types:**
```sql
select
  name,
  arn,
  tags,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:instance", "s3:bucket", "rds:db"]';
```

**Get tags for all resources in specific services:**
```sql
select
  name,
  arn,
  tags,
  region
from
  aws_tagging_resource
where
  resource_types = '["lambda", "dynamodb"]';
```

#### Common Resource Types

| Category | Resource Types |
|----------|----------------|
| Compute | `["ec2:instance", "lambda:function", "ecs:cluster"]` |
| Storage | `["s3:bucket", "ec2:volume", "efs:file-system"]` |
| Database | `["rds:db", "rds:cluster", "dynamodb:table"]` |
| Network | `["ec2:vpc", "ec2:security-group", "elbv2:loadbalancer"]` |

**Notes:**
- Resource types must be specified as a JSON array, even for single values
- Service and resource type names are case-sensitive and lowercase
- For a complete list, see the [AWS Resource Groups Tagging API documentation](https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_GetResources.html#API_GetResources_RequestParameters)
