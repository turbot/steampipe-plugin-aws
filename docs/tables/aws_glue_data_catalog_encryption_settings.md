---
title: "Steampipe Table: aws_glue_data_catalog_encryption_settings - Query AWS Glue Data Catalog using SQL"
description: "Allows users to query AWS Glue Data Catalog Encryption Settings."
folder: "Glue"
---

# Table: aws_glue_data_catalog_encryption_settings - Query AWS Glue Data Catalog using SQL

The AWS Glue Data Catalog is a fully managed, scalable, Apache Hive Metastore compatible, metadata repository. It provides a uniform repository where disparate systems can store and find metadata to keep track of data, and it makes this data available for ETL jobs and data queries. The Encryption Settings for the AWS Glue Data Catalog contain settings used to protect catalog resources with encryption.

## Table Usage Guide

The `aws_glue_data_catalog_encryption_settings` table in Steampipe provides you with information about the encryption settings of AWS Glue Data Catalogs. This table allows you, as a DevOps engineer or security analyst, to query encryption-specific details, including the encryption-at-rest settings and the return connection password encryption settings. You can utilize this table to gather insights on the encryption settings of your data catalogs, such as understanding the type of encryption used, the AWS KMS key ID used for encryption, and more. The schema outlines the various attributes of the encryption settings for you, including the catalog ID, create time, update time, and associated metadata.

## Examples

### Basic info
Analyze the settings to understand the encryption status and location of your AWS Glue Data Catalog. This is useful for maintaining data security and ensuring compliance with regional data regulations.

```sql+postgres
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings;
```

```sql+sqlite
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings;
```

### List settings where encryption at rest is disabled
Determine the areas in which encryption at rest is disabled to enhance security measures and protect sensitive data within your AWS Glue Data Catalog.

```sql+postgres
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

```sql+sqlite
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings
where
  json_extract(encryption_at_rest, '$.CatalogEncryptionMode') = 'DISABLED';
```

### List settings where connection password encryption is disabled
Discover the segments where connection password encryption is not enabled in the AWS Glue Data Catalog. This query is particularly useful for identifying potential security vulnerabilities related to password protection.

```sql+postgres
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

```sql+sqlite
select
  encryption_at_rest,
  connection_password_encryption,
  region,
  account_id
from
  aws_glue_data_catalog_encryption_settings
where
  json_extract(connection_password_encryption, '$.ReturnConnectionPasswordEncrypted') = 'false';
```

### List encryption at rest key details associated to settings
Identify the key details of encryption at rest associated with specific settings. This can help in assessing security measures and managing data protection strategies.

```sql+postgres
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

```sql+sqlite
select
  json_extract(encryption_at_rest, '$.SseAwsKmsKeyId') as key_arn,
  k.key_manager as key_manager,
  k.creation_date as key_creation_date,
  s.region,
  s.account_id
from
  aws_glue_data_catalog_encryption_settings s
  join aws_kms_key k on json_extract(s.encryption_at_rest, '$.SseAwsKmsKeyId') = k.arn
  and s.region = k.region;
```

### List connection password encryption key details associated to settings
Determine the areas in which the encryption key details are associated with certain settings, allowing for a comprehensive review of security measures across different regions and accounts. This query is particularly useful for understanding the management and creation date of encryption keys, contributing to enhanced data protection efforts.

```sql+postgres
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

```sql+sqlite
select
  json_extract(connection_password_encryption, '$.AwsKmsKeyId') as key_arn,
  k.key_manager as key_manager,
  k.creation_date as key_creation_date,
  s.region,
  s.account_id
from
  aws_glue_data_catalog_encryption_settings s
  join aws_kms_key k on json_extract(s.connection_password_encryption, '$.AwsKmsKeyId') = k.arn
  and s.region = k.region;
```