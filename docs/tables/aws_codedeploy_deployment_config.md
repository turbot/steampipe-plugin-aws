---
title: "Steampipe Table: aws_codedeploy_deployment_config - Query AWS CodeDeploy Deployment Configurations using SQL"
description: "Allows users to query AWS CodeDeploy Deployment Configurations to retrieve information about the deployment configurations within AWS CodeDeploy service."
folder: "CodeDeploy"
---

# Table: aws_codedeploy_deployment_config - Query AWS CodeDeploy Deployment Configurations using SQL

The AWS CodeDeploy Deployment Configurations is a feature of AWS CodeDeploy, a service that automates code deployments to any instance, including Amazon EC2 instances and servers hosted on-premise. Deployment configurations specify deployment rules and success/failure conditions used by AWS CodeDeploy when pushing out new application versions. This enables you to have a consistent, repeatable process for releasing new software, eliminating the complexity of updating applications and systems.

## Table Usage Guide

The `aws_codedeploy_deployment_config` table in Steampipe provides you with information about deployment configurations within AWS CodeDeploy. This table allows you as a DevOps engineer, developer, or system administrator to query deployment configuration details, including deployment configuration names, minimum healthy hosts, and compute platform. You can utilize this table to gather insights on configurations, such as those with specific compute platforms, minimum healthy host requirements, and more. The schema outlines the various attributes of the deployment configuration for you, including the deployment configuration ID, deployment configuration name, and the compute platform.

## Examples

### Basic info
Explore various configurations of your AWS CodeDeploy deployments to understand their compute platforms, creation times, and regions. This can help you manage and optimize your deployments effectively.

```sql+postgres
select
  arn,
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_deployment_config;
```

```sql+sqlite
select
  arn,
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_deployment_config;
```

### Get the configuration count for each compute platform
This query helps you understand the distribution of configurations across different compute platforms in your AWS CodeDeploy service. It's useful for gaining insights into how your deployment configurations are spread across different platforms, aiding in resource allocation and strategic planning.

```sql+postgres
select
  count(arn) as configuration_count,
  compute_platform
from
  aws_codedeploy_deployment_config
group by
  compute_platform;
```

```sql+sqlite
select
  count(arn) as configuration_count,
  compute_platform
from
  aws_codedeploy_deployment_config
group by
  compute_platform;
```

### List the user managed deployment configurations
Determine the areas in which user-managed deployment configurations have been set up. This is useful to understand where and when specific computing platforms were established, providing insights into the regional distribution and timeline of your deployment configurations.

```sql+postgres
select
  arn,
  deployment_config_id,
  deployment_config_name
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_deployment_config
where
  create_time is not null;
```

```sql+sqlite
select
  arn,
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_deployment_config
where
  create_time is not null;
```

### List the minimum healthy hosts required by each deployment configuration
Discover the segments that require the least number of healthy hosts for each deployment configuration. This can be useful in optimizing resource allocation and ensuring efficient application deployment.

```sql+postgres
select
  arn,
  deployment_config_id,
  deployment_config_name
  compute_platform,
  minimum_healthy_hosts ->> 'Type' as host_type,
  minimum_healthy_hosts ->> 'Value' as host_value,
  region
from
  aws_codedeploy_deployment_config
where
  create_time is not null;
```

```sql+sqlite
select
  arn,
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  json_extract(minimum_healthy_hosts, '$.Type') as host_type,
  json_extract(minimum_healthy_hosts, '$.Value') as host_value,
  region
from
  aws_codedeploy_deployment_config
where
  create_time is not null;
```

### Get traffic routing details for `TimeBasedCanary` deployment configurations
Determine the areas in which your AWS CodeDeploy configurations are utilizing TimeBasedCanary deployments. This can be useful for understanding how traffic is managed during deployments, and to assess the percentage and intervals of traffic being directed to your new service versions.

```sql+postgres
select
  arn,
  deployment_config_id,
  deployment_config_name,
  traffic_routing_config -> 'TimeBasedCanary' ->> 'CanaryInterval' as canary_interval,
  traffic_routing_config -> 'TimeBasedCanary' ->> 'CanaryPercentage' as canary_percentage
from
  aws_codedeploy_deployment_config
where
  traffic_routing_config ->> 'Type' = 'TimeBasedCanary';
```

```sql+sqlite
select
  arn,
  deployment_config_id,
  deployment_config_name,
  json_extract(traffic_routing_config, '$.TimeBasedCanary.CanaryInterval') as canary_interval,
  json_extract(traffic_routing_config, '$.TimeBasedCanary.CanaryPercentage') as canary_percentage
from
  aws_codedeploy_deployment_config
where
  json_extract(traffic_routing_config, '$.Type') = 'TimeBasedCanary';
```

### Get traffic routing details for `TimeBasedLinear` deployment configurations
Explore the intricacies of traffic routing for deployments using a 'TimeBasedLinear' configuration. This allows you to understand the rate of change over time, helping to optimize deployment strategies.

```sql+postgres
select
  arn,
  deployment_config_id,
  deployment_config_name,
  traffic_routing_config -> 'TimeBasedLinear' ->> 'LinearInterval' as linear_interval,
  traffic_routing_config -> 'TimeBasedLinear' ->> 'LinearPercentage' as linear_percentage
from
  aws_codedeploy_deployment_config
where
  traffic_routing_config ->> 'Type' = 'TimeBasedLinear';
```

```sql+sqlite
select
  arn,
  deployment_config_id,
  deployment_config_name,
  json_extract(traffic_routing_config, '$.TimeBasedLinear.LinearInterval') as linear_interval,
  json_extract(traffic_routing_config, '$.TimeBasedLinear.LinearPercentage') as linear_percentage
from
  aws_codedeploy_deployment_config
where
  json_extract(traffic_routing_config, '$.Type') = 'TimeBasedLinear';
```