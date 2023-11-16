---
title: "Table: aws_ecrpublic_repository - Query AWS Elastic Container Registry Public Repository using SQL"
description: "Allows users to query AWS Elastic Container Registry Public Repository to get detailed information about each ECR public repository within an AWS account."
---

# Table: aws_ecrpublic_repository - Query AWS Elastic Container Registry Public Repository using SQL

The `aws_ecrpublic_repository` table in Steampipe provides information about each ECR public repository within an AWS account. This table allows DevOps engineers to query repository-specific details, including the repository ARN, repository name, creation date, and associated metadata. This table can be used to gather insights on repositories, such as the number of images per repository, the status of each repository, and more. The schema outlines the various attributes of the ECR public repository, including the repository ARN, creation date, image tag mutability, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecrpublic_repository` table, you can use the `.inspect aws_ecrpublic_repository` command in Steampipe.

### Key columns:

- `repository_name`: The name of the repository. This can be used to join with other tables that provide additional details about the repository.
- `repository_arn`: The Amazon Resource Name (ARN) of the repository. This is a unique identifier for the repository and can be used to join with other tables that use ARN for identification.
- `created_at`: The date and time when the repository was created. This is useful for tracking the lifecycle of the repository.

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
