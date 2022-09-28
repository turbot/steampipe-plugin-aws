# Table: aws_codeartifact_domain

A domain contains one or more repositories. All assets and metadata in a domain's repositories are encrypted with one customer master key (CMK) stored in AWS Key Management Service (AWS KMS).

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
