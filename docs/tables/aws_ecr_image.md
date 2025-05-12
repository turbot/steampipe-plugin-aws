---
title: "Steampipe Table: aws_ecr_image - Query Amazon ECR Images using SQL"
description: "Allows users to query Amazon Elastic Container Registry (ECR) Images and retrieve detailed information about each image, including image tags, push timestamps, image sizes, and more."
folder: "ECR"
---

# Table: aws_ecr_image - Query Amazon ECR Images using SQL

The Amazon Elastic Container Registry (ECR) Images are Docker images that are stored within AWS's managed and highly available registry. ECR Images allow you to easily store, manage, and deploy Docker container images in a secure environment. They are integrated with AWS Identity and Access Management (IAM) for resource-level control and support for private Docker repositories.

## Table Usage Guide

The `aws_ecr_image` table in Steampipe provides you with information about Images within Amazon Elastic Container Registry (ECR). This table allows you, as a DevOps engineer, to query image-specific details, including image tags, push timestamps, image sizes, and associated metadata. You can utilize this table to gather insights on images, such as image scan findings, image vulnerability details, verification of image tags, and more. The schema outlines the various attributes of the ECR image for you, including the image digest, image tags, image scan status, and associated tags.

## Examples

### Basic info
Explore the details of your AWS Elastic Container Registry (ECR) images, like when they were last updated and their size, to better manage your resources. This can help in identifying outdated or oversized images, thus optimizing your ECR utilization.

```sql+postgres
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

```sql+sqlite
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
Identify instances where your repository images might have vulnerabilities by examining the severity of scan findings. This allows you to assess the security of your images and take necessary actions based on the severity of the findings.

```sql+postgres
select
  repository_name,
  image_scan_findings_summary ->> 'FindingSeverityCounts' as finding_severity_counts,
  image_scan_findings_summary ->> 'ImageScanCompletedAt' as image_scan_completed_at,
  image_scan_findings_summary ->> 'VulnerabilitySourceUpdatedAt' as vulnerability_source_updated_at
from
  aws_ecr_image;
```

```sql+sqlite
select
  repository_name,
  json_extract(image_scan_findings_summary, '$.FindingSeverityCounts') as finding_severity_counts,
  json_extract(image_scan_findings_summary, '$.ImageScanCompletedAt') as image_scan_completed_at,
  json_extract(image_scan_findings_summary, '$.VulnerabilitySourceUpdatedAt') as vulnerability_source_updated_at
from
  aws_ecr_image;
```

### List image tags for the images
Explore which image tags are associated with the images in your AWS ECR repositories. This can help you manage and organize your resources more effectively.

```sql+postgres
select
  repository_name,
  registry_id,
  image_digest,
  image_tags
from
  aws_ecr_image;
```

```sql+sqlite
select
  repository_name,
  registry_id,
  image_digest,
  image_tags
from
  aws_ecr_image;
```

### List images pushed in last 10 days for a repository
Determine the images that have been uploaded to a specific repository in the last 10 days. This is useful for tracking recent updates or additions to the repository.

```sql+postgres
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

```sql+sqlite
select
  repository_name,
  image_digest,
  image_pushed_at,
  image_size_in_bytes
from
  aws_ecr_image
where
  image_pushed_at >= datetime('now','-10 day')
and
  repository_name = 'test1';
```

### List images for repositories created in the last 20 days
Explore recently created repositories and the images they contain. This query is useful for keeping track of new content and managing resources within a 20-day timeframe.

```sql+postgres
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

```sql+sqlite
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
  r.created_at >= datetime('now', '-20 days');
```

### Get repository policy for each image's repository
Determine the access policies associated with each image's repository in AWS Elastic Container Registry (ECR). This can help to identify potential security risks, such as open access to sensitive images.

```sql+postgres
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

```sql+sqlite
select
  i.repository_name as repository_name,
  r.repository_uri as repository_uri,
  i.image_digest as image_digest,
  i.image_tags as image_tags,
  json_extract(s.value, '$.Effect') as effect,
  json_extract(s.value, '$.Action') as action,
  json_extract(s.value, '$.Condition') as condition,
  json_extract(s.value, '$.Principal') as principal
from
  aws_ecr_image as i,
  aws_ecr_repository as r,
  json_each(r.policy, '$.Statement') as s
where
  i.repository_name = r.repository_name;
```