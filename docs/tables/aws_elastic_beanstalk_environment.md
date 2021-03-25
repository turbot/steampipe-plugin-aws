# Table: aws_elastic_beanstalk_environment

AWS Elastic Beanstalk makes easy for you to create, deploy, and manage scalable, fault-tolerant applications running on the Amazon Web Services cloud.

## Examples

### Basic info

```sql
select
  environment_id,
  environment_name,
  application_name,
  environment_arn,
  tier
from
  aws_elastic_beanstalk_environment;
```

### List environments whose configuration update enable

```sql
select
  environment_name,
  abortable_operation_in_progress
from
  aws_elastic_beanstalk_environment
where
  abortable_operation_in_progress = 'true';
```

### List environments whose health is not responsive

```sql
select
  environment_name,
  application_name,
  environment_id,
  health
from
  aws_elastic_beanstalk_environment
where
  health = 'Red';
```

### List applications having health status disabled

```sql
select
  environment_name,
  health_status
from
  aws_elastic_beanstalk_environment
where
  health_status = 'Suspended';
```
