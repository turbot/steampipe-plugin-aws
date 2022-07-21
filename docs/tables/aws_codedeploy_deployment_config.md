# Table: aws_codedeploy_deployment_config

A deployment configuration is set of deployment rules and deployment success and failure conditions used by CodeDeploy during a deployment. If your deployment uses the EC2/On-Premises compute platform, you can specify the minimum number of healthy instances for the deployment. If your deployment uses the AWS Lambda or the Amazon ECS compute platform, you can specify how traffic is routed to your updated Lambda function or ECS task set.

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
  aws_codedeploy_deployment_config;
```


### Get total configs for deployement on each platform

```sql
select
  count(deployment_config_id) as config_count,
  compute_platform
from
  aws_codedeploy_deployment_config
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
  aws_codedeploy_deployment_config
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
  aws_codedeploy_deployment_config
where 
  create_time is not null;
```

### Get total configs for each traffic routing type

```sql
select
  count(deployment_config_id) as config_count,
  traffic_routing_config ->> 'Type'
from
  aws_codedeploy_deployment_config
where 
  traffic_routing_config is not null
group by
  traffic_routing_config ->> 'Type';
```