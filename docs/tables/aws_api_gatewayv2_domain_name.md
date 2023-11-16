---
title: "Table: aws_api_gatewayv2_domain_name - Query AWS API Gateway Domain Names using SQL"
description: "Allows users to query AWS API Gateway Domain Names and provides information about each domain name within the AWS API Gateway Service. This table can be used to query domain name details, including associated API mappings, security policy, and associated tags."
---

# Table: aws_api_gatewayv2_domain_name - Query AWS API Gateway Domain Names using SQL

The `aws_api_gatewayv2_domain_name` table in Steampipe provides information about each domain name within the AWS API Gateway Service. This table allows users to query domain name details, including associated API mappings, security policy, and associated tags. The schema outlines the various attributes of the domain name, including the domain name ARN, domain name, endpoint type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gatewayv2_domain_name` table, you can use the `.inspect aws_api_gatewayv2_domain_name` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of an API Gateway Domain Name. This can be used to join with other tables that return ARN information.
- `domain_name`: The name of the domain. This key column can be used to join with other tables that return domain name information.
- `security_policy`: The security policy of the domain. This column can be used to join with other tables that return security policy information.

## Examples

### Basic info

```sql
select
  domain_name,
  mutual_tls_authentication,
  tags,
  title,
  akas
from
  aws_api_gatewayv2_domain_name;
```

### List of all edge endpoint type domain name

```sql
select
  domain_name,
  config ->> 'EndpointType' as endpoint_type
from
  aws_api_gatewayv2_domain_name
  cross join jsonb_array_elements(domain_name_configurations) as config
where
  config ->> 'EndpointType' = 'EDGE';
```

### API gatewayv2 domain name configuration info

```sql
select
  domain_name,
  config ->> 'EndpointType' as endpoint_type,
  config ->> 'CertificateName' as certificate_name,
  config ->> 'CertificateArn' as certificate_arn,
  config ->> 'CertificateUploadDate' as certificate_upload_date,
  config ->> 'DomainNameStatus' as domain_name_status,
  config ->> 'DomainNameStatusMessage' as domain_name_status_message,
  config ->> 'ApiGatewayDomainName' as api_gateway_domain_name,
  config ->> 'HostedZoneId' as hosted_zone_id,
  config ->> 'OwnershipVerificationCertificateArn' as ownership_verification_certificate_arn,
  config -> 'SecurityPolicy' as security_policy
from
  aws_api_gatewayv2_domain_name
  cross join jsonb_array_elements(domain_name_configurations) as config;
```

### Get mutual TLS authentication configuration of each domain name

```sql
select
  domain_name,
  mutual_tls_authentication ->> 'TruststoreUri' as truststore_uri,
  mutual_tls_authentication ->> 'TruststoreVersion' as truststore_version,
  mutual_tls_authentication ->> 'TruststoreWarnings' as truststore_warnings
from
  aws_api_gatewayv2_domain_name;
```

### Get certificate details of each domain names

```sql
select
  d.domain_name,
  config ->> 'CertificateArn' as certificate_arn,
  c.certificate,
  c.certificate_transparency_logging_preference,
  c.created_at,
  c.imported_at,
  c.issuer,
  c.issued_at,
  c.key_algorithm
from
  aws_api_gatewayv2_domain_name AS d
  cross join jsonb_array_elements(d.domain_name_configurations) AS config
  left join aws_acm_certificate AS c ON c.certificate_arn = config ->> 'CertificateArn';
```