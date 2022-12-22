# Table: aws_securityhub_standards_subscription

Security Hub also generates its own findings as the result of running automated and continuous checks against the rules in a set of supported security standards. These checks provide a readiness score and identify specific accounts and resources that require attention.

## Examples

### Basic info

```sql
select
  name,
  standards_arn,
  description,
  region
from
  aws_securityhub_standards_subscription;
```


### List enabled security hub standards

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


### List standards whose status is not ready

```sql
select
  name,
  standards_arn,
  standards_subscription_arn,
  standards_status,
  standards_status_reason_code
from
  aws_securityhub_standards_subscription
where
  standards_status <> 'READY';
```

### List standards that are not managed by AWS

```sql
select
  name,
  standards_arn,
  standards_managed_by ->> 'Company' as standards_managed_by_company
from
  aws_securityhub_standards_subscription
where
  standards_managed_by ->> 'Company' <> 'AWS';
```