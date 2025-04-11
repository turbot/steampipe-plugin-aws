---
title: "Steampipe Table: aws_wellarchitected_milestone - Query AWS Well-Architected Tool Milestones using SQL"
description: "Allows users to query AWS Well-Architected Tool Milestones for detailed information about the milestones of a workload."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_milestone - Query AWS Well-Architected Tool Milestones using SQL

The AWS Well-Architected Tool Milestone is a resource within the AWS Well-Architected Tool service. It allows you to review and improve your workloads following AWS architectural best practices. The Milestones feature enables you to track changes to your workload over time, documenting architectural decisions and the reasons for those decisions.

## Table Usage Guide

The `aws_wellarchitected_milestone` table in Steampipe provides you with information about the milestones of a workload within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, architect, or developer, to query milestone-specific details, including the milestone name, date, and associated workload information. You can utilize this table to gather insights on milestones, such as the milestone history of a workload, changes made in each milestone, and more. The schema outlines the various attributes of the milestone for you, including the milestone number, record ID, workload ID, and associated tags.

## Examples

### Basic Info
Explore which milestones have been recorded in different regions for AWS workloads to better manage and optimize your cloud architecture. This can help in understanding the progress and geographical distribution of your workloads.

```sql+postgres
select
  workload_id,
  milestone_name,
  milestone_number,
  recorded_at,
  region
from
  aws_wellarchitected_milestone;
```

```sql+sqlite
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
Determine the most recent progress point for each workload in your AWS Well-Architected framework. This can help you track your workloads' evolution and understand where each one stands in terms of its development cycle.

```sql+postgres
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

```sql+sqlite
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
  aws_wellarchitected_milestone m
join
  latest_milestones l
on
  m.milestone_number = l.milestone_number
  and m.workload_id = l.workload_id;
```

### Get workload details associated with each milestone
Identify instances where specific workloads are associated with certain milestones within the AWS Well-Architected framework. This is useful for understanding the distribution of workloads across different milestones, providing insights into project progress and management.

```sql+postgres
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

```sql+sqlite
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
This example demonstrates how to pinpoint specific details related to a certain workload for a particular milestone. This can be beneficial in project management scenarios, where one may need to assess the environment, industry, and ownership aspects of a workload at a specific stage in the project lifecycle.

```sql+postgres
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

```sql+sqlite
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