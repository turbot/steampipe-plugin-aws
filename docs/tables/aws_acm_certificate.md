---
title: "Steampipe Table: aws_acm_certificate - Query AWS Certificate Manager certificates using SQL"
description: "Allows users to query AWS Certificate Manager certificates. This table provides information about each certificate, including the domain name, status, issuer, and more. It can be used to monitor certificate details, validity, and expiration data."
folder: "ACM"
---

# Table: aws_acm_certificate - Query AWS Certificate Manager certificates using SQL

The AWS Certificate Manager (ACM) is a service that lets you easily provision, manage, and deploy public and private Secure Sockets Layer/Transport Layer Security (SSL/TLS) certificates for use with AWS services and your internal connected resources. SSL/TLS certificates are used to secure network communications and establish the identity of websites over the Internet as well as resources on private networks. AWS Certificate Manager removes the time-consuming manual process of purchasing, uploading, and renewing SSL/TLS certificates.

## Table Usage Guide

The `aws_acm_certificate` table in Steampipe provides you with information about certificates within AWS Certificate Manager (ACM). This table allows you, as a DevOps engineer, to query certificate-specific details, including domain name, status, issuer, and expiration data. You can utilize this table to gather insights on certificates, such as certificate status, verification of issuer, and more. The schema outlines the various attributes of the ACM certificate for you, including the certificate ARN, creation date, domain name, and associated tags.



## Examples

### Basic info
Analyze the settings to understand the status and usage of your AWS Certificate Manager (ACM) certificates. This can help identify any issues with certificates, such as failure reasons, and see which domains they're associated with, aiding in efficient resource management and troubleshooting.

```sql+postgres
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

```sql+sqlite
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
Identify instances where your AWS certificates have expired. This allows you to maintain security by promptly replacing or renewing these certificates.

```sql+postgres
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  status = 'EXPIRED';
```

```sql+sqlite
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  status = 'EXPIRED';
```


### List certificates for which transparency logging is disabled
Discover the segments with disabled transparency logging in certificate settings to enhance security and compliance efforts. This allows for proactive mitigation of potential risks associated with non-transparent logging.

```sql+postgres
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  certificate_transparency_logging_preference <> 'ENABLED';
```

```sql+sqlite
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  certificate_transparency_logging_preference != 'ENABLED';
```


### List certificates without application tag key
Identify the certificates that are missing an application tag key. This can help in pinpointing areas where tagging conventions may not have been followed, aiding in better resource management.

```sql+postgres
select
  certificate_arn,
  tags
from
  aws_acm_certificate
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  certificate_arn,
  tags
from
  aws_acm_certificate
where
  json_extract(tags, '$.application') is null;
```