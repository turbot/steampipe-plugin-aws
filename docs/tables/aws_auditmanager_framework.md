# Table: aws_auditmanager_framework

The framework library is the central place from which you can access and manage frameworks in AWS Audit Manager.The framework library contains a catalog of standard and custom frameworks.

## Examples

### Basic info

```sql
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework;
```

### List custom audit manager frameworks

```sql
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework
where
  type = 'Custom';
```
