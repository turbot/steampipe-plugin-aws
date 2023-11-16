---
title: "Table: aws_kms_key - Query AWS KMS Key using SQL"
description: "Allows users to query AWS KMS Key data including cryptographic details, key usage, key state, and associated metadata."
---

# Table: aws_kms_key - Query AWS KMS Key using SQL

The `aws_kms_key` table in Steampipe provides information about Key Management Service (KMS) keys within AWS. This table allows DevOps engineers to query key-specific details, including cryptographic details, key usage, key state, and associated metadata. Users can utilize this table to gather insights on keys, such as keys rotation status, key type, key state, and more. The schema outlines the various attributes of the KMS key, including the key ARN, creation date, key state, key usage, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kms_key` table, you can use the `.inspect aws_kms_key` command in Steampipe.

**Key columns**:

- `key_id`: The unique identifier for the KMS key. It can be used to join this table with other tables.
- `arn`: The Amazon Resource Name (ARN) for the KMS key. It provides a unique identifier for the key across all of AWS.
- `key_state`: The state of the KMS key. It can be used to filter keys based on their current state.

## Examples

### Basic info

```sql
select
  id,
  title,
  arn,
  key_manager,
  creation_date
from
  aws_kms_key;
```

### List of KMS keys where key rotation is not enabled

```sql
select
  id,
  key_rotation_enabled
from
  aws_kms_key
where
  not key_rotation_enabled;
```


### List of KMS Customer Managed keys scheduled for deletion

```sql
select
  id,
  key_state,
  deletion_date
from
  aws_kms_key
where
  key_state = 'PendingDeletion';
```


### List of unused Customer Managed Keys

```sql
select
  id,
  enabled as key_enabled
from
  aws_kms_key
where
  not enabled;
```


### Count of AWS KMS keys by Key manager

```sql
select
  key_manager,
  count(key_manager) as count
from
  aws_kms_key
group by
  key_manager;
```