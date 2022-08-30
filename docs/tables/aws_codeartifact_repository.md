# Table: aws_codeartifact_repository

An AWS CodeBuild project configures how CodeBuild builds your source code. For example, it tells CodeBuild where to get the source code and which build environment to use.

## Examples

### Basic info

```sql
select
  arn,
  domain_name,
  domain_owner,
  upstreams,
  tags
from
  aws_codeartifact_repository;
```

### List all upstream repositories

```sql
select
  arn,
  domain_name,
  domain_owner,
  tags,
  repository_endpoint
from
  aws_codeartifact_repository
where
  repository_endpoint is not null;
```

### List of repository policy statements that grant external access

```sql
select
  arn,
  p as principal,
  a as action,
  s ->> 'Effect' as effect
from
  aws_codeartifact_repository,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  string_to_array(p, ':') as pa,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and (
    pa [5] != account_id
    or p = '*'
  );
```

### Get upstream package details associated to each repository

```sql
select
  arn,
  domain_name,
  domain_owner,
  u ->> 'RepositoryName' as upstream_repo_name
from
  aws_codeartifact_repository,
  jsonb_array_elements(upstreams) u;
```
