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
