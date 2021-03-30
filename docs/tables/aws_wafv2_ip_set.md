# Table: aws_wafv2_ip_set

An AWS WAFv2 IP Set contains one or more IP addresses or blocks of IP addresses (IPv4 and IPv6) specified in Classless Inter-Domain Routing (CIDR) notation.

## Examples

### Basic info

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set;
```


### List global (CLOUDFRONT) IP sets

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set
where
  scope = 'CLOUDFRONT';
```


### List IP sets having IPv4 address version

```sql
select
  name,
  description,
  arn,
  id,
  scope,
  addresses,
  ip_address_version,
  region
from
  aws_wafv2_ip_set
where
  ip_address_version = 'IPV4';
```


### List IP sets having a specific IP address

```sql
select
  name,
  description,
  arn,
  ip_address_version,
  region,
  address
from
  aws_wafv2_ip_set,
  jsonb_array_elements_text(addresses) as address
where
  address = '1.2.3.4/32';
```