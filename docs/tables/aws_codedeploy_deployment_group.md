# Table: aws_codedeploy_app

A deployment group is the AWS CodeDeploy entity for grouping EC2 instances or AWS Lambda functions in a CodeDeploy deployment. For EC2 deployments, it is a set of instances associated with an application that you target for a deployment.

## Examples

### Basic info

```sql
select
  arn,
  deployment_group_id,
  deployment_group_name,
  application_name,
  deployment_style,
  region
from
  aws_codedeploy_deployment_group;
```

### Get total deployment groups on each platform

```sql
select
  count(arn) as group_count,
  compute_platform
from
  aws_codedeploy_deployment_group
group by
  compute_platform;
```

### List the last successful deployment for each deployment group

```sql
select
  arn,
  deployment_group_id,
  last_successful_deployment
from
  aws_codedeploy_deployment_group;
```

### Get total deployment groups based on deployment style

```sql
select
  count(arn) as group_count,
  deployment_style
from
  aws_codedeploy_deployment_group
group by
  deployment_style;
```

### List the deployment groups having automatic rollback enabled

```sql
select
  arn,
  deployment_group_id,
  deployment_group_name,
  auto_rollback_configuration ->> 'Enabled' as auto_rollback_configuration_enabled
from
  aws_codedeploy_deployment_group
where
  auto_rollback_configuration ->> 'Enabled' = 'true' ;
```

### List all autoscaling groups in a particular deployment group for an application
```sql
select
  arn as group_arn,
  deployment_group_id,
  deployment_group_name,
  auto_scaling_groups ->> 'Hook' as auto_scaling_group_hook,
  auto_scaling_groups ->> 'Name' as auto_scaling_group_name
from
  aws_codedeploy_deployment_group
where
  application_name = 'abc' and deployment_group_name = 'def' ;
```

### List the deployment groups having automatic rollback enabled

```sql
select
  arn,
  deployment_group_id,
  deployment_group_name,
  alarm_configuration ->> 'Enabled' as alarm_configuration_enabled
from
  aws_codedeploy_deployment_group
where
  alarm_configuration ->> 'Enabled' = 'true' ;
```
