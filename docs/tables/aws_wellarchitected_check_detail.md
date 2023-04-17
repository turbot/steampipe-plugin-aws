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

### List total checks per the associated status per workload

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

### Get workload details for trusted advisor checks having errors

```sql
select
  d.question_id,
  d.choice_id,
  d.pillar_id,
  d.name,
  d.status,
  w.workload_name,
  w.workload_id,
  w.environment,
  w.industry,
  w.owner
from
  aws_wellarchitected_check_detail d,
  aws_wellarchitected_workload w
where
  d.workload_id = w.workload_id
  and d.status = 'ERROR';
```