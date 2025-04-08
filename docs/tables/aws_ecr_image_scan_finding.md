---
title: "Steampipe Table: aws_ecr_image_scan_finding - Query Amazon Elastic Container Registry (ECR) Image Scan Findings using SQL"
description: "Allows users to query Amazon ECR Image Scan Findings to retrieve detailed information about image scan findings, including attributes such as the severity of the finding, description, and package name where the vulnerability was found."
folder: "ECR"
---

# Table: aws_ecr_image_scan_finding - Query Amazon Elastic Container Registry (ECR) Image Scan Findings using SQL

The Amazon Elastic Container Registry (ECR) Image Scan Findings is a feature of AWS ECR that allows you to identify any software vulnerabilities in your Docker images. It uses the Common Vulnerabilities and Exposures (CVEs) database from the open-source Clair project. It provides detailed findings, severity levels, and a description of the vulnerabilities.

## Table Usage Guide

The `aws_ecr_image_scan_finding` table in Steampipe provides you with information about Image Scan Findings within Amazon Elastic Container Registry (ECR). This table allows you, as a DevOps engineer, to query specific details about image scan findings, including attributes such as the severity of the finding, description, and package name where the vulnerability was found. You can utilize this table to gather insights on image scan findings, such as identifying high-risk vulnerabilities, verifying package vulnerabilities, and more. The schema outlines the various attributes of the Image Scan Finding for you, including the repository name, image digest, finding severity, and associated metadata.

**Important Notes**
- You or your roles that have the AWS managed `ReadOnlyAccess` policy attached also need to attach the AWS managed `AmazonInspector2ReadOnlyAccess` policy to query this table.

## Examples

### List scan findings for an image
Identify potential vulnerabilities in a specific image within a repository. This assists in enhancing the security by highlighting areas of concern and providing insights into the severity and nature of the detected issues.

```sql+postgres
select
  repository_name,
  image_tag,
  name,
  severity,
  description,
  attributes,
  uri,
  image_scan_status,
  image_scan_completed_at,
  vulnerability_source_updated_at
from
  aws_ecr_image_scan_finding
where
  repository_name = 'my-repo'
  and image_tag = 'my-image-tag';
```

```sql+sqlite
select
  repository_name,
  image_tag,
  name,
  severity,
  description,
  attributes,
  uri,
  image_scan_status,
  image_scan_completed_at,
  vulnerability_source_updated_at
from
  aws_ecr_image_scan_finding
where
  repository_name = 'my-repo'
  and image_tag = 'my-image-tag';
```

### Get CVEs for all images pushed in the last 24 hours
Explore potential vulnerabilities in your system by identifying Common Vulnerabilities and Exposures (CVEs) in all images that have been pushed in the last 24 hours. This is particularly useful for maintaining system security and identifying areas that may need immediate attention or updates.

```sql+postgres
select
  f.repository_name,
  f.image_tag,
  f.name,
  f.severity,
  jsonb_pretty(f.attributes) as attributes
from
  (
    select
      repository_name,
      jsonb_array_elements_text(image_tags) as image_tag
    from
      aws_ecr_image as i
    where
      i.image_pushed_at > now() - interval '24' hour
  )
  images
  left outer join
    aws_ecr_image_scan_finding as f
    on images.repository_name = f.repository_name
    and images.image_tag = f.image_tag;
```

```sql+sqlite
select
  f.repository_name,
  f.image_tag,
  f.name,
  f.severity,
  f.attributes as attributes
from
  (
    select
      repository_name,
      json_each.value as image_tag
    from
      aws_ecr_image as i,
      json_each(i.image_tags)
    where
      i.image_pushed_at > datetime('now', '-24 hours')
  )
  images
  left outer join
    aws_ecr_image_scan_finding as f
    on images.repository_name = f.repository_name
    and images.image_tag = f.image_tag;
```