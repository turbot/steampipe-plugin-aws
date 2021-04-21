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


### List secrets whose automatic rotation is not enabled

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


### List secrets whose automatic rotation interval is more than 7 days

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


### List secrets which are not replicated in other regions

```sql
select
  name,
  created_date,
  description,
  replication_status
from
  aws_secretsmanager_secret
where
  replication_status is NULL;
```