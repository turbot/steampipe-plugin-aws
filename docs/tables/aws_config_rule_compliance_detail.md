---
title: "Steampipe Table: aws_config_rule_compliance_detail - Query AWS Config Rule Compliance Details using SQL"
description: "Allows users to query AWS Config Rule Compliance Details to get detailed evaluation results for AWS Config rules."
folder: "Config"
---

# Table: aws_config_rule_compliance_detail - Query AWS Config Rule Compliance Details using SQL

AWS Config Rule Compliance Details provide detailed evaluation results for a specific AWS Config rule. This includes which AWS resources were evaluated, when each resource was last evaluated, and whether each resource complies with the rule. This table is useful for in-depth compliance auditing and tracking individual resource compliance status.

## Table Usage Guide

The `aws_config_rule_compliance_detail` table in Steampipe provides detailed compliance evaluation results for AWS Config rules. This table allows you, as a DevOps engineer, compliance officer, or security auditor, to query resource-level compliance details, including evaluation results, timestamps, and compliance status for each resource evaluated by a Config rule. You can utilize this table to identify non-compliant resources, track compliance changes over time, and audit your AWS environment's adherence to compliance standards.

**Important:** This table uses `aws_config_rule` as a parent hydrate, which means:
- Without any filters, it will list compliance details for **all** Config rules in your account
- You can optionally filter by `config_rule_name` to get details for a specific rule
- You can filter by `compliance_type` to only see compliant or non-compliant resources

## Examples

### List all compliance details across all Config rules
Get compliance details for all resources evaluated by all Config rules in your account.

```sql+postgres
select
  config_rule_name,
  resource_type,
  resource_id,
  compliance_type,
  result_recorded_time
from
  aws_config_rule_compliance_detail
order by
  config_rule_name,
  compliance_type;
```

```sql+sqlite
select
  config_rule_name,
  resource_type,
  resource_id,
  compliance_type,
  result_recorded_time
from
  aws_config_rule_compliance_detail
order by
  config_rule_name,
  compliance_type;
```

### Get compliance details for a specific Config rule
Retrieve all evaluation results for a specific AWS Config rule to see which resources are compliant or non-compliant.

```sql+postgres
select
  config_rule_name,
  resource_type,
  resource_id,
  compliance_type,
  config_rule_invoked_time,
  result_recorded_time
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'required-tags-check';
```

```sql+sqlite
select
  config_rule_name,
  resource_type,
  resource_id,
  compliance_type,
  config_rule_invoked_time,
  result_recorded_time
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'required-tags-check';
```

### List all non-compliant resources across all Config rules
Identify all non-compliant resources across your entire AWS Config setup for prioritized remediation.

```sql+postgres
select
  config_rule_name,
  resource_type,
  resource_id,
  annotation,
  result_recorded_time
from
  aws_config_rule_compliance_detail
where
  compliance_type = 'NON_COMPLIANT'
order by
  result_recorded_time desc;
```

```sql+sqlite
select
  config_rule_name,
  resource_type,
  resource_id,
  annotation,
  result_recorded_time
from
  aws_config_rule_compliance_detail
where
  compliance_type = 'NON_COMPLIANT'
order by
  result_recorded_time desc;
```

### Find recently evaluated non-compliant resources
Identify resources that were recently found to be non-compliant for immediate action.

```sql+postgres
select
  config_rule_name,
  resource_type,
  resource_id,
  compliance_type,
  annotation,
  result_recorded_time
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'ec2-instance-managed-by-ssm'
  and compliance_type = 'NON_COMPLIANT'
  and result_recorded_time >= now() - interval '7 days'
order by
  result_recorded_time desc;
```

```sql+sqlite
select
  config_rule_name,
  resource_type,
  resource_id,
  compliance_type,
  annotation,
  result_recorded_time
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'ec2-instance-managed-by-ssm'
  and compliance_type = 'NON_COMPLIANT'
  and result_recorded_time >= datetime('now', '-7 days')
order by
  result_recorded_time desc;
```

### Group non-compliant resources by resource type
Understand which resource types have the most compliance issues for a specific rule.

```sql+postgres
select
  config_rule_name,
  resource_type,
  count(*) as non_compliant_count
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'required-tags-check'
  and compliance_type = 'NON_COMPLIANT'
group by
  config_rule_name,
  resource_type
order by
  non_compliant_count desc;
```

```sql+sqlite
select
  config_rule_name,
  resource_type,
  count(*) as non_compliant_count
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'required-tags-check'
  and compliance_type = 'NON_COMPLIANT'
group by
  config_rule_name,
  resource_type
order by
  non_compliant_count desc;
```

### Join with config rule table to get rule details
Combine compliance details with Config rule metadata for comprehensive reporting.

```sql+postgres
select
  d.config_rule_name,
  r.description as rule_description,
  d.resource_type,
  d.resource_id,
  d.compliance_type,
  d.annotation,
  d.result_recorded_time
from
  aws_config_rule_compliance_detail as d
  inner join aws_config_rule as r on d.config_rule_name = r.name
where
  d.config_rule_name = 's3-bucket-ssl-requests-only'
  and d.compliance_type = 'NON_COMPLIANT';
```

```sql+sqlite
select
  d.config_rule_name,
  r.description as rule_description,
  d.resource_type,
  d.resource_id,
  d.compliance_type,
  d.annotation,
  d.result_recorded_time
from
  aws_config_rule_compliance_detail as d
  inner join aws_config_rule as r on d.config_rule_name = r.name
where
  d.config_rule_name = 's3-bucket-ssl-requests-only'
  and d.compliance_type = 'NON_COMPLIANT';
```

### Track compliance changes over time
Analyze when resources were last evaluated to track compliance monitoring patterns.

```sql+postgres
select
  config_rule_name,
  resource_type,
  compliance_type,
  date_trunc('day', result_recorded_time) as evaluation_date,
  count(*) as evaluation_count
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'iam-password-policy'
  and result_recorded_time >= now() - interval '30 days'
group by
  config_rule_name,
  resource_type,
  compliance_type,
  date_trunc('day', result_recorded_time)
order by
  evaluation_date desc;
```

```sql+sqlite
select
  config_rule_name,
  resource_type,
  compliance_type,
  date(result_recorded_time) as evaluation_date,
  count(*) as evaluation_count
from
  aws_config_rule_compliance_detail
where
  config_rule_name = 'iam-password-policy'
  and result_recorded_time >= datetime('now', '-30 days')
group by
  config_rule_name,
  resource_type,
  compliance_type,
  date(result_recorded_time)
order by
  evaluation_date desc;
```
