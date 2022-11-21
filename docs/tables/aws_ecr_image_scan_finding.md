# Table: aws_ecr_image_scan_finding

Amazon Elastic Container Registry (Amazon ECR) stores Docker images and allows you to scan them on push, or periodically.
The corresponding CVE findings are available in this table for an image tag, in a repository.

**Note**: Users or roles that have the AWS managed `ReadOnlyAccess` policy attached also need to attach the AWS managed `AmazonInspector2ReadOnlyAccess` policy to query this table.

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
