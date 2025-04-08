---
title: "Steampipe Table: aws_dms_certificate - Query AWS DMS Certificates using SQL"
description: "Allows users to query AWS DMS (Database Migration Service) Certificates. This table provides information about SSL/TLS certificates used in AWS DMS for encrypting data during database migration tasks. Certificates play a crucial role in ensuring the security and integrity of data transferred between source and target databases."
folder: "DMS"
---

# Table: aws_dms_certificate - Query AWS DMS Certificates using SQL

AWS DMS (Database Migration Service) Certificate refers to an SSL/TLS certificate used in AWS DMS for encrypting data during the process of migrating databases. This certificate plays a crucial role in ensuring the security and integrity of the data as it is transferred between the source and target databases in a migration task.

## Table Usage Guide

The `aws_dms_certificate` table in Steampipe enables users to query information about AWS DMS Certificates. These certificates are used to secure the data during database migration tasks. Users can retrieve details such as the certificate identifier, ARN, certificate creation date, signing algorithm, valid-to date, and region. Additionally, the table allows users to filter certificates based on various criteria, such as expiration date, signing algorithm, ownership, and more.

## Examples

### Basic info
Retrieve basic information about AWS DMS Certificates, including their identifiers, ARNs, certificate creation dates, signing algorithms, valid-to dates, and regions. This query provides an overview of the certificates in your AWS environment.

```sql+postgres
select
  certificate_identifier,
  arn,
  certificate_creation_date,
  signing_algorithm,
  valid_to_date,
  region
from
  aws_dms_certificate;
```

```sql+sqlite
select
  certificate_identifier,
  arn,
  certificate_creation_date,
  signing_algorithm,
  valid_to_date,
  region
from
  aws_dms_certificate;
```

### List certificates expiring in next 10 days
Identify AWS DMS Certificates that are set to expire within the next 10 days. This query helps you proactively manage certificate renewals.

```sql+postgres
select
  certificate_identifier,
  arn,
  key_length,
  signing_algorithm,
  valid_to_date
from
  aws_dms_certificate
where
  valid_to_date <= current_date + interval '10' day;
```

```sql+sqlite
select
  certificate_identifier,
  arn,
  key_length,
  signing_algorithm,
  valid_to_date
from
  aws_dms_certificate
where
  valid_to_date <= date('now', '+10 day');
```

### List certificates with SHA256 signing algorithm
Retrieve AWS DMS Certificates that use the SHA256 with RSA signing algorithm. This query helps you identify certificates with specific security configurations.

```sql+postgres
select
  certificate_identifier,
  arn,
  signing_algorithm,
  key_length,
  certificate_owner
from
  aws_dms_certificate
where
  signing_algorithm = 'SHA256withRSA';
```

```sql+sqlite
select
  certificate_identifier,
  arn,
  signing_algorithm,
  key_length,
  certificate_owner
from
  aws_dms_certificate
where
  signing_algorithm = 'SHA256withRSA';
```

### List certificates not owned by the current account
Identify AWS DMS Certificates that are not owned by the current AWS account. This query helps you keep track of certificates associated with other accounts.

```sql+postgres
select
  certificate_identifier,
  arn,
  certificate_owner,
  account_id
from
  aws_dms_certificate
where
  certificate_owner <> account_id;
```

```sql+sqlite
select
  certificate_identifier,
  arn,
  certificate_owner,
  account_id
from
  aws_dms_certificate
where
  certificate_owner <> account_id;
```

### Get the number of days left until certificates expire
Retrieve AWS DMS Certificates along with the number of days left until they expire. This query helps you monitor certificate expiration dates.

```sql+postgres
select
  certificate_identifier,
  arn,
  certificate_owner,
  (valid_to_date - current_date) as days_left,
  region
from
  aws_dms_certificate;
```

```sql+sqlite
select
  certificate_identifier,
  arn,
  certificate_owner,
  (julianday(valid_to_date) - julianday('now')) as days_left,
  region
from
  aws_dms_certificate;
```