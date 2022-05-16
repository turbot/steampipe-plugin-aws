# Table: aws_securityhub_standards_control

Security Hub provides controls for the following standards.

- `CIS AWS Foundations`

- `Payment Card Industry Data Security Standard (PCI DSS)`

- `AWS Foundational Security Best Practices`

## Examples

### Basic info

```sql
select
  name,
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control;
```

### List disabled controls

```sql
select
  name,
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status = 'DISABLED';
```

### List controls with high severity

```sql
select
  name,
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'HIGH';
```
