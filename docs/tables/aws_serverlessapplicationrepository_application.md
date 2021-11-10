# Table: aws_serverlessapplicationrepository_application

The AWS Serverless Application Repository is a managed repository for serverless applications. It enables teams, organizations, and individual developers to store and share reusable applications, and easily assemble and deploy serverless architectures in powerful new ways. Using the Serverless Application Repository, you don't need to clone, build, package, or publish source code to AWS before deploying it. Instead, you can use pre-built applications from the Serverless Application Repository in your serverless architectures, helping you and your teams reduce duplicated work, ensure organizational best practices, and get to market faster.

A serverless application is a combination of Lambda functions, event sources, and other resources that work together to perform tasks.

## Examples

### Basic info

```sql
select
  name,
  application_id,
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
  application_id,
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
