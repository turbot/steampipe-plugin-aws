---
title: "Steampipe Table: aws_api_gateway_domain_name - Query AWS API Gateway Domain Names using SQL"
description: "Allows users to query AWS API Gateway Domain Names and retrieve details about each domain's configuration, certificate, and associated API."
folder: "API Gateway"
---

# Table: aws_api_gateway_domain_name - Query AWS API Gateway Domain Names using SQL

The AWS API Gateway Domain Name is a component of Amazon's API Gateway service that allows you to create, configure, and manage a custom domain name to maintain a consistent user experience. It enables routing of incoming requests to various backend services, including AWS Lambda functions, and provides features like SSL certificates for secure communication. This is crucial for providing a seamless and secure API communication channel for your applications.

## Table Usage Guide

The `aws_api_gateway_domain_name` table in Steampipe provides you with information about domain names within AWS API Gateway. This table allows you, as a DevOps engineer, to query domain-specific details, including the domain name, certificate details, and the associated API. You can utilize this table to gather insights on domains, such as the domain's endpoint configuration, the type of certificate used, and the API it's associated with. The schema outlines the various attributes of the domain name for you, including the domain name, certificate upload date, certificate ARN, and endpoint configuration.

## Examples

### Basic info
Determine the areas in which your API Gateway domain name configurations are operating in AWS. This can help you understand the status and ownership of your domain names, providing insights into their distribution and certificate details.

```sql+postgres
select
  domain_name,
  certificate_arn,
  distribution_domain_name,
  distribution_hosted_zone_id,
  domain_name_status,
  ownership_verification_certificate_arn
from
  aws_api_gateway_domain_name;
```

```sql+sqlite
select
  domain_name,
  certificate_arn,
  distribution_domain_name,
  distribution_hosted_zone_id,
  domain_name_status,
  ownership_verification_certificate_arn
from
  aws_api_gateway_domain_name;
```

### List available domain names
Determine the areas in which domain names are available for use in the AWS API Gateway. This is beneficial for identifying potential new domains for your applications.

```sql+postgres
select
  domain_name,
  certificate_arn,
  certificate_upload_date,
  regional_certificate_arn,
  domain_name_status
from
  aws_api_gateway_domain_name
where
  domain_name_status = 'AVAILABLE';
```

```sql+sqlite
select
  domain_name,
  certificate_arn,
  certificate_upload_date,
  regional_certificate_arn,
  domain_name_status
from
  aws_api_gateway_domain_name
where
  domain_name_status = 'AVAILABLE';
```

### Get certificate details of each domain name
Discover the segments that provide detailed insights about the certificates associated with each domain name. This is useful in understanding the security measures in place and their configurations, aiding in better management of your web assets.

```sql+postgres
select
  d.domain_name,
  d.regional_certificate_arn,
  c.certificate,
  c.certificate_transparency_logging_preference,
  c.created_at,
  c.imported_at,
  c.issuer,
  c.issued_at,
  c.key_algorithm
from
  aws_api_gateway_domain_name as d,
  aws_acm_certificate as c
where
  c.certificate_arn = d.regional_certificate_arn;
```

```sql+sqlite
select
  d.domain_name,
  d.regional_certificate_arn,
  c.certificate,
  c.certificate_transparency_logging_preference,
  c.created_at,
  c.imported_at,
  c.issuer,
  c.issued_at,
  c.key_algorithm
from
  aws_api_gateway_domain_name as d,
  aws_acm_certificate as c
where
  c.certificate_arn = d.regional_certificate_arn;
```

### Get endpoint configuration details of each domain
Determine the configuration details of each domain in your AWS API Gateway to better understand the types of endpoints used and identify any associated Virtual Private Cloud (VPC) endpoints.

```sql+postgres
select
  domain_name,
  endpoint_configuration -> 'Types' as endpoint_types,
  endpoint_configuration -> 'VpcEndpointIds' as vpc_endpoint_ids
from
  aws_api_gateway_domain_name;
```

```sql+sqlite
select
  domain_name,
  json_extract(endpoint_configuration, '$.Types') as endpoint_types,
  json_extract(endpoint_configuration, '$.VpcEndpointIds') as vpc_endpoint_ids
from
  aws_api_gateway_domain_name;
```

### Get mutual TLS authentication configuration of each domain name
This query can be used to analyze the mutual TLS authentication settings for each domain name in an AWS API Gateway. It provides insights into the truststore details, which can be beneficial for improving security configurations and troubleshooting potential issues.

```sql+postgres
select
  domain_name,
  mutual_tls_authentication ->> 'TruststoreUri' as truststore_uri,
  mutual_tls_authentication ->> 'TruststoreVersion' as truststore_version,
  mutual_tls_authentication ->> 'TruststoreWarnings' as truststore_warnings
from
  aws_api_gateway_domain_name;
```

```sql+sqlite
select
  domain_name,
  json_extract(mutual_tls_authentication, '$.TruststoreUri') as truststore_uri,
  json_extract(mutual_tls_authentication, '$.TruststoreVersion') as truststore_version,
  json_extract(mutual_tls_authentication, '$.TruststoreWarnings') as truststore_warnings
from
  aws_api_gateway_domain_name;
```