# Table: aws_ecr_image_finding

Amazon Elastic Container Registry (Amazon ECR) stores Docker images and allows you to scan them on push, or periodically.
The corresponding CVE findings are available in this table for an image tag, in a repository.

## Examples

### Get CVEs count from an image tag scan, group by severity

```sql
select
  repository_name,
  image_tag,
  image_digest,
  image_scan_status,
  image_scan_completed_at,
  vulnerability_source_updated_at,
  severity,
  count(*) as nb 
from
  aws_ecr_image_scan_finding 
where
  repository_name = 'my-repo' 
  and image_tag = 'my-image-tag' 
group by
  1,
  2,
  3,
  4,
  5,
  6,
  7
```

### Get details from an image tag scan

```sql
select
  repository_name,
  image_tag,
  image_digest,
  image_scan_status,
  image_scan_completed_at,
  vulnerability_source_updated_at,
  severity,
  name,
  description,
  attributes,
  uri 
from
  aws_ecr_image_scan_finding 
where
  repository_name = 'my-repo' 
  and image_tag = 'my-image-tag'
```

### Get CVEs for all images pushed in the last 24 hours

```sql
select
  findings.* 
from
  (
    select
      repository_name,
      jsonb_array_elements_text(image_tags) as image_tag 
    from
      aws_ecr_image i 
    where
      i.image_pushed_at > now() - interval '24' hour 
  )
  images 
  left outer join
    aws_ecr_image_scan_finding findings 
    on images.repository_name = findings.repository_name 
    and images.image_tag = findings.image_tag
```
