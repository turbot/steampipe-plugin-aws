---
title: "Table: aws_vpc_eip_address_transfer - Query AWS VPC Elastic IP Address Transfers using SQL"
description: "Allows users to query Elastic IP Address Transfers in AWS VPC."
---

# Table: aws_vpc_eip_address_transfer - Query AWS VPC Elastic IP Address Transfers using SQL

The `aws_vpc_eip_address_transfer` table in Steampipe provides information about Elastic IP Address Transfers within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers to query transfer-specific details, including associated Elastic IPs, transfer status, and associated metadata. Users can utilize this table to gather insights on transfer events, such as tracking the movement of Elastic IPs, monitoring the status of IP transfers, and more. The schema outlines the various attributes of the Elastic IP Address transfer, including the transfer ID, allocation ID, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_eip_address_transfer` table, you can use the `.inspect aws_vpc_eip_address_transfer` command in Steampipe.

**Key columns**:

- `transfer_id`: This is the unique identifier for the transfer. It is useful for tracking specific transfers and can be used to join with other tables that track transfer events.
- `public_ip`: This column contains the Elastic IP that is being transferred. This is essential for identifying which IPs are involved in transfers.
- `status`: This column indicates the status of the transfer. It is useful for monitoring the progress of IP transfers.

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