---
title: "Steampipe Table: aws_inspector2_finding - Query AWS Inspector findings using SQL"
description: "Allows users to query AWS Inspector findings to gain insights into potential security issues and vulnerabilities within their AWS resources. The table provides detailed information about each finding, including its severity, status, and associated resources."
folder: "Inspector2"
---

# Table: aws_inspector2_finding - Query AWS Inspector findings using SQL

The AWS Inspector is a security assessment service that helps improve the security and compliance of applications deployed on AWS. It automatically assesses applications for vulnerabilities or deviations from best practices, including impacted networks, and insecure configurations. Inspector findings provide detailed information about security vulnerabilities and recommendations for remediation.

## Table Usage Guide

The `aws_inspector2_finding` table in Steampipe provides you with information about findings identified by AWS Inspector. This table allows you, as a security analyst, to query finding-specific details, including their severity, status, and the resources they are associated with. You can utilize this table to gather insights on potential security issues and vulnerabilities within your AWS resources. The schema outlines the various attributes of the findings, including the finding ARN, creation date, severity, status, and associated resources.

When you run an assessment with AWS Inspector, it analyzes your target resources such as EC2 instances, ECS clusters, or RDS databases and generates findings that highlight security vulnerabilities, potential misconfigurations, and other security-related issues. These findings provide you with detailed information about the identified vulnerabilities, including severity levels, affected resources, and recommended remediation steps.

## Examples

### Basic info
Explore which security vulnerabilities exist in your AWS infrastructure and determine their severity and whether fixes are available. This can help prioritize remediation efforts and improve overall security posture.

```sql+postgres
select
  arn,
  description,
  fix_available,
  inspector_score,
  severity,
  finding_account_id
from
  aws_inspector2_finding;
```

```sql+sqlite
select
  arn,
  description,
  fix_available,
  inspector_score,
  severity,
  finding_account_id
from
  aws_inspector2_finding;
```

### List findings with high severity
Discover the segments that have high severity findings in your AWS Inspector data. This is useful for prioritizing issues that require immediate attention due to their potential impact on your AWS resources.

```sql+postgres
select
  arn,
  source,
  vendor_severity,
  status,
  severity
from
  aws_inspector2_finding
where
  severity = 'HIGH';
```

```sql+sqlite
select
  arn,
  source,
  vendor_severity,
  status,
  severity
from
  aws_inspector2_finding
where
  severity = 'HIGH';
```

### Count the number of findings by severity
Analyze the severity of findings from AWS Inspector to understand the distribution and frequency of issues. This can help prioritize remediation efforts based on the severity of identified problems.

```sql+postgres
select
  severity,
  count(severity)
from
  aws_inspector2_finding
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
  aws_inspector2_finding
group by
  severity
order by
  severity;
```

### List findings in last 10 days
Discover the segments that have been identified as potential issues within the past 10 days through AWS Inspector. This is beneficial in maintaining system security by allowing you to promptly address any recent findings.

```sql+postgres
select
  title,
  arn,
  severity
from
  aws_inspector2_finding
where
  last_observed_at >= now() - interval '10' day;
```

```sql+sqlite
select
  title,
  arn,
  severity
from
  aws_inspector2_finding
where
  last_observed_at >= datetime('now', '-10 days');
```

### List suppressed findings
Discover the segments that have been marked as 'suppressed' within your AWS Inspector findings. This can be particularly useful for identifying and managing potential security vulnerabilities that are currently not being addressed.

```sql+postgres
select
  arn,
  status,
  type,
  resources,
  vulnerable_packages
from
  aws_inspector2_finding
where
  status = 'SUPPRESSED';
```

```sql+sqlite
select
  arn,
  status,
  type,
  resources,
  vulnerable_packages
from
  aws_inspector2_finding
where
  status = 'SUPPRESSED';
```

### List package vulnerability findings
Identify instances where software packages have vulnerabilities in your AWS environment. This helps in proactively addressing potential security risks in your system.

```sql+postgres
select
  arn,
  status,
  type,
  resources,
  vulnerable_packages
from
  aws_inspector2_finding
where
  type = 'PACKAGE_VULNERABILITY';
```

```sql+sqlite
select
  arn,
  status,
  type,
  resources,
  vulnerable_packages
from
  aws_inspector2_finding
where
  type = 'PACKAGE_VULNERABILITY';
```

