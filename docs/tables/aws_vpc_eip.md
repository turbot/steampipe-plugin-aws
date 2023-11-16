---
title: "Table: aws_vpc_eip - Query AWS VPC Elastic IP Addresses using SQL"
description: "Allows users to query AWS VPC Elastic IP Addresses"
---

# Table: aws_vpc_eip - Query AWS VPC Elastic IP Addresses using SQL

The `aws_vpc_eip` table in Steampipe provides information about Elastic IP Addresses (EIP) within AWS Virtual Private Cloud (VPC). This table allows DevOps engineers to query EIP-specific details, including allocation IDs, association IDs, domain, public IP, and associated metadata. Users can utilize this table to gather insights on Elastic IP Addresses, such as their allocation status, associated instances, network interface details, and more. The schema outlines the various attributes of the Elastic IP Address, including the public IPv4 address, allocation ID, association ID, domain, instance ID, network interface ID, private IP address, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_eip` table, you can use the `.inspect aws_vpc_eip` command in Steampipe.

**Key columns**:

- `public_ip`: This column stores the public IP address of the Elastic IP Address. This is an important column because it can be used to join with other tables that also contain public IP addresses.
- `allocation_id`: This column stores the allocation ID of the Elastic IP Address. This is useful because it can be used to join with other tables that also contain allocation IDs.
- `instance_id`: This column stores the ID of the instance that the Elastic IP Address is associated with. This is important because it can be used to join with other tables that also contain instance IDs.

## Examples

### List of unused elastic IPs

```sql
select
  public_ip,
  domain association_id
from
  aws_vpc_eip
where
  association_id is null;
```


### Count of elastic IPs by instance Ids

```sql
select
  public_ipv4_pool,
  count(public_ip) as elastic_ips
from
  aws_vpc_eip
group by
  public_ipv4_pool;
```
