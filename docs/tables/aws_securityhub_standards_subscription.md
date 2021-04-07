# Table: aws_securityhub_standards_subscription

Security Hub also generates its own findings as the result of running automated and continuous checks against the rules in a set of supported security standards. These checks provide a readiness score and identify specific accounts and resources that require attention.

## Examples

### Basic info

```sql
select
  name,
  standards_arn,
  description region
from
  aws_securityhub_standards_subscription;
```


### List enabled standards in security hub

```sql
select
  name,
  standards_arn,
  enabled_by_default
from
  aws_securityhub_standards_subscription
where
  enabled_by_default;
```


### List standards where standard status is not ready

```sql
select
  name,
  standards_arn,
  standards_subscription_arn,
  standards_status
from
  aws_securityhub_standards_subscription
where
  standards_status <> 'READY';
```