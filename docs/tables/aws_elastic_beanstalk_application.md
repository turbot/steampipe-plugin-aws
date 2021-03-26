# Table: aws_elastic_beanstalk_application

AWS Elastic Beanstalk offers the fastest and simplest way to deploy developer's applications without worrying about the infrastructure while maintaining high availability.

## Examples

### Basic info

```sql
select
  name,
  arn,
  description,
  date_created,
  date_updated,
  versions
from
  aws_elastic_beanstalk_application;
```


### Get resource life cycle configuration details for each application

```sql
select
  name,
  resource_lifecycle_config ->> 'ServiceRole' as role,
  resource_lifecycle_config -> 'VersionLifecycleConfig' ->> 'MaxAgeRule' as max_age_rule,
  resource_lifecycle_config -> 'VersionLifecycleConfig' ->> 'MaxCountRule' as max_count_rule
from
  aws_elastic_beanstalk_application;
```

 