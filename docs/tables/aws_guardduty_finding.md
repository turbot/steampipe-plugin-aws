# Table: aws_guardduty_finding

A GuardDuty finding represents a potential security issue detected within your network. GuardDuty generates a finding whenever it detects unexpected and potentially malicious activity in your AWS environment.

## Examples

### Basic info

```sql
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding;
```

### List findings that are not archived

```sql
select
  id,
  detector_id,
  arn,
  created_at
from
  aws_guardduty_finding
where
  service ->> 'Archived' = 'false';
```
