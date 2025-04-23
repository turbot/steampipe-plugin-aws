---
title: "Steampipe Table: aws_codeartifact_domain - Query AWS CodeArtifact Domains using SQL"
description: "Allows users to query AWS CodeArtifact Domains for details such as domain ownership, encryption key, and policy information."
folder: "CodeArtifact"
---

# Table: aws_codeartifact_domain - Query AWS CodeArtifact Domains using SQL

The AWS CodeArtifact Domain is a fundamental resource within the AWS CodeArtifact service, which is a fully managed artifact repository service. It enables you to easily store, publish, and share software packages in a scalable and secure manner. Each domain allows for the management and organization of your package assets across multiple repositories.

## Table Usage Guide

The `aws_codeartifact_domain` table in Steampipe provides you with information about domains within AWS CodeArtifact. This table allows you, as a DevOps engineer, to query domain-specific details, including domain ownership, encryption key, and associated policy information. You can utilize this table to gather insights on domains, such as who owns a domain, what encryption key is used, and what policies are applied. The schema outlines the various attributes of the AWS CodeArtifact domain for you, including the domain ARN, domain owner, encryption key, and associated policies.

## Examples

### Basic info
Discover the segments that provide insights into the creation, ownership, and status of AWS CodeArtifact domains, in order to better understand and manage your resources. This could be beneficial for maintaining security protocols and efficient resource allocation.

```sql+postgres
select
  arn,
  created_time,
  encryption_key,
  status,
  owner,
  tags
from
  aws_codeartifact_domain;
```

```sql+sqlite
select
  arn,
  created_time,
  encryption_key,
  status,
  owner,
  tags
from
  aws_codeartifact_domain;
```

### List unencrypted domains
Identify instances where AWS CodeArtifact domains are unencrypted, providing a useful method to highlight potential security vulnerabilities within your AWS infrastructure. This can aid in enhancing data protection measures by pinpointing areas that require encryption implementation.

```sql+postgres
select
  arn,
  created_time,
  status,
  s3_bucket_arn,
  tags
from
  aws_codeartifact_domain
where
  encryption_key is null;
```

```sql+sqlite
select
  arn,
  created_time,
  status,
  s3_bucket_arn,
  tags
from
  aws_codeartifact_domain
where
  encryption_key is null;
```

### List inactive domains
Determine the areas in which domains are not actively used within the AWS CodeArtifact service. This can be useful in identifying unused resources, potentially helping to reduce costs and optimize resource management.

```sql+postgres
select
  arn,
  created_time,
  status,
  s3_bucket_arn,
  tags
from
  aws_codeartifact_domain
where
  status != 'Active';
```

```sql+sqlite
select
  arn,
  created_time,
  status,
  s3_bucket_arn,
  tags
from
  aws_codeartifact_domain
where
  status != 'Active';
```

### List domain policy statements that grant external access
Explore which domain policy statements in your AWS CodeArtifact domain allow external access. This is useful to identify potential security vulnerabilities and ensure that only authorized entities have access to your domain.

```sql+postgres
select
  arn,
  p as principal,
  a as action,
  s ->> 'Effect' as effect
from
  aws_codeartifact_domain,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  string_to_array(p, ':') as pa,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and (
    pa [5] != account_id
    or p = '*'
  );
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```

### Get S3 bucket details associated with each domain
Determine the areas in which S3 bucket details are linked with each domain to assess the elements within the domain's encryption key and public bucket policy. This can be useful to gain insights into the security configuration of your AWS CodeArtifact domains and associated S3 buckets.

```sql+postgres
select
  d.arn as domain_arn,
  b.arn as bucket_arn,
  d.encryption_key domain_encryption_key,
  bucket_policy_is_public
from
  aws_codeartifact_domain d
  join aws_s3_bucket b on d.s3_bucket_arn = b.arn;
```

```sql+sqlite
select
  d.arn as domain_arn,
  b.arn as bucket_arn,
  d.encryption_key as domain_encryption_key,
  bucket_policy_is_public
from
  aws_codeartifact_domain d
  join aws_s3_bucket b on d.s3_bucket_arn = b.arn;
```

### Get KMS key details associated with each the domain
Explore which domains are associated with specific KMS keys to gain insights into their encryption status and management. This can help in assessing the security configuration of your AWS CodeArtifact domains.

```sql+postgres
select
  d.arn as domain_arn,
  d.encryption_key domain_encryption_key,
  key_manager,
  key_state
from
  aws_codeartifact_domain d
  join aws_kms_key k on d.encryption_key = k.arn;
```

```sql+sqlite
select
  d.arn as domain_arn,
  d.encryption_key as domain_encryption_key,
  key_manager,
  key_state
from
  aws_codeartifact_domain d
  join aws_kms_key k on d.encryption_key = k.arn;
```

### List domains using customer managed encryption
Discover the segments that use customer-managed encryption in your AWS CodeArtifact domains. This can be beneficial for assessing your security protocols and identifying areas where you're maintaining direct control over your encryption keys.

```sql+postgres
select
  d.arn as domain_arn,
  d.encryption_key domain_encryption_key,
  key_manager,
  key_state
from
  aws_codeartifact_domain d
  join aws_kms_key k on d.encryption_key = k.arn
where 
  key_manager = 'CUSTOMER';
```

```sql+sqlite
select
  d.arn as domain_arn,
  d.encryption_key as domain_encryption_key,
  key_manager,
  key_state
from
  aws_codeartifact_domain as d
  join aws_kms_key as k on d.encryption_key = k.arn
where 
  key_manager = 'CUSTOMER';
```


