# Table: aws_elasticache_cluster

A cluster is a collection of one or more cache nodes, all of which run an instance of the Redis cache engine software.

## Examples

### List of unencrypted elasticache clusters

```sql
select
	cache_cluster_id,
	cache_node_type,
	at_rest_encryption_enabled
from
	aws_elasticache_cluster
where
	not at_rest_encryption_enabled;
```


### List of elasticache clusters whose availability zone count is less than 2

```sql
select
	cache_cluster_id,
	preferred_availability_zone
from
	aws_elasticache_cluster
where
	preferred_availability_zone <> 'Multiple';
```


### List of elasticache cluster that DO NOT enforce encryption in transit

```sql
select
	cache_cluster_id,
	cache_node_type,
	transit_encryption_enabled
from
	aws_elasticache_cluster
where
	not transit_encryption_enabled;
```


### Count of elasticache cluster provisioned with undesired(for example cache.m5.large and cache.m4.4xlarge is desired) node type(s).

```sql
select
  cache_node_type,
  count(*) as count
from
  aws_elasticache_cluster
where
  cache_node_type not in ('cache.m5.large', 'cache.m4.4xlarge')
group by
  cache_node_type;
```


### List of elasticache cluster which has inactive notification configuration topics

```sql
select
	cache_cluster_id,
	cache_cluster_status,
	notification_configuration ->> 'TopicArn' as topic_arn,
	notification_configuration ->> 'TopicStatus' as topic_status
from
	aws_elastic_cache_cluster
where
	notification_configuration ->> 'TopicStatus' = 'inactive';
```


### Security group details attached with the cache cluster

```sql
select
	cache_cluster_id,
	sg ->> 'SecurityGroupId' as security_group_id,
	sg ->> 'Status' as status
from
	aws_elasticache_cluster,
	jsonb_array_elements(security_groups) as sg;
```