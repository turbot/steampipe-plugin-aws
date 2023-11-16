---
title: "Table: aws_glacier_vault - Query AWS Glacier Vaults using SQL"
description: "Allows users to query AWS Glacier Vaults for detailed information on each vault, including the vault's name, ARN, creation date, number of archives, size of archives, and more."
---

# Table: aws_glacier_vault - Query AWS Glacier Vaults using SQL

The `aws_glacier_vault` table in Steampipe provides information about Vaults within AWS Glacier. This table allows DevOps engineers to query vault-specific details, including vault names, ARNs, creation dates, number of archives, size of archives, and more. Users can utilize this table to gather insights on vaults, such as the total size of all archives in the vault, the number of archives in the vault, and the date the vault was last accessed. The schema outlines the various attributes of the Glacier Vault, including the vault ARN, creation date, last inventory date, number of archives, size of archives, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glacier_vault` table, you can use the `.inspect aws_glacier_vault` command in Steampipe.

**Key columns**:

- `vault_name`: The name of the vault. This can be used to join with other tables that contain vault names.
- `vault_arn`: The Amazon Resource Name (ARN) for the vault. This can be used to join with other tables that contain vault ARNs.
- `account_id`: The AWS account ID associated with the vault. This can be used to join with other tables that contain AWS account IDs.

## Examples

### Basic info

```sql
select
  vault_name,
  creation_date,
  last_inventory_date,
  number_of_archives,
  size_in_bytes
from
  aws_glacier_vault;
```


### List vaults that grant full access to the resource

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_glacier_vault,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and a in ('*', 'glacier:*');
```


### List vaults that grant anonymous access to the resource

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_glacier_vault,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  p = '*'
  and s ->> 'Effect' = 'Allow';
```


### Get the archival age in days before deletion for each vault

```sql
select
  title,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' -> 'NumericLessThan' ->> 'glacier:archiveageindays' as archive_age_in_days
from
  aws_glacier_vault,
  jsonb_array_elements(vault_lock_policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Action') as a;
```


### List vaults without owner tag key

```sql
select
  vault_name,
  tags
from
  aws_glacier_vault
where
  not tags :: JSONB ? 'owner';
```

### List vaults with notifications enabled

```sql
select
  vault_name,
  vault_notification_config ->> 'SNSTopic' as sns_topic,
  vault_notification_config ->> 'Events' as notification_events
from
  aws_glacier_vault
where
  vault_notification_config is not null;
```
