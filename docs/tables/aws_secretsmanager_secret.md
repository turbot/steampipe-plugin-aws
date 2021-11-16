# Table: aws_secretsmanager_secret

The AWS Secrets Manager Secret resource creates a secret and stores it in Secrets Manager.

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
