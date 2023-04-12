# Table: aws_wellarchitected_milestone

A milestone records the state of a workload at a particular point in time. Save a milestone after you initially complete all the questions associated with a workload. As you change your workload based on items in your improvement plan, you can save additional milestones to measure progress.

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
with latest_milestones as (
  select
    max(milestone_number) as milestone_number,
    workload_id
  from
    aws_wellarchitected_milestone
  group by
    workload_id
) select
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

### Get workload details associated to each milestone

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