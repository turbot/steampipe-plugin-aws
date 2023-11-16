---
title: "Table: aws_codeartifact_domain - Query AWS CodeArtifact Domains using SQL"
description: "Allows users to query AWS CodeArtifact Domains for details such as domain ownership, encryption key, and policy information."
---

# Table: aws_codeartifact_domain - Query AWS CodeArtifact Domains using SQL

The `aws_codeartifact_domain` table in Steampipe provides information about domains within AWS CodeArtifact. This table allows DevOps engineers to query domain-specific details, including domain ownership, encryption key, and associated policy information. Users can utilize this table to gather insights on domains, such as who owns a domain, what encryption key is used, and what policies are applied. The schema outlines the various attributes of the AWS CodeArtifact domain, including the domain ARN, domain owner, encryption key, and associated policies.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_codeartifact_domain` table, you can use the `.inspect aws_codeartifact_domain` command in Steampipe.

### Key columns:

- `name`: The name of the domain. It is unique to an AWS account and useful for joining with other tables that reference AWS CodeArtifact domains by name.
- `arn`: The Amazon Resource Number (ARN) of the domain. This is a unique identifier for the domain and can be used for joining with other tables that reference AWS CodeArtifact domains by ARN.
- `domain_owner`: The AWS account ID that owns the domain. This can be used for joining with other tables that reference AWS account IDs.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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

### Get S3 bucket details associated with each domain

```sql
select
  d.arn as domain_arn,
  b.arn as bucket_arn,
  d.encryption_key domain_encryption_key,
  bucket_policy_is_public
from
  aws_codeartifact_domain d
  join aws_s3_bucket b on d.s3_bucket_arn = b.arn;
```

### Get KMS key details associated with each the domain

```sql
select
  d.arn as domain_arn,
  d.encryption_key domain_encryption_key,
  key_manager,
  key_state
from
  aws_codeartifact_domain d
  join aws_kms_key k on d.encryption_key = k.arn;
```

### List domains using customer managed encryption

```sql
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