### Get resource details of findings
Explore the specific details of identified resources within findings. This enables a deeper understanding of each finding's context and aids in subsequent decision-making processes.

```sql+postgres
select
  f.arn as finding_arn,
  r ->> 'Id' as resource_id,
  r ->> 'Type' as resource_type,
  r ->> 'Details' as resource_details,
  r ->> 'Partition' as partition,
  r ->> 'Tags' as resource_tags
from
  aws_inspector2_finding as f,
  jsonb_array_elements(resources) as r;
```

```sql+sqlite
select
  f.arn as finding_arn,
  json_extract(r.value, '$.Id') as resource_id,
  json_extract(r.value, '$.Type') as resource_type,
  json_extract(r.value, '$.Details') as resource_details,
  json_extract(r.value, '$.Partition') as partition,
  json_extract(r.value, '$.Tags') as resource_tags
from
  aws_inspector2_finding as f,
  json_each(f.resources) as r;
```

### Get vulnerable package details of findings
Discover the segments that are vulnerable within your system by analyzing the details of problematic packages. This query is useful in identifying potential areas of risk and planning for appropriate remediation measures.

```sql+postgres
select
  f.arn,
  f.vulnerability_id,
  v ->> 'Name' as vulnerability_package_name,
  v ->> 'Version' as vulnerability_package_version,
  v ->> 'Arch' as vulnerability_package_arch,
  v ->> 'Epoch' as vulnerability_package_epoch,
  v ->> 'FilePath' as vulnerability_package_file_path,
  v ->> 'FixedInVersion' as vulnerability_package_fixed_in_version,
  v ->> 'PackageManager' as vulnerability_package_package_manager,
  v ->> 'Release' as vulnerability_package_release,
  v ->> 'Remediation' as vulnerability_package_remediation,
  v ->> 'SourceLambdaLayerArn' as source_lambda_layer_arn,
  v ->> 'Name' as source_layer_hash
from
  aws_inspector2_finding as f,
  jsonb_array_elements(vulnerable_packages) as v;
```

```sql+sqlite
select
  f.arn,
  f.vulnerability_id,
  json_extract(v.value, '$.Name') as vulnerability_package_name,
  json_extract(v.value, '$.Version') as vulnerability_package_version,
  json_extract(v.value, '$.Arch') as vulnerability_package_arch,
  json_extract(v.value, '$.Epoch') as vulnerability_package_epoch,
  json_extract(v.value, '$.FilePath') as vulnerability_package_file_path,
  json_extract(v.value, '$.FixedInVersion') as vulnerability_package_fixed_in_version,
  json_extract(v.value, '$.PackageManager') as vulnerability_package_package_manager,
  json_extract(v.value, '$.Release') as vulnerability_package_release,
  json_extract(v.value, '$.Remediation') as vulnerability_package_remediation,
  json_extract(v.value, '$.SourceLambdaLayerArn') as source_lambda_layer_arn,
  json_extract(v.value, '$.Name') as source_layer_hash
from
  aws_inspector
```

### List exploit available findings
Identify instances where potential vulnerabilities in your AWS infrastructure have known exploits available. This allows you to prioritize urgent threats and address them promptly.

```sql+postgres
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available
from
  aws_inspector2_finding
where
  exploit_available = 'YES';
```

```sql+sqlite
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available
from
  aws_inspector2_finding
where
  exploit_available = 'YES';
```

### List findings that have fixes available through a version update
Identify potential security issues within your system that can be resolved through a version update. This is beneficial for maintaining system integrity and staying ahead of potential vulnerabilities.

```sql+postgres
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available
from
  aws_inspector2_finding
where
  fix_available = 'YES';
```

```sql+sqlite
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available
from
  aws_inspector2_finding
where
  fix_available = 'YES';
```

### List top 5 findings by inspector score
Identify instances where there are critical security findings by ranking them based on the severity of the inspector score. This is useful for prioritizing remediation efforts in your AWS environment.

```sql+postgres
select
  arn,
  inspector_score,
  first_observed_at,
  last_observed_at
  inspector_score_details
from
  aws_inspector2_finding
order by
  inspector_score desc;
```

```sql+sqlite
select
  arn,
  inspector_score,
  first_observed_at,
  last_observed_at,
  inspector_score_details
from
  aws_inspector2_finding
order by
  inspector_score desc;
```

### Get inspector score details of findings
Gain insights into the severity and source of potential security vulnerabilities by analyzing the adjusted CVSS scores of inspection findings. This allows for prioritizing the resolution of the most critical issues and identifying the sources of these vulnerabilities.

