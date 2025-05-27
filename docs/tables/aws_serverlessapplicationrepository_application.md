---
title: "Steampipe Table: aws_serverlessapplicationrepository_application - Query AWS Serverless Application Repository Applications using SQL"
description: "Allows users to query AWS Serverless Application Repository Applications to fetch details like application name, status, author, description, labels, license URL, creation time, and more."
folder: "Serverless Application Repository"
---

# Table: aws_serverlessapplicationrepository_application - Query AWS Serverless Application Repository Applications using SQL

The AWS Serverless Application Repository is a managed repository for serverless applications. It enables teams, developers, and organizations to discover, configure, and deploy serverless applications and components on AWS. It simplifies the management of serverless applications by providing a mechanism to store and share applications, and to easily configure and deploy them in AWS environments.

## Table Usage Guide

The `aws_serverlessapplicationrepository_application` table in Steampipe provides you with information about Applications within AWS Serverless Application Repository. This table enables you, as a DevOps engineer, to query application-specific details, including application name, status, author, description, labels, license URL, creation time, and more. You can utilize this table to gather insights on applications, such as applications by specific authors, applications with certain labels, applications under certain licenses, and more. The schema outlines the various attributes of the Application for you, including the application ID, home page URL, semantic version, and associated tags.

## Examples

### Basic info
Discover the segments that use AWS serverless applications and gain insights into their authors and creation times. This can be useful in understanding the distribution and usage of serverless applications across your AWS environment.

```sql+postgres
select
  name,
  arn,
  author,
  creation_time,
  description
from
  aws_serverlessapplicationrepository_application;
```

```sql+sqlite
select
  name,
  arn,
  author,
  creation_time,
  description
from
  aws_serverlessapplicationrepository_application;
```


### List applications created by verified author
Discover the segments that consist of applications created by verified authors, which can provide a level of trust and assurance in the application's functionality and security. This is particularly useful when assessing the credibility of applications within your AWS Serverless Application Repository.

```sql+postgres
select
  name,
  arn,
  author,
  is_verified_author
from
  aws_serverlessapplicationrepository_application
where
  is_verified_author;
```

```sql+sqlite
select
  name,
  arn,
  author,
  is_verified_author
from
  aws_serverlessapplicationrepository_application
where
  is_verified_author = 1;
```

### List application policy details
Determine the specifics of application policies within your AWS Serverless Application Repository. This query is useful for understanding the actions, principal organization IDs, and principals associated with each policy, providing valuable insight for policy management and security audits.

```sql+postgres
select
  name,
  jsonb_pretty(statement -> 'Actions') as actions,
  jsonb_pretty(statement -> 'PrincipalOrgIDs') as principal_org_ids,
  jsonb_pretty(statement -> 'Principals') as principals,
  statement ->> 'StatementId' as statement_id
from
  aws_serverlessapplicationrepository_application,
  jsonb_array_elements(statements) as statement;
```

```sql+sqlite
select
  name,
  json_extract(statement.value, '$.Actions') as actions,
  json_extract(statement.value, '$.PrincipalOrgIDs') as principal_org_ids,
  json_extract(statement.value, '$.Principals') as principals,
  json_extract(statement.value, '$.StatementId') as statement_id
from
  aws_serverlessapplicationrepository_application,
  json_each(statements) as statement;
```