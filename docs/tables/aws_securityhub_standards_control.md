---
title: "Steampipe Table: aws_securityhub_standards_control - Query AWS Security Hub Standards Control using SQL"
description: "Allows users to query AWS Security Hub Standards Control data including details about each security standard control available in an AWS account."
folder: "RDS"
---

# Table: aws_securityhub_standards_control - Query AWS Security Hub Standards Control using SQL

The AWS Security Hub Standards Control is a feature of AWS Security Hub that provides a comprehensive view of the security alerts and security posture across your AWS accounts. This includes continuous monitoring and automated compliance checks against standards such as CIS AWS Foundations Benchmark. It aggregates, organizes, and prioritizes your security alerts, or findings, from multiple AWS services, as well as from AWS partner solutions.

## Table Usage Guide

The `aws_securityhub_standards_control` table in Steampipe provides you with information about each security standard control available in your AWS account. This table allows you, as a DevOps engineer, security analyst, or other professional, to query control-specific details, including its status, related AWS resources, severity, and compliance status. You can utilize this table to gather insights on controls, such as controls that are currently non-compliant, controls that have a high severity level, and more. The schema outlines the various attributes of the standards control for you, including the control ID, control status, related AWS resources, severity, and compliance status.

## Examples

### Basic info
Gain insights into the status and severity rating of various controls in your AWS SecurityHub to ensure your security standards are met. This is crucial for maintaining a robust security posture and promptly addressing any potential vulnerabilities or non-compliance issues.

```sql+postgres
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control;
```

```sql+sqlite
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control;
```

### List disabled controls
Identify instances where certain security controls within AWS Security Hub are disabled, allowing you to assess potential vulnerabilities and take corrective action.

```sql+postgres
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status = 'DISABLED';
```

```sql+sqlite
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status = 'DISABLED';
```

### Count the number of controls by severity
Assess the distribution of security controls based on their severity in your AWS Security Hub. This can help prioritize actions depending on the severity of security controls.

```sql+postgres
select
  severity_rating,
  count(severity_rating)
from
  aws_securityhub_standards_control
group by
  severity_rating
order by
  severity_rating;
```

```sql+sqlite
select
  severity_rating,
  count(severity_rating)
from
  aws_securityhub_standards_control
group by
  severity_rating
order by
  severity_rating;
```

### List controls with high severity
Discover the segments that have high severity in your security controls. This can be particularly useful for prioritizing security measures and addressing the most critical vulnerabilities first.

```sql+postgres
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'HIGH';
```

```sql+sqlite
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'HIGH';
```

### List controls which were updated in the last 30 days
Determine the areas in which your Security Hub controls have been updated recently. This can be useful for keeping track of changes to your security posture and identifying any potential vulnerabilities that need to be addressed.

```sql+postgres
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status_updated_at >= (now() - interval '30' day);
```

```sql+sqlite
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status_updated_at >= datetime('now', '-30 day');
```

### List CIS AWS foundations benchmark controls with critical severity
Determine the areas in which critical severity controls are found within the AWS foundations benchmark. This could be useful for prioritizing areas of concern in your security strategy.

```sql+postgres
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'CRITICAL'
  and arn like '%cis-aws-foundations-benchmark%';
```

```sql+sqlite
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  severity_rating = 'CRITICAL'
  and arn like '%cis-aws-foundations-benchmark%';
```

### List related requirements benchmark for S3 controls
Discover the segments that are related to specific S3 controls in AWS Security Hub. This allows you to better manage and enhance your security posture by understanding the interdependencies between different requirements and controls.

```sql+postgres
select
  control_id,
  r as related_requirements
from
  aws_securityhub_standards_control,
  jsonb_array_elements_text(related_requirements) as r
where
  control_id like '%S3%'
group by
  control_id, r
order by
  control_id, r;
```

```sql+sqlite
select
  control_id,
  json_extract(r.value, '$') as related_requirements
from
  aws_securityhub_standards_control,
  json_each(related_requirements) as r
where
  control_id like '%S3%'
group by
  control_id, json_extract(r.value, '$')
order by
  control_id, json_extract(r.value, '$');
```

### List controls which require PCI DSS benchmark
Discover the segments that have security controls aligned with the PCI DSS benchmark, a crucial step in ensuring your AWS services comply with these important data security standards. This query assists in identifying these specific controls, aiding in the process of regulatory compliance.

```sql+postgres
select
  r as related_requirements,
  control_id
from
  aws_securityhub_standards_control,
  jsonb_array_elements_text(related_requirements) as r
where
  r like '%PCI%'
group by
  r, control_id
order by
  r, control_id;
```

```sql+sqlite
select
  json_extract(r.value, '$') as related_requirements,
  control_id
from
  aws_securityhub_standards_control,
  json_each(related_requirements) as r
where
  json_extract(r.value, '$') like '%PCI%'
group by
  json_extract(r.value, '$'), control_id
order by
  json_extract(r.value, '$'), control_id;
```