# Table: aws_elasticache_parameter_group

Parameter Groups are a collection of parameters which control the behavior of the ElastiCache cluster.

### Basic parameter group info

```sql
select
  cache_parameter_group_name,
  description,
  cache_parameter_group_family,
  description,
  is_global
from
  aws_elasticache_parameter_group;
```


### List elasticache parameter groups in the undesired (for example, redis5.0 and memcached1.5 are desired) family types.

```sql
select
  cache_parameter_group_family,
  count(*) as count
from
  aws_elasticache_parameter_group
where
  cache_parameter_group_family not in ('redis5.0', 'memcached1.5')
group by
  cache_parameter_group_family;
```
