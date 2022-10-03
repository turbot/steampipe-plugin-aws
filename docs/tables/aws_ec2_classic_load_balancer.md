# Table: aws_ec2_classic_load_balancer

Classic Load Balancer provides basic load balancing across multiple Amazon EC2 instances and operates at both the request level and connection level. Classic Load Balancer is intended for applications that are built within the EC2-Classic network.

## Examples

### Instances associated with classic load balancers

```sql
select
  name,
  instances
from
  aws_ec2_classic_load_balancer;
```


### List of classic load balancers whose logging is not enabled

```sql
select
  name,
  access_log_enabled
from
  aws_ec2_classic_load_balancer
where
  access_log_enabled = 'false';
```


### Security groups attached to each classic load balancer

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as sg
from
  aws_ec2_classic_load_balancer;
```


### Classic load balancers listener info

```sql
select
  name,
  listener_description -> 'Listener' ->> 'InstancePort' as instance_port,
  listener_description -> 'Listener' ->> 'InstanceProtocol' as instance_protocol,
  listener_description -> 'Listener' ->> 'LoadBalancerPort' as load_balancer_port,
  listener_description -> 'Listener' ->> 'Protocol' as load_balancer_protocol,
  listener_description -> 'SSLCertificateId' ->> 'SSLCertificateId' as ssl_certificate,
  listener_description -> 'Listener' ->> 'PolicyNames' as policy_names
from
  aws_ec2_classic_load_balancer
  cross join jsonb_array_elements(listener_descriptions) as listener_description;
```


### Health check info

```sql
select
  name,
  healthy_threshold,
  health_check_interval,
  health_check_target,
  health_check_timeout,
  unhealthy_threshold
from
  aws_ec2_classic_load_balancer;
```