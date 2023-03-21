# Table: aws_codedeploy_app

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

### Get the total configurations for each compute platform

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

### Display canary for the deployment configurations where traffic routing type is 'TimeBasedCanary'

```sql
select
  arn,
  deployment_config_id,
  deployment_config_name,
  traffic_routing_config ->> 'TimeBasedCanary' as cranary,
  region
from
  aws_codedeploy_deployment_config
where
  traffic_routing_config ->> 'Type' = 'TimeBasedCanary';
```

### Display canary for the deployment configurations where traffic routing type is 'TimeBasedLinear'

```sql
select
  arn,
  deployment_config_id,
  deployment_config_name,
  traffic_routing_config ->> 'TimeBasedLinear' as cranary,
  region
from
  aws_codedeploy_deployment_config
where
  traffic_routing_config ->> 'Type' = 'TimeBasedLinear';
```