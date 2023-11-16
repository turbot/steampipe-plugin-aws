---
title: "Table: aws_vpc_peering_connection - Query AWS VPC Peering Connections using SQL"
description: "Allows users to query VPC Peering Connections in Amazon Virtual Private Cloud (VPC)."
---

# Table: aws_vpc_peering_connection - Query AWS VPC Peering Connections using SQL

The `aws_vpc_peering_connection` table in Steampipe provides information about VPC Peering Connections within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers, security teams, and system administrators to query peering connection-specific details, including peering statuses, VPC IDs, region, and associated metadata. Users can utilize this table to gather insights on peering connections, such as connection status, verification of peering options, and more. The schema outlines the various attributes of the VPC peering connection, including the peering connection ID, creation date, requester VPC info, accepter VPC info, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_peering_connection` table, you can use the `.inspect aws_vpc_peering_connection` command in Steampipe.

**Key columns**:

- `pcx_id`: The ID of the VPC peering connection. This is the primary key of the table and can be used to join this table with other tables.
- `requester_vpc_id`: The ID of the requester VPC. This can be used to join this table with the `aws_vpc` table to get more details about the requester VPC.
- `accepter_vpc_id`: The ID of the accepter VPC. This can be used to join this table with the `aws_vpc` table to get more details about the accepter VPC.

## Examples

### Basic Info

```sql
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

```sql
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

```sql
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

### List accepter VPC connection details

```sql
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

### List VPC peering connections by specific VPC peering connection IDs

```sql
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

```sql
select
  id,
  jsonb_pretty(tags) as tags,
  jsonb_pretty(tags_src) as tags_src
from
  aws_vpc_peering_connection;
```

```sql
select
  id,
  tags ->> 'Name' as name
from
  aws_vpc_peering_connection;
```

### List VPC peering connections by specific tag's key

```sql
select
  id,
  jsonb_pretty(tags) as tags
from
  aws_vpc_peering_connection,
  jsonb_each(tags)
where
  key = 'turbot:TurbotCreatedPeeringConnection';
```

### List VPC peering connections by specific tag's key & value

```sql
select
  id,
  jsonb_pretty(tags) as tags
from
  aws_vpc_peering_connection
where
  tags @> '{"Name": "vpc-0639e12347e5b6bfb <=> vpc-8e1234f5"}';
```
