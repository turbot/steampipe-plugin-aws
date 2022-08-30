# Table: aws_ecr_image

Amazon Elastic Container Registry (Amazon ECR) stores Docker images, Open Container Initiative (OCI) images, and OCI compatible artifacts in private repositories. You can use the Docker CLI, or your preferred client, to push and pull images to and from your repositories.

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