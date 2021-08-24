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

### List of resources where resource is compliant with the effective tag policy

```sql
select
  name,
  arn,
  tags,
  compliance_status
from
  aws_tagging_resource
where
  compliance_status;
```