---
title: "Steampipe Table: aws_glue_security_configuration - Query AWS Glue Security Configurations using SQL"
description: "Allows users to query AWS Glue Security Configurations and gain insights into the security configurations of Glue resources."
folder: "Config"
---

# Table: aws_glue_security_configuration - Query AWS Glue Security Configurations using SQL

The AWS Glue Security Configuration is a feature within AWS Glue service that allows you to specify the security settings that are used for Glue ETL jobs. This includes settings for data encryption, Amazon CloudWatch Logs encryption, and job bookmark encryption. It helps in maintaining the security and privacy of your data during the ETL (extract, transform, and load) process.

## Table Usage Guide

The `aws_glue_security_configuration` table in Steampipe provides you with information about security configurations within AWS Glue. This table allows you, as a DevOps engineer, to query security configuration-specific details, including encryption settings, CloudWatch encryption settings, Job Bookmarks encryption settings, and S3 encryption settings. You can utilize this table to gather insights on security configurations, such as the status of encryption settings, the type of encryption used, and more. The schema outlines the various attributes of the Glue security configuration for you, including the name, creation time, and encryption settings.

## Examples

### Basic info
Explore the security configurations of your AWS Glue service to assess the encryption status of different components such as Cloud Watch, job bookmarks, and S3. This can be useful in maintaining data security and compliance by ensuring appropriate encryption is in place.

```sql+postgres
select
  name,
  created_time_stamp,
  cloud_watch_encryption,
  job_bookmarks_encryption,
  s3_encryption
from
  aws_glue_security_configuration;
```

```sql+sqlite
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
Explore the encryption details of your CloudWatch logs to ensure data security. This is particularly useful for identifying instances where encryption has not been disabled, thereby providing an additional layer of security for your data.

```sql+postgres
select
  name,
  cloud_watch_encryption ->> 'CloudWatchEncryptionMode' as encyption_mode,
  cloud_watch_encryption ->> 'KmsKeyArn' as kms_key_arn
from
  aws_glue_security_configuration
where
  cloud_watch_encryption ->> 'CloudWatchEncryptionMode' != 'DISABLED';
```

```sql+sqlite
select
  name,
  json_extract(cloud_watch_encryption, '$.CloudWatchEncryptionMode') as encyption_mode,
  json_extract(cloud_watch_encryption, '$.KmsKeyArn') as kms_key_arn
from
  aws_glue_security_configuration
where
  json_extract(cloud_watch_encryption, '$.CloudWatchEncryptionMode') != 'DISABLED';
```

### List job bookmarks encryption details
Explore the encryption status of job bookmarks in AWS Glue Security Configurations, focusing on those with active encryption modes. This can be useful for maintaining data security and compliance by ensuring sensitive information is properly encrypted.

```sql+postgres
select
  name,
  job_bookmarks_encryption ->> 'JobBookmarksEncryptionMode' as encyption_mode,
  job_bookmarks_encryption ->> 'KmsKeyArn' as kms_key_arn
from
  aws_glue_security_configuration
where
  job_bookmarks_encryption ->> 'JobBookmarksEncryptionMode' != 'DISABLED';
```

```sql+sqlite
select
  name,
  json_extract(job_bookmarks_encryption, '$.JobBookmarksEncryptionMode') as encyption_mode,
  json_extract(job_bookmarks_encryption, '$.KmsKeyArn') as kms_key_arn
from
  aws_glue_security_configuration
where
  json_extract(job_bookmarks_encryption, '$.JobBookmarksEncryptionMode') != 'DISABLED';
```

### List s3 encryption details
Discover the segments that are using encryption within your AWS S3 storage. This is useful for maintaining security standards and ensuring sensitive data is properly protected.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(e.value, '$.S3EncryptionMode') as encyption_mode,
  json_extract(e.value, '$.KmsKeyArn') as kms_key_arn
from
  aws_glue_security_configuration,
  json_each(s3_encryption) as e
where
  json_extract(e.value, '$.S3EncryptionMode') != 'DISABLED';
```