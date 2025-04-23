---
title: "Steampipe Table: aws_glacier_vault - Query AWS Glacier Vaults using SQL"
description: "Allows users to query AWS Glacier Vaults for detailed information on each vault, including the vault's name, ARN, creation date, number of archives, size of archives, and more."
folder: "Glacier"
---

# Table: aws_glacier_vault - Query AWS Glacier Vaults using SQL

AWS Glacier Vaults are a component of the Amazon Glacier service, designed for long-term, secure and durable storage of data for archiving and backup purposes. They provide an extremely low-cost storage solution that ensures data is kept safe for extended periods of time. Vaults also allow for the control of access through the use of resource-based policies.

## Table Usage Guide

The `aws_glacier_vault` table in Steampipe provides you with information about Vaults within AWS Glacier. This table allows you, as a DevOps engineer, to query vault-specific details, including vault names, ARNs, creation dates, number of archives, size of archives, and more. You can utilize this table to gather insights on vaults, such as the total size of all archives in the vault, the number of archives in the vault, and the date the vault was last accessed. The schema outlines the various attributes of the Glacier Vault for you, including the vault ARN, creation date, last inventory date, number of archives, size of archives, and associated tags.

## Examples

### Basic info
Explore the historical data and storage size of your AWS Glacier Vaults. This query is particularly useful for tracking the growth and usage of your vaults over time.

```sql+postgres
select
  vault_name,
  creation_date,
  last_inventory_date,
  number_of_archives,
  size_in_bytes
from
  aws_glacier_vault;
```

```sql+sqlite
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
Determine the areas in which full access to the resource is granted through vaults. This is useful for identifying potential security risks and ensuring appropriate access controls are in place.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(principal.value, '$') as principal,
  json_extract(action.value, '$') as action,
  json_extract(statement.value, '$.Effect') as effect,
  json_extract(statement.value, '$.Condition') as conditions
from
  aws_glacier_vault,
  json_each(policy_std, '$.Statement') as statement,
  json_each(json_extract(statement.value, '$.Principal.AWS')) as principal,
  json_each(json_extract(statement.value, '$.Action')) as action
where
  json_extract(statement.value, '$.Effect') = 'Allow'
  and (
    json_extract(action.value, '$') = '*'
    or json_extract(action.value, '$') = 'glacier:*'
  );
```

### List vaults that grant anonymous access to the resource
Determine the areas in which your data vaults may be vulnerable by identifying any instances that allow anonymous access. This is particularly useful for enhancing security measures and ensuring data privacy.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(p.value, '$') as principal,
  json_extract(a.value, '$') as action,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Condition') as conditions
from
  aws_glacier_vault,
  json_each(policy_std, '$.Statement') as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a
where
  p.value = '*'
  and json_extract(s.value, '$.Effect') = 'Allow';
```


### Get the archival age in days before deletion for each vault
This query is used to identify the number of days before each AWS Glacier vault is scheduled for deletion. This helps in managing data lifecycle and ensuring timely archival or deletion of data to optimize storage costs.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(action.value, '$') as action,
  json_extract(statement.value, '$.Effect') as effect,
  json_extract(statement.value, '$.Condition.NumericLessThan."glacier:archiveageindays"') as archive_age_in_days
from
  aws_glacier_vault,
  json_each(vault_lock_policy_std, '$.Statement') as statement,
  json_each(json_extract(statement.value, '$.Action')) as action;
```


### List vaults without owner tag key
Identify instances where AWS Glacier vaults lack an 'owner' tag. This can help in managing and organizing your resources effectively by ensuring every vault has an owner assigned.

```sql+postgres
select
  vault_name,
  tags
from
  aws_glacier_vault
where
  not tags :: JSONB ? 'owner';
```

```sql+sqlite
select
  vault_name,
  tags
from
  aws_glacier_vault
where
  json_extract(tags, '$.owner') is null;
```

### List vaults with notifications enabled
Discover the segments that have enabled notifications within your vaults. This is particularly useful for keeping track of important events and updates in real-time.

```sql+postgres
select
  vault_name,
  vault_notification_config ->> 'SNSTopic' as sns_topic,
  vault_notification_config ->> 'Events' as notification_events
from
  aws_glacier_vault
where
  vault_notification_config is not null;
```

```sql+sqlite
select
  vault_name,
  json_extract(vault_notification_config, '$.SNSTopic') as sns_topic,
  json_extract(vault_notification_config, '$.Events') as notification_events
from
  aws_glacier_vault
where
  vault_notification_config is not null;
```