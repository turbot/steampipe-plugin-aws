---
title: "Table: aws_directory_service_certificate - Query AWS Directory Service Certificates using SQL"
description: "Allows users to query AWS Directory Service Certificates to gather information about the certificates associated with AWS Managed Microsoft AD and Simple AD directories."
---

# Table: aws_directory_service_certificate - Query AWS Directory Service Certificates using SQL

The `aws_directory_service_certificate` table in Steampipe provides information about the certificates associated with AWS Managed Microsoft AD and Simple AD directories. This table allows IT administrators and security professionals to query certificate-specific details, including certificate state, expiry date, and associated metadata. Users can utilize this table to gather insights on certificates, such as active certificates, expired certificates, and certificates nearing expiry. The schema outlines the various attributes of the Directory Service Certificate, including the certificate ID, common name, expiry date, registered date, and the state of the certificate.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_directory_service_certificate` table, you can use the `.inspect aws_directory_service_certificate` command in Steampipe.

**Key columns**:

- `certificate_id`: The identifier of the AWS Managed Microsoft AD directory certificate. It can be used to join with other tables to gather more specific information about a particular certificate.
- `directory_id`: The identifier of the AWS Directory Service. This can be used to join with other AWS Directory Service tables to gather more information about the directory associated with the certificate.
- `common_name`: The fully qualified domain name (FQDN) of the certificate. This can be used to join with other tables that may contain information about the domain associated with the certificate.

## Examples

### Basic Info

```sql
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

```sql
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

```sql
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

```sql
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

### Get client certificate auth settings of each certificate

```sql
select
  directory_id,
  certificate_id,
  common_name,
  client_cert_auth_settings -> 'OCSPUrl' as ocsp_url
from
  aws_directory_service_certificate;
```

### Retrieve the number of certificates registered in each directory

```sql
select
  directory_id,
  count(*) as certificate_count
from
  aws_directory_service_certificate
group by
  directory_id;
```

### List all certificates that were registered more than a year ago and have not been deregistered

```sql
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

### Find the certificate with the latest registration date in each AWS partition

```sql
select
  distinct partition,
  registered_date_time
from
  aws_directory_service_certificate
order by
  partition,
  registered_date_time desc;
```
