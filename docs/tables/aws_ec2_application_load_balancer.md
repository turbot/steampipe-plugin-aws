# Table: aws_ec2_application_load_balancer

Application load balancers support content-based routing, and supports applications that run in containers and also provide additional visibility into the health of the target instances and containers.

## Examples

### Security group attached to the application load balancers

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as attached_security_group
from
  aws_ec2_application_load_balancer;
```


### Availability zone information

```sql
select
  name,
  az ->> 'LoadBalancerAddresses' as load_balancer_addresses,
  az ->> 'OutpostId' as outpost_id,
  az ->> 'SubnetId' as subnet_id,
  az ->> 'ZoneName' as zone_name
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(availability_zones) as az;
```


### List of application load balancers whose availability zone count is less than 1

```sql
select
  name,
  count(az ->> 'ZoneName') < 2 as zone_count_1
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(availability_zones) as az
group by
  name;
```


### List of application load balancers whose logging is not enabled

```sql
select
  name,
  lb ->> 'Key' as logging_key,
  lb ->> 'Value' as logging_value
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'access_logs.s3.enabled'
  and lb ->> 'Value' = 'false';
```


### List of application load balancers whose deletion protection is not enabled

```sql
select
  name,
  lb ->> 'Key' as deletion_protection_key,
  lb ->> 'Value' as deletion_protection_value
from
  aws_ec2_application_load_balancer
  cross join jsonb_array_elements(load_balancer_attributes) as lb
where
  lb ->> 'Key' = 'deletion_protection.enabled'
  and lb ->> 'Value' = 'false';
```
