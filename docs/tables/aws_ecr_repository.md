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

### List of ECR repositories which are not using Customer Managed Keys(CMK) for encryption
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

### Count of image details in each Repository
```sql
select
  count(image_details) as count
from
  aws_ecr_repository
where
  repository_name = 'steampipe-12';
```

### List of unencrypted ECR repositories
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

