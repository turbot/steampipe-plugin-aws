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

### Filter resources by a single tag key
Find resources that have the `Environment` tag set, regardless of its value.

```sql+postgres
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters = '[{"key": "Environment"}]';
```

```sql+sqlite
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters = '[{"key": "Environment"}]';
```

### Filter resources by one of several tag values
Find resources where the `Environment` tag is set to either `prod` or `staging`.

```sql+postgres
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters = '[{"key": "Environment", "values": ["prod", "staging"]}]';
```

```sql+sqlite
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters = '[{"key": "Environment", "values": ["prod", "staging"]}]';
```

### Filter resources that match multiple tags
Find resources that have the `Environment` tag set to `prod`, and also have an `Owner` tag set (regardless of value). Multiple entries in the array are ANDed together.

```sql+postgres
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters = '[{"key": "Environment", "values": ["prod"]}, {"key": "Owner"}]';
```

```sql+sqlite
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters = '[{"key": "Environment", "values": ["prod"]}, {"key": "Owner"}]';
```

### Filter resources using alternative tag filters (unioned)
Find resources matching either of two independent `tag_filters`, using an `in` list to union the results.

```sql+postgres
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters in (
    '[{"key": "Team", "values": ["platform"]}]',
    '[{"key": "Team", "values": ["data"]}]'
  );
```

```sql+sqlite
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tag_filters in (
    '[{"key": "Team", "values": ["platform"]}]',
    '[{"key": "Team", "values": ["data"]}]'
  );
```