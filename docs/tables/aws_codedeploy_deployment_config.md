---
title: "Table: aws_codedeploy_deployment_config - Query AWS CodeDeploy Deployment Configurations using SQL"
description: "Allows users to query AWS CodeDeploy Deployment Configurations to retrieve information about the deployment configurations within AWS CodeDeploy service."
---

# Table: aws_codedeploy_deployment_config - Query AWS CodeDeploy Deployment Configurations using SQL

The `aws_codedeploy_deployment_config` table in Steampipe provides information about deployment configurations within AWS CodeDeploy. This table allows DevOps engineers, developers and system administrators to query deployment configuration details, including deployment configuration names, minimum healthy hosts and compute platform. Users can utilize this table to gather insights on configurations, such as those with specific compute platforms, minimum healthy host requirements, and more. The schema outlines the various attributes of the deployment configuration, including the deployment configuration ID, deployment configuration name, and the compute platform.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codedeploy_deployment_config` table, you can use the `.inspect aws_codedeploy_deployment_config` command in Steampipe.

Key columns:

- `deployment_config_id`: The unique ID of a deployment configuration. This can be used to join with other tables when there is a need to correlate more information based on deployment configuration ID.
- `deployment_config_name`: The name of a deployment configuration. This is useful when querying specific deployment configurations by name.
- `compute_platform`: The name of the compute platform. This is important when filtering or joining tables based on the compute platform.

## Examples

### Basic info

```sql
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

```sql
select
  count(arn) as configuration_count,
  compute_platform
from
  aws_codedeploy_deployment_config
group by
  compute_platform;
```

### List the user managed deployment configurations

```sql
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

### List the minimum healthy hosts required by each deployment configuration

```sql
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

### Get traffic routing details for `TimeBasedCanary` deployment configurations

```sql
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

### Get traffic routing details for `TimeBasedLinear` deployment configurations

```sql
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