---
title: "Steampipe Table: aws_vpc_eip - Query AWS VPC Elastic IP Addresses using SQL"
description: "Allows users to query AWS VPC Elastic IP Addresses"
folder: "VPC"
---

# Table: aws_vpc_eip - Query AWS VPC Elastic IP Addresses using SQL

An AWS VPC Elastic IP Address is a static, public IPv4 address designed for dynamic cloud computing. It allows you to mask the failure of an instance or software by rapidly remapping the address to another instance in your account. With an Elastic IP address, you can manage your own IP addressing, connectivity, security, and delivery control in the AWS Cloud.

## Table Usage Guide

The `aws_vpc_eip` table in Steampipe provides you with information about Elastic IP Addresses (EIP) within AWS Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query EIP-specific details, including allocation IDs, association IDs, domain, public IP, and associated metadata. You can utilize this table to gather insights on Elastic IP Addresses, such as their allocation status, associated instances, network interface details, and more. The schema outlines the various attributes of the Elastic IP Address for you, including the public IPv4 address, allocation ID, association ID, domain, instance ID, network interface ID, private IP address, and associated tags.

## Examples

### List of unused elastic IPs
Identify instances where Elastic IPs are not currently in use within your AWS VPC environment. This can help in resource management and cost reduction by releasing unnecessary IPs.

```sql+postgres
select
  public_ip,
  domain association_id
from
  aws_vpc_eip
where
  association_id is null;
```

```sql+sqlite
select
  public_ip,
  domain association_id
from
  aws_vpc_eip
where
  association_id is null;
```


### Count of elastic IPs by instance Ids
Discover the segments that have differing numbers of elastic IPs within your AWS VPC environment. This allows for efficient allocation and management of your public IPv4 resources.

```sql+postgres
select
  public_ipv4_pool,
  count(public_ip) as elastic_ips
from
  aws_vpc_eip
group by
  public_ipv4_pool;
```

```sql+sqlite
select
  public_ipv4_pool,
  count(public_ip) as elastic_ips
from
  aws_vpc_eip
group by
  public_ipv4_pool;
```