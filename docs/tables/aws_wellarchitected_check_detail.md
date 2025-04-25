---
title: "Steampipe Table: aws_wellarchitected_check_detail - Query AWS Well-Architected Tool Check Details using SQL"
description: "Allows users to query AWS Well-Architected Tool Check Details for information on individual checks within a workload. The table provides data on the check status, risk, reason for risk, improvement plan, and other related details."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_check_detail - Query AWS Well-Architected Tool Check Details using SQL

The AWS Well-Architected Tool is a service that helps you review the state of your workloads and compares them to the latest AWS architectural best practices. The tool generates a report detailing areas where your architecture aligns with AWS best practices, and suggests areas for improvement. The Check Details feature specifically provides more granular information about individual checks within your workload's review.

## Table Usage Guide

The `aws_wellarchitected_check_detail` table in Steampipe provides you with information about individual checks within a workload in AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query check-specific details, including check status, risk, reason for risk, and improvement plan. You can utilize this table to gather insights on risk management, workload optimization, and improvement planning. The schema outlines the various attributes of the check detail for you, including the workload ID, lens alias, check ID, and associated metadata.

## Examples

### Basic info
This query allows users to examine the details of their AWS Well-Architected Framework checks, which can provide insights into the status and configuration of their AWS workloads. This can be beneficial for maintaining best practices, identifying potential issues, and ensuring optimal performance of AWS services.

```sql+postgres
select
  workload_id,
  lens_arn,
  pillar_id,
  question_id,
  choice_id,
  id,
  name,
  description,
  status
from
  aws_wellarchitected_check_detail;
```

```sql+sqlite
select
  workload_id,
  lens_arn,
  pillar_id,
  question_id,
  choice_id,
  id,
  name,
  description,
  status
from
  aws_wellarchitected_check_detail;
```

### List total checks per associated status per workload
Discover the segments that contain different workloads, and understand how many checks are associated with each workload status. This can help in assessing the workload's overall health and status efficiently.

```sql+postgres
select
  workload_id,
  status,
  count(id) as checks
from
  aws_wellarchitected_check_detail
group by
  workload_id,
  status;
```

```sql+sqlite
select
  workload_id,
  status,
  count(id) as checks
from
  aws_wellarchitected_check_detail
group by
  workload_id,
  status;
```

### Get check details for security pillar
Explore the specifics of security checks within your AWS architecture. This can help identify areas that require improvements or adjustments to enhance overall security.

```sql+postgres
select
  workload_id,
  lens_arn,
  pillar_id,
  question_id,
  choice_id,
  id,
  name,
  description,
  status
from
  aws_wellarchitected_check_detail
where 
  pillar_id = 'security';
```

```sql+sqlite
select
  workload_id,
  lens_arn,
  pillar_id,
  question_id,
  choice_id,
  id,
  name,
  description,
  status
from
  aws_wellarchitected_check_detail
where 
  pillar_id = 'security';
```

### Get trusted advisor checks with errors
Identify instances where the AWS Trusted Advisor checks have resulted in errors. This can help in promptly addressing problematic areas within your AWS environment, thereby improving system performance and security.

```sql+postgres
select
  id,
  choice_id,
  name,
  pillar_id,
  question_id,
  flagged_resources,
  updated_at
from
  aws_wellarchitected_check_detail
where 
  status = 'ERROR';
```

```sql+sqlite
select
  id,
  choice_id,
  name,
  pillar_id,
  question_id,
  flagged_resources,
  updated_at
from
  aws_wellarchitected_check_detail
where 
  status = 'ERROR';
```

### Get workload details for trusted advisor checks with errors
Identify the workloads in your AWS Well-Architected Framework that have checks with errors. This can help you pinpoint areas that need attention to improve your system's reliability, performance efficiency, cost optimization, operational excellence, and security.

```sql+postgres
select
  w.workload_name,
  w.workload_id,
  w.environment,
  w.industry,
  w.owner,
  d.name as check_name,
  d.flagged_resources,
  d.pillar_id
from
  aws_wellarchitected_check_detail d,
  aws_wellarchitected_workload w
where
  d.workload_id = w.workload_id
  and d.status = 'ERROR';
```

```sql+sqlite
select
  w.workload_name,
  w.workload_id,
  w.environment,
  w.industry,
  w.owner,
  d.name as check_name,
  d.flagged_resources,
  d.pillar_id
from
  aws_wellarchitected_check_detail d,
  aws_wellarchitected_workload w
where
  d.workload_id = w.workload_id
  and d.status = 'ERROR';
```

### Get trusted advisor check details for well-architected lens in a particular workload
Explore the status and details of trusted advisor checks for a specific workload in the well-architected framework. This can help identify potential issues and areas for improvement in your AWS infrastructure.

```sql+postgres
select
  id,
  choice_id,
  name,
  pillar_id,
  question_id,
  flagged_resources,
  status,
  updated_at
from
  aws_wellarchitected_check_detail
where
  lens_arn = 'arn:aws:wellarchitected::aws:lens/wellarchitected'
  and workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```

```sql+sqlite
select
  id,
  choice_id,
  name,
  pillar_id,
  question_id,
  flagged_resources,
  status,
  updated_at
from
  aws_wellarchitected_check_detail
where
  lens_arn = 'arn:aws:wellarchitected::aws:lens/wellarchitected'
  and workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```