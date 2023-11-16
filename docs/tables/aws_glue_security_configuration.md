---
title: "Table: aws_glue_security_configuration - Query AWS Glue Security Configurations using SQL"
description: "Allows users to query AWS Glue Security Configurations and gain insights into the security configurations of Glue resources."
---

# Table: aws_glue_security_configuration - Query AWS Glue Security Configurations using SQL

The `aws_glue_security_configuration` table in Steampipe provides information about security configurations within AWS Glue. This table allows DevOps engineers to query security configuration-specific details, including encryption settings, CloudWatch encryption settings, Job Bookmarks encryption settings, and S3 encryption settings. Users can utilize this table to gather insights on security configurations, such as the status of encryption settings, the type of encryption used, and more. The schema outlines the various attributes of the Glue security configuration, including the name, creation time, and encryption settings.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_security_configuration` table, you can use the `.inspect aws_glue_security_configuration` command in Steampipe.

### Key columns:

- `name`: The name of the security configuration. It can be used as a unique identifier for the security configuration and can be used to join with other tables that reference the security configuration by name.
- `encryption_configuration_s3_encryption_s3_encryption_mode`: The encryption mode used for S3 data. This can be used to understand the level of data protection applied to S3 data.
- `encryption_configuration_cloud_watch_encryption_cloud_watch_encryption_mode`: The encryption mode used for CloudWatch logs. This can be used to understand the level of data protection applied to CloudWatch logs.

## Examples

### Basic info

```sql
select
  name,
  created_time_stamp,
  cloud_watch_encryption,
  job_bookmarks_encryption,
  s3_encryption
from
  aws_glue_security_configuration;
```

### List cloud watch encryption details

```sql
select
  name,
  cloud_watch_encryption ->> 'CloudWatchEncryptionMode' as encyption_mode,
  cloud_watch_encryption ->> 'KmsKeyArn' as kms_key_arn
from
  aws_glue_security_configuration
where
  cloud_watch_encryption ->> 'CloudWatchEncryptionMode' != 'DISABLED';
```

### List job bookmarks encryption details

```sql
select
  name,
  job_bookmarks_encryption ->> 'JobBookmarksEncryptionMode' as encyption_mode,
  job_bookmarks_encryption ->> 'KmsKeyArn' as kms_key_arn
from
  aws_glue_security_configuration
where
  job_bookmarks_encryption ->> 'JobBookmarksEncryptionMode' != 'DISABLED';
```

### List s3 encryption details

```sql
select
  name,
  e ->> 'S3EncryptionMode' as encyption_mode,
  e ->> 'KmsKeyArn' as kms_key_arn
from
  aws_glue_security_configuration,
  jsonb_array_elements(s3_encryption) e
where
  e ->> 'S3EncryptionMode' != 'DISABLED';
```