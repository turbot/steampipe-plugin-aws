# Table: aws_route53_resolver_endpoint

The AWS Route 53 Resolver endpoints in a Virtual Private Cloud (VPC) that is used for DNS management.

## Examples

### Basic info

```sql
select
  name,
  id,
  direction,
  ip_address_count,
  status
from
  aws_route53_resolver_endpoint;
```

### List unhealthy endpoints

```sql
select
  name,
  id,
  direction,
  status,
  status_message
from
  aws_route53_resolver_endpoint
where
  status = 'ACTION_NEEDED';
```

### Get IP address details for each endpoint

```sql
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Status' as status,
  p ->> 'SubnetId' as subnet_id
from
  aws_route53_resolver_endpoint,
  jsonb_array_elements(ip_addresses) as p;
```
