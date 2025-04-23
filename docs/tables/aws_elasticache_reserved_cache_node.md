---
title: "Steampipe Table: aws_elasticache_reserved_cache_node - Query AWS ElastiCache Reserved Cache Nodes using SQL"
description: "Allows users to query AWS ElastiCache Reserved Cache Nodes to gather details such as the reservation status, start time, duration, and associated metadata."
folder: "ElastiCache"
---

# Table: aws_elasticache_reserved_cache_node - Query AWS ElastiCache Reserved Cache Nodes using SQL

AWS ElastiCache Reserved Cache Nodes are a type of node that you can purchase for a one-time, upfront payment in order to reserve capacity for future use. These nodes provide you with a significant discount compared to standard on-demand cache node pricing. They are ideal for applications with steady-state or predictable usage and can be used in any available AWS region.

## Table Usage Guide

The `aws_elasticache_reserved_cache_node` table in Steampipe provides you with information about the reserved cache nodes within AWS ElastiCache. This table allows you, as a DevOps engineer, to query reserved cache node-specific details, including the reservation status, start time, and duration. You can utilize this table to gather insights on reserved cache nodes, such as their current status, the time at which the reservation started, the duration of the reservation, and more. The schema outlines the various attributes of the reserved cache node for you, including the reserved cache node ID, cache node type, start time, duration, fixed price, usage price, cache node count, product description, offering type, state, recurring charges, and associated tags.

## Examples

### Basic info
Explore which AWS ElastiCache reserved nodes are currently active, and gain insights into their type and associated offering IDs. This can help in managing resources and planning for future capacity needs.

```sql+postgres
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node;
```

```sql+sqlite
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
Identify the reserved cache nodes that have been fully paid for upfront. This can help to manage costs and understand the financial commitment made for these resources.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which cache nodes are reserved for the longest duration. This can help prioritize which nodes to investigate for potential cost savings or performance improvements.

```sql+postgres
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

```sql+sqlite
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
Identify the reserved cache nodes within your AWS ElastiCache service, organized by their usage price. This can help prioritize cost management efforts by highlighting the most expensive nodes first.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which reserved cache nodes are not currently active, allowing for an assessment of resources that may be underutilized or potentially misallocated within your AWS ElastiCache service.

```sql+postgres
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

```sql+sqlite
select
  reserved_cache_node_id,
  arn,
  reserved_cache_nodes_offering_id,
  state,
  cache_node_type
from
  aws_elasticache_reserved_cache_node
where
  state != 'active';
```