---
title: "Steampipe Table: aws_backup_recovery_point - Query AWS Backup Recovery Points using SQL"
description: "Allows users to query AWS Backup Recovery Points to gather comprehensive information about each recovery point within an AWS Backup vault."
folder: "Backup"
---

# Table: aws_backup_recovery_point - Query AWS Backup Recovery Points using SQL

The AWS Backup Recovery Point is a component of AWS Backup, a fully managed backup service that makes it easy to centralize and automate the backup of data across AWS services. This resource, the recovery point, is an entity that contains all the metadata that AWS Backup needs to recover a protected resource, such as an Amazon RDS database, an Amazon EBS volume, or an Amazon S3 bucket. The recovery point is created after a successful backup of a protected resource.

## Table Usage Guide

The `aws_backup_recovery_point` table in Steampipe provides you with information about each recovery point within an AWS Backup vault. This table allows you, as a DevOps engineer or system administrator, to query recovery point-specific details, including the backup vault where the recovery point is stored, the source of the backup, the state of the recovery point, and associated metadata. You can utilize this table to gather insights on recovery points, such as identifying unencrypted recovery points, verifying backup completion status, and more. The schema outlines the various attributes of the recovery point for you, including the recovery point ARN, creation date, backup size, and associated tags.

Note: The value in the `tags` column will be populated only if its resource type has a checkmark for [Full AWS Backup management](https://docs.aws.amazon.com/aws-backup/latest/devguide/whatisbackup.html#full-management) as per AWS Backup docs. This means the recovery point ARN must match the pattern `arn:aws:backup:[a-z0-9\-]+:[0-9]{12}:recovery-point:.*`

## Examples

### Basic Info
Discover the segments that are significant in your AWS backup recovery points. This can be beneficial for assessing the status and type of resources within your backup vaults, which can help in managing your backup strategy effectively.

```sql+postgres
select
  backup_vault_name,
  recovery_point_arn,
  resource_type,
  status
from
  aws_backup_recovery_point;
```

```sql+sqlite
select
  backup_vault_name,
  recovery_point_arn,
  resource_type,
  status
from
  aws_backup_recovery_point;
```

### List encrypted recovery points
Identify instances where your recovery points are encrypted to ensure data security and compliance. This query is useful to maintain a secure and compliant data backup system by pinpointing the specific locations where encryption is applied.

```sql+postgres
select
  backup_vault_name,
  recovery_point_arn,
  resource_type,
  status,
  is_encrypted
from
  aws_backup_recovery_point
where
  is_encrypted;
```

```sql+sqlite
select
  backup_vault_name,
  recovery_point_arn,
  resource_type,
  status,
  is_encrypted
from
  aws_backup_recovery_point
where
  is_encrypted = 1;
```

### Get associated tags for the targeted Recovery Points EC2, EBS and S3 resource types
Retrieving metadata, in the form of tags, for recovery points associated with three resource types - EC2 instances, EBS volumes, and S3 buckets. Tags are key-value pairs that provide valuable information about AWS resources.

```sql+postgres
select
  r.backup_vault_name as backup_vault_name,
  r.recovery_point_arn as recovery_point_arn,
  r.resource_type as resource_type,
case
    when r.resource_type = 'EBS' then (
      select tags from aws_ebs_snapshot where arn = concat(
        (string_to_array(r.recovery_point_arn, '::'))[1],
        ':',
        r.account_id,
        ':',
        (string_to_array(r.recovery_point_arn, '::'))[2]
      )
    )
    when r.resource_type = 'EC2' then (
      select tags from aws_ec2_ami where image_id = (string_to_array(r.recovery_point_arn, '::image/'))[2]
    )
    when r.resource_type in ('S3', 'EFS') then r.tags
end as tags,
  r.region,
  r.account_id
from
  aws_backup_recovery_point as r;
```

```sql+sqlite
select
  r.backup_vault_name as backup_vault_name,
  r.recovery_point_arn as recovery_point_arn,
  r.resource_type as resource_type,
  case
    when r.resource_type = 'EBS' then (
      select tags from aws_ebs_snapshot where arn = substr(r.recovery_point_arn, instr(r.recovery_point_arn, '::') + 2)
    )
    when r.resource_type = 'EC2' then (
      select tags from aws_ec2_ami where image_id = substr(r.recovery_point_arn, instr(r.recovery_point_arn, '::image/') + 8)
    )
    when r.resource_type in ('S3', 'EFS') then r.tags
  end as tags,
  r.region,
  r.account_id
from
  aws_backup_recovery_point as r;
```