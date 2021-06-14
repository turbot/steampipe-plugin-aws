# Table: aws_ec2_ssl_policy

An SSL security policy is a combination of protocols and ciphers. Elastic Load Balancing uses a Secure Socket Layer (SSL) negotiation configuration, known as a security policy, to negotiate SSL connections between a client and the load balancer.

## Examples

### Basic info

```sql
select
  name,
  ssl_protocols
from
  aws_ec2_ssl_policy;
```


### List load balancer listeners that use an SSL policy with weak ciphers

```sql
select
  arn,
  ssl_policy
from
  aws_ec2_load_balancer_listener listener
join 
  aws_ec2_ssl_policy ssl_policy
on
  listener.ssl_policy = ssl_policy.Name
where
  ssl_policy.ciphers @> '[{"Name":"DES-CBC3-SHA"}]';
```
