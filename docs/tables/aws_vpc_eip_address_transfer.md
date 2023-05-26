# Table: aws_vpc_eip_address_transfer

An Elastic IP Address Transfer is a static, public IPv4 address that transfers from an AWS account (source account) to any other AWS account in the same AWS Region.

## Examples

### Basic info

```sql
select
  allocation_id,
  address_transfer_status,
  public_ip,
  transfer_account_id,
  transfer_offer_accepted_timestamp
from
  aws_vpc_eip_address_transfer;
```

### List address transfers accepted in last 30 days

```sql
select
  allocation_id,
  address_transfer_status,
  public_ip,
  transfer_account_id,
  transfer_offer_accepted_timestamp
from
  aws_vpc_eip_address_transfer
where
  transfer_offer_accepted_timestamp >= now() - interval '30' day;
```

### List address transfers expiring in 10 days

```sql
select
  allocation_id,
  address_transfer_status,
  public_ip,
  transfer_account_id,
  transfer_offer_expiration_timestamp
from
  aws_vpc_eip_address_transfer
where
  transfer_offer_expiration_timestamp >= now() - interval '10' day;
```

### Get VPC details for the elastic IP address transfer

```sql
select
  t.allocation_id,
  t.address_transfer_status,
  t.transfer_account_id,
  i.vpc_id,
  v.cidr_block,
  v.state,
  v.is_default
from
  aws_vpc_eip eip,
  aws_ec2_instance i,
  aws_vpc_eip_address_transfer t,
  aws_vpc v
where
  eip.instance_id = i.instance_id
  and t.allocation_id = eip.allocation_id
  and v.vpc_id = i.vpc_id;
```