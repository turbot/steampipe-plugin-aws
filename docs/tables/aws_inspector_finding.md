---
title: "Steampipe Table: aws_inspector_finding - Query AWS Inspector Findings using SQL"
description: "Allows users to query AWS Inspector Findings to obtain detailed information about the Amazon Inspector findings that are generated during the assessment of the target applications."
folder: "Inspector"
---

# Table: aws_inspector_finding - Query AWS Inspector Findings using SQL

The AWS Inspector Finding is a resource within AWS Inspector service that allows you to identify potential security issues and deviations from best practices. It provides detailed descriptions of the security findings and offers recommendations on how to fix them. AWS Inspector automatically assesses applications for vulnerabilities or deviations from best practices, including impacted networks, instances, and attached storage.

## Table Usage Guide

The `aws_inspector_finding` table in Steampipe provides you with information about AWS Inspector findings. AWS Inspector is an automated security assessment service that helps you improve the security and compliance of applications deployed on AWS. This table allows you, as a security analyst, developer, or DevOps engineer, to query finding-specific details, including the finding ARN, severity, title, description, recommendation, and associated metadata. You can utilize this table to gather insights on findings, such as findings with high severity, findings associated with a specific rule package, verification of recommendations, and more. The schema outlines the various attributes of the AWS Inspector finding for you, including the finding ARN, severity, title, description, recommendation, and associated tags.

## Examples

### Basic info
Explore which AWS Inspector findings have been identified, focusing on their severity and confidence levels. This can help in prioritizing remediation efforts based on the severity and confidence of the findings.

```sql+postgres
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

```sql+sqlite
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
Identify instances where there are high severity findings in the AWS Inspector. This is useful in prioritizing security issues that need immediate attention.

```sql+postgres
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

```sql+sqlite
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
Analyze the severity levels of AWS inspector findings to understand how many issues fall into each category. This can help prioritize remediation efforts by focusing on the most severe findings first.

```sql+postgres
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

```sql+sqlite
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
Identify instances where security findings have been recorded in the last 10 days. This allows you to stay updated on recent security issues and take necessary actions.

```sql+postgres
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

```sql+sqlite
select
  title,
  id,
  confidence,
  severity
from
  aws_inspector_finding
where
  created_at >= datetime('now','-10 days');
```

### List attributes for each finding
Determine the characteristics of each identified issue within your AWS Inspector service. This can help in understanding the nature of the problems and strategizing appropriate solutions.

```sql+postgres
select
  title,
  id, 
  jsonb_pretty(attributes) as attributes
from
  aws_inspector_finding;
```

```sql+sqlite
select
  title,
  id, 
  attributes
from
  aws_inspector_finding;
```

### Get asset attributes for each finding
This query is used to uncover the details of each asset's attributes associated with a specific finding in AWS Inspector. This can help in identifying instances where anomalies or issues have been detected, providing insights into potential areas of risk or concern within your AWS environment.

```sql+postgres
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

```sql+sqlite
select
  id,
  title,
  json_extract(asset_attributes, '$.AgentId') as agent_id,
  json_extract(asset_attributes, '$.AmiId') as ami_id,
  json_extract(asset_attributes, '$.Hostname') as hostname,
  json_extract(asset_attributes, '$.Tags') as tags
from
  aws_inspector_finding;
```

### List EC2 instances with high severity
Discover the segments that are operating Amazon EC2 instances with high severity findings. This is useful for identifying potential security vulnerabilities and risks in your AWS infrastructure.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which specific service attributes are linked to each finding, enabling a more comprehensive understanding of the findings in AWS Inspector. This can assist in better assessment planning and rule package selection for future inspections.

```sql+postgres
select
  id,
  title,
  service_attributes ->> 'AssessmentRunArn' as assessment_run_arn,
  service_attributes ->> 'RulesPackageArn' as rules_package_arn,
  service_attributes ->> 'SchemaVersion' as schema_version
from
  aws_inspector_finding;
```

```sql+sqlite
select
  id,
  title,
  json_extract(service_attributes, '$.AssessmentRunArn') as assessment_run_arn,
  json_extract(service_attributes, '$.RulesPackageArn') as rules_package_arn,
  json_extract(service_attributes, '$.SchemaVersion') as schema_version
from
  aws_inspector_finding;
```

### Get assessment run details for findings
This query is used to analyze the details of assessment runs linked to specific findings in AWS Inspector. It's useful for identifying potential security vulnerabilities and understanding the scope of any issues identified during the assessment runs.

```sql+postgres
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

```sql+sqlite
select
  f.id,
  r.title,
  json_extract(f.service_attributes, '$.AssessmentRunArn') as assessment_run_arn,
  r.assessment_template_arn,
  r.finding_counts
from
  aws_inspector_finding as f
join
  aws_inspector_assessment_run as r
on
  json_extract(f.service_attributes, '$.AssessmentRunArn') = r.arn;
```

### List findings order by confidence
Explore which AWS Inspector findings are most reliable by sorting them according to their confidence levels. This can help prioritize remediation efforts by focusing first on findings with the highest confidence.

```sql+postgres
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

```sql+sqlite
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