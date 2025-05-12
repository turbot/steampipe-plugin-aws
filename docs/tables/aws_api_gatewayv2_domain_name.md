---
title: "Steampipe Table: aws_api_gatewayv2_domain_name - Query AWS API Gateway Domain Names using SQL"
description: "Allows users to query AWS API Gateway Domain Names and provides information about each domain name within the AWS API Gateway Service. This table can be used to query domain name details, including associated API mappings, security policy, and associated tags."
folder: "API Gateway"
---

# Table: aws_api_gatewayv2_domain_name - Query AWS API Gateway Domain Names using SQL

The AWS API Gateway Domain Name is a component of Amazon API Gateway that you associate with a DNS hostname. It's utilized to provide a custom domain for an API that you deploy through the service. The custom domain name can be used to route requests to the API, providing a more user-friendly URL for your API endpoints.

## Table Usage Guide

The `aws_api_gatewayv2_domain_name` table in Steampipe provides you with information about each domain name within the AWS API Gateway Service. This table allows you to query domain name details, including associated API mappings, security policy, and associated tags. The schema outlines the various attributes of the domain name for you, including the domain name ARN, domain name, endpoint type, and associated tags.

## Examples

### Basic info
Explore the security and metadata aspects of your AWS API Gateway domain names. This query is useful to gain insights into the mutual TLS authentication status, associated tags, title, and alternative names of your domain names, crucial for maintaining secure and organized API management.Analyze the settings to understand the security measures and metadata associated with different domains in your AWS API Gateway. This query can help you assess the use of mutual TLS authentication and keep track of domains through their tags, titles, and alternate names.


```sql+postgres
select
  domain_name,
  mutual_tls_authentication,
  tags,
  title,
  akas
from
  aws_api_gatewayv2_domain_name;
```

```sql+sqlite
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
Identify instances where the endpoint type of a domain name in AWS API Gateway is 'EDGE'. This query is useful in understanding and managing your API Gateway configurations, especially when dealing with edge-optimized API setups.Analyze the settings to understand the distribution of edge endpoint types within your AWS API Gateway domain names. This can help optimize your API's performance by identifying areas that may benefit from a different endpoint type.


```sql+postgres
select
  domain_name,
  config ->> 'EndpointType' as endpoint_type
from
  aws_api_gatewayv2_domain_name
  cross join jsonb_array_elements(domain_name_configurations) as config
where
  config ->> 'EndpointType' = 'EDGE';
```

```sql+sqlite
select
  domain_name,
  json_extract(config.value, '$.EndpointType') as endpoint_type
from
  aws_api_gatewayv2_domain_name,
  json_each(domain_name_configurations) as config
where
  json_extract(config, '$.EndpointType') = 'EDGE';
```

### API gatewayv2 domain name configuration info
Determine the configuration details of your API Gateway's domain name to understand its security policy, certificate details, and status. This information can be useful when troubleshooting issues or assessing the security posture of your API Gateway."Explore the configuration details of your API Gateway domain names to understand their current status, security policies, and associated certificates. This can help in managing your domain names and ensuring their secure and optimal operation."


```sql+postgres
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

```sql+sqlite
select
  domain_name,
  json_extract(config.value, '$.EndpointType') as endpoint_type,
  json_extract(config.value, '$.CertificateName') as certificate_name,
  json_extract(config.value, '$.CertificateArn') as certificate_arn,
  json_extract(config.value, '$.CertificateUploadDate') as certificate_upload_date,
  json_extract(config.value, '$.DomainNameStatus') as domain_name_status,
  json_extract(config.value, '$.DomainNameStatusMessage') as domain_name_status_message,
  json_extract(config.value, '$.ApiGatewayDomainName') as api_gateway_domain_name,
  json_extract(config.value, '$.HostedZoneId') as hosted_zone_id,
  json_extract(config.value, '$.OwnershipVerificationCertificateArn') as ownership_verification_certificate_arn,
  json_extract(config.value, '$.SecurityPolicy') as security_policy
from
  aws_api_gatewayv2_domain_name,
  json_each(domain_name_configurations) as config;
```

### Get mutual TLS authentication configuration of each domain name
Explore the setup of mutual TLS authentication for each domain name, focusing on the truststore details. This can be beneficial for understanding the security measures in place and identifying any potential warnings or issues.Explore the configuration of mutual TLS authentication for each domain name, which can help you identify potential security issues and ensure that your domains are properly secured. This can be particularly useful for maintaining compliance and identifying any domains that may require additional security measures.


```sql+postgres
select
  domain_name,
  mutual_tls_authentication ->> 'TruststoreUri' as truststore_uri,
  mutual_tls_authentication ->> 'TruststoreVersion' as truststore_version,
  mutual_tls_authentication ->> 'TruststoreWarnings' as truststore_warnings
from
  aws_api_gatewayv2_domain_name;
```

```sql+sqlite
select
  domain_name,
  json_extract(mutual_tls_authentication, '$.TruststoreUri') as truststore_uri,
  json_extract(mutual_tls_authentication, '$.TruststoreVersion') as truststore_version,
  json_extract(mutual_tls_authentication, '$.TruststoreWarnings') as truststore_warnings
from
  aws_api_gatewayv2_domain_name;
```

### Get certificate details of each domain names
Determine the specifics of certificates associated with each domain name, including their creation and issuance details, key algorithm, and transparency logging preferences. This can help in managing and maintaining the security aspects of your domain names.This query allows you to examine the details of certificates associated with each domain name. It's useful for understanding the security measures in place for your domains, such as the issuing authority, creation date, and key algorithm.


```sql+postgres
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

```sql+sqlite
select
  d.domain_name,
  json_extract(config.value, '$.CertificateArn') as certificate_arn,
  c.certificate,
  c.certificate_transparency_logging_preference,
  c.created_at,
  c.imported_at,
  c.issuer,
  c.issued_at,
  c.key_algorithm
from
  aws_api_gatewayv2_domain_name AS d,
  json_each(d.domain_name_configurations) AS config
  left join aws_acm_certificate AS c ON c.certificate_arn = json_extract(config.value, '$.CertificateArn');
```