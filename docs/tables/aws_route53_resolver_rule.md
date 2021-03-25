# Table: aws_route53_resolver_rule

The AWS Route 53 Resolver rule specifies how to route DNS queries out of the VPC.

## Examples

### List all rules

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


### Get a specific rule

```sql
select
  name,
  domain_name owner_id,
  resolver_endpoint_id,
  rule_type,
  share_status,
  status
from
  aws_route53_resolver_rule
where
  id = 'rslvr-rr-389d2ef50c092349a';
```


### List all rules where VPCs are not associated

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


### List the IP addresses enabled for outbound DNS queries

```sql
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Port' as port
from
  aws_route53_resolver_rule,
  jsonb_array_elements(target_ips) as p;
```


### List of resolver rule shared with another account

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