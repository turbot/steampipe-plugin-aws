# Table: aws_ecr_repository

Amazon Elastic Container Registry (Amazon ECR) is a managed container image registry service.

## Examples

### SSM association basic info

```sql
select
  repository_name,
  registry_id,
  repository_arn,
  repository_uri,
  created_at,
  region,
  account_id
from
  aws_ecr_repository;
```


### Image vulnerability count by severity. 

```sql
select
  repository_name,
  detail -> 'ImageScanFindingsSummary' -> 'FindingSeverityCounts' ->> 'INFORMATIONAL' as informational_severity_counts,
  detail -> 'ImageScanFindingsSummary' -> 'FindingSeverityCounts' ->> 'LOW' as low_severity_counts,
  detail -> 'ImageScanFindingsSummary' -> 'FindingSeverityCounts' ->> 'MEDIUM' as medium_severity_counts
from
  aws_ecr_repository,
  jsonb_array_elements(image_details) as details,
  jsonb(details) as detail;
```


### List of repositories whose image scanning has failed

```sql
select
  repository_name,
  detail -> 'ImageScanStatus' ->> 'Status' as scan_status
from
  aws_ecr_repository,
  jsonb_array_elements(image_details) as details,
  jsonb(details) as detail
where
  detail -> 'ImageScanStatus' ->> 'Status' = 'FAILED';
```


### List of ECR repositories which are not using Customer Managed Keys(CMK) for encryption

```sql
select
  repository_name,
  encryption_configuration ->> 'EncryptionType' as encryption_type,
  encryption_configuration ->> 'KmsKey' as kms_key
from
  aws_ecr_repository
where
  encryption_configuration ->> 'EncryptionType' = 'AES256';
```


### List Repositories whose tag immutability is disabled

```sql
select
  repository_name,
  image_tag_mutability
from
  aws_ecr_repository
where
  image_tag_mutability = 'IMMUTABLE';
```


### List of repositories where automatic image scanning is disabled

```sql
select
  repository_name,
  image_scanning_configuration ->> 'ScanOnPush' as scan_on_push
from
  aws_ecr_repository
where
  image_scanning_configuration ->> 'ScanOnPush' = 'false';
```
