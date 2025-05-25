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

The `aws_tagging_resource` table supports filtering by resource types to help you focus on specific AWS services or resource types. This is particularly useful when you need to:
- Audit tags on specific resource categories
- Generate compliance reports for particular services
- Analyze costs for specific resource types
- Troubleshoot issues with particular AWS services

#### How Resource Type Filtering Works

The filter uses the `resource_types` column, which accepts a JSON array of strings. Each string can be specified in two formats:

| Format | Description | Example | Matches |
|--------|-------------|---------|---------|
| `service` | All resources from a specific AWS service | `"ec2"` | All EC2 resources (instances, volumes, security groups, etc.) |
| `service:resourceType` | Specific resource type within a service | `"ec2:instance"` | Only EC2 instances |

#### Basic Examples

**Filter by a single resource type:**
```sql
select
  name,
  arn,
  resource_types,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:instance"]';
```

**Filter by multiple specific resource types:**
```sql
select
  name,
  arn,
  resource_types,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:instance", "s3:bucket", "rds:db"]';
```

**Filter by entire service (all resource types within a service):**
```sql
select
  name,
  arn,
  resource_types,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2", "s3"]';
```

#### Advanced Examples

**Mix service-level and resource-type-level filters:**
```sql
-- Get all Lambda resources and only RDS database instances
select
  name,
  arn,
  resource_types,
  region
from
  aws_tagging_resource
where
  resource_types = '["lambda", "rds:db"]';
```

**Common resource types for compliance auditing:**
```sql
-- Focus on compute and storage resources
select
  name,
  arn,
  compliance_status,
  tags,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:instance", "ec2:volume", "s3:bucket", "rds:db", "ecs:cluster"]'
  and compliance_status = false;
```

**Network and security resources:**
```sql
-- Audit network-related resources
select
  name,
  arn,
  tags,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:security-group", "ec2:vpc", "ec2:subnet", "elbv2:loadbalancer"]';
```

#### Common Resource Type Values

Here are some frequently used resource type combinations:

| Use Case | Resource Types |
|----------|----------------|
| Compute resources | `["ec2:instance", "lambda:function", "ecs:cluster"]` |
| Storage resources | `["s3:bucket", "ec2:volume", "efs:file-system"]` |
| Database resources | `["rds:db", "rds:cluster", "dynamodb:table"]` |
| Network resources | `["ec2:vpc", "ec2:subnet", "ec2:security-group"]` |
| All EC2 resources | `["ec2"]` |

#### Important Notes

- **JSON Format**: The `resource_types` value must be a valid JSON array of strings, even for a single resource type
- **Case Sensitivity**: Service and resource type names are case-sensitive and should be lowercase
- **Performance**: Filtering by resource types can significantly improve query performance by reducing the number of resources scanned
- **Reference**: For a complete list of valid service names and resource type names, refer to the [AWS Resource Groups Tagging API documentation](https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_GetResources.html#API_GetResources_RequestParameters)
- **Discovery**: Use the [`aws_resource_explorer_supported_resource_type`](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) table to discover available resource types, though note that some supported services may not appear in that API response
