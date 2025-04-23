---
title: "Steampipe Table: aws_kms_alias - Query AWS Key Management Service (KMS) alias using SQL"
description: "Allows users to query AWS KMS aliases and retrieve information about their associated keys, including the key ID, alias name, and alias ARN."
folder: "KMS"
---

# Table: aws_kms_alias - Query AWS Key Management Service (KMS) alias using SQL

The AWS Key Management Service (KMS) alias is a user-friendly identifier for a KMS key. These aliases allow you to simplify cryptographic workflows by referring to a key by a consistent name throughout its lifecycle. KMS aliases can be used to manage cryptographic keys, enabling secure access to services and applications.

## Table Usage Guide

The `aws_kms_alias` table in Steampipe provides you with information about aliases within AWS Key Management Service (KMS). This table allows you, as a DevOps engineer, to query alias-specific details, including the alias name, alias ARN, and the key it is associated with. You can utilize this table to gather insights on aliases, such as the keys they are associated with and the ARNs of the aliases. The schema outlines the various attributes of the KMS alias for you, including the alias name, alias ARN, and associated key ID.

## Examples

### Basic info
Discover the segments that have been created within your AWS Key Management Service (KMS), including their unique identifiers and creation dates. This can help in identifying and managing your encryption keys, ensuring they are correctly configured and up to date.

```sql+postgres
select
  alias_name,
  title,
  arn,
  target_key_id,
  creation_date
from
  aws_kms_alias;
```

```sql+sqlite
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
Discover the segments where key rotation is disabled in AWS Key Management Service. This is useful in identifying potential security risks, as disabling key rotation can make cryptographic keys more susceptible to compromise.

```sql+postgres
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

```sql+sqlite
select
  k.id as key_id,
  k.key_rotation_enabled as key_rotation_enabled,
  a.alias_name as alias_name,
  a.arn as alias_arn
from
  aws_kms_key as k,
  aws_kms_alias as a
where
  k.id = a.target_key_id and key_rotation_enabled = 0;
```

### List of KMS Customer Managed key alias that is scheduled for deletion
Determine the areas in which the AWS Key Management Service (KMS) has scheduled customer-managed keys for deletion. This allows you to proactively manage your encryption keys and mitigate potential security risks.

```sql+postgres
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

```sql+sqlite
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
Determine the number of aliases associated with each unique key to understand the utilization and management of keys within your AWS Key Management Service.

```sql+postgres
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

```sql+sqlite
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