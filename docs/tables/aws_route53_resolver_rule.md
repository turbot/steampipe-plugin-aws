---
title: "Steampipe Table: aws_route53_resolver_rule - Query AWS Route 53 Resolver Rule using SQL"
description: "Allows users to query AWS Route 53 Resolver Rule to obtain data on DNS resolver rules configured in an AWS account."
folder: "Route 53"
---

# Table: aws_route53_resolver_rule - Query AWS Route 53 Resolver Rule using SQL

The AWS Route 53 Resolver Rule is a feature of Amazon Route 53, a highly available and scalable cloud Domain Name System (DNS) web service. It allows you to specify how DNS queries from your VPC are routed to your network. This rule can help to simplify DNS operations and enhance security by directing queries to your managed DNS service.

## Table Usage Guide

The `aws_route53_resolver_rule` table in Steampipe provides you with information about DNS resolver rules within AWS Route 53. This table allows you, as a DevOps engineer, to query resolver rule-specific details, including rule action, domain name, rule type, and associated metadata. You can utilize this table to gather insights on resolver rules, such as rule configuration, rule status, and rule action. The schema outlines the various attributes of the resolver rule for you, including the rule ID, resolver endpoint ID, rule action, and associated tags.

## Examples

### Basic info
Explore which domain names are associated with specific resolver rules in AWS Route53. This can help identify areas where rules may need to be updated or shared differently for optimal network routing.

```sql+postgres
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

```sql+sqlite
select
  name,
  domain_name as owner_id,
  resolver_endpoint_id,
  rule_type,
  share_status,
  status
from
  aws_route53_resolver_rule;
```


### List rules that are not associated with VPCs
Discover the segments that are not connected to any Virtual Private Networks (VPNs). This is useful for identifying potential security risks or unused resources within your network infrastructure.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which specific rules allow outbound DNS queries by IP address. This can help assess the elements within your network security setup, providing insights into potential vulnerabilities or areas for optimization.

```sql+postgres
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Port' as port
from
  aws_route53_resolver_rule,
  jsonb_array_elements(target_ips) as p;
```

```sql+sqlite
select
  name,
  json_extract(p.value, '$.Ip') as ip,
  json_extract(p.value, '$.Port') as port
from
  aws_route53_resolver_rule,
  json_each(target_ips) as p;
```


### List resolver rules shared with another account
Identify instances where AWS Route53 resolver rules are shared with another account, in order to better manage and understand your shared resources.

```sql+postgres
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

```sql+sqlite
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