# Table: aws_ecrpublic_repository

A public repository is open to publicly pull images from and is visible on the Amazon ECR Public Gallery. When creating a public repository you specify catalog data which helps users find and use your images.

## Examples

### Basic info

```sql
select
  repository_name,
  registry_id,
  arn,
  repository_uri,
  created_at,
  region,
  account_id
from
  aws_ecrpublic_repository;
```

### List repository policy statements that grant full access for each repository

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_ecrpublic_repository,
  jsonb_array_elements(policy -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and a in ('*', 'ecr-public:*');
```
