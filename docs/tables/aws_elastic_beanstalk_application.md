# Table: aws_elastic_beanstalk_application

AWS Elastic Beanstalk is an easy-to-use service. It offers the fastest and simplest way to deploy developer's applications without worrying about the infrastructure while maintaining high availability.

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


### Details of resource life cycle config for application

```sql
select
  name,
  jsonb_pretty(resource_lifecycle_config)
from
  aws_elastic_beanstalk_application;
```

 