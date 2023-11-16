---
title: "Table: aws_backup_protected_resource - Query AWS Backup Protected Resources using SQL"
description: "Allows users to query AWS Backup Protected Resources to retrieve detailed information about the resources that are backed up by AWS Backup service."
---

# Table: aws_backup_protected_resource - Query AWS Backup Protected Resources using SQL

The `aws_backup_protected_resource` table in Steampipe provides information about the resources that are backed up by AWS Backup service. This table allows DevOps engineers, security analysts, and system administrators to query resource-specific details, including resource ARN, type, backup plan ID, and the last backup time. Users can utilize this table to gather insights on backed up resources, such as retrieving the last backup time, identifying resources that are not backed up, verifying the backup plan associated with each resource, and more. The schema outlines the various attributes of the backed up resource, including the resource ARN, resource type, backup plan ID, and last backup time.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_protected_resource` table, you can use the `.inspect aws_backup_protected_resource` command in Steampipe.

Key columns:

- `resource_arn`: The Amazon Resource Name (ARN) of the backed up resource. This can be used to join with other tables that contain resource ARN, such as the `aws_ec2_instance` table.
- `resource_type`: The type of the AWS resource. This can be used to filter resources based on their type.
- `backup_plan_id`: The backup plan ID that is used to backup the resource. This can be used to join with the `aws_backup_plan` table to retrieve detailed information about the backup plan.

## Examples

### Basic Info

```sql
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource;
```

### List EBS volumes that are backed up

```sql
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource
where
  resource_type = 'EBS';
```
