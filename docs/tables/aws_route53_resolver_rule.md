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

### List of ip addresses that an outbound endpoint forwards DNS queries to

```sql
select
	name,
	p ->> 'Ip' as ip
from
	aws_route53_resolver_rule,
	jsonb_array_elements(target_ips) as p;
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


