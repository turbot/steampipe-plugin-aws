---
title: "Table: aws_serverlessapplicationrepository_application - Query AWS Serverless Application Repository Applications using SQL"
description: "Allows users to query AWS Serverless Application Repository Applications to fetch details like application name, status, author, description, labels, license URL, creation time, and more."
---

# Table: aws_serverlessapplicationrepository_application - Query AWS Serverless Application Repository Applications using SQL

The `aws_serverlessapplicationrepository_application` table in Steampipe provides information about Applications within AWS Serverless Application Repository. This table allows DevOps engineers to query application-specific details, including application name, status, author, description, labels, license URL, creation time, and more. Users can utilize this table to gather insights on applications, such as applications by specific authors, applications with certain labels, applications under certain licenses, and more. The schema outlines the various attributes of the Application, including the application ID, home page URL, semantic version, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_serverlessapplicationrepository_application` table, you can use the `.inspect aws_serverlessapplicationrepository_application` command in Steampipe.

### Key columns:

- `application_id`: The Amazon Resource Name (ARN) of the application. It is a unique identifier and can be used to join this table with other tables.
- `author`: The name of the author publishing the app. It can be used to filter applications by a specific author.
- `creation_time`: The date/time the application was created. It can be used to filter applications based on their creation time.

## Examples

### Basic info

```sql
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

```sql
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

### List application policy details

```sql
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
