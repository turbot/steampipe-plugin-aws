# Table: aws_vpc_internet_gateway

An internet gateway is a horizontally scaled, redundant, and highly available VPC component that allows communication between VPC and the internet.

## Examples

### List unattached internet gateways

```sql
select
  internet_gateway_id,
  attachments
from
  aws_vpc_internet_gateway
where
  attachments is null;
```


### Find VPCs attached to the internet gateways

```sql
select
  internet_gateway_id,
  att ->> 'VpcId' as vpc_id
from
  aws_vpc_internet_gateway
  cross join jsonb_array_elements(attachments) as att;
```
