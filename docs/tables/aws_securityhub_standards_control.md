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
  severity_rating
order by
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
  control_status_updated_at >= (now() - interval '30' day);
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
  severity_rating = 'CRITICAL'
  and arn like '%cis-aws-foundations-benchmark%';
```

### List related requirements benchmark for S3 controls

```sql
select
  control_id,
  r as related_requirements
from
  aws_securityhub_standards_control,
  jsonb_array_elements_text(related_requirements) as r
where
  control_id like '%S3%'
group by
  control_id, r
order by
  control_id, r;
```

### List controls which require PCI DSS benchmark

```sql
select
  r as related_requirements,
  control_id
from
  aws_securityhub_standards_control,
  jsonb_array_elements_text(related_requirements) as r
where
  r like '%PCI%'
group by
  r, control_id
order by
  r, control_id;
```
