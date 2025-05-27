---
title: "Steampipe Table: aws_app_runner_service - Query AWS App Runner Services using SQL"
description: "Allows users to query AWS App Runner services, providing detailed information on service configurations, scaling, and network settings."
folder: "App Runner"
---

# Table: aws_app_runner_service - Query AWS App Runner Services using SQL

AWS App Runner is a fully managed service that makes it easy to build, deploy, and run containerized web applications and APIs at scale. The `aws_app_runner_service` table in Steampipe allows you to query information about your App Runner services in AWS, including their configurations, scaling policies, and network settings.

## Table Usage Guide

The `aws_app_runner_service` table enables cloud administrators and DevOps engineers to gather detailed insights into their App Runner services. You can query various aspects of the services, such as their scaling configurations, network settings, health checks, and service URLs. This table is particularly useful for monitoring service health, managing configurations, and ensuring that your applications are running efficiently.

## Examples

### Basic service information
Retrieve basic information about your AWS App Runner services, including their name, ARN, and region.

```sql+postgres
select
  service_name,
  arn,
  region,
  created_at,
  updated_at
from
  aws_app_runner_service;
```

```sql+sqlite
select
  service_name,
  arn,
  region,
  created_at,
  updated_at
from
  aws_app_runner_service;
```

### List services with specific network configurations
Identify services that are configured with a specific network configuration, such as VPC or public network settings.

```sql+postgres
select
  service_name,
  arn,
  network_configuration
from
  aws_app_runner_service
where
  (network_configuration -> 'EgressConfiguration' ->> 'VpcConnectorArn') is not null;
```

```sql+sqlite
select
  service_name,
  arn,
  network_configuration
from
  aws_app_runner_service
where
  json_extract(network_configuration, '$.VpcConfiguration') is not null;
```

### List services with auto-scaling configurations
Retrieve information about services that have specific auto-scaling configurations.

```sql+postgres
select
  service_name,
  arn,
  auto_scaling_configuration_summary
from
  aws_app_runner_service
where
  jsonb_path_exists(auto_scaling_configuration_summary, '$.AutoScalingConfigurationArn');
```

```sql+sqlite
select
  service_name,
  arn,
  auto_scaling_configuration_summary
from
  aws_app_runner_service
where
  json_extract(auto_scaling_configuration_summary, '$.AutoScalingConfigurationArn') is not null;
```

### List services with specific observability configurations
Identify services that have observability features enabled, such as logging or tracing.

```sql+postgres
select
  service_name,
  arn,
  observability_configuration
from
  aws_app_runner_service
where
  (observability_configuration ->> 'ObservabilityConfigurationArn') is not null;
```

```sql+sqlite
select
  service_name,
  arn,
  observability_configuration
from
  aws_app_runner_service
where
  json_extract(observability_configuration, '$.ObservabilityConfigurationArn') is not null;
```

### List services created within a specific time frame
Fetch services that were created within a specific date range, which can be useful for auditing purposes.

```sql+postgres
select
  service_name,
  arn,
  created_at
from
  aws_app_runner_service
where
  created_at >= '2023-01-01T00:00:00Z' and created_at <= '2023-12-31T23:59:59Z';
```

```sql+sqlite
select
  service_name,
  arn,
  created_at
from
  aws_app_runner_service
where
  created_at >= '2023-01-01T00:00:00Z' and created_at <= '2023-12-31T23:59:59Z';
```

### Get service URLs for all services
Retrieve the service URLs for all App Runner services, which can be useful for accessing or sharing service endpoints.

```sql+postgres
select
  service_name,
  arn,
  service_url
from
  aws_app_runner_service;
```

```sql+sqlite
select
  service_name,
  arn,
  service_url
from
  aws_app_runner_service;
```