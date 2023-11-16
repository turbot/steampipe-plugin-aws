---
title: "Table: aws_guardduty_filter - Query AWS GuardDuty Filter using SQL"
description: "Allows users to query AWS GuardDuty Filters to retrieve information about existing filters, their conditions, actions, and associated metadata."
---

# Table: aws_guardduty_filter - Query AWS GuardDuty Filter using SQL

The `aws_guardduty_filter` table in Steampipe provides information about filters within AWS GuardDuty. This table allows security analysts to query filter-specific details, including filter conditions, actions, and associated metadata. Users can utilize this table to gather insights on filters, such as filter actions, conditions, and the detector ID to which the filter is associated. The schema outlines the various attributes of the GuardDuty filter, including the filter name, detector ID, rank, description, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_guardduty_filter` table, you can use the `.inspect aws_guardduty_filter` command in Steampipe.

### Key columns:

- `name`: The unique name of the filter. This can be used to join this table with other tables that reference GuardDuty filters by their names.
- `detector_id`: The unique ID of the detector that the filter is associated with. This is useful for joining with tables that contain detector-specific information.
- `action`: The action that GuardDuty takes when a filter match occurs. This is crucial for understanding the implications of a filter match.

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
