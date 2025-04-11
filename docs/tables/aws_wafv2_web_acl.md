---
title: "Steampipe Table: aws_wafv2_web_acl - Query AWS WAFv2 WebACLs using SQL"
description: "Allows users to query AWS WAFv2 WebACLs to retrieve and manage information about WebACL resources within AWS WAFv2."
folder: "WAFv2"
---

# Table: aws_wafv2_web_acl - Query AWS WAFv2 WebACLs using SQL

AWS WAFv2 WebACLs are a key component of AWS WAF, a web application firewall that helps protect your web applications or APIs against common web exploits that may affect availability, compromise security, or consume excessive resources. The WebACLs allow you to manage a collection of rules that use the same settings. These rules can identify patterns of malicious behavior and take action to block, allow, or count web requests.

## Table Usage Guide

The `aws_wafv2_web_acl` table in Steampipe provides you with information about WebACL resources within AWS WAFv2. This table allows you, as a DevOps engineer, to query WebACL-specific details, including associated rules, actions, visibility configurations, and associated metadata. You can utilize this table to gather insights on WebACLs, such as rules associated with each WebACL, actions for each rule, and the scope of the WebACL. The schema outlines for you the various attributes of the WebACL, including the ARN, capacity, default action, description, and associated tags.

## Examples

### Basic info
Identify instances where your AWS WAFv2 web access control lists (ACLs) are managed by the firewall manager to understand your current security posture and capacity. This could be useful in assessing potential vulnerabilities and planning for capacity management.

```sql+postgres
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

```sql+sqlite
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
Explore the status and metrics of your web access control lists (ACLs) to understand their operational efficiency and performance. This can help in monitoring and managing the security of your web applications.

```sql+postgres
select
  name,
  id,
  visibility_config ->> 'CloudWatchMetricsEnabled' as cloud_watch_metrics_enabled,
  visibility_config ->> 'MetricName' as metric_name
from
  aws_wafv2_web_acl;
```

```sql+sqlite
select
  name,
  id,
  json_extract(visibility_config, '$.CloudWatchMetricsEnabled') as cloud_watch_metrics_enabled,
  json_extract(visibility_config, '$.MetricName') as metric_name
from
  aws_wafv2_web_acl;
```


### List web ACLs whose sampled requests are not enabled
Identify the Web Access Control Lists (ACLs) within your AWS infrastructure where the sampling of requests is disabled. This could be useful for security audits or ensuring optimal configuration settings.

```sql+postgres
select
  name,
  id,
  visibility_config ->> 'SampledRequestsEnabled' as sampled_requests_enabled
from
  aws_wafv2_web_acl
where
  visibility_config ->> 'SampledRequestsEnabled' = 'false';
```

```sql+sqlite
select
  name,
  id,
  json_extract(visibility_config, '$.SampledRequestsEnabled') as sampled_requests_enabled
from
  aws_wafv2_web_acl
where
  json_extract(visibility_config, '$.SampledRequestsEnabled') = 'false';
```


### Get the attack patterns defined in each rule for each web ACL
Identify the specific attack patterns defined in each rule for each web application firewall. This can be useful in understanding the security measures in place and potentially identifying areas for improvement or adjustment.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```


### List regional web ACLs
Determine the areas in which regional web access control lists (ACLs) are in use. This information can be useful for understanding the geographical distribution of your web security measures.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have disabled logging in their web ACLs to understand potential security blind spots in your AWS WAF configuration. This can be used to enhance security measures by ensuring all activities are properly logged and monitored.

```sql+postgres
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

```sql+sqlite
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
This query allows you to pinpoint the specific Application Load Balancers associated with each Web Access Control List. This can be particularly useful in understanding your network's security configuration or identifying potential vulnerabilities.

```sql+postgres
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

```sql+sqlite
select
  lb.name as application_load_balancer_name,
  w.name as web_acl_name,
  w.id as web_acl_id,
  w.scope as web_acl_scope,
  lb.type as application_load_balancer_type
from
  aws_ec2_application_load_balancer as lb,
  aws_wafv2_web_acl as w,
  json_each(associated_resources) as arns
where
  lb.arn = arns.value;
```