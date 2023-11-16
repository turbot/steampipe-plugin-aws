---
title: "Table: aws_elasticache_reserved_cache_node - Query AWS ElastiCache Reserved Cache Nodes using SQL"
description: "Allows users to query AWS ElastiCache Reserved Cache Nodes to gather details such as the reservation status, start time, duration, and associated metadata."
---

# Table: aws_elasticache_reserved_cache_node - Query AWS ElastiCache Reserved Cache Nodes using SQL

The `aws_elasticache_reserved_cache_node` table in Steampipe provides information about the reserved cache nodes within AWS ElastiCache. This table allows DevOps engineers to query reserved cache node-specific details, including the reservation status, start time, and duration. Users can utilize this table to gather insights on reserved cache nodes, such as their current status, the time at which the reservation started, the duration of the reservation, and more. The schema outlines the various attributes of the reserved cache node, including the reserved cache node ID, cache node type, start time, duration, fixed price, usage price, cache node count, product description, offering type, state, recurring charges, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_reserved_cache_node` table, you can use the `.inspect aws_elasticache_reserved_cache_node` command in Steampipe.

**Key columns**:

- `reserved_cache_node_id`: The unique identifier for the reservation. This can be used to join with other tables that contain information about the reserved cache node.

- `cache_node_type`: The type of the cache node. This can be useful when joining with other tables that contain information about the type of cache nodes.

- `state`: The current state of the reservation. This can be useful when joining with other tables that contain information about the state of the reservations.

## Examples

### Basic info

```sql
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node;
```

### List reserved cache nodes with offering type `All Upfront`

```sql
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node
where
  offering_type = 'All Upfront';
```

### List reserved cache nodes order by duration

```sql
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node
order by
  duration desc;
```

### List reserved cache nodes order by usage price

```sql
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node
order by
  usage_price desc;
```

### List reserved cache nodes which are not active

```sql
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node
where
  state <> 'active';
```
