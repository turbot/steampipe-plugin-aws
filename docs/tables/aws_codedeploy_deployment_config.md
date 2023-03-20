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

### Get total configurations deployed on each platform

```sql
select
  count(arn) as configuration_count,
  compute_platform
from
  aws_codedeploy_deployment_config
group by
  compute_platform;
```

### List the user defined configurations

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

### Find the minimum healthy hosts of a particular deployment configuration

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
