# Table: aws_vpc_eip

An Elastic IP Address Transfer is a static, public IPv4 address that transfer from AWS account (source account) to any other AWS account in the same AWS Region.

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

### List address transfers expires in 10 days

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