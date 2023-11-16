---
title: "Table: aws_api_gateway_domain_name - Query AWS API Gateway Domain Names using SQL"
description: "Allows users to query AWS API Gateway Domain Names and retrieve details about each domain's configuration, certificate, and associated API."
---

# Table: aws_api_gateway_domain_name - Query AWS API Gateway Domain Names using SQL

The `aws_api_gateway_domain_name` table in Steampipe provides information about domain names within AWS API Gateway. This table allows DevOps engineers to query domain-specific details, including the domain name, certificate details, and the associated API. Users can utilize this table to gather insights on domains, such as the domain's endpoint configuration, the type of certificate used, and the API it's associated with. The schema outlines the various attributes of the domain name, including the domain name, certificate upload date, certificate ARN, and endpoint configuration.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gateway_domain_name` table, you can use the `.inspect aws_api_gateway_domain_name` command in Steampipe.

**Key columns**:

- `domain_name`: The name of the domain. This can be used to join with other tables that contain information about the domain.
- `certificate_arn`: The Amazon Resource Name (ARN) of an AWS-managed certificate. This can be used to join with other tables that contain certificate information.
- `api_mapping_selection_expression`: The version of the API Gateway REST API that the stage should point to. This can be used to join with other tables that contain API version information.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
select
  domain_name,
  aws_api_gateway_domain_name -> 'Types' as endpoint_types,
  aws_api_gateway_domain_name -> 'VpcEndpointIds' as vpc_endpoint_ids,
from
  aws_api_gateway_domain_name;
```

### Get mutual TLS authentication configuration of each domain name

```sql
select
  domain_name,
  mutual_tls_authentication ->> 'TruststoreUri' as truststore_uri,
  mutual_tls_authentication ->> 'TruststoreVersion' as truststore_version,
  mutual_tls_authentication ->> 'TruststoreWarnings' as truststore_warnings
from
  aws_api_gateway_domain_name;
```