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

You can filter the results to retrieve only resources of specific AWS services or particular resource types within those services. This is useful for focusing your queries on particular categories of AWS resources and improving the relevance of your results.

The filter uses the `resource_types` column, which accepts a JSON array of strings. Each string should be in the format:
- `service` — to match all resources of a service (e.g., `"ec2"`)
- `service:resourceType` — to match a specific resource type in a service (e.g., `"ec2:instance"`)

**Examples:**

- To retrieve all EC2 instances, S3 buckets, and any AWS Audit Manager resource:

```sql
select
  name,
  arn,
  resource_types,
  region
from
  aws_tagging_resource
where
  resource_types = '["ec2:instance", "s3:bucket", "auditmanager"]';
```

- To filter by a single resource type, such as EC2 instances:

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

**Note:**
- For a complete list of valid service names and resource type names, refer to the [AWS documentation](https://docs.aws.amazon.com/resourcegroupstagging/latest/APIReference/API_GetResources.html#API_GetResources_RequestParameters). The [`aws_resource_explorer_supported_resource_type`](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_resource_explorer_supported_resource_type) can help provide a set of valid input values, but a couple of AWS services that are actually supported are not returned from that API.
