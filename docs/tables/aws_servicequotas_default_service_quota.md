# Table: aws_servicequotas_default_service_quota

AWS account has default quotas, formerly referred to as limits, for each AWS service. Unless otherwise noted, each quota is Region-specific. You can request increases for some quotas, and other quotas cannot be increased.

## Examples

### Basic info

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_default_service_quota;
```

### List global default service quotas

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_default_service_quota
where
  global_quota;
```

### List default service quotas for a specific service

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_default_service_quota
where
  service_code = 'athena';
```