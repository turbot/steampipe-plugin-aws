# Table: aws_trusted_advisor_check_summary

A Trusted Advisor check is a specific evaluation or assessment performed by Trusted Advisor in different categories. These checks cover various areas, including cost optimization, security, performance, and fault tolerance. Each check examines a specific aspect of your AWS resources and provides recommendations for improvement.

## Examples

### Basi info

```sql
select
  name,
  check_id,
  category,
  description,
  status,
  timestamp,
  resources_flagged
from
  aws_trusted_advisor_check_summary;
```

### Get error check summaries

```sql
select
  name,
  check_id,
  category,
  status
from
  aws_trusted_advisor_check_summary
where
  status = 'error';
```

### Get last 5days check summaries

```sql
select
  name,
  check_id,
  description,
  status,
  timestamp
from
  aws_trusted_advisor_check_summary
where
  timestamp >= now() - interval '5 day';
```

### Get resource summaries of each check

```sql
select
  name,
  check_id,
  resources_flagged,
  resources_ignored,
  resources_processed,
  resources_suppressed
from
  aws_trusted_advisor_check_summary;
```
