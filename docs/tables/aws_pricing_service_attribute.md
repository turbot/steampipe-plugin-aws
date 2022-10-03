# Table: aws_pricing_service_attribute

Returns the metadata for one service or a list of the metadata for all services.

## Examples

### Basic info

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute;
```

### List attribute details of AWS Backup service

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup';
```

### List supported attribute values of AWS Backup service for termType attribute

```sql
select
  service_code,
  attribute_name,
  attribute_values
from
  aws_pricing_service_attribute
where
  service_code = 'AWSBackup' and attribute_name = 'termType';
```
