# Table: aws_secrets_manager_secret

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
  aws_secrets_manager_secret;
```


### List secrets whose automatic rotation is not enabled

```sql
select
  name,
  created_date,
  description,
  rotation_enabled
from
  aws_secrets_manager_secret
where
  rotation_enabled = 'false';
```


### List secrets whose automatic rotation interval is greater tha 7 days

```sql
select
  name,
  created_date,
  description,
  rotation_enabled,
  rotation_rules
from
  aws_secrets_manager_secret
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
  aws_secrets_manager_secret
where
  replication_status is NULL;
```