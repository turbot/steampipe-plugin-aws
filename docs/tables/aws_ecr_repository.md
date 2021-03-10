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
### List Repositories whose imageTag set as MUTABLE

```sql
select
	repository_name,
	image_tag_mutability
from
	aws_ecr_repository
where
	image_tag_mutability = 'MUTABLE';
```

### List Repositories who has image scan configuration with both CRITICAL findings 

```sql
select
  repository_name,
	image_scanning_configuration as scan
from
	aws_ecr_repository
where
	finding-severity-counts= 'CRITICAL';
```
### Count of image ids in a requested Repository

```sql
select
	count(image_ids) as count
from
	aws_ecr_repository
where
	repository_name = 'steampipe-12';
```
### List of unencrypted ECR repository

```sql
select
	repository_name,
	encryption_configuration ->> 'EncryptionType' as encryption_type,
	encryption_configuration ->> 'KmsKey' as kms_key
from
	aws_ecr_repository
where
	encryption_configuration ->> 'KmsKey' = 'null';
```

### List of ECR repository which are not using Customer Managed Keys(CMK)

```sql
select
	repository_name,
	encryption_configuration ->> 'EncryptionType' as encryption_type,
	encryption_configuration ->> 'KmsKey' as kms_key
from
	aws_ecr_repository
where
	encryption_configuration ->> 'KmsKey' = 'KMS';
```