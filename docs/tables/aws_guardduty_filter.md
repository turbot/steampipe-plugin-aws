# Table: aws_guardduty_filter

A GuardDuty filter allows you to view findings that match the criteria you specify and filter out any unmatched findings. 

## Examples

### Basic info

```sql
select
  name,
  detector_id,
  action,
  rank
from
  aws_guardduty_filter;
```

### List filters that will archive the findings

```sql
select
  name,
  detector_id,
  action,
  rank
from
  aws_guardduty_filter
where
  action = 'ARCHIVE';
```

### Get the filter which will be applied first to the findings

```sql
select
  name,
  region,
  detector_id,
  action,
  rank
from
  aws_guardduty_filter
where
  rank = 1;
```

### Get the criteria details for a filter

```sql
select
  name,
  jsonb_pretty(finding_criteria) as finding_criteria
from
  aws_guardduty_filter
where
  name = 'filter-1';
```

### Count the number of filters by region and detector

```sql
select
  region,
  detector_id,
  count(name)
from
  aws_guardduty_filter
group by
  region,
  detector_id
order by
  count desc;
```
