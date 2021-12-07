# Table: aws_vpc_peering_connection

A VPC peering connection is a networking connection between two VPCs that enables you to route traffic between them using private IPv4 addresses or IPv6 addresses. Instances in either VPC can communicate with each other as if they are within the same network. You can create a VPC peering connection between your own VPCs, or with a VPC in another AWS account. The VPCs can be in different regions (also known as an inter-region VPC peering connection).

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
