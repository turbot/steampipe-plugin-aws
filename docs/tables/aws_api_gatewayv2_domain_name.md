# Table: aws_api_gatewayv2_domain_name

In AWS API Gateway, Domain Names allow you to customize the URL for your API endpoints. Instead of using the default API Gateway domain, you can associate your own custom domain name with your API. This provides a more branded and user-friendly URL for your API consumers.

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