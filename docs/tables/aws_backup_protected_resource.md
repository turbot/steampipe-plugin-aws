---
title: "Steampipe Table: aws_backup_protected_resource - Query AWS Backup Protected Resources using SQL"
description: "Allows users to query AWS Backup Protected Resources to retrieve detailed information about the resources that are backed up by AWS Backup service."
folder: "Backup"
---

# Table: aws_backup_protected_resource - Query AWS Backup Protected Resources using SQL

AWS Backup Protected Resources are the critical data, system configurations, and applications that are safeguarded by AWS Backup. This service provides a fully managed, policy-based backup solution, simplifying the process of backing up data across AWS services. It offers a centralized place to manage backups, audit and monitor activities, and apply retention policies, thus enhancing data protection and compliance.

## Table Usage Guide

The `aws_backup_protected_resource` table in Steampipe provides you with information about the resources that are backed up by AWS Backup service. This table allows you, as a DevOps engineer, security analyst, or system administrator, to query resource-specific details, including resource ARN, type, backup plan ID, and the last backup time. You can utilize this table to gather insights on backed up resources, such as retrieving the last backup time, identifying resources that are not backed up, verifying the backup plan associated with each resource, and more. The schema outlines the various attributes of the backed up resource, including the resource ARN, resource type, backup plan ID, and last backup time for you.

## Examples

### Basic Info
Discover the segments that are protected by AWS Backup service and when they were last backed up. This is useful for maintaining data recovery readiness and ensuring that critical resources are sufficiently protected.

```sql+postgres
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource;
```

```sql+sqlite
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource;
```

### List EBS volumes that are backed up
Determine the areas in which EBS volumes are backed up, allowing you to understand the reach of your backup strategy and ensure no critical data is left unprotected.

```sql+postgres
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource
where
  resource_type = 'EBS';
```

```sql+sqlite
select
  resource_arn,
  resource_type,
  last_backup_time
from
  aws_backup_protected_resource
where
  resource_type = 'EBS';
```