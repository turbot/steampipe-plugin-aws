---
title: "Steampipe Table: aws_vpc_eip_address_transfer - Query AWS VPC Elastic IP Address Transfers using SQL"
description: "Allows users to query Elastic IP Address Transfers in AWS VPC."
folder: "VPC"
---

# Table: aws_vpc_eip_address_transfer - Query AWS VPC Elastic IP Address Transfers using SQL

The AWS VPC Elastic IP Address Transfer is a feature within Amazon's Virtual Private Cloud (VPC) service. It allows for the allocation, association, or release of Elastic IP addresses in your VPC. This feature is crucial for managing network connections and ensuring high availability of services, by allowing re-mapping of addresses to instances in case of instance or Availability Zone failures.

## Table Usage Guide

The `aws_vpc_eip_address_transfer` table in Steampipe provides you with information about Elastic IP Address Transfers within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query transfer-specific details, including associated Elastic IPs, transfer status, and associated metadata. You can utilize this table to gather insights on transfer events, such as tracking the movement of Elastic IPs, monitoring the status of IP transfers, and more. The schema outlines the various attributes of the Elastic IP Address transfer for you, including the transfer ID, allocation ID, status, and associated tags.

## Examples

### Basic info
Uncover the details of Elastic IP address transfers within your AWS VPC to gain insights into their statuses and the accounts involved in the process. This can be helpful in tracking IP address allocation and managing network resources effectively.

```sql+postgres
select
  allocation_id,
  address_transfer_status,
  public_ip,
  transfer_account_id,
  transfer_offer_accepted_timestamp
from
  aws_vpc_eip_address_transfer;
```

```sql+sqlite
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
Determine the instances where IP address transfers were accepted in the past month. This is useful for keeping track of recent changes in your network infrastructure.

```sql+postgres
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

```sql+sqlite
select
  allocation_id,
  address_transfer_status,
  public_ip,
  transfer_account_id,
  transfer_offer_accepted_timestamp
from
  aws_vpc_eip_address_transfer
where
  transfer_offer_accepted_timestamp >= datetime('now','-30 day');
```

### List address transfers expiring in 10 days
Determine the areas in which Elastic IP address transfers within your AWS VPC are set to expire within the next 10 days. This can assist in effectively managing your IP resources and avoiding unexpected address expiration.

```sql+postgres
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

```sql+sqlite
select
  allocation_id,
  address_transfer_status,
  public_ip,
  transfer_account_id,
  transfer_offer_expiration_timestamp
from
  aws_vpc_eip_address_transfer
where
  transfer_offer_expiration_timestamp >= datetime('now', '-10 days');
```

### Get VPC details for the elastic IP address transfer
Discover the segments that are involved in the transfer of an elastic IP address within a Virtual Private Cloud (VPC). This is beneficial for managing network configurations and understanding the status and details of the IP address transfer.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```