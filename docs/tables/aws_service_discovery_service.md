---
title: "Table: aws_service_discovery_service - Query AWS Service Discovery Service using SQL"
description: "Allows users to query AWS Service Discovery Service to retrieve detailed information about AWS resources that are registered with AWS Cloud Map."
---

# Table: aws_service_discovery_service - Query AWS Service Discovery Service using SQL

The `aws_service_discovery_service` table in Steampipe provides information about AWS resources that are registered with AWS Cloud Map. This table allows DevOps engineers to query service-specific details, including service ID, ARN, name, and associated metadata. Users can utilize this table to gather insights on services, such as service configurations, health check configurations, and more. The schema outlines the various attributes of the service, including the service ID, ARN, name, description, DNS configurations, health check configurations, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_service_discovery_service` table, you can use the `.inspect aws_service_discovery_service` command in Steampipe.

### Key columns:

- `id`: The ID of the service. This is a unique identifier and can be used to join this table with other tables.
- `arn`: The Amazon Resource Name (ARN) of the service. This is a unique identifier for the service across all of AWS and can be used to join this table with other tables.
- `name`: The name of the service. This can be used to filter specific services based on their names.

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
