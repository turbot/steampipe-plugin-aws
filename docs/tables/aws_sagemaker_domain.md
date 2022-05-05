# Table: aws_sagemaker_domain

An Amazon SageMaker Domain consists of an associated Amazon Elastic File System (Amazon EFS) volume; a list of authorized users; and a variety of security, application, policy, and Amazon Virtual Private Cloud (Amazon VPC) configurations.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_domain;
```

### List sagemaker domains where efs volume is unencrypted

```sql
select
  name,
  creation_time,
  home_efs_file_system_id,
  kms_key_id
from
  aws_sagemaker_domain
where 
  kms_key_id is null;
```

### List publicly accessible sagemaker domains

```sql
select
  name,
  arn,
  creation_time,
  app_network_access_type
from
  aws_sagemaker_domain
where 
  app_network_access_type = 'PublicInternetOnly';
  ```
  