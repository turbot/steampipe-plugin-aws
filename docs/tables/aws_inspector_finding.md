---
title: "Table: aws_inspector_finding - Query AWS Inspector Findings using SQL"
description: "Allows users to query AWS Inspector Findings to obtain detailed information about the Amazon Inspector findings that are generated during the assessment of the target applications."
---

# Table: aws_inspector_finding - Query AWS Inspector Findings using SQL

The `aws_inspector_finding` table in Steampipe provides information about AWS Inspector findings. AWS Inspector is an automated security assessment service that helps improve the security and compliance of applications deployed on AWS. This table allows security analysts, developers, and DevOps engineers to query finding-specific details, including the finding ARN, severity, title, description, recommendation, and associated metadata. Users can utilize this table to gather insights on findings, such as findings with high severity, findings associated with a specific rule package, verification of recommendations, and more. The schema outlines the various attributes of the AWS Inspector finding, including the finding ARN, severity, title, description, recommendation, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_inspector_finding` table, you can use the `.inspect aws_inspector_finding` command in Steampipe.

**Key columns**:

- `arn`: The ARN of the finding. This can be used to join with other tables that contain AWS Inspector findings.
- `severity`: The severity of the finding. This is important as it allows users to prioritize findings based on their impact.
- `title`: The title of the finding. This is useful for identifying the type of issue detected by AWS Inspector.

## Examples

### Basic info

```sql
select
  id,
  arn,
  agent_id as instance_id,
  asset_type,
  confidence,
  severity
from
  aws_inspector_finding;
```

### List findings with high severity

```sql
select
  id,
  arn,
  agent_id as instance_id,
  asset_type,
  confidence,
  severity
from
  aws_inspector_finding
where
  severity = 'High';
```

### Count the number of findings by severity

```sql
select
  severity,
  count(severity)
from
  aws_inspector_finding
group by
  severity
order by
  severity;
```

### List last 10 days findings

```sql
select
  title,
  id,
  confidence,
  severity
from
  aws_inspector_finding
where
  created_at >= now() - interval '10' day;
```

### List attributes for each finding

```sql
select
  title,
  id, 
  jsonb_pretty(attributes) as attributes
from
  aws_inspector_finding;
```

### Get asset attributes for each finding

```sql
select
  id,
  title,
  asset_attributes ->> 'AgentId' as agent_id,
  asset_attributes ->> 'AmiId' as ami_id,
  asset_attributes ->> 'Hostname' as hostname,
  asset_attributes ->> 'Tags' as tags
from
  aws_inspector_finding;
```

### List EC2 instances with high severity

```sql
select
  distinct i.instance_id,
  i.instance_state,
  i.instance_type,
  f.title,
  f.service,
  f.severity,
  f.confidence
from
  aws_ec2_instance as i,
  aws_inspector_finding as f
where
  severity = 'High'
and
  i.instance_id = f.agent_id;
```

### Get service attributes for each finding

```sql
select
  id,
  title,
  service_attributes ->> 'AssessmentRunArn' as assessment_run_arn,
  service_attributes ->> 'RulesPackageArn' as rules_package_arn,
  service_attributes ->> 'SchemaVersion' as schema_version,
from
  aws_inspector_finding;
```

### Get assessment run details for findings

```sql
select
  f.id,
  r.title,
  f.service_attributes ->> 'AssessmentRunArn' as assessment_run_arn,
  r.assessment_template_arn,
  r.finding_counts
from
  aws_inspector_finding as f,
  aws_inspector_assessment_run as r
where
  f.service_attributes ->> 'AssessmentRunArn' = r.arn;
```

### List findings order by confidence

```sql
select
  id,
  arn,
  agent_id as instance_id,
  asset_type,
  confidence,
  severity
from
  aws_inspector_finding
order by
  confidence;
```
