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
```sql+postgres
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

```sql+sqlite
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
```sql+postgres
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

```sql+sqlite
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
```sql+postgres
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

```sql+sqlite
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
| Compute | `["ec2:instance", "lambda:function", "ecs:cluster", "eks:cluster"]` |
| Storage | `["s3:bucket", "ec2:volume", "elasticfilesystem:file-system"]` |
| Database | `["rds:db", "rds:cluster", "dynamodb:table"]` |
| Network | `["ec2:vpc", "ec2:subnet", "ec2:security-group", "elasticloadbalancing:loadbalancer"]` |
| Security | `["iam:role", "iam:policy", "kms:key"]` |
| Monitoring | `["logs:log-group", "cloudwatch:alarm", "cloudwatch:dashboard"]` |

### API Behavior and Resource Type Discovery

#### Automatic Batching

The AWS Resource Groups Tagging API limits each request to 100 resource type filters. Steampipe handles this limitation transparently:

1. **Automatic splitting**: Resource type lists exceeding 100 items are automatically split into batches
2. **Sequential execution**: Each batch is processed as a separate API request
3. **Result aggregation**: All results are combined and deduplicated by ARN
4. **Seamless streaming**: Results are returned as a single, unified dataset

This means you can query hundreds of resource types without manual batching:

```sql+postgres
-- This works seamlessly even though it exceeds the 100-item API limit
select name, arn, tags
from aws_tagging_resource
where resource_types = '[
  "ec2:instance", "ec2:volume", "ec2:snapshot", "ec2:image", "ec2:security-group",
  "s3:bucket", "lambda:function", "rds:db", "rds:cluster", "dynamodb:table",
  -- ... add as many as needed
]';
```

```sql+sqlite
-- This works seamlessly even though it exceeds the 100-item API limit
select name, arn, tags
from aws_tagging_resource
where resource_types = '[
  "ec2:instance", "ec2:volume", "ec2:snapshot", "ec2:image", "ec2:security-group",
  "s3:bucket", "lambda:function", "rds:db", "rds:cluster", "dynamodb:table",
  -- ... add as many as needed
]';
```

#### Discovering Available Resource Types

To find which resource types you can query, use these approaches:

**1. Query AWS Resource Explorer for supported types:**
```sql+postgres
-- List all available resource types
select
  service,
  resource_type,
  service || ':' || resource_type as full_resource_type
from
  aws_resource_explorer_supported_resource_type
order by
  service, resource_type;
```

```sql+sqlite
-- List all available resource types
select
  service,
  resource_type,
  service || ':' || resource_type as full_resource_type
from
  aws_resource_explorer_supported_resource_type
order by
  service, resource_type;
```

**2. Find resource types for a specific service:**
```sql+postgres
-- Example: Find all EC2 resource types
select
  service,
  resource_type,
  service || ':' || resource_type as full_resource_type
from
  aws_resource_explorer_supported_resource_type
where
  service = 'ec2'
order by
  resource_type;
```

```sql+sqlite
-- Example: Find all EC2 resource types
select
  service,
  resource_type,
  service || ':' || resource_type as full_resource_type
from
  aws_resource_explorer_supported_resource_type
where
  service = 'ec2'
order by
  resource_type;
```

#### Important Notes

- **JSON array format**: Resource types must always be specified as a JSON array, even for single values: `'["ec2:instance"]'`
- **Service vs. resource type**: Use `"ec2"` to query all EC2 resources, or `"ec2:instance"` for specific types
- **Case sensitivity**: Resource type filters are case-sensitive and must match AWS conventions
- **Performance**: While batching is automatic, querying many resource types may take longer due to multiple API calls
- **Regional data**: Results are returned for the region specified in your connection configuration

For the complete list of supported services and resource types, refer to:
- [AWS Resource Groups Tagging API supported services](https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/supported-services.html)
- The [`aws_resource_explorer_supported_resource_type`](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) table in your Steampipe instance
