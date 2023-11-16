---
title: "Table: aws_route53_resolver_rule - Query AWS Route 53 Resolver Rule using SQL"
description: "Allows users to query AWS Route 53 Resolver Rule to obtain data on DNS resolver rules configured in an AWS account."
---

# Table: aws_route53_resolver_rule - Query AWS Route 53 Resolver Rule using SQL

The `aws_route53_resolver_rule` table in Steampipe provides information about DNS resolver rules within AWS Route 53. This table allows DevOps engineers to query resolver rule-specific details, including rule action, domain name, rule type, and associated metadata. Users can utilize this table to gather insights on resolver rules, such as rule configuration, rule status, and rule action. The schema outlines the various attributes of the resolver rule, including the rule ID, resolver endpoint ID, rule action, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_resolver_rule` table, you can use the `.inspect aws_route53_resolver_rule` command in Steampipe.

**Key columns**:

- `id`: The ID of the resolver rule. This is a unique identifier that can be used to join this table with other tables.
- `domain_name`: The domain name that the resolver rule forwards requests for. This can be useful in identifying where DNS queries are being directed.
- `rule_action`: The action that DNS queries take when they match the resolver rule. This provides insight into the behavior of the resolver rule.

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
