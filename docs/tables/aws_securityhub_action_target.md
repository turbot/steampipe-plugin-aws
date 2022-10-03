# Table: aws_securityhub_action_target

You can use custom actions on findings and insights in Security Hub to trigger target actions in Amazon CloudWatch Events.

## Examples

### Basic info

```sql
select
  name,
  arn,
  region
from
  aws_securityhub_action_target;
```

### Get details of a specific action target

```sql
select
  name,
  arn,
  region
from
  aws_securityhub_action_target
where
  arn = 'arn:aws:securityhub:ap-south-1:*****:action/custom/test';
```
