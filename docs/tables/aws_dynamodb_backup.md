# Table: aws_dynamodb_backup

DynamoDB backup is used to create full backups of  tables for long-term retention and archival for regulatory compliance needs.

## Examples

### List backups with their corresponding tables

```sql
select
  name,
  table_name,
  table_id
from
  aws_dynamodb_backup;
```


### Basic backup info

```sql
select
  name,
  backup_status,
  backup_type,
  backup_expiry_datetime,
  backup_size_bytes
from
  aws_dynamodb_backup;
```