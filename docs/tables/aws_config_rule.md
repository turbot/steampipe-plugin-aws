---
title: "Table: aws_config_rule - Query AWS Config Rules using SQL"
description: "Allows users to query Config Rules in AWS Config service. It provides information about each Config Rule, including its name, ARN, description, scope, and compliance status."
---

# Table: aws_config_rule - Query AWS Config Rules using SQL

The `aws_config_rule` table in Steampipe provides information about Config Rules within AWS Config service. This table allows DevOps engineers to query rule-specific details, including the rule name, ARN, description, scope, and compliance status. Users can utilize this table to gather insights on Config Rules, such as rules that are non-compliant, rules applied to specific resources, and more. The schema outlines the various attributes of the Config Rule, including the rule ARN, creation date, input parameters, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_config_rule` table, you can use the `.inspect aws_config_rule` command in Steampipe.

### Key columns:

- `name`: The name of the AWS Config rule. This is the primary key of the table and can be used to join with other tables that contain Config rule names.
- `arn`: The Amazon Resource Number (ARN) of the Config rule. This can be used to join with any other table that contains Config rule ARNs.
- `compliance.status`: The compliance status of the Config rule. This column is useful for filtering rules based on their compliance status.

## Examples

### Basic info

```sql
select
  name,
  rule_id,
  arn,
  rule_state,
  created_by,
  scope
from
  aws_config_rule;
```

### List inactive rules

```sql
select
  name,
  rule_id,
  arn,
  rule_state
from
  aws_config_rule
where
  rule_state <> 'ACTIVE';
```

### List active rules for S3 buckets

```sql
select
  name,
  rule_id,
  tags
from
  aws_config_rule
where
  name Like '%s3-bucket%';
```

### List complaince details by config rule

```sql
select
  jsonb_pretty(compliance_by_config_rule) as compliance_info
from
  aws_config_rule
where
  name = 'approved-amis-by-id';
```

### List complaince types by config rule

```sql
select
  name as config_rule_name,
  compliance_status -> 'Compliance' -> 'ComplianceType' as compliance_type
from
  aws_config_rule,
  jsonb_array_elements(compliance_by_config_rule) as compliance_status;
```

### List config rules that run in proactive mode

```sql
select
  name as config_rule_name,
  c ->> 'Mode' as evaluation_mode
from
  aws_config_rule,
  jsonb_array_elements(evaluation_modes) as c
where
  c ->> 'Mode' = 'PROACTIVE';
```
