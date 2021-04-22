# Table: aws_ecrpublic_repository

A public repository is open to publicly pull images from and is visible on the Amazon ECR Public Gallery. When creating a public repository you specify catalog data which helps users find and use your images.

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
  aws_ecrpublic_repository;
```

### List repositories whose image scanning has failed

```sql
select
  repository_name,
  detail -> 'ImageScanStatus' ->> 'Status' as scan_status
from
  aws_ecrpublic_repository,
  jsonb_array_elements(image_details) as details,
  jsonb(details) as detail
where
  detail -> 'ImageScanStatus' ->> 'Status' = 'FAILED';
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
  aws_ecrpublic_repository,
  jsonb_array_elements(policy -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and a in ('*', 'ecr-public:*');
```

### Get repository image vulnerability count by severity for each repository

```sql
select
  repository_name,
  detail -> 'ImageScanFindingsSummary' -> 'FindingSeverityCounts' ->> 'INFORMATIONAL' as informational_severity_counts,
  detail -> 'ImageScanFindingsSummary' -> 'FindingSeverityCounts' ->> 'LOW' as low_severity_counts,
  detail -> 'ImageScanFindingsSummary' -> 'FindingSeverityCounts' ->> 'MEDIUM' as medium_severity_counts
from
  aws_ecrpublic_repository,
  jsonb_array_elements(image_details) as details,
  jsonb(details) as detail;
```
