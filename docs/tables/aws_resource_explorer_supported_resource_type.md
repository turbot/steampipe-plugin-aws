# Table: aws_resource_explorer_supported_resource_type

This table retrieves all resource types currently supported by AWS Resource Explorer.

## Examples

### Basic info

```sql
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type;
```

### List iam service resource types supported by AWS resource explorer

```sql
select
  service,
  resource_type
from
  aws_resource_explorer_supported_resource_type
where
  service = 'iam';
```
