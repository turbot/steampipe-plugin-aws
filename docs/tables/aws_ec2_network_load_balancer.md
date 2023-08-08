# Table: aws_ec2_network_load_balancer

AWS Network Load Balancer (NLB) distributes end user traffic across multiple cloud resources to ensure low latency and high throughput for applications. When a target becomes slow or unavailable, the Network Load Balancer routes traffic to another target.

## Examples

### Count of AZs registered with network load balancers

```sql
select
  name,
  count(az ->> 'ZoneName') as zone_count
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```


### List of network load balancers where Cross-Zone Load Balancing is enabled

```sql
select
  name,
  lb ->> 'Key' as cross_zone,
  lb ->> 'Value' as cross_zone_value
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'load_balancing.cross_zone.enabled'
  and lb ->> 'Value' = 'false';
```


### List of network load balancers where logging is not enabled

```sql
select
  name,
  lb ->> 'Key' as logging_key,
  lb ->> 'Value' as logging_value
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'access_logs.s3.enabled'
  and lb ->> 'Value' = 'false';
```


### List of network load balancers where deletion protection is not enabled

```sql
select
  name,
  lb ->> 'Key' as deletion_protection_key,
  lb ->> 'Value' as deletion_protection_value
from
  aws_ec2_network_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'deletion_protection.enabled'
  and lb ->> 'Value' = 'false';
```
