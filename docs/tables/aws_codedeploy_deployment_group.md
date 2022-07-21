# Table: aws_codedeploy_deployment_group

A deployment group is a set of individual instances. A deployment group contains individually tagged instances, Amazon EC2 instances in Amazon EC2 Auto Scaling groups, or both. The deployment group contains settings and configurations used during the deployment. Most deployment group settings depend on the compute platform used by your application. Some settings, such as rollbacks, triggers, and alarms can be configured for deployment groups for any compute platform.

## Examples

### Basic info

```sql
select
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  create_time,
  minimum_healthy_hosts,
  region
from
  aws_codedeploy_deployment_group;
```


### Get total configs for deployement on each platform

```sql
select
  count(deployment_config_id) as config_count,
  compute_platform
from
  aws_codedeploy_deployment_group
group by
  compute_platform;
```


### List all configs having minimum healthy instances configurations

```sql
select
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  minimum_healthy_hosts ->> 'Type' as minimum_healthy_host_type,
  minimum_healthy_hosts ->> 'Value' as minimum_healthy_host_value,
  region
from
  aws_codedeploy_deployment_group
where
  minimum_healthy_hosts is not null;
```

### List all user managed deployment configs

```sql
select
  deployment_config_id,
  deployment_config_name,
  compute_platform,
  create_time,
  region
from
  aws_codedeploy_deployment_group
where 
  create_time is not null;
```

### Get total configs for each traffic routing type

```sql
select
  count(deployment_config_id) as config_count,
  traffic_routing_config ->> 'Type'
from
  aws_codedeploy_deployment_group
where 
  traffic_routing_config is not null
group by
  traffic_routing_config ->> 'Type';
```