---
title: "Steampipe Table: aws_s3_directory_bucket - Query AWS S3 Directory Buckets using SQL"
description: "Allows users to query AWS S3 directory buckets for detailed information about their configuration and properties."
folder: "S3"
---

# Table: aws_s3_directory_bucket - Query AWS S3 Directory Buckets using SQL

An AWS S3 Directory Bucket is a high-performance storage option for applications that need consistent single-digit millisecond data access. Directory buckets are designed for applications that require high throughput and low latency, such as machine learning training, real-time analytics, and high-performance computing workloads.

## Table Usage Guide

The `aws_s3_directory_bucket` table in Steampipe provides you with information about S3 directory buckets within Amazon Simple Storage Service (S3). This table allows you, as a DevOps engineer, to query directory bucket-specific details, including their names, creation dates, and regional information. You can utilize this table to gather insights on directory buckets, such as their distribution across regions and creation timelines.

## Examples

### Basic info
Explore which AWS S3 directory buckets exist in different regions to understand your high-performance storage distribution across your AWS infrastructure.

```sql+postgres
select
  name,
  region,
  account_id,
  creation_date
from
  aws_s3_directory_bucket;
```

```sql+sqlite
select
  name,
  region,
  account_id,
  creation_date
from
  aws_s3_directory_bucket;
```

### List directory buckets created in the last 30 days
Discover the segments that were created recently in your Amazon S3 directory buckets. This could be useful for tracking new high-performance storage deployments.

```sql+postgres
select
  name,
  region,
  account_id,
  creation_date
from
  aws_s3_directory_bucket
where
  creation_date >= now() - interval '30 days';
```

```sql+sqlite
select
  name,
  region,
  account_id,
  creation_date
from
  aws_s3_directory_bucket
where
  creation_date >= datetime('now', '-30 days');
```

### Count directory buckets by region
Analyze the distribution of your S3 directory buckets across different AWS regions to understand your high-performance storage footprint.

```sql+postgres
select
  region,
  count(*) as bucket_count
from
  aws_s3_directory_bucket
group by
  region
order by
  bucket_count desc;
```

```sql+sqlite
select
  region,
  count(*) as bucket_count
from
  aws_s3_directory_bucket
group by
  region
order by
  bucket_count desc;
```

### List directory buckets with specific naming patterns
Identify directory buckets that follow specific naming conventions, which can be useful for organizational and security audits.

```sql+postgres
select
  name,
  region,
  creation_date
from
  aws_s3_directory_bucket
where
  name like 'prod-%' or name like 'dev-%';
```

```sql+sqlite
select
  name,
  region,
  creation_date
from
  aws_s3_directory_bucket
where
  name like 'prod-%' or name like 'dev-%';
```

### List directory buckets with encryption configuration
Find directory buckets that have server-side encryption configured for enhanced security.

```sql+postgres
select
  name,
  region,
  server_side_encryption_configuration
from
  aws_s3_directory_bucket
where
  server_side_encryption_configuration is not null;
```

```sql+sqlite
select
  name,
  region,
  server_side_encryption_configuration
from
  aws_s3_directory_bucket
where
  server_side_encryption_configuration is not null;
```

### List directory buckets with lifecycle rules
Identify directory buckets that have lifecycle management policies configured.

```sql+postgres
select
  name,
  region,
  lifecycle_rules
from
  aws_s3_directory_bucket
where
  lifecycle_rules is not null;
```

```sql+sqlite
select
  name,
  region,
  lifecycle_rules
from
  aws_s3_directory_bucket
where
  lifecycle_rules is not null;
```

### List directory buckets with bucket policies
Identify directory buckets that have bucket policies configured for access control.

```sql+postgres
select
  name,
  region,
  policy
from
  aws_s3_directory_bucket
where
  policy is not null;
```

```sql+sqlite
select
  name,
  region,
  policy
from
  aws_s3_directory_bucket
where
  policy is not null;
```

### List directory buckets with all lifecycle rules
Extract all lifecycle rules from directory buckets, creating a row for each rule.

```sql+postgres
select
  name,
  region,
  jsonb_array_elements(lifecycle_rules) ->> 'ID' as lifecycle_rule_id,
  jsonb_array_elements(lifecycle_rules) ->> 'Status' as lifecycle_status,
  jsonb_array_elements(lifecycle_rules) -> 'Expiration' ->> 'Days' as expiration_days,
  jsonb_array_elements(lifecycle_rules) -> 'Expiration' ->> 'Date' as expiration_date
from
  aws_s3_directory_bucket
where
  lifecycle_rules is not null;
```

