# Table: aws_codebuild_source_credential

Source repository credentials contain credential information for an AWS CodeBuild project that has its source code stored in GitHub, GitHub Enterprise, or Bitbucket repository.

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
