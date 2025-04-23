---
title: "Steampipe Table: aws_appconfig_application - Query AWS AppConfig Applications using SQL"
description: "Allows users to query AWS AppConfig Applications to gather detailed information about each application, including its name, description, associated environments, and more."
folder: "AppConfig"
---

# Table: aws_appconfig_application - Query AWS AppConfig Applications using SQL

The AWS AppConfig Application is a feature of AWS AppConfig, which is a service that enables you to create, manage, and quickly deploy application configurations. It is designed to use AWS Lambda, Amazon ECS, Amazon S3, and other AWS services. AWS AppConfig Application helps you manage the configurations of your Amazon Web Services applications in a centralized manner, reducing error and increasing speed in deployment.

## Table Usage Guide

The `aws_appconfig_application` table in Steampipe provides you with information about AWS AppConfig Applications. This table allows you, as a DevOps engineer or other technical professional, to query application-specific details, including its ID, name, description, and associated environments. You can utilize this table to gather insights on applications, such as their deployment strategies, associated configurations, and more. The schema outlines the various attributes of the AppConfig application for you, including the application ID, name, description, and associated tags.

## Examples

### Basic info
Explore which AWS AppConfig applications are currently in use. This can help you manage and monitor your applications effectively, ensuring they're configured correctly and align with your operational requirements.

```sql+postgres
select
  arn,
  id,
  name,
  description,
  tags
from
  aws_appconfig_application;
```

```sql+sqlite
select
  arn,
  id,
  name,
  description,
  tags
from
  aws_appconfig_application;
```