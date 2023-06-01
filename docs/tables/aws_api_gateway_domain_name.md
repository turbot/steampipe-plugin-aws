# Table: aws_api_gateway_domain_name

In AWS API Gateway, Domain Names allow you to customize the URL for your API endpoints. Instead of using the default API Gateway domain, you can associate your own custom domain name with your API. This provides a more branded and user-friendly URL for your API consumers.

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

### Get certificate details of each domain names

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