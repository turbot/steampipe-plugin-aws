---
title: "Table: aws_tagging_resource - Query AWS Resource Tagging API using SQL"
description: "Allows users to query AWS Resource Tagging API to get details about resources and their associated tags."
---

# Table: aws_tagging_resource - Query AWS Resource Tagging API using SQL

The `aws_tagging_resource` table in Steampipe provides information about resources and their associated tags in AWS. This table allows DevOps engineers to query resource-specific details, including resource ARN, resource type, and associated tags. Users can utilize this table to gather insights on resources, such as resources with specific tags, resources of a certain type, and more. The schema outlines the various attributes of the AWS resource, including the resource ARN, resource type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_tagging_resource` table, you can use the `.inspect aws_tagging_resource` command in Steampipe.

**Key columns**:

- `arn`: This is the Amazon Resource Name (ARN) of the resource. It is a unique identifier for the resource and can be used to join this table with other tables that also contain resource ARNs.
- `resource_type`: This column specifies the type of the resource (like EC2 instance, S3 bucket, etc.). It is useful for filtering resources by type.
- `tags`: This column contains the tags associated with the resource. It is useful for filtering resources based on specific tag values.

## Examples

### Basic info

```sql
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

```sql
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