```sql+sqlite
select
  name,
  region,
  json_extract(value, '$.ID') as lifecycle_rule_id,
  json_extract(value, '$.Status') as lifecycle_status,
  json_extract(value, '$.Expiration.Days') as expiration_days,
  json_extract(value, '$.Expiration.Date') as expiration_date
from
  aws_s3_directory_bucket,
  json_each(lifecycle_rules)
where
  lifecycle_rules is not null;
```

### List directory buckets with all encryption configurations
Extract all encryption configuration details from directory buckets.

```sql+postgres
select
  name,
  region,
  jsonb_array_elements(server_side_encryption_configuration) ->> 'BucketKeyEnabled' as bucket_key_enabled,
  jsonb_array_elements(server_side_encryption_configuration) -> 'ApplyServerSideEncryptionByDefault' ->> 'SSEAlgorithm' as encryption_algorithm,
  jsonb_array_elements(server_side_encryption_configuration) -> 'ApplyServerSideEncryptionByDefault' ->> 'KMSMasterKeyID' as kms_key_id
from
  aws_s3_directory_bucket
where
  server_side_encryption_configuration is not null;
```

```sql+sqlite
select
  name,
  region,
  json_extract(value, '$.BucketKeyEnabled') as bucket_key_enabled,
  json_extract(value, '$.ApplyServerSideEncryptionByDefault.SSEAlgorithm') as encryption_algorithm,
  json_extract(value, '$.ApplyServerSideEncryptionByDefault.KMSMasterKeyID') as kms_key_id
from
  aws_s3_directory_bucket,
  json_each(server_side_encryption_configuration)
where
  server_side_encryption_configuration is not null;
```

### List directory buckets with all policy statements
Extract all policy statements from bucket policies, creating a row for each statement.

```sql+postgres
select
  name,
  region,
  jsonb_array_elements(policy -> 'Statement') ->> 'Sid' as policy_statement_id,
  jsonb_array_elements(policy -> 'Statement') ->> 'Effect' as policy_effect,
  jsonb_array_elements(policy -> 'Statement') -> 'Action' as policy_actions,
  jsonb_array_elements(policy -> 'Statement') -> 'Resource' as policy_resources
from
  aws_s3_directory_bucket
where
  policy is not null;
```

```sql+sqlite
select
  name,
  region,
  json_extract(value, '$.Sid') as policy_statement_id,
  json_extract(value, '$.Effect') as policy_effect,
  json_extract(value, '$.Action') as policy_actions,
  json_extract(value, '$.Resource') as policy_resources
from
  aws_s3_directory_bucket,
  json_each(policy, '$.Statement')
where
  policy is not null;
```

### Comprehensive directory bucket analysis
Extract and analyze all configuration details from directory buckets in a single query.

```sql+postgres
with bucket_configs as (
  select
    name,
    region,
    jsonb_array_elements(lifecycle_rules) as lifecycle_rule,
    jsonb_array_elements(server_side_encryption_configuration) as encryption_config,
    jsonb_array_elements(policy -> 'Statement') as policy_statement
  from
    aws_s3_directory_bucket
  where
    lifecycle_rules is not null
    and server_side_encryption_configuration is not null
    and policy is not null
)
select
  name,
  region,
  lifecycle_rule ->> 'ID' as lifecycle_rule_id,
  lifecycle_rule ->> 'Status' as lifecycle_status,
  lifecycle_rule -> 'Expiration' ->> 'Days' as expiration_days,
  encryption_config ->> 'BucketKeyEnabled' as bucket_key_enabled,
  encryption_config -> 'ApplyServerSideEncryptionByDefault' ->> 'SSEAlgorithm' as encryption_algorithm,
  policy_statement ->> 'Sid' as policy_statement_id,
  policy_statement ->> 'Effect' as policy_effect
from
  bucket_configs;
```

```sql+sqlite
with bucket_configs as (
  select
    name,
    region,
    json_extract(lifecycle_rule.value, '$.ID') as lifecycle_rule_id,
    json_extract(lifecycle_rule.value, '$.Status') as lifecycle_status,
    json_extract(lifecycle_rule.value, '$.Expiration.Days') as expiration_days,
    json_extract(encryption_config.value, '$.BucketKeyEnabled') as bucket_key_enabled,
    json_extract(encryption_config.value, '$.ApplyServerSideEncryptionByDefault.SSEAlgorithm') as encryption_algorithm,
    json_extract(policy_statement.value, '$.Sid') as policy_statement_id,
    json_extract(policy_statement.value, '$.Effect') as policy_effect
  from
    aws_s3_directory_bucket,
    json_each(lifecycle_rules) as lifecycle_rule,
    json_each(server_side_encryption_configuration) as encryption_config,
    json_each(policy, '$.Statement') as policy_statement
  where
    lifecycle_rules is not null
    and server_side_encryption_configuration is not null
    and policy is not null
)
select * from bucket_configs;
```
