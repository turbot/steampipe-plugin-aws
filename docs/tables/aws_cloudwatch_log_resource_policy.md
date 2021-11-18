# Table: aws_cloudwatch_log_resource_policy

A log resource policy is a policy that enables one or more entities to put logs to a log group in the account.

## Examples

### Basic Info

```sql
select
  policy_name,
  last_updated_time,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_cloudwatch_log_resource_policy;
```
