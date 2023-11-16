---
title: "Table: aws_elastic_beanstalk_application - Query AWS Elastic Beanstalk Applications using SQL"
description: "Allows users to query AWS Elastic Beanstalk Applications to obtain details about their configurations, versions, environment, and other metadata."
---

# Table: aws_elastic_beanstalk_application - Query AWS Elastic Beanstalk Applications using SQL

The `aws_elastic_beanstalk_application` table in Steampipe provides information about applications within AWS Elastic Beanstalk. This table allows DevOps engineers to query application-specific details, including application ARN, description, date created, date updated and associated metadata. Users can utilize this table to gather insights on applications, such as application versions, configurations, associated environments, and more. The schema outlines the various attributes of the Elastic Beanstalk application, including the resource lifecycles, configurations, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elastic_beanstalk_application` table, you can use the `.inspect aws_elastic_beanstalk_application` command in Steampipe.

**Key columns**:

- `application_name`: The name of the application. This can be used to join with other tables that require application name as input.
- `application_arn`: The Amazon Resource Name (ARN) of the application. This can be used to join with other tables that require the ARN.
- `date_created`: The date when the application was created. This can be useful for tracking application lifecycle and auditing purposes.

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
