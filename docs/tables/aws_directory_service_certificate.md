---
title: "Steampipe Table: aws_directory_service_certificate - Query AWS Directory Service Certificates using SQL"
description: "Allows users to query AWS Directory Service Certificates to gather information about the certificates associated with AWS Managed Microsoft AD and Simple AD directories."
folder: "Directory Service"
---

# Table: aws_directory_service_certificate - Query AWS Directory Service Certificates using SQL

The AWS Directory Service Certificate is a component of the AWS Directory Service, which simplifies the setup and management of Windows and Linux directories in the cloud. These certificates are used to establish secure LDAP communications between your applications and your AWS managed directories. They provide an extra layer of security by encrypting your data and establishing a secure connection.

## Table Usage Guide

The `aws_directory_service_certificate` table in Steampipe provides you with information about the certificates associated with AWS Managed Microsoft AD and Simple AD directories. This table allows you as an IT administrator or security professional to query certificate-specific details, including certificate state, expiry date, and associated metadata. You can utilize this table to gather insights on certificates, such as active certificates, expired certificates, and certificates nearing expiry. The schema outlines the various attributes of the Directory Service Certificate for you, including the certificate ID, common name, expiry date, registered date, and the state of the certificate.

## Examples

### Basic Info
Determine the status and validity of your AWS Directory Service's security certificates. This is particularly useful for maintaining system security by ensuring certificates are up-to-date and appropriately configured.

```sql+postgres
select
  directory_id,
  certificate_id,
  common_name,
  type,
  state,
  expiry_date_time
from
  aws_directory_service_certificate;
```

```sql+sqlite
select
  directory_id,
  certificate_id,
  common_name,
  type,
  state,
  expiry_date_time
from
  aws_directory_service_certificate;
```

### List 'MicrosoftAD' type directories
Determine the areas in which 'MicrosoftAD' type directories are being used. This query can be useful to gain insights into the distribution and application of these directories within your AWS environment.

```sql+postgres
select
  c.certificate_id,
  c.common_name,
  c.directory_id,
  c.type as certificate_type,
  d.name as directory_name,
  d.type as directory_type
from
  aws_directory_service_certificate c,
  aws_directory_service_directory d
where
  d.type = 'MicrosoftAD';
```

```sql+sqlite
select
  c.certificate_id,
  c.common_name,
  c.directory_id,
  c.type as certificate_type,
  d.name as directory_name,
  d.type as directory_type
from
  aws_directory_service_certificate c,
  aws_directory_service_directory d
where
  d.type = 'MicrosoftAD';
```

### List deregistered certificates
Identify instances where certificates have been deregistered within the AWS directory service. This can be useful in understanding the history of your security configuration and tracking changes over time.

```sql+postgres
select
  common_name,
  directory_id,
  type,
  state
from
  aws_directory_service_certificate
where
  state = 'Deregistered';
```

```sql+sqlite
select
  common_name,
  directory_id,
  type,
  state
from
  aws_directory_service_certificate
where
  state = 'Deregistered';
```

### List certificates that will expire in the coming 7 days
Identify the certificates that are due to expire in the next week. This allows you to proactively manage and renew them before they lapse, ensuring continuous and secure operations.

```sql+postgres
select
  directory_id,
  certificate_id,
  common_name,
  type,
  state,
  expiry_date_time
from
  aws_directory_service_certificate
where
  expiry_date_time >= now() + interval '7' day;
```

```sql+sqlite
select
  directory_id,
  certificate_id,
  common_name,
  type,
  state,
  expiry_date_time
from
  aws_directory_service_certificate
where
  expiry_date_time >= datetime('now', '+7 day');
```

### Get client certificate auth settings of each certificate
Analyze the authentication settings of each certificate to understand the Online Certificate Status Protocol (OCSP) URL's configuration. This can help in ensuring the certificates are correctly configured for client authentication, thereby enhancing security.

```sql+postgres
select
  directory_id,
  certificate_id,
  common_name,
  client_cert_auth_settings -> 'OCSPUrl' as ocsp_url
from
  aws_directory_service_certificate;
```

```sql+sqlite
select
  directory_id,
  certificate_id,
  common_name,
  json_extract(client_cert_auth_settings, '$.OCSPUrl') as ocsp_url
from
  aws_directory_service_certificate;
```

### Retrieve the number of certificates registered in each directory
Determine the distribution of certificates across various directories to understand their allocation and manage resources more effectively.

```sql+postgres
select
  directory_id,
  count(*) as certificate_count
from
  aws_directory_service_certificate
group by
  directory_id;
```

```sql+sqlite
select
  directory_id,
  count(*) as certificate_count
from
  aws_directory_service_certificate
group by
  directory_id;
```

### List all certificates that were registered more than a year ago and have not been deregistered
Pinpoint the specific instances where certificates have been registered for over a year and have not yet been deregistered. This can be useful for maintaining security standards and ensuring outdated certificates are properly managed.

```sql+postgres
select
  common_name,
  directory_id,
  type,
  state
from
  aws_directory_service_certificate
where
  registered_date_time <= now() - interval '1 year'
  and state not like 'Deregister%';
```

```sql+sqlite
select
  common_name,
  directory_id,
  type,
  state
from
  aws_directory_service_certificate
where
  registered_date_time <= datetime('now', '-1 year')
  and state not like 'Deregister%';
```

### Find the certificate with the latest registration date in each AWS partition
Discover the segments that have the most recent certificate registrations within each AWS partition. This can be useful for maintaining up-to-date security practices and ensuring compliance within your AWS infrastructure.

```sql+postgres
select
  distinct partition,
  registered_date_time
from
  aws_directory_service_certificate
order by
  partition,
  registered_date_time desc;
```

```sql+sqlite
select
  distinct partition,
  registered_date_time
from
  aws_directory_service_certificate
order by
  partition,
  registered_date_time desc;
```