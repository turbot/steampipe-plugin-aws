---
title: "Steampipe Table: aws_inspector_assessment_target - Query AWS Inspector Assessment Targets using SQL"
description: "Allows users to query AWS Inspector Assessment Targets. The `aws_inspector_assessment_target` table in Steampipe provides information about assessment targets within AWS Inspector. This table allows DevOps engineers to query target-specific details, including ARN, name, and associated resource group ARN. Users can utilize this table to gather insights on assessment targets, such as their creation time, last updated time, and more. The schema outlines the various attributes of the assessment target, including the target ARN, creation date, and associated tags."
folder: "Inspector"
---

# Table: aws_inspector_assessment_target - Query AWS Inspector Assessment Targets using SQL

The AWS Inspector Assessment Target is a part of the AWS Inspector service, which is an automated security assessment service that helps improve the security and compliance of applications deployed on AWS. Assessment Targets are specifically used in AWS Inspector to define the Amazon EC2 instances that are included in the assessment. They help in identifying potential security issues, vulnerabilities, or deviations from best practices.

## Table Usage Guide

The `aws_inspector_assessment_target` table in Steampipe provides you with information about assessment targets within AWS Inspector. This table allows you, as a DevOps engineer, to query target-specific details, including ARN, name, and associated resource group ARN. You can utilize this table to gather insights on assessment targets, such as their creation time, last updated time, and more. The schema outlines the various attributes of the assessment target for you, including the target ARN, creation date, and associated tags.

## Examples

### Basic info
Explore which AWS Inspector Assessment Targets are available and when they were created or updated. This query is useful for monitoring changes and managing resources across different AWS regions.

```sql+postgres
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target;
```

```sql+sqlite
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target;
```


### List assessment targets created within the last 7 days
Identify recently created assessment targets within the past week to stay updated on new additions in the AWS Inspector service. This can help in monitoring and managing your resources effectively.

```sql+postgres
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at > (current_date - interval '7' day);
```

```sql+sqlite
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at > date('now','-7 day');
```


### List assessment targets that were updated after creation
Discover the segments that have undergone changes since their initial creation in the AWS Inspector service. This is useful for tracking modifications and ensuring that updates have been successfully implemented.

```sql+postgres
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at != updated_at;
```

```sql+sqlite
select
  name,
  arn,
  resource_group_arn,
  created_at,
  updated_at,
  region
from
  aws_inspector_assessment_target
where
  created_at <> updated_at;
```