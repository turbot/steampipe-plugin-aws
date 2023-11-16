---
title: "Table: aws_wellarchitected_check_detail - Query AWS Well-Architected Tool Check Details using SQL"
description: "Allows users to query AWS Well-Architected Tool Check Details for information on individual checks within a workload. The table provides data on the check status, risk, reason for risk, improvement plan, and other related details."
---

# Table: aws_wellarchitected_check_detail - Query AWS Well-Architected Tool Check Details using SQL

The `aws_wellarchitected_check_detail` table in Steampipe provides information about individual checks within a workload in AWS Well-Architected Tool. This table allows DevOps engineers to query check-specific details, including check status, risk, reason for risk, and improvement plan. Users can utilize this table to gather insights on risk management, workload optimization, and improvement planning. The schema outlines the various attributes of the check detail, including the workload ID, lens alias, check ID, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_check_detail` table, you can use the `.inspect aws_wellarchitected_check_detail` command in Steampipe.

**Key columns**:

- `workload_id`: This is the unique identifier for a workload. It can be used to join this table with other tables that contain workload-specific information.
- `lens_alias`: This represents the lens through which the check was made. It can be used to join with other tables that contain lens-specific data.
- `check_id`: This is the unique identifier for a check. It can be used to join this table with other tables that contain check-specific information.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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

```sql
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

```sql
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

```sql
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