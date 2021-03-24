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
  id = 'rslvr-rr-389d2ef50c094970b';
```

### List of associations that were created between Resolver rules and VPCs

```sql
select
  name,
  p ->> 'Id' as id,
  p ->> 'Status' as status,
  p ->> 'VPCId' as vpc_id
from
  aws_route53_resolver_rule,
  jsonb_array_elements(resolver_rule_associations) as p;
```
### List of IP addresses and ports that an outbound endpoint forwards DNS queries

```sql
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Port' as port
from
  aws_route53_resolver_rule,
  jsonb_array_elements(target_ips) as p;
```
### List of resolver rule not shared with another account

```sql
select
  name,
  id,
  share_status,
  rule_type
from
  aws_route53_resolver_rule
where
  share_status = 'NOT_SHARED';
```
