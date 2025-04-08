---
title: "Steampipe Table: aws_service_discovery_service - Query AWS Service Discovery Service using SQL"
description: "Allows users to query AWS Service Discovery Service to retrieve detailed information about AWS resources that are registered with AWS Cloud Map."
folder: "Service Discovery"
---

# Table: aws_service_discovery_service - Query AWS Service Discovery Service using SQL

The AWS Service Discovery Service is a managed solution that makes it easy for microservices and containerized applications to connect with each other over the network. It automatically manages the registration and health checks of services on your behalf, making it easier to discover and connect services across multiple AWS accounts and AWS Regions. Thus, it simplifies the process of building and scaling microservices, service-oriented, and containerized applications.

## Table Usage Guide

The `aws_service_discovery_service` table in Steampipe provides you with information about AWS resources that are registered with AWS Cloud Map. This table enables you, as a DevOps engineer, to query service-specific details, including service ID, ARN, name, and associated metadata. You can utilize this table to gather insights on services, such as service configurations, health check configurations, and more. The schema outlines the various attributes of the service for you, including the service ID, ARN, name, description, DNS configurations, health check configurations, and associated tags.

## Examples

### Basic info
Analyze the settings to understand the fundamental details of your AWS Service Discovery resources, such as their names, identifiers, and regions. This can help manage and organize your resources effectively, ensuring optimal utilization and control.

```sql+postgres
select
  name,
  id,
  arn,
  type,
  region
from
  aws_service_discovery_service;
```

```sql+sqlite
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
Determine the areas in which DNS services are being utilized within your AWS environment. This analysis can help you manage and optimize your cloud resources effectively.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  arn,
  type,
  create_date
from
  aws_service_discovery_service
where
  type like '%dns%';
```

### List HTTP type services
Explore the services in your AWS environment that are configured for HTTP type communication. This can help in understanding the network communication setup and identifying potential areas for security improvement.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which new services have been added to your AWS environment in the past month. This is useful for tracking recent changes and understanding the evolution of your service architecture.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  description,
  create_date
from
  aws_service_discovery_service
where
  create_date >= datetime('now', '-30 day');
```

### Count services by type
Explore the distribution of different service types within your AWS Service Discovery to better manage and optimize your resources.

```sql+postgres
select
  type,
  count(type)
from
  aws_service_discovery_service
group by
  type;
```

```sql+sqlite
select
  type,
  count(type)
from
  aws_service_discovery_service
group by
  type;
```

### Get health check config details of services
Determine the areas in which the health check configurations of services are failing. This query is useful in identifying problematic services that might need attention or adjustment.

```sql+postgres
select
  name,
  id,
  health_check_config ->> 'Type' as health_check_type,
  health_check_config ->> 'FailureThreshold' as failure_threshold,
  health_check_config ->> 'ResourcePath' as resource_path
from
  aws_service_discovery_service;
```

```sql+sqlite
select
  name,
  id,
  json_extract(health_check_config, '$.Type') as health_check_type,
  json_extract(health_check_config, '$.FailureThreshold') as failure_threshold,
  json_extract(health_check_config, '$.ResourcePath') as resource_path
from
  aws_service_discovery_service;
```

### Get custom health check config details of services
Analyze the configuration of services to understand the custom health check details. This information could be useful in identifying services that may have a higher failure threshold, thus requiring more attention or adjustment.

```sql+postgres
select
  name,
  id,
  health_check_custom_config ->> 'FailureThreshold' as failure_threshold
from
  aws_service_discovery_service;
```

```sql+sqlite
select
  name,
  id,
  json_extract(health_check_custom_config, '$.FailureThreshold') as failure_threshold
from
  aws_service_discovery_service;
```

### Get namespace details of services
Identify the details of specific services within a namespace to understand the number and type of services running, along with their DNS properties. This can be useful in managing and optimizing your AWS service discovery.

```sql+postgres
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

```sql+sqlite
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