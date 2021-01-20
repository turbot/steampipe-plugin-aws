# Table: aws_vpc_egress_only_internet_gateway

An egress-only internet gateway is a horizontally scaled, redundant, and highly available VPC component that allows outbound communication over IPv6 from instances in your VPC to the internet, and prevents the internet from initiating an IPv6 connection with your instances

## Examples

### Egress only internet gateway basic info

```sql
select
  id,
  att ->> 'State' as state,
  att ->> 'VpcId' as vpc_id,
  tags,
  region
from
  aws_vpc_egress_only_internet_gateway
  cross join jsonb_array_elements(attachments) as att;
```


### List unattached egress only gateways

```sql
select
  id,
  attachments
from
  aws_vpc_egress_only_internet_gateway
where
  attachments is null;
```


### List all the egress only gateways attached to default VPC

```sql
select
  id,
  vpc.is_default
from
  aws_vpc_egress_only_internet_gateway
  cross join jsonb_array_elements(attachments) as i
  join aws_vpc vpc on i ->> 'VpcId' = vpc.vpc_id
where
  vpc.is_default = true;
```
