---
title: "Table: aws_kms_alias - Query AWS Key Management Service (KMS) alias using SQL"
description: "Allows users to query AWS KMS aliases and retrieve information about their associated keys, including the key ID, alias name, and alias ARN."
---

# Table: aws_kms_alias - Query AWS Key Management Service (KMS) alias using SQL

The `aws_kms_alias` table in Steampipe provides information about aliases within AWS Key Management Service (KMS). This table allows DevOps engineers to query alias-specific details, including the alias name, alias ARN, and the key it is associated with. Users can utilize this table to gather insights on aliases, such as the keys they are associated with and the ARNs of the aliases. The schema outlines the various attributes of the KMS alias, including the alias name, alias ARN, and associated key ID.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kms_alias` table, you can use the `.inspect aws_kms_alias` command in Steampipe.

*Key columns:*

- `alias_name`: This is the alias name. It is a unique identifier for the alias and can be used to join this table with other tables.
- `alias_arn`: This is the Amazon Resource Name (ARN) of the alias. It provides a unique identifier for the alias across all AWS accounts and can be used for joining with other tables that contain alias ARNs.
- `target_key_id`: This is the ID of the key associated with the alias. It is useful for joining with other tables that contain key IDs, allowing for queries that return information about the keys associated with aliases.

## Examples

### Basic info

```sql
select
  alias_name,
  title,
  arn,
  target_key_id,
  creation_date
from
  aws_kms_alias;
```

### List of KMS key alias where key rotation disabled on the key

```sql
select
  k.id as key_id,
  k.key_rotation_enabled as key_rotation_enabled,
  a.alias_name as alias_name,
  a.arn as alias_arn
from
  aws_kms_key as k,
  aws_kms_alias as a
where
  k.id = a.target_key_id and not key_rotation_enabled;
```

### List of KMS Customer Managed key alias that is scheduled for deletion

```sql
select
  a.alias_name as alias_name,
  k.id as key_id,
  k.key_state as key_state,
  k.deletion_date as key_deletion_date
from
  aws_kms_key as k,
  aws_kms_alias as a
where
  k.id = a.target_key_id and key_state = 'PendingDeletion';
```

### Count of aliases by key id

```sql
select
  k.id as key_id,
  count(a.alias_name) as count
from
  aws_kms_key as k
  left join aws_kms_alias as a
  on k.id = a.target_key_id
group by
  key_id;
```
