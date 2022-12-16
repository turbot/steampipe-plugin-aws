# Table: aws_drs_job

AWS Elastic Disaster Recovery (AWS DRS) Job provides information about currently running jobs, and completed jobs. Jobs are created when you launch Drill or Recovery instances, terminate launched instances, or perform a fallback.

## Examples

### Basic Info

```sql
select
  title,
  arn,
  status,
  initiated_by
from
  aws_drs_job;
```

### List jobs that are in pending state

```sql
select
  title,
  arn,
  status,
  initiated_by,
  creation_date_time
from
  aws_drs_job
where
  status = 'PENDING';
```

### List jobs that were started in past 30 days

```sql
select
  title,
  arn,
  status,
  initiated_by,
  type,
  creation_date_time,
  end_date_time
from
  aws_drs_job
where
  creation_date_time >= now() - interval '30' day;
```
