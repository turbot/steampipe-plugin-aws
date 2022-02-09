# Table: aws_servicequotas_service_quota

Service Quotas is an AWS service that helps you manage your quotas for many AWS services, from one location. Along with looking up the quota values, you can also request a quota increase from the Service Quotas console.

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
  aws_servicequotas_service_quota;
```

### List global service quotas

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_service_quota;
where
  global_quota;
```

### List service quotas for a specific service

```sql
select
  quota_name,
  quota_code,
  quota_arn,
  service_name,
  service_code,
  value
from
  aws_servicequotas_service_quota
where
  service_code = 'athena';
```