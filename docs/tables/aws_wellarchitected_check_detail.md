# Table: aws_wellarchitected_check_detail

Account details for a Well-Architected best practice in relation to Trusted Advisor checks.

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