# Table: aws_api_gatewayv2_domain_name

The domain name is a component of a uniform resource locator (URL) used to access web sites,

## Examples

### List of all edge endpoint type domain name

```sql
select
  domain_name,
  config ->> 'endpointType' as endpoint_type
from
  aws_api_gatewayv2_domain_name
  cross join jsonb_array_elements(domain_name_configurations) as config
where
  config ->> 'endpointType' = 'EDGE';
```


### API gatewayv2 domain name configuration info

```sql
select
  domain_name,
  config ->> 'endpointType' as endpoint_type,
  config ->> 'certificateName' as certificate_name,
  config ->> 'apiGatewayDomainName' as api_gateway_domain_name,
  config ->> 'hostedZoneId' as hosted_zone_id,
  config -> 'securityPolicy' as security_policy
from
  aws_api_gatewayv2_domain_name
  cross join jsonb_array_elements(domain_name_configurations) as config;
```
