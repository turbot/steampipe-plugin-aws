---
title: "Table: aws_securityhub_standards_control - Query AWS Security Hub Standards Control using SQL"
description: "Allows users to query AWS Security Hub Standards Control data including details about each security standard control available in an AWS account."
---

# Table: aws_securityhub_standards_control - Query AWS Security Hub Standards Control using SQL

The `aws_securityhub_standards_control` table in Steampipe provides information about each security standard control available in an AWS account. This table allows DevOps engineers, security analysts, and other professionals to query control-specific details, including its status, related AWS resources, severity, and compliance status. Users can utilize this table to gather insights on controls, such as controls that are currently non-compliant, controls that have a high severity level, and more. The schema outlines the various attributes of the standards control, including the control ID, control status, related AWS resources, severity, and compliance status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_standards_control` table, you can use the `.inspect aws_securityhub_standards_control` command in Steampipe.

**Key columns**:

- `standards_control_arn`: The ARN of the standards control. This can be useful for joining with other tables that reference the standards control by ARN.
- `control_id`: The ID of the control. This is a unique identifier for the control and can be used for specific queries or joins with other tables.
- `compliance_status`: The compliance status of the control. This can be useful for identifying controls that are non-compliant and may need attention.

## Examples

### Basic info

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control;
```

### List disabled controls

```sql
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

```sql
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

```sql
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

```sql
select
  control_id,
  control_status,
  severity_rating
from
  aws_securityhub_standards_control
where
  control_status_updated_at >= (now() - interval '30' day);
```

### List CIS AWS foundations benchmark controls with critical severity

```sql
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

```sql
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

### List controls which require PCI DSS benchmark

```sql
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
