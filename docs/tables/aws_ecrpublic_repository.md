---
title: "Steampipe Table: aws_ecrpublic_repository - Query AWS Elastic Container Registry Public Repository using SQL"
description: "Allows users to query AWS Elastic Container Registry Public Repository to get detailed information about each ECR public repository within an AWS account."
folder: "ECR Public"
---

# Table: aws_ecrpublic_repository - Query AWS Elastic Container Registry Public Repository using SQL

The AWS Elastic Container Registry Public Repository is a service that allows you to store, manage, and deploy Docker images. It eliminates the need to operate your own container repositories or worry about scaling the underlying infrastructure. It is a fully-managed service that makes it easy to store, manage, share, and deploy your container images and artifacts anywhere.

## Table Usage Guide

The `aws_ecrpublic_repository` table in Steampipe provides you with information about each ECR public repository within your AWS account. This table allows you, as a DevOps engineer, to query repository-specific details, including the repository ARN, repository name, creation date, and associated metadata. You can use this table to gather insights on repositories, such as the number of images per repository, the status of each repository, and more. The schema outlines the various attributes of the ECR public repository for you, including the repository ARN, creation date, image tag mutability, and associated tags.

## Examples

### Basic info
Explore which public repositories are available in your AWS Elastic Container Registry. This can help you manage and track your container images, understand their origins and creation times, and identify the specific regions and accounts associated with each repository.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which repository policy statements are granting full access. This is useful for security audits and ensuring that access permissions are correctly configured.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(p.value, '$') as principal,
  json_extract(a.value, '$') as action,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Condition') as conditions
from
  aws_ecrpublic_repository,
  json_each(json_extract(policy, '$.Statement')) as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a
where
  json_extract(s.value, '$.Effect') = 'Allow'
  and json_extract(a.value, '$') in ('*', 'ecr-public:*');
```