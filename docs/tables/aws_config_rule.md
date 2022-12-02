# Table: aws_config_rule

An AWS Config rule represents an AWS Lambda function that you create for a custom rule or a predefined function for an AWS managed rule. The function evaluates configuration items to assess whether your AWS resources comply with your desired configurations. This function can run when AWS Config detects a configuration change to an AWS resource and at a periodic frequency that you choose (for example, every 24 hours).

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

### List config rules that runs in proactive mode

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