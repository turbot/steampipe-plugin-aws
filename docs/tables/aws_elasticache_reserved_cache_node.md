# Table: aws_elasticache_reserved_cache_node

Amazon ElastiCache Reserved Nodes give you the option to make a low, one-time payment for each cache node you want to reserve and in turn receive a significant discount on the hourly charge for that Node. Amazon ElastiCache provides three ElastiCache Reserved Node types (All Upfront, No Upfront, and Partial Upfront Reserved Instances) that enable you to balance the amount you pay upfront with your effective hourly price.

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
