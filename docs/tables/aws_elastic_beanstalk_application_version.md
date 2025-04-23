---
title : "Steampipe Table: aws_elastic_beanstalk_application_version - Query AWS Elastic Beanstalk Application Versions using SQL"
description: "Allows users to query AWS Elastic Beanstalk Application Versions to obtain details about their configurations, environments, and other metadata."
folder: "Elastic Beanstalk"
---

# Table: aws_elastic_beanstalk_application_version - Query AWS Elastic Beanstalk Application Versions using SQL

The AWS Elastic Beanstalk Application Version is a component of AWS's platform-as-a-service (PaaS) offering, Elastic Beanstalk. It allows developers to deploy and manage applications in multiple languages without worrying about the infrastructure that runs those applications. The Elastic Beanstalk Application Version handles capacity provisioning, load balancing, and automatic scaling, among other tasks, enabling developers to focus on their application code.

## Table Usage Guide

The `aws_elastic_beanstalk_application_version` table in Steampipe provides you with information about application versions within AWS Elastic Beanstalk. This table enables you, as a DevOps engineer, to query application version-specific details, including application version ARN, description, date created, date updated, and associated metadata. You can utilize this table to gather insights on application versions, such as application version configurations, associated environments, and more. The schema outlines for you the various attributes of the Elastic Beanstalk application version, including the resource lifecycles, configurations, and associated tags.

## Examples

### Basic info

Explore the application versions in your AWS Elastic Beanstalk environment to understand their creation and update timeline, as well as the different configurations available. This can help in managing the application versions better by keeping track of their configurations and update history.

```sql+postgres
select
  application_name,
  application_version_arn,
  version_label,
  description,
  date_created,
  date_updated,
  source_bundle
from
  aws_elastic_beanstalk_application_version;
```

```sql+sqlite
select
  application_name,
  application_version_arn,
  version_label,
  description,
  date_created,
  date_updated,
  source_bundle
from
  aws_elastic_beanstalk_application_version;
```

### List the recently updated application versions

Identify the application versions that have been recently updated in your AWS Elastic Beanstalk environment. This can help in tracking the recent changes made to the application versions and understanding the impact of these changes on the environment.

```sql+postgres
select
  application_name,
  application_version_arn,
  version_label,
  date_updated
from
  aws_elastic_beanstalk_application_version
order by
  date_updated desc;
```

```sql+sqlite
select
  application_name,
  application_version_arn,
  version_label,
  date_updated
from
  aws_elastic_beanstalk_application_version
order by
  date_updated desc;
```

### List the application versions which are 'Processed'

Identify the application versions that are in the 'Processed' state in your AWS Elastic Beanstalk environment. This can help in understanding the status of the application versions and their readiness for deployment.

```sql+postgres
select
  application_name,
  application_version_arn,
  version_label,
  status
from
  aws_elastic_beanstalk_application_version
where
  status = 'Processed';
```

```sql+sqlite
select
  application_name,
  application_version_arn,
  version_label,
  status
from
  aws_elastic_beanstalk_application_version
where
  status = 'Processed';
```

### List the application versions of a specific application

Identify the application versions of a specific application in your AWS Elastic Beanstalk environment. This can help in understanding the different versions available for a specific application and their configurations.

```sql+postgres
select
  application_name,
  application_version_arn,
  version_label,
  description,
  date_created,
  date_updated,
  source_bundle
from
  aws_elastic_beanstalk_application_version
where
  application_name = 'my-application';
```

```sql+sqlite
select
  application_name,
  application_version_arn,
  version_label,
  description,
  date_created,
  date_updated,
  source_bundle
from
  aws_elastic_beanstalk_application_version
where
  application_name = 'my-application';
```

### List the application versions with specific tags

Identify the application versions with specific tags in your AWS Elastic Beanstalk environment. This can help in understanding the tags associated with the application versions and their metadata.

```sql+postgres
select
  application_name,
  application_version_arn,
  version_label,
  tags
from
  aws_elastic_beanstalk_application_version
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  application_name,
  application_version_arn,
  version_label,
  tags
from
  aws_elastic_beanstalk_application_version
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List the application versions where the source repository is stored in CodeCommit

Identify the application versions where the source repository is stored in AWS CodeCommit in your AWS Elastic Beanstalk environment. This can help in understanding the source repository of the application versions and their configurations.

```sql+postgres
select
  application_name,
  application_version_arn,
  version_label
from
  aws_elastic_beanstalk_application_version
where
  source_build_information ->> 'SourceRepository' = 'CodeCommit';
```

```sql+sqlite
select
  application_name,
  application_version_arn,
  version_label
from
  aws_elastic_beanstalk_application_version
where
  json_extract(source_build_information, '$.SourceRepository') = 'CodeCommit';
```
