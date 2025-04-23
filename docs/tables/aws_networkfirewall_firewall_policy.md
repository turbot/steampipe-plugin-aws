---
title: "Steampipe Table: aws_networkfirewall_firewall - Query AWS Network Firewall Policy using SQL"
description: "Allows users to query AWS Network Firewall Policy for detailed information about each firewall's properties and settings."
folder: "Network Firewall"
---

# Table: aws_networkfirewall_firewall_policy - Query AWS Network Firewall Policy using SQL

An AWS Network Firewall Policy defines the behavior of a firewall in a particular stateless or stateful rule group. It sets the actions that are taken on a packet when it matches the rule criteria. The policy can be tailored to fit any network security needs, offering granular control over the traffic passing through the firewall.

## Table Usage Guide

The `aws_networkfirewall_firewall_policy` table in Steampipe provides you with information about each firewall policy in AWS Network Firewall. This table enables you, as a network administrator, security analyst, or DevOps engineer, to query specific details about firewalls, including firewall policies, subnet mappings, and associated VPCs. You can utilize this table to gain insights into firewall configurations, such as firewall policy ARNs, VPC IDs, subnet IDs, and more. The schema outlines the various attributes of the firewall for you, including the firewall ARN, firewall name, firewall policy ARN, VPC ID, subnet mapping, delete protection status, and associated tags.

## Examples

### Basic info
Explore the status and regional distribution of your AWS Network Firewall policies. This allows you to understand the overall security posture and manage resources effectively across various regions.

```sql+postgres
select
  arn,
  name,
  description,
  firewall_policy_status,
  region,
  tags
from
  aws_networkfirewall_firewall_policy;
```

```sql+sqlite
select
  arn,
  name,
  description,
  firewall_policy_status,
  region,
  tags
from
  aws_networkfirewall_firewall_policy;
```

### List policies using default encryption
Determine the areas in your AWS network firewall policies where the default encryption is being used. This is useful for assessing your network's security measures and identifying any potential areas for improvement.

```sql+postgres
select
  arn,
  name,
  description,
  firewall_policy_status,
  encryption_configuration
from
  aws_networkfirewall_firewall_policy
where 
  encryption_configuration ->> 'Type' = 'aws_OWNED_KMS_KEY';
```

```sql+sqlite
select
  arn,
  name,
  description,
  firewall_policy_status,
  encryption_configuration
from
  aws_networkfirewall_firewall_policy
where 
  json_extract(encryption_configuration, '$.Type') = 'aws_OWNED_KMS_KEY';
```

### List inactive policies
Identify instances where certain firewall policies within your AWS Network Firewall are not active. This can help in assessing the security posture of your network and ensure that all necessary policies are in effect.

```sql+postgres
select
  arn,
  name,
  description,
  firewall_policy_status,
  region,
  tags
from
  aws_networkfirewall_firewall_policy
where
  firewall_policy_status != 'ACTIVE';
```

```sql+sqlite
select
  arn,
  name,
  description,
  firewall_policy_status,
  region,
  tags
from
  aws_networkfirewall_firewall_policy
where
  firewall_policy_status != 'ACTIVE';
```

### Get policy's default stateful actions and rule group details
Determine the default actions and rule group details of a policy within a network firewall. This can be useful in understanding the policy's behavior and configuration, which is crucial for managing network security.

```sql+postgres
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatefulDefaultActions' as stateful_default_actions,
  firewall_policy -> 'StatefulRuleGroupReferences' as stateful_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

```sql+sqlite
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  json_extract(firewall_policy, '$.StatefulDefaultActions') as stateful_default_actions,
  json_extract(firewall_policy, '$.StatefulRuleGroupReferences') as stateful_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

### Get policy's default stateless actions and rule group details for full packets
Explore the default actions and rule group details for full packets in a policy to better understand the firewall's configuration and status. This can help in assessing the security measures in place and identifying areas for improvement.

```sql+postgres
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatelessDefaultActions' as stateless_default_actions,
  firewall_policy -> 'StatelessRuleGroupReferences' as stateless_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

```sql+sqlite
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  json_extract(firewall_policy, '$.StatelessDefaultActions') as stateless_default_actions,
  json_extract(firewall_policy, '$.StatelessRuleGroupReferences') as stateless_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

### Get policy's default stateless actions and rule group details for fragmented packets
This query allows you to examine the default actions and rule group details for fragmented packets within a firewall policy. It's particularly useful for understanding your network firewall's behavior and configuration when handling fragmented packets.

```sql+postgres
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatelessFragmentDefaultActions' as stateless_default_actions,
  firewall_policy -> 'StatelessRuleGroupReferences' as stateless_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

```sql+sqlite
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  json_extract(firewall_policy, '$.StatelessFragmentDefaultActions') as stateless_default_actions,
  json_extract(firewall_policy, '$.StatelessRuleGroupReferences') as stateless_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

### Get policy's custom stateless actions
This query is useful for understanding the custom actions associated with your network firewall policies in AWS. It allows you to assess the stateless actions configured and their definitions, enabling you to manage your security measures more effectively.

```sql+postgres
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatelessRuleGroupReferences' ->> 'ActionName' as custom_action_name,
  firewall_policy -> 'StatelessRuleGroupReferences' ->> 'ActionDefinition' as custom_action_definition
from
  aws_networkfirewall_firewall_policy;
```

```sql+sqlite
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  json_extract(firewall_policy, '$.StatelessRuleGroupReferences.ActionName') as custom_action_name,
  json_extract(firewall_policy, '$.StatelessRuleGroupReferences.ActionDefinition') as custom_action_definition
from
  aws_networkfirewall_firewall_policy;
```