---
title: "Table: aws_glue_data_catalog_encryption_settings - Query AWS Glue Data Catalog using SQL"
description: "Allows users to query AWS Glue Data Catalog Encryption Settings."
---

# Table: aws_glue_data_catalog_encryption_settings - Query AWS Glue Data Catalog using SQL

The `aws_glue_data_catalog_encryption_settings` table in Steampipe provides information about the encryption settings of AWS Glue Data Catalogs. This table allows DevOps engineers and security analysts to query encryption-specific details, including the encryption-at-rest settings and the return connection password encryption settings. Users can utilize this table to gather insights on the encryption settings of their data catalogs, such as understanding the type of encryption used, the AWS KMS key ID used for encryption, and more. The schema outlines the various attributes of the encryption settings, including the catalog ID, create time, update time, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_data_catalog_encryption_settings` table, you can use the `.inspect aws_glue_data_catalog_encryption_settings` command in Steampipe.

### Key columns:

- `catalog_id`: The ID of the data catalog. This is a key column that can be used to join this table with other tables to get more detailed information about the data catalog.
- `create_time`: The time at which the encryption settings were created. This can be useful for auditing and tracking changes over time.
- `update_time`: The time at which the encryption settings were last updated. This can be useful for monitoring and ensuring that encryption settings are up-to-date.

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