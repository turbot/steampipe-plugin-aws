---
title: "Steampipe Table: aws_config_rule - Query AWS Config Rules using SQL"
description: "Allows users to query Config Rules in AWS Config service. It provides information about each Config Rule, including its name, ARN, description, scope, and compliance status."
folder: "Config"
---

# Table: aws_config_rule - Query AWS Config Rules using SQL

AWS Config Rules is a service that enables you to automate the evaluation of recorded configurations against the desired configurations. With Config Rules, you can review changes to configurations and relationships between AWS resources, dive into detailed resource configuration histories, and determine your overall compliance against the configurations specified in your internal guidelines. This enables you to simplify compliance auditing, security analysis, change management, and operational troubleshooting.

## Table Usage Guide

The `aws_config_rule` table in Steampipe provides you with information about Config Rules within the AWS Config service. This table allows you, as a DevOps engineer, to query rule-specific details, including the rule name, ARN, description, scope, and compliance status. You can utilize this table to gather insights on Config Rules, such as rules that are non-compliant, rules applied to specific resources, and more. The schema outlines the various attributes of the Config Rule for you, including the rule ARN, creation date, input parameters, and associated tags.

## Examples

### Basic info
Explore which AWS configuration rules are in place to gain insights into the current security and compliance state of your AWS resources. This can help identify potential areas of risk or non-compliance.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that consist of inactive rules within your AWS configuration to help identify potential areas for optimization or deletion. This could be useful in maintaining a clean and efficient system by removing or updating unused elements.

```sql+postgres
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

```sql+sqlite
select
  name,
  rule_id,
  arn,
  rule_state
from
  aws_config_rule
where
  rule_state != 'ACTIVE';
```

### List active rules for S3 buckets
Discover the segments that contain active rules for your S3 buckets to better manage and monitor your AWS resources. This is particularly useful for ensuring compliance and security within your cloud storage environment.

```sql+postgres
select
  name,
  rule_id,
  tags
from
  aws_config_rule
where
  name Like '%s3-bucket%';
```

```sql+sqlite
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
Determine the compliance status of a specific AWS Config rule. This is useful to ensure that your AWS resources are following the set rules for approved Amazon Machine Images (AMIs), thereby maintaining a secure and compliant environment.

```sql+postgres
select
  jsonb_pretty(compliance_by_config_rule) as compliance_info
from
  aws_config_rule
where
  name = 'approved-amis-by-id';
```

```sql+sqlite
select
  compliance_by_config_rule
from
  aws_config_rule
where
  name = 'approved-amis-by-id';
```

### List complaince types by config rule
Determine the areas in which your AWS configuration rules are compliant or non-compliant. This can help you identify potential issues and ensure your configurations align with best practices.

```sql+postgres
select
  name as config_rule_name,
  compliance_status -> 'Compliance' -> 'ComplianceType' as compliance_type
from
  aws_config_rule,
  jsonb_array_elements(compliance_by_config_rule) as compliance_status;
```

```sql+sqlite
select
  name as config_rule_name,
  json_extract(compliance_status.value, '$.Compliance.ComplianceType') as compliance_type
from
  aws_config_rule,
  json_each(compliance_by_config_rule) as compliance_status;
```

### List config rules that run in proactive mode
Identify instances where configuration rules are set to operate in proactive mode, which allows for continuous monitoring and automated compliance checks of your system.

```sql+postgres
select
  name as config_rule_name,
  c ->> 'Mode' as evaluation_mode
from
  aws_config_rule,
  jsonb_array_elements(evaluation_modes) as c
where
  c ->> 'Mode' = 'PROACTIVE';
```

```sql+sqlite
select
  name as config_rule_name,
  json_extract(c.value, '$.Mode') as evaluation_mode
from
  aws_config_rule,
  json_each(evaluation_modes) as c
where
  json_extract(c.value, '$.Mode') = 'PROACTIVE';
```