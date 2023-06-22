# Table: aws_service_discovery_service

The AWS Service Discovery service is a fully managed service that allows you to easily discover, register, and connect services within your AWS infrastructure. It simplifies the process of building and managing complex distributed applications by providing a way for services to locate and communicate with each other dynamically.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  type,
  region
from
  aws_service_discovery_service;
```

### List DNS services

```sql
select
  name,
  id,
  arn,
  type,
  create_date
from
  aws_service_discovery_service
where
  type ilike '%dns%';
```

### List HTTP type services

```sql
select
  name,
  id,
  arn,
  type,
  description
from
  aws_service_discovery_service
where
  type = 'HTTP';
```

### List services created in the last 30 days

```sql
select
  name,
  id,
  description,
  create_date
from
  aws_service_discovery_service
where
  create_date >= now() - interval '30' day;
```

### Count services by type

```sql
select
  type,
  count(type)
from
  aws_service_discovery_service
group by
  type;
```

### Get health check config details of services

```sql
select
  name,
  id,
  health_check_config ->> 'Type' as health_check_type,
  health_check_config ->> 'FailureThreshold' as failure_threshold,
  health_check_config ->> 'ResourcePath' as resource_path
from
  aws_service_discovery_service;
```

### Get custom health check config details of services

```sql
select
  name,
  id,
  health_check_custom_config ->> 'FailureThreshold' as failure_threshold
from
  aws_service_discovery_service;
```

### Get namespace details of services

```sql
select
  s.name,
  s.id,
  s.namespace_id,
  n.service_count,
  n.type as namespace_type,
  n.dns_properties
from
  aws_service_discovery_service as s,
  aws_service_discovery_namespace as n;
```
