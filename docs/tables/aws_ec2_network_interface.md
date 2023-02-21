# Table: aws_ec2_network_interface

An AWS EC2 Network interface represents an elastic network interface (ENI) in AWS.

## Examples

### Basic IP address info

```sql
select
  network_interface_id,
  interface_type,
  description,
  private_ip_address,
  association_public_ip,
  mac_address
from
  aws_ec2_network_interface;
```

### Find all ENIs with private IPs that are in a given subnet (10.66.0.0/16)

```sql
select
  network_interface_id,
  interface_type,
  description,
  private_ip_address,
  association_public_ip,
  mac_address
from
  aws_ec2_network_interface
where
  private_ip_address :: cidr <<= '10.66.0.0/16';
```

### Count of ENIs by interface type

```sql
select
  interface_type,
  count(interface_type) as count
from
  aws_ec2_network_interface
group by
  interface_type
order by
  count desc;
```

### Security groups attached to each ENI

```sql
select
  network_interface_id as eni,
  sg ->> 'GroupId' as "security group id",
  sg ->> 'GroupName' as "security group name"
from
  aws_ec2_network_interface
  cross join jsonb_array_elements(groups) as sg
order by
  eni;
```

### Get network details for each ENI

```sql
select
  e.network_interface_id,
  v.vpc_id,
  v.is_default,
  v.cidr_block,
  v.state,
  v.account_id,
  v.region
from
  aws_ec2_network_interface e,
  aws_vpc v
where 
  e.vpc_id = v.vpc_id;
```