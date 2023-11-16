---
title: "Table: aws_secretsmanager_secret - Query AWS Secrets Manager Secret using SQL"
description: "Allows users to query AWS Secrets Manager Secret data, including metadata, versions, rotation configuration, and more."
---

# Table: aws_secretsmanager_secret - Query AWS Secrets Manager Secret using SQL

The `aws_secretsmanager_secret` table in Steampipe provides information about secrets within AWS Secrets Manager. This table allows DevOps engineers to query secret-specific details, including metadata, versions, rotation configuration, and more. Users can utilize this table to gather insights on secrets, such as secret rotation status, associated resource policies, and more. The schema outlines the various attributes of the secret, including the secret ARN, name, description, rotation rules, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_secretsmanager_secret` table, you can use the `.inspect aws_secretsmanager_secret` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the secret. This is a unique identifier that is used to join this table with other tables.
- `name`: The friendly name of the secret. This can be used to join with other tables that also contain secret names.
- `rotation_enabled`: Indicates whether automatic rotation is enabled for the secret. This is useful for understanding the security posture of your secrets.

## Examples

### Basic info

```sql
select
  name,
  created_date,
  description,
  last_accessed_date
from
  aws_secretsmanager_secret;
```


### List secrets that do not automatically rotate

```sql
select
  name,
  created_date,
  description,
  rotation_enabled
from
  aws_secretsmanager_secret
where
  not rotation_enabled;
```


### List secrets that automatically rotate every 7 days

```sql
select
  name,
  created_date,
  description,
  rotation_enabled,
  rotation_rules
from
  aws_secretsmanager_secret
where
  rotation_rules -> 'AutomaticallyAfterDays' > '7';
```


### List secrets that are not replicated in other regions

```sql
select
  name,
  created_date,
  description,
  replication_status
from
  aws_secretsmanager_secret
where
  replication_status is null;
```

### List policy details for the secrets

```sql
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_secretsmanager_secret;
```
