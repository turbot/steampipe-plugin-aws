# Table: aws_elasticache_subnet_group

A subnet group is a collection of subnets that you can designate for your clusters running in an Amazon Virtual Private Cloud (VPC) environment.

## Examples

### Basic Subnet Group info

```sql
select
  cache_subnet_group_name,
  cache_subnet_group_description,
  region,
  account_id
from
  aws_elasticache_subnet_group;
```


### VPCs and Subnets info of each subnet group

```sql
select
  vpc_id,
  sub -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  sub ->> 'SubnetIdentifier' as subnet_identifier,
  sub ->> 'SubnetOutpost' as subnet_outpost
from
  aws_elasticache_subnet_group,
  jsonb_array_elements(subnets) as sub;
```


### Elasticache clusters in each subnet group

```sql
select
  cache_cluster_id,
  cache_subnet_group_name,
  vpc_id
from
  aws_elasticache_subnet_group as sb
  join aws_elasticache_cluster as c on sb.cache_subnet_group_name = c.cache_subnet_group_name;
```
