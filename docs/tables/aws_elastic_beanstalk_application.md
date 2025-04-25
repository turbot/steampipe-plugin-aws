---
title: "Steampipe Table: aws_elastic_beanstalk_application - Query AWS Elastic Beanstalk Applications using SQL"
description: "Allows users to query AWS Elastic Beanstalk Applications to obtain details about their configurations, versions, environment, and other metadata."
folder: "Elastic Beanstalk"
---

# Table: aws_elastic_beanstalk_application - Query AWS Elastic Beanstalk Applications using SQL

The AWS Elastic Beanstalk Application is a component of AWS's platform-as-a-service (PaaS) offering, Elastic Beanstalk. It allows developers to deploy and manage applications in multiple languages without worrying about the infrastructure that runs those applications. The Elastic Beanstalk Application handles capacity provisioning, load balancing, and automatic scaling, among other tasks, enabling developers to focus on their application code.

## Table Usage Guide

The `aws_elastic_beanstalk_application` table in Steampipe provides you with information about applications within AWS Elastic Beanstalk. This table enables you, as a DevOps engineer, to query application-specific details, including application ARN, description, date created, date updated, and associated metadata. You can utilize this table to gather insights on applications, such as application versions, configurations, associated environments, and more. The schema outlines for you the various attributes of the Elastic Beanstalk application, including the resource lifecycles, configurations, and associated tags.

## Examples

### Basic info
Explore the applications in your AWS Elastic Beanstalk environment to understand their creation and update timeline, as well as the different versions available. This can help in managing the applications better by keeping track of their versions and update history.

```sql+postgres
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

```sql+sqlite
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
Determine the life cycle configurations of your applications to understand the roles assigned and the rules set for version management. This can help in optimizing resource usage and maintaining application health.

```sql+postgres
select
  name,
  resource_lifecycle_config ->> 'ServiceRole' as role,
  resource_lifecycle_config -> 'VersionLifecycleConfig' ->> 'MaxAgeRule' as max_age_rule,
  resource_lifecycle_config -> 'VersionLifecycleConfig' ->> 'MaxCountRule' as max_count_rule
from
  aws_elastic_beanstalk_application;
```

```sql+sqlite
select
  name,
  json_extract(resource_lifecycle_config, '$.ServiceRole') as role,
  json_extract(resource_lifecycle_config, '$.VersionLifecycleConfig.MaxAgeRule') as max_age_rule,
  json_extract(resource_lifecycle_config, '$.VersionLifecycleConfig.MaxCountRule') as max_count_rule
from
  aws_elastic_beanstalk_application;
```