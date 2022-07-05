# Table: aws_service_attribute

Returns the metadata for one service or a list of the metadata for all services.

## Examples

### Basic info

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_service_attribute;
```

### List attribute details of AWS Backup service

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_service_attribute
where
  service_code = 'AWSBackup';
```
