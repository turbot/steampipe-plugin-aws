# Table: aws_networkfirewall_firewall_policy

An AWS Network Firewall firewall policy defines the monitoring and protection behavior for a firewall. The details of the behavior are defined in the rule groups that you add to your policy, and in some policy default settings. To use a firewall policy, you associate it with one or more firewalls.

## Examples

### Basic info

```sql
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

```sql
select
  arn,
  name,
  description,
  firewall_policy_status,
  encryption_configuration
from
  aws_networkfirewall_firewall_policy
where 
  encryption_configuration ->> 'Type' = `aws_OWNED_KMS_KEY';
```

### List inactive policies

```sql
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

```sql
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatefulDefaultActions' as stateful_default_actions,
  firewall_policy -> 'StatefulRuleGroupReferences' as stateful_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

### Get policy's default stateless actions and rule group details for full packets

```sql
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatelessDefaultActions' as stateless_default_actions,
  firewall_policy -> 'StatelessRuleGroupReferences' as stateless_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

### Get policy's default stateless actions and rule group details for fragmented packets

```sql
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatelessFragmentDefaultActions' as stateless_default_actions,
  firewall_policy -> 'StatelessRuleGroupReferences' as stateless_rule_group_references
from
  aws_networkfirewall_firewall_policy;
```

### Get policy's custom stateless actions

```sql
select
  arn,
  name as firewall_policy_name,
  firewall_policy_status,
  firewall_policy -> 'StatelessRuleGroupReferences' ->> 'ActionName' as custom_action_name,
  firewall_policy -> 'StatelessRuleGroupReferences' ->> 'ActionDefinition' as custom_action_definition
from
  aws_networkfirewall_firewall_policy;
```