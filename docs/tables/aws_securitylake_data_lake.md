---
title: "Steampipe Table: aws_securitylake_data_lake - Query AWS Lake Formation Data Lakes using SQL"
description: "Allows users to query AWS Lake Formation Data Lakes for information such as the Data Lake name, creation time, last modified time, and more."
folder: "Security Lake"
---

# Table: aws_securitylake_data_lake - Query AWS Lake Formation Data Lakes using SQL

The AWS Lake Formation is a service that makes it easy to set up, secure, and manage your data lakes. It simplifies the process of data ingestion, cataloging, transformation, and security. With Lake Formation, you can query your data using SQL, making it accessible for analysis and decision-making processes.

## Table Usage Guide

The `aws_securitylake_data_lake` table in Steampipe provides you with information about Data Lakes within AWS Lake Formation. This table allows you, as a DevOps engineer, to query Data Lake-specific details, including the Data Lake name, creation time, last modified time, and more. You can utilize this table to gather insights on Data Lakes, such as their creation times, last modified times, and other associated metadata. The schema outlines the various attributes of the Data Lake for you, including the Data Lake name, creation time, last modified time, and more.

## Examples

### Basic info
Determine the areas in which your AWS Security Lake data is being replicated and stored. This allows you to assess the status and security measures applied to your data storage and replication processes.

```sql+postgres
select
  kms_key_id,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake;
```

```sql+sqlite
select
  kms_key_id,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake;
```

### Get S3 bucket details of each data lake
Explore the security configurations of your data lakes by identifying the ones stored in public S3 buckets. This allows for a quick assessment of potential vulnerabilities and helps maintain proper data privacy standards.

```sql+postgres
select
  distinct b.name as bucket_name,
  l.s3_bucket_arn,
  b.creation_date,
  b.bucket_policy_is_public,
  b.versioning_enabled,
  b.block_public_acls
from
  aws_securitylake_data_lake as l,
  aws_s3_bucket as b
where
  l.s3_bucket_arn = b.arn;
```

```sql+sqlite
select
  distinct b.name as bucket_name,
  l.s3_bucket_arn,
  b.creation_date,
  b.bucket_policy_is_public,
  b.versioning_enabled,
  b.block_public_acls
from
  aws_securitylake_data_lake as l,
  aws_s3_bucket as b
where
  l.s3_bucket_arn = b.arn;
```

### Get retention setting details of data lake
Determine the status of your data lake's security settings, including encryption, replication, and storage details. This is particularly useful for managing data retention and ensuring optimal storage class configurations.

```sql+postgres
select
  l.replication_role_arn,
  l.s3_bucket_arn,
  l.status,
  r ->> 'RetentionPeriod' as retention_period,
  r ->> 'StorageClass' as storage_class
from
  aws_securitylake_data_lake as l,
  jsonb_array_elements(retention_settings) as r;
```

```sql+sqlite
select
  l.replication_role_arn,
  l.s3_bucket_arn,
  l.status,
  json_extract(r.value, '$.RetentionPeriod') as retention_period,
  json_extract(r.value, '$.StorageClass') as storage_class
from
  aws_securitylake_data_lake as l,
  json_each(l.retention_settings) as r;
```

### List data lakes where the configuration operation is in a pending state
Determine the areas in which data lakes are yet to complete their configuration process. This is beneficial in identifying and resolving potential delays in the setup of your data lakes.

```sql+postgres
select
  kms_key_id,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake
where
  status = 'PENDING';
```

```sql+sqlite
select
  kms_key_id,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake
where
  status = 'PENDING';
```