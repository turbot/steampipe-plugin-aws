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