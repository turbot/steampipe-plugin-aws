# Table: aws_vpc_customer_gateway

A customer gateway is a resource that is installed on the customer side and is often linked to the provider side.

## Examples

### Customer gateway basic detail

```sql
select
  customer_gateway_id,
  type,
  state,
  bgp_asn,
  certificate_arn,
  device_name,
  ip_address
from
  aws_vpc_customer_gateway;
```


### Count of customer gateways by certificate_arn

```sql
select
  type,
  count(customer_gateway_id) as customer_gateway_id_count
from
  aws_vpc_customer_gateway
group by
  type;
```
