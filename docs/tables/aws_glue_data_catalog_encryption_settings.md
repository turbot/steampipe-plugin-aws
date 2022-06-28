# Table: aws_glue_data_catalog_encryption_settings

The Data Catalog Encryption Settings in AWS Glue maintains Data Catalog resources security. 

## Examples

### Basic info

```sql
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings;
```

### List settings where encryption at rest is disabled

```sql
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings
where
  encryption_at_rest ->> 'CatalogEncryptionMode' = 'DISABLED';
```

### List settings where connection password encryption is disabled

```sql
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings
where
  connection_password_encryption ->> 'ReturnConnectionPasswordEncrypted' = 'false';
```

### List encryption at rest key details associated to settings

```sql
select
  encryption_at_rest ->> 'SseAwsKmsKeyId' as key_arn,
  k.key_manager as key_manager,
  k.creation_date as key_creation_date,
  s.region,
  s.account_id
from
  aws_glue_data_catalog_encryption_settings s
  join aws_kms_key k on s.encryption_at_rest ->> 'SseAwsKmsKeyId' = k.arn
  and s.region = k.region;
```

### List connection password encryption key details associated to settings

```sql
select
  connection_password_encryption ->> 'AwsKmsKeyId' as key_arn,
  k.key_manager as key_manager,
  k.creation_date as key_creation_date,
  s.region,
  s.account_id
from
  aws_glue_data_catalog_encryption_settings s
  join aws_kms_key k on s.connection_password_encryption ->> 'AwsKmsKeyId' = k.arn
  and s.region = k.region;
```