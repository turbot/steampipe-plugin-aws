# Table: aws_codedeploy_deployment_config

A deployment configuration is a set of rules and success and failure conditions used by CodeDeploy during a deployment. These rules and conditions are different, depending on whether you deploy to an EC2/On-Premises compute platform, AWS Lambda compute platform, or Amazon ECS compute platform.

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