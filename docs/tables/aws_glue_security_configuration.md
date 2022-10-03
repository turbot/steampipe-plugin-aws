# Table: aws_glue_security_configuration

A security configuration in AWS Glue contains the properties that are needed when you write encrypted data. You create security configurations on the AWS Glue console to provide the encryption properties that are used by crawlers, jobs, and development endpoints.

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