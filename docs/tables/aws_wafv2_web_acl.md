---
title: "Table: aws_wafv2_web_acl - Query AWS WAFv2 WebACLs using SQL"
description: "Allows users to query AWS WAFv2 WebACLs to retrieve and manage information about WebACL resources within AWS WAFv2."
---

# Table: aws_wafv2_web_acl - Query AWS WAFv2 WebACLs using SQL

The `aws_wafv2_web_acl` table in Steampipe provides information about WebACL resources within AWS WAFv2. This table allows DevOps engineers to query WebACL-specific details, including associated rules, actions, visibility configurations, and associated metadata. Users can utilize this table to gather insights on WebACLs, such as rules associated with each WebACL, actions for each rule, and the scope of the WebACL. The schema outlines the various attributes of the WebACL, including the ARN, capacity, default action, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wafv2_web_acl` table, you can use the `.inspect aws_wafv2_web_acl` command in Steampipe.

### Key columns:

- `arn`: The Amazon Resource Name (ARN) of the WebACL. This can be used to join with other tables that provide WebACL-specific information.
- `id`: The unique identifier of the WebACL. This can be used to join with other tables that provide WebACL-specific information.
- `name`: The name of the WebACL. This can be used to join with other tables that provide WebACL-specific information.

## Examples

### Basic info

```sql
select
  name,
  id,
  scope,
  description,
  capacity,
  managed_by_firewall_manager
from
  aws_wafv2_web_acl;
```


### Get CloudWatch metrics details for each web ACL

```sql
select
  name,
  id,
  visibility_config ->> 'CloudWatchMetricsEnabled' as cloud_watch_metrics_enabled,
  visibility_config ->> 'MetricName' as metric_name
from
  aws_wafv2_web_acl;
```


### List web ACLs whose sampled requests are not enabled

```sql
select
  name,
  id,
  visibility_config ->> 'SampledRequestsEnabled' as sampled_requests_enabled
from
  aws_wafv2_web_acl
where
  visibility_config ->> 'SampledRequestsEnabled' = 'false';
```


### Get the attack patterns defined in each rule for each web ACL

```sql
select
  name,
  id,
  r ->> 'Name' as name,
  r -> 'Statement' ->> 'AndStatement' as and_statement,
  r -> 'Statement' ->> 'ByteMatchStatement' as byte_match_statement,
  r -> 'Statement' ->> 'GeoMatchStatement' as geo_match_statement,
  r -> 'Statement' ->> 'IPSetReferenceStatement' as ip_set_reference_statement,
  r -> 'Statement' ->> 'NotStatement' as not_statement,
  r -> 'Statement' ->> 'OrStatement' as or_statement,
  r -> 'Statement' ->> 'RateBasedStatement' as rate_based_statement,
  r -> 'Statement' ->> 'RegexPatternSetReferenceStatement' as regex_pattern_set_reference_statement,
  r -> 'Statement' ->> 'RuleGroupReferenceStatement' as rule_group_reference_statement,
  r -> 'Statement' ->> 'SizeConstraintStatement' as size_constraint_statement,
  r -> 'Statement' ->> 'SqliMatchStatement' as sql_match_statement,
  r -> 'Statement' ->> 'XssMatchStatement' as xss_match_statement
from
  aws_wafv2_web_acl,
  jsonb_array_elements(rules) as r;
```


### List regional web ACLs

```sql
select
  name,
  id,
  scope,
  region
from
  aws_wafv2_web_acl
where
  scope = 'REGIONAL';
```


### List web ACLs with logging disabled

```sql
select
  name,
  id,
  scope,
  region
from
  aws_wafv2_web_acl
where
  logging_configuration is null;
```

### Get details for ALBs associated with each web ACL

```sql
select
  lb.name as application_load_balancer_name,
  w.name as web_acl_name,
  w.id as web_acl_id,
  w.scope as web_acl_scope,
  lb.type as application_load_balancer_type
from
  aws_ec2_application_load_balancer as lb,
  aws_wafv2_web_acl as w,
  jsonb_array_elements_text(associated_resources) as arns
where
  lb.arn = arns;
```