---
title: "Steampipe Table: aws_acmpca_certificate_authority - Query AWS ACM PCA Certificate Authorities using SQL"
description: "Allows users to query AWS ACM PCA Certificate Authorities. It can be used to monitor certificate authorities details, validity, usage mode and expiration data."
folder: "ACM"
---

# Table: aws_acmpca_certificate_authority - Query AWS ACM PCA Certificate Authorities using SQL

The `aws_acmpca_certificate_authority` table provides detailed information about AWS Certificate Manager Private Certificate Authority (ACM PCA) certificate authorities. These entities enable you to securely issue and manage your private certificates. This table allows for querying configurations, statuses, key storage standards, and more for each certificate authority within your AWS account.

## Table Usage Guide

This table can be utilized to monitor the configuration and operational health of your private certificate authorities managed through AWS ACM PCA. It enables security analysts, compliance auditors, and cloud administrators to assess the certificate authorities' compliance with policies, investigate issuance metadata, and understand the security standards being applied.

## Examples

### Basic information
Retrieve basic details about your ACM PCA Certificate Authorities.

```sql+postgres
select
  arn,
  status,
  created_at,
  not_before,
  not_after,
  key_storage_security_standard,
  failure_reason
from
  aws_acmpca_certificate_authority;
```

```sql+sqlite
select
  arn,
  status,
  datetime(created_at) AS created_at,
  datetime(not_before) AS not_before,
  datetime(not_after) AS not_after,
  key_storage_security_standard,
  failure_reason
from
  aws_acmpca_certificate_authority;
```

### Certificate authorities with specific key storage security standards
List certificate authorities that comply with a specific key storage security standard.

```sql+postgres
select
  arn,
  status,
  key_storage_security_standard
from
  aws_acmpca_certificate_authority
where
  key_storage_security_standard = 'FIPS_140_2_LEVEL_3_OR_HIGHER';
```

```sql+sqlite
select
  arn,
  status,
  key_storage_security_standard
from
  aws_acmpca_certificate_authority
where
  key_storage_security_standard = 'FIPS_140_2_LEVEL_3_OR_HIGHER';
```

### Certificate authorities by status
Find certificate authorities by their operational status, e.g., `ACTIVE`, `DISABLED`.

```sql+postgres
select
  arn,
  status,
  created_at,
  last_state_change_at
from
  aws_acmpca_certificate_authority
where
  status = 'ACTIVE';
```

```sql+sqlite
select
  arn,
  status,
  datetime(created_at) AS created_at,
  datetime(last_state_change_at) AS last_state_change_at
from
  aws_acmpca_certificate_authority
where
  status = 'ACTIVE';
```

### Tagged Certificate Authorities
Identify certificate authorities tagged with specific key-value pairs for organizational purposes.

```sql+postgres
select
  arn,
  tags
from
  aws_acmpca_certificate_authority
where
  (tags ->> 'Project') = 'MyProject';
```

```sql+sqlite
select
  arn,
  json_extract(tags, '$.Project') AS project_tag
from
  aws_acmpca_certificate_authority
where
  project_tag = 'MyProject';
```
