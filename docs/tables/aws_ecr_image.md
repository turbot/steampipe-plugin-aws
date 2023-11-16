---
title: "Table: aws_ecr_image - Query Amazon ECR Images using SQL"
description: "Allows users to query Amazon Elastic Container Registry (ECR) Images and retrieve detailed information about each image, including image tags, push timestamps, image sizes, and more."
---

# Table: aws_ecr_image - Query Amazon ECR Images using SQL

The `aws_ecr_image` table in Steampipe provides information about Images within Amazon Elastic Container Registry (ECR). This table allows DevOps engineers to query image-specific details, including image tags, push timestamps, image sizes, and associated metadata. Users can utilize this table to gather insights on images, such as image scan findings, image vulnerability details, verification of image tags, and more. The schema outlines the various attributes of the ECR image, including the image digest, image tags, image scan status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecr_image` table, you can use the `.inspect aws_ecr_image` command in Steampipe.

### Key columns:

- `repository_name`: The name of the repository that the image resides in. This column is useful for joining with the `aws_ecr_repository` table.
- `image_digest`: The sha256 digest of the image manifest. This column is important as it uniquely identifies an image within a repository.
- `image_tags`: The list of image tags associated with the image. This column is useful for filtering and categorizing images.

## Examples

### Basic info

```sql
select
  repository_name,
  image_digest,
  image_pushed_at,
  image_size_in_bytes,
  registry_id,
  image_scan_status,
  image_tags
from
  aws_ecr_image;
```

### List image scan findings

```sql
select
  repository_name,
  image_scan_findings_summary ->> 'FindingSeverityCounts' as finding_severity_counts,
  image_scan_findings_summary ->> 'ImageScanCompletedAt' as image_scan_completed_at,
  image_scan_findings_summary ->> 'VulnerabilitySourceUpdatedAt' as vulnerability_source_updated_at
from
  aws_ecr_image;
```

### List image tags for the images

```sql
select
  repository_name,
  registry_id,
  image_digest,
  image_tags
from
  aws_ecr_image;
```

### List images pushed in last 10 days for a repository

```sql
select
  repository_name,
  image_digest,
  image_pushed_at,
  image_size_in_bytes
from
  aws_ecr_image
where
  image_pushed_at >= now() - interval '10' day
and
  repository_name = 'test1';
```

### List images for repositories created in the last 20 days

```sql
select
  i.repository_name as repository_name,
  r.repository_uri as repository_uri,
  i.image_digest as image_digest,
  i.image_tags as image_tags
from
  aws_ecr_image as i,
  aws_ecr_repository as r
where
  i.repository_name = r.repository_name
and
  r.created_at >= now() - interval '20' day;
```

### Get repository policy for each image's repository

```sql
select
  i.repository_name as repository_name,
  r.repository_uri as repository_uri,
  i.image_digest as image_digest,
  i.image_tags as image_tags,
  s ->> 'Effect' as effect,
  s ->> 'Action' as action,
  s ->> 'Condition' as condition,
  s ->> 'Principal' as principal
from
  aws_ecr_image as i,
  aws_ecr_repository as r,
  jsonb_array_elements(r.policy -> 'Statement') as s
where
  i.repository_name = r.repository_name;
```

### Scan images with trivy for a particular repository

```sql
select
  artifact_name,
  artifact_type,
  metadata,
  results
from
  trivy_scan_artifact as a,
  aws_ecr_image as i
where
  artifact_name = image_uri
  and repository_name = 'hello';
```
