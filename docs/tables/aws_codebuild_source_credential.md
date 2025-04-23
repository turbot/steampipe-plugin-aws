---
title: "Steampipe Table: aws_codebuild_source_credential - Query AWS CodeBuild Source Credentials using SQL"
description: "Allows users to query AWS CodeBuild Source Credentials"
folder: "CodeBuild"
---

# Table: aws_codebuild_source_credential - Query AWS CodeBuild Source Credentials using SQL

The AWS CodeBuild Source Credentials are used to interact with external code repositories. They store the authentication information required to access private repositories in GitHub, BitBucket, and AWS CodeCommit. This feature enables secure connection to these repositories, allowing AWS CodeBuild to read the source code for build operations.

## Table Usage Guide

The `aws_codebuild_source_credential` table in Steampipe provides you with information about source credentials within AWS CodeBuild. This table allows you as a DevOps engineer to query specific details about source credentials, including the ARN, server type, authentication type, and token. You can utilize this table to gather insights on source credentials, such as identifying the server types, verifying the authentication types, and more. The schema outlines the various attributes of the source credential for you, including the ARN, server type, authentication type, and token.

## Examples

### Basic info
Determine the areas in which authentication types and server types are used across different regions. This can provide useful insights for managing and optimizing the use of AWS CodeBuild source credentials.

```sql+postgres
select
  arn,
  server_type,
  auth_type,
  region
from
  aws_codebuild_source_credential;
```

```sql+sqlite
select
  arn,
  server_type,
  auth_type,
  region
from
  aws_codebuild_source_credential;
```


### List projects using OAuth to access GitHub source repository
This query helps identify projects that are utilizing OAuth for accessing GitHub as their source repository. This could be useful for auditing purposes, ensuring the correct authorization method is being used for accessing code repositories.

```sql+postgres
select
  p.arn as project_arn,
  p.source ->> 'Location' as source_repository, 
  p.source ->> 'Type' as source_repository_type,
  c.auth_type as authorization_type
from
  aws_codebuild_project as p
  join aws_codebuild_source_credential as c on (p.region = c.region and p.source ->> 'Type' = c.server_type)
where
  p.source ->> 'Type' = 'GITHUB'
  and c.auth_type = 'OAUTH';
```

```sql+sqlite
select
  p.arn as project_arn,
  json_extract(p.source, '$.Location') as source_repository, 
  json_extract(p.source, '$.Type') as source_repository_type,
  c.auth_type as authorization_type
from
  aws_codebuild_project as p
  join aws_codebuild_source_credential as c on (p.region = c.region and json_extract(p.source, '$.Type') = c.server_type)
where
  json_extract(p.source, '$.Type') = 'GITHUB'
  and c.auth_type = 'OAUTH';
```