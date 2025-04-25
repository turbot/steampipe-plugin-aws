---
title: "Steampipe Table: aws_backup_vault - Query AWS Backup Vaults using SQL"
description: "Allows users to query AWS Backup Vaults, providing detailed information about each backup vault, including its name, ARN, recovery points, and more."
folder: "Backup"
---

# Table: aws_backup_vault - Query AWS Backup Vaults using SQL

The AWS Backup Vault is a secured place where AWS Backup stores backup data. It provides a scalable, fully managed, policy-based resource for managing and protecting data across AWS services. It is designed to simplify data protection, enable regulatory compliance, and save costs by eliminating the need to create and manage custom scripts and manual processes.

## Table Usage Guide

The `aws_backup_vault` table in Steampipe provides you with information about backup vaults within AWS Backup. This table allows you, as a DevOps engineer, to query vault-specific details, including the vault name, ARN, number of recovery points, and associated metadata. You can utilize this table to gather insights on backup vaults, such as the number of recovery points for each vault, the creation date of each vault, and more. The schema outlines the various attributes of the backup vault for you, including the vault name, ARN, creation date, last resource backup time, and associated tags.

## Examples

### Basic Info
Uncover the details of your AWS backup vaults, including their names, unique identifiers, and the dates they were created. This can be particularly useful for auditing purposes, allowing you to keep track of your resources and their creation timelines.

```sql+postgres
select
  name,
  arn,
  creation_date
from
  aws_backup_vault;
```

```sql+sqlite
select
  name,
  arn,
  creation_date
from
  aws_backup_vault;
```

### List vaults older than 90 days
Identify backup vaults that have been established for over 90 days. This can be beneficial in assessing long-standing storage resources that may require maintenance or review.

```sql+postgres
select
  name,
  arn,
  creation_date
from
  aws_backup_vault
where
  creation_date <= (current_date - interval '90' day)
order by
  creation_date;
```

```sql+sqlite
select
  name,
  arn,
  creation_date
from
  aws_backup_vault
where
  creation_date <= date('now','-90 day')
order by
  creation_date;
```

### List vaults that do not prevent the deletion of backups in the backup vault
Determine the areas in which your backup vaults may be at risk, specifically those that do not have policies in place to prevent the deletion of backups. This query is useful in identifying potential vulnerabilities and ensuring the safety of your data.

```sql+postgres
select
  name
from
  aws_backup_vault,
  jsonb_array_elements(policy -> 'Statement') as s
where
  s ->> 'Principal' = '*'
  and s ->> 'Effect' != 'Deny'
  and s ->> 'Action' like '%DeleteBackupVault%';
```

```sql+sqlite
select
  name
from
  aws_backup_vault
where
  json_extract(policy, '$.Statement[*].Principal') = '*'
  and json_extract(policy, '$.Statement[*].Effect') != 'Deny'
  and json_extract(policy, '$.Statement[*].Action') like '%DeleteBackupVault%';
```

### List policy details for backup vaults
Determine the areas in which your AWS backup vault policies are applied. This helps in understanding the security measures in place for your backup vaults, assisting in maintaining data integrity and safety.

```sql+postgres
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_backup_vault;
```

```sql+sqlite
select
  name,
  policy,
  policy_std
from
  aws_backup_vault;
```