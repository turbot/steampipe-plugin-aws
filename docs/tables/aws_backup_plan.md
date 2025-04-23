---
title: "Steampipe Table: aws_backup_plan - Query AWS Backup Plan using SQL"
description: "Allows users to query AWS Backup Plan data, providing detailed information about each backup plan created within an AWS account. Useful for DevOps engineers to monitor and manage backup strategies and ensure data recovery processes are in place."
folder: "Backup"
---

# Table: aws_backup_plan - Query AWS Backup Plan using SQL

The AWS Backup Plan is a policy-based solution for defining, scheduling, and automating the backup activities of AWS resources. It enables you to centralize and automate data protection across AWS services, simplifying management and reducing operational costs. With AWS Backup, you can customize where and how you backup your resources, providing flexibility and control over your data protection strategy.

## Table Usage Guide

The `aws_backup_plan` table in Steampipe provides you with information about each backup plan within AWS Backup. This table allows you, as a DevOps engineer, to query backup plan-specific details, including backup options, creation and version details, and associated metadata. You can utilize this table to gather insights on backup plans, such as the backup frequency, backup window, lifecycle of the backup, and more. The schema outlines the various attributes of the backup plan for you, including the backup plan ARN, creation date, version, and associated tags.

## Examples

### Basic Info
Assess the elements within your AWS backup plans to understand when they were created and when they were last executed. This can help in monitoring and managing your backup strategies effectively.

```sql+postgres
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan;
```

```sql+sqlite
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan;
```

### List plans older than 90 days
Determine the areas in which backup plans have been inactive for more than 90 days. This can aid in identifying outdated or potentially unnecessary backup plans, facilitating better resource management.

```sql+postgres
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan
where
  creation_date <= (current_date - interval '90' day)
order by
  creation_date;
```

```sql+sqlite
select
  name,
  backup_plan_id,
  arn,
  creation_date,
  last_execution_date
from
  aws_backup_plan
where
  creation_date <= date('now','-90 day')
order by
  creation_date;
```

### List plans that were deleted in the last 7 days
Determine the areas in which backup plans were recently removed within the AWS environment to keep track of changes and maintain security standards.

```sql+postgres
select
  name,
  arn,
  creation_date,
  deletion_date
from
  aws_backup_plan
where
  deletion_date > current_date - 7
order by
  deletion_date;
```

```sql+sqlite
select
  name,
  arn,
  creation_date,
  deletion_date
from
  aws_backup_plan
where
  deletion_date > date('now','-7 day')
order by
  deletion_date;
```