# Table: aws_securityhub_standards_control

Security Hub provides controls for the following standards.

- `CIS AWS Foundations`

- `Payment Card Industry Data Security Standard (PCI DSS)`

- `AWS Foundational Security Best Practices`

## Examples

### Basic info

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control;
```

### List disabled controls

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status = 'DISABLED';
```

### Count the number of controls by severity

```sql
select
  severity_rating,
  count(severity_rating)
from
  aws_securityhub_standards_control
group by
  severity_rating;
```

### List controls with high severity

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'HIGH';
```

### List controls which were updated in the last 30 days

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  date_part('day', now() - control_status_updated_at) <= 30;
```

### List CIS AWS foundations benchmark controls with critical severity

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'CRITICAL' and standards_control_arn like '%cis-aws-foundations-benchmark%';
```
