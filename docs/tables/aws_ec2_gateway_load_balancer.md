# Table: aws_ec2_gateway_load_balancer

Gateway Load Balancer makes it easy to deploy, scale, and run third-party virtual networking appliances. Providing load balancing and auto scaling for fleets of third-party appliances, Gateway Load Balancer is transparent to the source and destination of traffic. This capability makes it well suited for working with third-party appliances for security, network analytics, and other use cases

## Examples

### Basic gateway load balancer  info

```sql
select
  name,
  arn,
  type,
  state_code,
  vpc_id,
  availability_zones
from
  aws_ec2_gateway_load_balancer

```

### Availability zone information of all the gateway load balancer

```sql
select
  name,
  az ->> 'LoadBalancerAddresses' as load_balancer_addresses,
  az ->> 'OutpostId' as outpost_id,
  az ->> 'SubnetId' as subnet_id,
  az ->> 'ZoneName' as zone_name
from
  aws_ec2_gateway_load_balancer
  cross join jsonb_array_elements(availability_zones) as az;

```

### List of gateway load balancers whose availability zone count is less than 1

```sql
select
  name,
  count(az ->> 'ZoneName') < 2 as zone_count_1
from
  aws_ec2_gateway_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```

### List of application load balancers whose deletion protection is not enabled

```sql
select
  name,
  lb ->> 'Key' as deletion_protection_key,
  lb ->> 'Value' as deletion_protection_value
from
  aws_ec2_gateway_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'deletion_protection.enabled'
  and lb ->> 'Value' = 'false';
```


### List of gateway load balancers whose load balancing cross zone is enabled

```sql
select
  name,
  lb ->> 'Key' as load_balancing_cross_zone_key,
  lb ->> 'Value' as load_balancing_cross_zone_value
from
  aws_ec2_gateway_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'load_balancing.cross_zone.enabled'
  and lb ->> 'Value' = 'true';
```

### Security group attached to the gateway load balancers

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as attached_security_group
from
  aws_ec2_gateway_load_balancer;
```

### List of gateway load balancerw with state other than active

```sql
select
  name,
  state_code
from
  aws_ec2_gateway_load_balancer
where
 state_code <> 'active'
```