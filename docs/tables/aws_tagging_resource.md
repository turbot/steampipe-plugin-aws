# Table: aws_tagging_resource

A patch baseline defines which patches are approved for installation on your instances.

## Examples

### Basic info

```sql
select
  name,
  arn,
  compliance_status,
  tags,
  region
from
  aws_tagging_resource;
```

### List of resources with tag key test-resource

```sql
select
  name,
  arn,
  tags
from
  aws_tagging_resource
where
  tags -> 'test-resource' is not null;
```