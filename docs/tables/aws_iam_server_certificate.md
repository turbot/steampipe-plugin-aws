---
title: "Table: aws_iam_server_certificate - Query AWS IAM Server Certificates using SQL"
description: "Allows users to query AWS IAM Server Certificates"
---

# Table: aws_iam_server_certificate - Query AWS IAM Server Certificates using SQL

The `aws_iam_server_certificate` table in Steampipe provides information about server certificates within AWS Identity and Access Management (IAM). This table allows DevOps engineers to query certificate-specific details, including the certificate body, certificate chain, and associated metadata. Users can utilize this table to gather insights on certificates, such as certificates' expiration dates, the path of the certificate, and more. The schema outlines the various attributes of the IAM server certificate, including the server certificate name, certificate ID, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_iam_server_certificate` table, you can use the `.inspect aws_iam_server_certificate` command in Steampipe.

### Key columns:
- `server_certificate_name`: The name that identifies the server certificate. This can be used to join this table with other tables that require server certificate name.
- `server_certificate_id`: The ID for the server certificate. This can be used to join this table with other tables that require server certificate ID.
- `arn`: The Amazon Resource Name (ARN) specifying the server certificate. This can be used to join this table with other tables that require the ARN of a server certificate.

## Examples

### Basic info

```sql
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

```sql
select
  name,
  arn,
  expiration
from
  aws_iam_server_certificate
where
  expiration < now()::timestamp;
```
