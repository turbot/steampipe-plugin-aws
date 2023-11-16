---
title: "Table: aws_ecr_image_scan_finding - Query Amazon Elastic Container Registry (ECR) Image Scan Findings using SQL"
description: "Allows users to query Amazon ECR Image Scan Findings to retrieve detailed information about image scan findings, including attributes such as the severity of the finding, description, and package name where the vulnerability was found."
---

# Table: aws_ecr_image_scan_finding - Query Amazon Elastic Container Registry (ECR) Image Scan Findings using SQL

The `aws_ecr_image_scan_finding` table in Steampipe provides information about Image Scan Findings within Amazon Elastic Container Registry (ECR). This table allows DevOps engineers to query specific details about image scan findings, including attributes such as the severity of the finding, description, and package name where the vulnerability was found. Users can utilize this table to gather insights on image scan findings, such as identifying high-risk vulnerabilities, verifying package vulnerabilities, and more. The schema outlines the various attributes of the Image Scan Finding, including the repository name, image digest, finding severity, and associated metadata.

**Note**: Users or roles that have the AWS managed `ReadOnlyAccess` policy attached also need to attach the AWS managed `AmazonInspector2ReadOnlyAccess` policy to query this table.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecr_image_scan_finding` table, you can use the `.inspect aws_ecr_image_scan_finding` command in Steampipe.

**Key columns**:

- `repository_name`: This column is useful as it contains the name of the ECR repository. This can be used to join this table with other tables that contain repository-specific information.
- `image_digest`: This column is important because it contains the image digest for the scanned image. This can be used to join this table with other tables that contain image-specific information.
- `finding_severity`: This column is useful as it contains the severity of the finding. This can be used to filter or sort the findings based on their severity.

## Examples

### List scan findings for an image

```sql
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
  and image_tag = 'my-image-tag'
```

### Get CVEs for all images pushed in the last 24 hours

```sql
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
    and images.image_tag = f.image_tag
```
