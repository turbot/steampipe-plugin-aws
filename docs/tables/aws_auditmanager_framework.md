---
title: "Steampipe Table: aws_auditmanager_framework - Query AWS Audit Manager Framework using SQL"
description: "Allows users to query AWS Audit Manager Frameworks"
folder: "Audit Manager"
---

# Table: aws_auditmanager_framework - Query AWS Audit Manager Framework using SQL

The AWS Audit Manager Framework is a feature of AWS Audit Manager that helps you continuously audit your AWS usage to simplify your compliance with regulations and industry standards. It automates evidence collection to enable you to scale your audit capability in AWS, reducing the effort needed to assess risk and compliance. This feature is especially useful for organizations that need to maintain a consistent audit process across various AWS services.

## Table Usage Guide

The `aws_auditmanager_framework` table in Steampipe provides you with information about frameworks within AWS Audit Manager. This table allows you, as a DevOps engineer, to query framework-specific details, including the framework's ARN, ID, type, and associated metadata. You can utilize this table to gather insights on frameworks, such as the number of controls associated with each framework, the compliance type, and more. The schema outlines the various attributes of the Audit Manager Framework for you, including the framework ARN, creation date, last updated date, and associated tags.

## Examples

### Basic info
Explore which audit frameworks are currently implemented in your AWS environment. This can help in assessing your existing auditing strategies and identifying areas for improvement.

```sql+postgres
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework;
```

```sql+sqlite
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework;
```

### List custom audit manager frameworks
Uncover the details of your custom audit frameworks within AWS Audit Manager. This query is useful for understanding the scope and details of your custom configurations, aiding in the management and review of your audit frameworks.

```sql+postgres
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework
where
  type = 'Custom';
```

```sql+sqlite
select
  name,
  arn,
  id,
  type
from
  aws_auditmanager_framework
where
  type = 'Custom';
```