---
title: "Table: aws_ecr_repository - Query AWS ECR Repositories using SQL"
description: "Allows users to query AWS Elastic Container Registry (ECR) Repositories and retrieve detailed information about each repository."
---

# Table: aws_ecr_repository - Query AWS ECR Repositories using SQL

The `aws_ecr_repository` table in Steampipe provides information about repositories within AWS Elastic Container Registry (ECR). This table allows DevOps engineers to query repository-specific details, including repository ARN, repository URI, and creation date. Users can utilize this table to gather insights on repositories, such as repository policies, image scanning configurations, image tag mutability, and more. The schema outlines the various attributes of the ECR repository, including the repository name, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecr_repository` table, you can use the `.inspect aws_ecr_repository` command in Steampipe.

Key columns:

- `repository_name`: The name of the repository. This column can be used to join this table with others that contain repository-specific information.
- `repository_arn`: The Amazon Resource Number (ARN) of the repository. This column is useful for joining with other tables that use ARN for resource identification.
- `repository_uri`: The URI of the repository. This column is important as it provides the location of the repository, which can be used for joining with other tables that contain URI-based resource data.

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
  aws_ecr_repository;
```

### List repositories which are not using Customer Managed Keys (CMK) for encryption

```sql
select
  repository_name,
  encryption_configuration ->> 'EncryptionType' as encryption_type,
  encryption_configuration ->> 'KmsKey' as kms_key
from
  aws_ecr_repository
where
  encryption_configuration ->> 'EncryptionType' = 'AES256';
```

### List repositories with automatic image scanning disabled

```sql
select
  repository_name,
  image_scanning_configuration ->> 'ScanOnPush' as scan_on_push
from
  aws_ecr_repository
where
  image_scanning_configuration ->> 'ScanOnPush' = 'false';
```

### List images for each repository

```sql
select
  r.repository_name as repository_name,
  i.image_digest as image_digest,
  i.image_tags as image_tags,
  i.image_pushed_at as image_pushed_at,
  i.image_size_in_bytes as image_size_in_bytes,
  i.last_recorded_pull_time as last_recorded_pull_time,
  i.registry_id as registry_id,
  i.image_scan_status as image_scan_status
from
  aws_ecr_repository as r,
  aws_ecr_image as i
where
  r.repository_name = i.repository_name;
```

### List images with failed scans

```sql
select
  r.repository_name as repository_name,
  i.image_digest as image_digest,
  i.image_scan_status as image_scan_status
from
  aws_ecr_repository as r,
  aws_ecr_image as i
where
  r.repository_name = i.repository_name
  and i.image_scan_status ->> 'Status' = 'FAILED';
```

### List repositories whose tag immutability is disabled

```sql
select
  repository_name,
  image_tag_mutability
from
  aws_ecr_repository
where
  image_tag_mutability = 'IMMUTABLE';
```

### List repositories whose lifecycle policy rule is not configured to remove untagged and old images

```sql
select
  repository_name,
  r -> 'selection' ->> 'tagStatus' as tag_status,
  r -> 'selection' ->> 'countType' as count_type
from
  aws_ecr_repository,
  jsonb_array_elements(lifecycle_policy -> 'rules') as r
where
  (
    (r -> 'selection' ->> 'tagStatus' <> 'untagged')
    and (
      r -> 'selection' ->> 'countType' <> 'sinceImagePushed'
    )
  );
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
  aws_ecr_repository,
  jsonb_array_elements(policy -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and a in ('*', 'ecr:*');
```

### List repository scanning configuration settings

```sql
select
  repository_name,
  r ->> 'AppliedScanFilters' as applied_scan_filters,
  r ->> 'RepositoryArn' as repository_arn,
  r ->> 'ScanFrequency' as scan_frequency,
  r ->> 'ScanOnPush' as scan_on_push
from
  aws_ecr_repository,
  jsonb_array_elements(repository_scanning_configuration -> 'ScanningConfigurations') as r;

```

### List repositories where the scanning frequency is set to manual

```sql
select
  repository_name,
  r ->> 'RepositoryArn' as repository_arn,
  r ->> 'ScanFrequency' as scan_frequency
from
  aws_ecr_repository,
  jsonb_array_elements(repository_scanning_configuration -> 'ScanningConfigurations') as r
where
  r ->> 'ScanFrequency' = 'MANUAL';
```

### List repositories with scan-on-push is disabled

```sql
select
  repository_name,
  r ->> 'RepositoryArn' as repository_arn,
  r ->> 'ScanOnPush' as scan_on_push
from
  aws_ecr_repository,
  jsonb_array_elements(repository_scanning_configuration -> 'ScanningConfigurations') as r
where
 r ->> 'ScanOnPush' = 'false';
```
