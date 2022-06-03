# Table: aws_servicequotas_service_quota_change_request

For adjustable quotas, you can request a quota increase. Smaller increases are automatically approved, and larger requests are submitted to AWS Support. You can track your request case in the AWS Support console. Requests to increase service quotas do not receive priority support.

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
  id,
  case_id,
  status,
  quota_name,
  quota_code,
  desired_value
from
  aws_servicequotas_service_quota_change_request
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