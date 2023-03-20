# Table: aws_codedeploy_app

A deployment group is the AWS CodeDeploy entity for grouping EC2 instances or AWS Lambda functions in a CodeDeploy deployment. For EC2 deployments, it is a set of instances associated with an application that you target for a deployment.

## Examples

### Basic info

```sql
select
  arn,
  deployment_group_id,
  deployment_group_name
  application_name,
  deployment_style,
  region
from
  aws_codedeploy_deployment_group
where
  application_name = 'abc';
```

### Get total deployment groups on each platform

```sql
select
  count(arn) as group_count,
  compute_platform
from
  aws_codedeploy_deployment_group
where
  application_name = 'abc'
group by
  compute_platform;
```

### List the last successful deployment of a deployment group
```sql
select
  arn,
  deployment_group_id,
  last_successful_deployment
from
  aws_codedeploy_deployment_group
where
  application_name = 'abc';
```
