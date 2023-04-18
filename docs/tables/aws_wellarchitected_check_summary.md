# Table: aws_wellarchitected_check_summary

Trusted Advisor check summary.

## Examples

### Basic info

```sql
select
  id,
  name,
  description,
  jsonb_pretty(account_summary) as account_summary,
  choice_id,
  lens_arn,
  pillar_id,
  question_id,
  status,
  region,
  workload_id
from
  aws_wellarchitected_check_summary;
```

### Get summarized trusted advisor check report for a workload

```sql
select
  workload_id,
  id,
  name,
  jsonb_pretty(account_summary) as account_summary,
  status,
  choice_id,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```

### List trusted advisor checks with errors

```sql
select
  workload_id,
  id,
  name,
  jsonb_pretty(account_summary) as account_summary,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  status = 'ERROR';
```

### Get account summary for trusted advisor checks

```sql
select
  workload_id,
  id,
  name,
  account_summary ->> 'ERROR' as errors,
  account_summary ->> 'FETCH_FAILED' as fetch_failed,
  account_summary ->> 'NOT_AVAILABLE' as not_available,
  account_summary ->> 'OKAY' as okay,
  account_summary ->> 'WARNING' as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary;
```

### Get account summary for trusted advisor checks for wellarchitected lens in a particular workload

```sql
select
  workload_id,
  id,
  name,
  account_summary ->> 'ERROR' as errors,
  account_summary ->> 'FETCH_FAILED' as fetch_failed,
  account_summary ->> 'NOT_AVAILABLE' as not_available,
  account_summary ->> 'OKAY' as okay,
  account_summary ->> 'WARNING' as warnings,
  pillar_id,
  question_id
from
  aws_wellarchitected_check_summary
where
  lens_arn = 'arn:aws:wellarchitected::aws:lens/wellarchitected'
  and workload_id = 'abcdc851ac1d8d9d5b9938615da016ce';
```