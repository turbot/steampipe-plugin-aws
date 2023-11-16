---
title: "Table: aws_wellarchitected_milestone - Query AWS Well-Architected Tool Milestones using SQL"
description: "Allows users to query AWS Well-Architected Tool Milestones for detailed information about the milestones of a workload."
---

# Table: aws_wellarchitected_milestone - Query AWS Well-Architected Tool Milestones using SQL

The `aws_wellarchitected_milestone` table in Steampipe provides information about the milestones of a workload within AWS Well-Architected Tool. This table allows DevOps engineers, architects, and developers to query milestone-specific details, including the milestone name, date, and associated workload information. Users can utilize this table to gather insights on milestones, such as milestone history of a workload, changes made in each milestone, and more. The schema outlines the various attributes of the milestone, including the milestone number, record ID, workload ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_milestone` table, you can use the `.inspect aws_wellarchitected_milestone` command in Steampipe.

**Key columns**:

- `milestone_number`: This is the milestone number. It can be used to identify a specific milestone of a workload.
- `workload_id`: This is the ID of the workload that the milestone is associated with. This can be used to join this table with the workload table.
- `record_id`: This is the unique ID of the milestone record. It can be used to uniquely identify a milestone in a workload.

## Examples

### Basic Info

```sql
select
  workload_id,
  milestone_name,
  milestone_number,
  recorded_at,
  region
from
  aws_wellarchitected_milestone;
```

### Get the latest milestone for each workload

```sql
with latest_milestones as 
(
  select
    max(milestone_number) as milestone_number,
    workload_id
  from
    aws_wellarchitected_milestone
  group by
    workload_id
) 
select
  m.workload_id,
  m.milestone_name,
  m.milestone_number as latest_milestone_number,
  m.recorded_at,
  m.region
from
  aws_wellarchitected_milestone m,
  latest_milestones l
where
  m.milestone_number = l.milestone_number
  and m.workload_id = l.workload_id;
```

### Get workload details associated with each milestone

```sql
select
  m.milestone_name,
  m.milestone_number,
  w.workload_name,
  w.workload_id,
  w.environment,
  w.industry,
  w.owner
from
  aws_wellarchitected_workload w,
  aws_wellarchitected_milestone m
where
  w.workload_id = m.workload_id;
```

### Get workload details for a particular milestone

```sql
select
  m.milestone_name,
  m.milestone_number,
  w.workload_name,
  w.workload_id,
  w.environment,
  w.industry,
  w.owner
from
  aws_wellarchitected_workload w,
  aws_wellarchitected_milestone m
where
  w.workload_id = m.workload_id
  and milestone_number = 1
  and w.workload_id = 'abcdec851ac1d8d9d5b9938615da016ce';
```