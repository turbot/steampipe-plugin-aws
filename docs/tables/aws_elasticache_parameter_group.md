# Table: aws_elasticache_parameter_group

Parameter Groups are a collection of parameters which control the behaviour of the ElastiCache cluster.

### List of global elasticache parameter groups

```sql
select
  cache_parameter_group_name,
  cache_parameter_group_family,
  description,
  is_global
from
  aws_elasticache_parameter_group
where
  is_global;
```


### Count of elasticache parameter groups in the undesired (for example redis5.0 and memcached1.5 is desired) family type.

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
