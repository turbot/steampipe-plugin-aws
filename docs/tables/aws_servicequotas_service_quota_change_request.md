# Table: aws_servicequotas_service_quota_change_request

Service Quotas is an AWS service that helps you manage your quotas for many AWS services, from one location. Along with looking up the quota values, you can also request a quota increase from the Service Quotas console.

## Examples

### Basic info

```sql
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request;
```

### List denied service quota change requests

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
  status = 'DENIED';
```

### List service quota change requests for a specific service

```sql
select
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request
where
  service_code = 'athena';
```