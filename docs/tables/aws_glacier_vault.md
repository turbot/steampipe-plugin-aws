# Table: aws_glacier_vault

A vault is a way to group archives together in Amazon S3 Glacier.

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
