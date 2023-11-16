---
title: "Table: aws_acm_certificate - Query AWS Certificate Manager certificates using SQL"
description: "Allows users to query AWS Certificate Manager certificates. This table provides information about each certificate, including the domain name, status, issuer, and more. It can be used to monitor certificate details, validity, and expiration data."
---

# Table: aws_acm_certificate - Query AWS Certificate Manager certificates using SQL

The `aws_acm_certificate` table in Steampipe provides information about certificates within AWS Certificate Manager (ACM). This table allows DevOps engineers to query certificate-specific details, including domain name, status, issuer, and expiration data. Users can utilize this table to gather insights on certificates, such as certificate status, verification of issuer, and more. The schema outlines the various attributes of the ACM certificate, including the certificate ARN, creation date, domain name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_acm_certificate` table, you can use the `.inspect aws_acm_certificate` command in Steampipe.

**Key columns**:

- `certificate_arn`: The Amazon Resource Name (ARN) of the certificate. This is a unique identifier and can be used to join this table with other tables.
- `domain_name`: The fully qualified domain name (FQDN), such as www.example.com, for which ACM issued the certificate. This can be used to filter certificates by domain.
- `status`: The status of the ACM certificate. This can be used to monitor and manage certificate statuses.

## Examples

### Basic info

```sql
select
  certificate_arn,
  domain_name,
  failure_reason,
  in_use_by,
  status,
  key_algorithm
from
  aws_acm_certificate;
```


### List of expired certificates

```sql
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  status = 'expired';
```


### List certificates for which transparency logging is disabled

```sql
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  certificate_transparency_logging_preference <> 'ENABLED';
```


### List certificates without application tag key

```sql
select
  certificate_arn,
  turbot_tags
from
  aws_acm_certificate
where
  not turbot_tags :: JSONB ? 'application';
```
