---
title: "Steampipe Table: aws_vpc_peering_connection - Query AWS VPC Peering Connections using SQL"
description: "Allows users to query VPC Peering Connections in Amazon Virtual Private Cloud (VPC)."
folder: "VPC"
---

# Table: aws_vpc_peering_connection - Query AWS VPC Peering Connections using SQL

The AWS VPC Peering Connection is a networking connection between two Virtual Private Clouds (VPCs) that enables you to route traffic between them using private IPv4 addresses or IPv6 addresses. It's a one-to-one relationship, and doesn't require a gateway, VPN connection, or separate network hardware. This connection can be made between your own VPCs, or with a VPC in another AWS account within a single region.

## Table Usage Guide

The `aws_vpc_peering_connection` table in Steampipe provides you with information about VPC Peering Connections within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, security team member, or system administrator, to query peering connection-specific details, including peering statuses, VPC IDs, region, and associated metadata. You can utilize this table to gather insights on peering connections, such as connection status, verification of peering options, and more. The schema outlines the various attributes of the VPC peering connection for you, including the peering connection ID, creation date, requester VPC info, accepter VPC info, and associated tags.

## Examples

### Basic Info
Determine the areas in which Virtual Private Cloud (VPC) peering connections are established between different AWS accounts and regions. This information can help you manage network access, improve security, and optimize resource allocation.

```sql+postgres
select
  id,
  accepter_owner_id,
  accepter_region,
  accepter_vpc_id,
  expiration_time,
  requester_owner_id,
  requester_region,
  requester_vpc_id
from
  aws_vpc_peering_connection;
```

```sql+sqlite
select
  id,
  accepter_owner_id,
  accepter_region,
  accepter_vpc_id,
  expiration_time,
  requester_owner_id,
  requester_region,
  requester_vpc_id
from
  aws_vpc_peering_connection;
```

### List VPC peering connections by approval status
Determine the areas in which Virtual Private Cloud (VPC) peering connections are still awaiting approval. This is particularly useful in managing network accessibility and ensuring secure and efficient data transfer between different VPCs.

```sql+postgres
select
  id,
  accepter_vpc_id,
  requester_vpc_id,
  status_code,
  status_message
from
  aws_vpc_peering_connection
where
  status_code = 'pending-acceptance';
```

```sql+sqlite
select
  id,
  accepter_vpc_id,
  requester_vpc_id,
  status_code,
  status_message
from
  aws_vpc_peering_connection
where
  status_code = 'pending-acceptance';
```

### List requester VPC connection details
Discover the segments that involve details of VPC connections requested in an AWS environment, enabling you to understand who is requesting connections, from which regions, and what their specific peering options are.

```sql+postgres
select
  id,
  requester_cidr_block,
  requester_owner_id,
  requester_region,
  requester_vpc_id,
  jsonb_pretty(requester_cidr_block_set) as requester_cidr_block_set,
  jsonb_pretty(requester_ipv6_cidr_block_set) as requester_ipv6_cidr_block_set,
  jsonb_pretty(requester_peering_options) as requester_peering_options
from
  aws_vpc_peering_connection;
```

```sql+sqlite
select
  id,
  requester_cidr_block,
  requester_owner_id,
  requester_region,
  requester_vpc_id,
  requester_cidr_block_set,
  requester_ipv6_cidr_block_set,
  requester_peering_options
from
  aws_vpc_peering_connection;
```

### List accepter VPC connection details
Explore the details of accepted VPC connections to understand their configurations, ownership, and regional distribution. This can aid in managing network access and ensuring secure data transfers within your AWS environment.

```sql+postgres
select
  id,
  accepter_cidr_block,
  accepter_owner_id,
  accepter_region,
  accepter_vpc_id,
  jsonb_pretty(accepter_cidr_block_set) as accepter_cidr_block_set,
  jsonb_pretty(accepter_ipv6_cidr_block_set) as accepter_ipv6_cidr_block_set,
  jsonb_pretty(accepter_peering_options) as accepter_peering_options
from
  aws_vpc_peering_connection;
```

```sql+sqlite
select
  id,
  accepter_cidr_block,
  accepter_owner_id,
  accepter_region,
  accepter_vpc_id,
  accepter_cidr_block_set,
  accepter_ipv6_cidr_block_set,
  accepter_peering_options
from
  aws_vpc_peering_connection;
```

### List VPC peering connections by specific VPC peering connection IDs
This query is useful to identify specific VPC peering connections by their IDs. It allows you to gain insights into the ownership, region, and associated VPC details of both the accepter and requester, which can be beneficial for network management and troubleshooting tasks.

```sql+postgres
select
  id,
  accepter_owner_id,
  accepter_region,
  accepter_vpc_id,
  expiration_time,
  requester_owner_id,
  requester_region,
  requester_vpc_id
from
  aws_vpc_peering_connection
where
  id in ('pcx-0a0403619dd2f3b24', 'pcx-048825e2c43ffd99e');
```

```sql+sqlite
select
  id,
  accepter_owner_id,
  accepter_region,
  accepter_vpc_id,
  expiration_time,
  requester_owner_id,
  requester_region,
  requester_vpc_id
from
  aws_vpc_peering_connection
where
  id in ('pcx-0a0403619dd2f3b24', 'pcx-048825e2c43ffd99e');
```

### List VPC peering connections with tag details
Discover the segments that have VPC peering connections along with their respective tags. This can be beneficial for gaining insights into the connections between different VPCs and their associated metadata, aiding in network management and security audits.

```sql+postgres
select
  id,
  jsonb_pretty(tags) as tags,
  jsonb_pretty(tags_src) as tags_src
from
  aws_vpc_peering_connection;
```

```sql+sqlite
select
  id,
  tags,
  tags_src
from
  aws_vpc_peering_connection;
```

```sql+postgres
select
  id,
  tags ->> 'Name' as name
from
  aws_vpc_peering_connection;
```

```sql+sqlite
select
  id,
  json_extract(tags, '$.Name') as name
from
  aws_vpc_peering_connection;
```

### List VPC peering connections by specific tag's key
Explore which Virtual Private Cloud (VPC) peering connections have been specifically marked with the 'turbot:TurbotCreatedPeeringConnection' tag. This could be useful in understanding and managing connections that were automatically created by Turbot, a cloud governance platform.

```sql+postgres
select
  v.id,
  jsonb_pretty(tags) as tags
from
  aws_vpc_peering_connection as v,
  jsonb_each(tags)
where
  key = 'turbot:TurbotCreatedPeeringConnection';
```

```sql+sqlite
select
  v.id,
  json_extract(t.value, '$') as tags
from
  aws_vpc_peering_connection as v,
  json_each(tags) as t
where
  json_extract(t.value, '$.key') = 'turbot:TurbotCreatedPeeringConnection';
```

### List VPC peering connections by specific tag's key & value
Discover the segments that have specific peering connections within a virtual private cloud (VPC) network using specific tags. This is useful for managing and organizing your network connections based on their assigned tags.

```sql+postgres
select
  id,
  jsonb_pretty(tags) as tags
from
  aws_vpc_peering_connection
where
  tags @> '{"Name": "vpc-0639e12347e5b6bfb <=> vpc-8e1234f5"}';
```

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```
