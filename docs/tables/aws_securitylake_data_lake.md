---
title: "Table: aws_securitylake_data_lake - Query AWS Lake Formation Data Lakes using SQL"
description: "Allows users to query AWS Lake Formation Data Lakes for information such as the Data Lake name, creation time, last modified time, and more."
---

# Table: aws_securitylake_data_lake - Query AWS Lake Formation Data Lakes using SQL

The `aws_securitylake_data_lake` table in Steampipe provides information about Data Lakes within AWS Lake Formation. This table allows DevOps engineers to query Data Lake-specific details, including Data Lake name, creation time, last modified time, and more. Users can utilize this table to gather insights on Data Lakes, such as Data Lake creation times, last modified times, and other associated metadata. The schema outlines the various attributes of the Data Lake, including the Data Lake name, creation time, last modified time, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securitylake_data_lake` table, you can use the `.inspect aws_securitylake_data_lake` command in Steampipe.

### Key columns:

- `name`: The name of the Data Lake. This is a key identifier and can be used to join this table with other tables that also contain Data Lake names.
- `created_time`: The time when the Data Lake was created. This can be useful for tracking the age of Data Lakes and identifying any that may be outdated or unused.
- `last_modified_time`: The last time the Data Lake was modified. This can be useful for tracking recent changes or modifications to the Data Lake.

## Examples

### Basic info

```sql
select
  encryption_key,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake;
```

### Get S3 bucket details of each data lake

```sql
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

```sql
select
  l.encryption_key,
  l.replication_role_arn,
  l.s3_bucket_arn,
  l.status,
  r ->> 'RetentionPeriod' as retention_period,
  r ->> 'StorageClass' as storage_class
from
  aws_securitylake_data_lake as l,
  jsonb_array_elements(retention_settings) as r;
```

### List data lakes where the configuration operation is in a pending state

```sql
select
  encryption_key,
  replication_role_arn,
  s3_bucket_arn,
  status
from
  aws_securitylake_data_lake
where
  status = 'PENDING';
```
