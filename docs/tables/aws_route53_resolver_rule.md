# Table: aws_route53_resolver_rule

The AWS Route 53 Resolver rule specifies how to route DNS queries out of the VPC.

## Examples

### Basic info

```sql
select
  name,
  domain_name owner_id,
  resolver_endpoint_id,
  rule_type,
  share_status,
  status
from
  aws_route53_resolver_rule;
```


### List rules that are not associated with VPCs

```sql
select
  name,
  id,
  arn,
  resolver_rule_associations
from
  aws_route53_resolver_rule
Where
  resolver_rule_associations = '[]';
```


### List the IP addresses enabled for outbound DNS queries for each rule

```sql
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Port' as port
from
  aws_route53_resolver_rule,
  jsonb_array_elements(target_ips) as p;
```


### List resolver rules shared with another account

```sql
select
  name,
  id,
  share_status,
  rule_type
from
  aws_route53_resolver_rule
where
  share_status = 'SHARED';
```
