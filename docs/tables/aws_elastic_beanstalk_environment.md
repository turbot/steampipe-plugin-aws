# Table: aws_elastic_beanstalk_environment

AWS Elastic Beanstalk makes easy for you to create, deploy, and manage scalable, fault-tolerant applications running on the Amazon Web Services cloud.

## Examples

### Basic ElasticBeanstalk Environment info

```sql
select
  application_name,
  environment_id,
  environment_name,
  environment_arn,
  tier
from
  aws_elastic_beanstalk_environment;
```


### list of Environments whose AbortableOperationInProgress is set enable

```sql
select
  environment_name
from
  aws_elastic_beanstalk_environment
where
  abortable_operation_in_progress = 'true';
```


### list the Environments whose health is not responsive

```sql
select
  environment_name
from
  aws_elastic_beanstalk_environment
where
  health = 'Red';
```


### list of applications running in environment whose health status disabled

```sql
select
  environment_name
from
  aws_elastic_beanstalk_environment
where
  health_status = 'Suspended';
```