```sql+postgres
select
  arn,
  inspector_score_details -> 'AdjustedCvss' ->> 'Score' as adjusted_cvss_score,
  inspector_score_details -> 'AdjustedCvss' ->> 'ScScoreSourceore' as adjusted_cvss_source_score,
  inspector_score_details -> 'AdjustedCvss' ->> 'ScoScoringVectorre' as adjusted_cvss_scoring_vector,
  inspector_score_details -> 'AdjustedCvss' ->> 'Version' as adjusted_cvss_version,
  inspector_score_details -> 'AdjustedCvss' -> 'Adjustments' as adjusted_cvss_adjustments,
  inspector_score_details -> 'AdjustedCvss' ->> 'CvssSource' as adjusted_cvss_cvss_source
from
  aws_inspector2_finding;
```

```sql+sqlite
select
  arn,
  json_extract(inspector_score_details, '$.AdjustedCvss.Score') as adjusted_cvss_score,
  json_extract(inspector_score_details, '$.AdjustedCvss.ScScoreSourceore') as adjusted_cvss_source_score,
  json_extract(inspector_score_details, '$.AdjustedCvss.ScoScoringVectorre') as adjusted_cvss_scoring_vector,
  json_extract(inspector_score_details, '$.AdjustedCvss.Version') as adjusted_cvss_version,
  json_extract(inspector_score_details, '$.AdjustedCvss.Adjustments') as adjusted_cvss_adjustments,
  json_extract(inspector_score_details, '$.AdjustedCvss.CvssSource') as adjusted_cvss_cvss_source
from
  aws_inspector2_finding;
```

### Get network reachability details of findings
Discover the segments that are reachable within your network and the open port ranges. This can help to identify potential vulnerabilities or areas for improvement in network security.

```sql+postgres
select
  arn,
  network_reachability_details -> 'NetworkPath' -> 'Steps' as network_pathsteps,
  network_reachability_details -> 'OpenPortRange' ->> 'Begin' as open_port_range_begin,
  network_reachability_details -> 'OpenPortRange' ->> 'End' as open_port_range_end,
  network_reachability_details -> 'Protocol' as protocol
from
  aws_inspector2_finding;
```

```sql+sqlite
select
  arn,
  json_extract(network_reachability_details, '$.NetworkPath.Steps') as network_pathsteps,
  json_extract(network_reachability_details, '$.OpenPortRange.Begin') as open_port_range_begin,
  json_extract(network_reachability_details, '$.OpenPortRange.End') as open_port_range_end,
  json_extract(network_reachability_details, '$.Protocol') as protocol
from
  aws_inspector2_finding;
```

### List findings by resource tags
Determine the areas in which security findings are linked to specific resource tags. This is particularly useful for identifying potential vulnerabilities within your 'Dev' and 'Prod' environments.

```sql+postgres
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available,
  resource_tags
from
  aws_inspector2_finding
where
  resource_tags = '[{"key": "Name", "value": "Dev"}, {"key": "Name", "value": "Prod"}]';
```

```sql+sqlite
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available,
  resource_tags
from
  aws_inspector2_finding
where
  resource_tags = '[{"key": "Name", "value": "Dev"}, {"key": "Name", "value": "Prod"}]';
```

### List findings by vulnerable packages
Discover the segments that are vulnerable and have potential fixes available. This query is useful to assess the security status of your system and understand where immediate action is required.

```sql+postgres
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available,
  vulnerable_package
from
  aws_inspector2_finding
where
  vulnerable_package = '[{"architecture": "arc", "epoch": "231321", "name": "myVulere", "release": "v0.2.0", "sourceLambdaLayerArn": "arn:aws:lambda:us-west-2:123456789012:layer:my-layer:1", "sourceLayerHash": "dbasjkhda872", "version": "v0.1.0"}]';
```

```sql+sqlite
select
  arn,
  finding_account_id,
  first_observed_at,
  fix_available,
  exploit_available,
  vulnerable_package
from
  aws_inspector2_finding
where
  vulnerable_package = '[{"architecture": "arc", "epoch": "231321", "name": "myVulere", "release": "v0.2.0", "sourceLambdaLayerArn": "arn:aws:lambda:us-west-2:123456789012:layer:my-layer:1", "sourceLayerHash": "dbasjkhda872", "version": "v0.1.0"}]';
```