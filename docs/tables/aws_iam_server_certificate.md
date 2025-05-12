---
title: "Steampipe Table: aws_iam_server_certificate - Query AWS IAM Server Certificates using SQL"
description: "Allows users to query AWS IAM Server Certificates"
folder: "IAM"
---

# Table: aws_iam_server_certificate - Query AWS IAM Server Certificates using SQL

The AWS IAM Server Certificate is a resource in AWS Identity and Access Management (IAM) that you upload to deploy an SSL/TLS-based app on AWS. It contains a public key certificate, a private key, and an optional certificate chain, which is an ordered list of certificates that includes the root certificate and intermediate certificates. This enables secure connections from a client, such as a web browser, to an AWS service like a load balancer.

## Table Usage Guide

The `aws_iam_server_certificate` table in Steampipe provides you with information about server certificates within AWS Identity and Access Management (IAM). This table allows you, as a DevOps engineer, to query certificate-specific details, including the certificate body, certificate chain, and associated metadata. You can utilize this table to gather insights on certificates, such as certificates' expiration dates, the path of the certificate, and more. The schema outlines the various attributes of the IAM server certificate for you, including the server certificate name, certificate ID, creation date, and associated tags.

## Examples

### Basic info
Gain insights into your AWS server certificates, including their names, ARNs, and IDs, as well as their upload and expiration dates. This can help manage your certificates, ensuring they're up-to-date and preventing potential security issues.

```sql+postgres
select
  name,
  arn,
  server_certificate_id,
  upload_date,
  expiration
from
  aws_iam_server_certificate;
```

```sql+sqlite
select
  name,
  arn,
  server_certificate_id,
  upload_date,
  expiration
from
  aws_iam_server_certificate;
```

### List expired certificates
Determine the areas in which your AWS IAM server certificates have expired. This is useful to ensure your system's security by replacing or renewing those certificates promptly.

```sql+postgres
select
  name,
  arn,
  expiration
from
  aws_iam_server_certificate
where
  expiration < now()::timestamp;
```

```sql+sqlite
select
  name,
  arn,
  expiration
from
  aws_iam_server_certificate
where
  expiration < datetime('now');
```