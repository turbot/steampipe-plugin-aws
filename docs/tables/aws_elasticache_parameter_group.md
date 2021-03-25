# Table: aws_elasticache_parameter_group

Parameter Groups are a collection of parameters which control the behavior of the ElastiCache cluster.

### Basic info

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


### List parameter groups that are not compatible with redis 5.0 and memcached 1.5

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
