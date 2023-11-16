---
title: "Table: aws_codebuild_source_credential - Query AWS CodeBuild Source Credentials using SQL"
description: "Allows users to query AWS CodeBuild Source Credentials"
---

# Table: aws_codebuild_source_credential - Query AWS CodeBuild Source Credentials using SQL

The `aws_codebuild_source_credential` table in Steampipe provides information about source credentials within AWS CodeBuild. This table allows DevOps engineers to query specific details about source credentials, including the ARN, server type, authentication type, and token. Users can utilize this table to gather insights on source credentials, such as identifying the server types, verifying the authentication types, and more. The schema outlines the various attributes of the source credential, including the ARN, server type, authentication type, and token.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codebuild_source_credential` table, you can use the `.inspect aws_codebuild_source_credential` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the source credential. This can be used as a unique identifier for the source credential, and can be used to join this table with other tables that reference the ARN.
- `server_type`: The type of the source server. This can be used to filter source credentials based on the server type.
- `auth_type`: The type of authentication used for the source credential. This can be used to filter source credentials based on the authentication type.

## Examples

### Basic info

```sql
select
  arn,
  server_type,
  auth_type,
  region
from
  aws_codebuild_source_credential;
```


### List projects using OAuth to access GitHub source repository

```sql
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
