---
title: "Steampipe Table: aws_ecr_repository - Query AWS ECR Repositories using SQL"
description: "Allows users to query AWS Elastic Container Registry (ECR) Repositories and retrieve detailed information about each repository."
folder: "ECR"
---

# Table: aws_ecr_repository - Query AWS ECR Repositories using SQL

The AWS ECR Repository is a managed docker container registry service provided by Amazon Web Services. It makes it easy for developers to store, manage, and deploy Docker container images. Amazon ECR eliminates the need to operate your own container repositories or worry about scaling the underlying infrastructure.

## Table Usage Guide

The `aws_ecr_repository` table in Steampipe provides you with information about repositories within AWS Elastic Container Registry (ECR). This table allows you, as a DevOps engineer, to query repository-specific details, including repository ARN, repository URI, and creation date. You can utilize this table to gather insights on repositories, such as repository policies, image scanning configurations, image tag mutability, and more. The schema outlines the various attributes of the ECR repository for you, including the repository name, creation date, and associated tags.

## Examples

### Basic info
Explore which Elastic Container Registry (ECR) repositories are available in your AWS account and determine their associated details such as creation date and region. This can be beneficial in managing your repositories and understanding their distribution across different regions.

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
  aws_ecr_repository;
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
  aws_ecr_repository;
```

### List repositories which are not using Customer Managed Keys (CMK) for encryption
Determine the areas in which repositories are not utilizing Customer Managed Keys for encryption. This is useful for enhancing security measures by identifying potential vulnerabilities in your encryption methods.

```sql+postgres
select
  repository_name,
  encryption_configuration ->> 'EncryptionType' as encryption_type,
  encryption_configuration ->> 'KmsKey' as kms_key
from
  aws_ecr_repository
where
  encryption_configuration ->> 'EncryptionType' = 'AES256';
```

```sql+sqlite
select
  repository_name,
  json_extract(encryption_configuration, '$.EncryptionType') as encryption_type,
  json_extract(encryption_configuration, '$.KmsKey') as kms_key
from
  aws_ecr_repository
where
  json_extract(encryption_configuration, '$.EncryptionType') = 'AES256';
```

### List repositories with automatic image scanning disabled
Identify instances where automatic image scanning is disabled in repositories. This is useful to ensure security measures are consistently applied across all repositories.

```sql+postgres
select
  repository_name,
  image_scanning_configuration ->> 'ScanOnPush' as scan_on_push
from
  aws_ecr_repository
where
  image_scanning_configuration ->> 'ScanOnPush' = 'false';
```

```sql+sqlite
select
  repository_name,
  json_extract(image_scanning_configuration, '$.ScanOnPush') as scan_on_push
from
  aws_ecr_repository
where
  json_extract(image_scanning_configuration, '$.ScanOnPush') = 'false';
```

### List images for each repository
Determine the images associated with each repository to understand their size, push time, last pull time, and scan status. This can help in managing repository resources, tracking image usage, and ensuring security compliance.

```sql+postgres
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

```sql+sqlite
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
  aws_ecr_repository as r
join
  aws_ecr_image as i
on
  r.repository_name = i.repository_name;
```

### List images with failed scans
Identify instances where image scans have failed in your AWS ECR repositories. This can help in diagnosing and rectifying issues related to image scanning, thereby improving the security and reliability of your container images.

```sql+postgres
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

```sql+sqlite
select
  r.repository_name as repository_name,
  i.image_digest as image_digest,
  json_extract(i.image_scan_status, '$.Status') as image_scan_status
from
  aws_ecr_repository as r
join
  aws_ecr_image as i
on
  r.repository_name = i.repository_name
where
  json_extract(i.image_scan_status, '$.Status') = 'FAILED';
```

### List repositories whose tag immutability is disabled
Determine the areas in which image tag immutability is disabled within your repositories. This allows you to identify and manage potential vulnerabilities in your AWS Elastic Container Registry.

```sql+postgres
select
  repository_name,
  image_tag_mutability
from
  aws_ecr_repository
where
  image_tag_mutability = 'IMMUTABLE';
```

```sql+sqlite
select
  repository_name,
  image_tag_mutability
from
  aws_ecr_repository
where
  image_tag_mutability = 'IMMUTABLE';
```

### List repositories whose lifecycle policy rule is not configured to remove untagged and old images
Determine the areas in which repositories are not configured to automatically clean up untagged and old images. This can help in managing storage and avoiding unnecessary costs associated with unused or outdated images.

```sql+postgres
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

```sql+sqlite
select
  repository_name,
  json_extract(r.value, '$.selection.tagStatus') as tag_status,
  json_extract(r.value, '$.selection.countType') as count_type
from
  aws_ecr_repository,
  json_each(lifecycle_policy, 'rules') as r
where
  (
    (json_extract(r.value, '$.selection.tagStatus') <> 'untagged')
    and (
      json_extract(r.value, '$.selection.countType') <> 'sinceImagePushed'
    )
  );
```

### List repository policy statements that grant full access for each repository
Identify instances where full access has been granted to each repository. This is useful to review and manage access permissions, ensuring optimal security and control over your data repositories.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(p.value, '$') as principal,
  json_extract(a.value, '$') as action,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Condition') as conditions
from
  aws_ecr_repository,
  json_each(policy, '$.Statement') as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as p,
  json_each(json_extract(s.value, '$.Action')) as a
where
  json_extract(s.value, '$.Effect') = 'Allow'
  and (
    json_extract(a.value, '$') = '*'
    or json_extract(a.value, '$') = 'ecr:*'
  );
```

### List repository scanning configuration settings
Determine the frequency and triggers for scanning within your repositories to optimize security checks and resource management. This enables you to understand the efficiency and effectiveness of your scanning configurations.

```sql+postgres
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

```sql+sqlite
select
  repository_name,
  json_extract(r.value, '$.AppliedScanFilters') as applied_scan_filters,
  json_extract(r.value, '$.RepositoryArn') as repository_arn,
  json_extract(r.value, '$.ScanFrequency') as scan_frequency,
  json_extract(r.value, '$.ScanOnPush') as scan_on_push
from
  aws_ecr_repository,
  json_each(repository_scanning_configuration, '$.ScanningConfigurations') as r;
```

### List repositories where the scanning frequency is set to manual
Determine the areas in your AWS ECR repositories where the scanning frequency is manually set. This allows you to identify instances where automated scanning is not enabled, potentially leaving your repositories vulnerable to undetected issues.

```sql+postgres
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

```sql+sqlite
select
  repository_name,
  json_extract(r.value, '$.RepositoryArn') as repository_arn,
  json_extract(r.value, '$.ScanFrequency') as scan_frequency
from
  aws_ecr_repository,
  json_each(repository_scanning_configuration, '$.ScanningConfigurations') as r
where
  json_extract(r.value, '$.ScanFrequency') = 'MANUAL';
```

### List repositories with scan-on-push is disabled
Identify instances where the scan-on-push feature is disabled in your repositories. This can help improve your security measures by ensuring all repositories are scanned for vulnerabilities upon each push.

```sql+postgres
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

```sql+sqlite
select
  repository_name,
  json_extract(r.value, '$.RepositoryArn') as repository_arn,
  json_extract(r.value, '$.ScanOnPush') as scan_on_push
from
  aws_ecr_repository,
  json_each(repository_scanning_configuration, 'ScanningConfigurations') as r
where
 json_extract(r.value, '$.ScanOnPush') = 'false';
```