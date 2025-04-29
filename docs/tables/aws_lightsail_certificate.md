---
title: "Steampipe Table: aws_lightsail_certificate - Query AWS Lightsail Certificates using SQL"
description: "Allows users to query AWS Lightsail Certificates for detailed information about SSL/TLS certificates, including their status, validity periods, and associated metadata."
folder: "Lightsail"
---

# Table: aws_lightsail_certificate - Query AWS Lightsail Certificates using SQL

The AWS Lightsail Certificate is a feature of Amazon Lightsail that enables you to manage SSL/TLS certificates for your Lightsail resources. These certificates are used to secure connections to your Lightsail resources, such as load balancers and content delivery network (CDN) distributions. AWS Lightsail provides both automatic and manual certificate management options.

## Table Usage Guide

The `aws_lightsail_certificate` table in Steampipe provides you with information about AWS Lightsail certificates. This table allows you as a DevOps engineer to query certificate-specific details, including the certificate name, domain name, status, validity periods, and associated metadata. You can utilize this table to gather insights on certificates, such as validation status, expiration dates, and any failures that occurred during certificate validation. The schema outlines the various attributes of the AWS Lightsail certificate for you, including the certificate ARN, domain validation records, and associated tags.

## Examples

### Basic info
Explore the basic information about your Lightsail certificates, including their names, domain names, and status. This can help you understand the current state of your certificates and identify any that might need attention.

```sql+postgres
select
  name,
  domain_name,
  status,
  not_after,
  not_before,
  issuer_ca
from
  aws_lightsail_certificate;
```

```sql+sqlite
select
  name,
  domain_name,
  status,
  not_after,
  not_before,
  issuer_ca
from
  aws_lightsail_certificate;
```

### List certificates that are expiring in the next 30 days
Identify certificates that are approaching expiration to ensure timely renewal and prevent service disruptions.

```sql+postgres
select
  name,
  domain_name,
  not_after,
  status
from
  aws_lightsail_certificate
where
  not_after < now() + interval '30 days'
  and status = 'ISSUED';
```

```sql+sqlite
select
  name,
  domain_name,
  not_after,
  status
from
  aws_lightsail_certificate
where
  not_after < datetime('now', '+30 days')
  and status = 'ISSUED';
```

### List certificates with validation failures
Find certificates that failed validation, along with their failure details. This can help you troubleshoot issues and ensure your certificates are properly configured.

```sql+postgres
select
  name,
  domain_name,
  status,
  request_failure_reason
from
  aws_lightsail_certificate
where
  status = 'FAILED';
```

```sql+sqlite
select
  name,
  domain_name,
  status,
  request_failure_reason
from
  aws_lightsail_certificate
where
  status = 'FAILED';
```

### List certificates by status
Analyze the distribution of certificates by their current status to understand your certificate management landscape.

```sql+postgres
select
  status,
  count(*) as certificate_count
from
  aws_lightsail_certificate
group by
  status
order by
  certificate_count desc;
```

```sql+sqlite
select
  status,
  count(*) as certificate_count
from
  aws_lightsail_certificate
group by
  status
order by
  certificate_count desc;
```

### List certificates with their associated resources
Find certificates that are currently in use by Lightsail resources to understand certificate dependencies.

```sql+postgres
select
  name,
  domain_name,
  in_use_resource_count,
  status
from
  aws_lightsail_certificate
where
  in_use_resource_count > 0;
```

```sql+sqlite
select
  name,
  domain_name,
  in_use_resource_count,
  status
from
  aws_lightsail_certificate
where
  in_use_resource_count > 0;
```

### List certificates with specific tags
Find certificates that have specific tags associated with them to help organize and manage your certificates based on custom criteria.

```sql+postgres
select
  name,
  domain_name,
  tags
from
  aws_lightsail_certificate
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  name,
  domain_name,
  tags
from
  aws_lightsail_certificate
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List certificates with domain validation records
View the domain validation records for certificates to help with DNS configuration and validation.

```sql+postgres
select
  name,
  domain_name,
  domain_validation_records
from
  aws_lightsail_certificate
where
  domain_validation_records is not null;
```

```sql+sqlite
select
  name,
  domain_name,
  domain_validation_records
from
  aws_lightsail_certificate
where
  domain_validation_records is not null;
```

### List certificates with subject alternative names
Find certificates that have additional domain names associated with them.

```sql+postgres
select
  name,
  domain_name,
  subject_alternative_names
from
  aws_lightsail_certificate
where
  subject_alternative_names is not null;
```

```sql+sqlite
select
  name,
  domain_name,
  subject_alternative_names
from
  aws_lightsail_certificate
where
  subject_alternative_names is not null;
```

### List certificates by key algorithm
Analyze the distribution of certificates by their key algorithm to ensure compliance with security standards.

```sql+postgres
select
  key_algorithm,
  count(*) as certificate_count
from
  aws_lightsail_certificate
group by
  key_algorithm
order by
  certificate_count desc;
```

```sql+sqlite
select
  key_algorithm,
  count(*) as certificate_count
from
  aws_lightsail_certificate
group by
  key_algorithm
order by
  certificate_count desc;
```
