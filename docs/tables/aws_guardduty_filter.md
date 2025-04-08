---
title: "Steampipe Table: aws_guardduty_filter - Query AWS GuardDuty Filter using SQL"
description: "Allows users to query AWS GuardDuty Filters to retrieve information about existing filters, their conditions, actions, and associated metadata."
folder: "GuardDuty"
---

# Table: aws_guardduty_filter - Query AWS GuardDuty Filter using SQL

The AWS GuardDuty Filter is a feature of AWS GuardDuty that allows you to manage and define conditions for the findings that AWS GuardDuty includes in its threat detection reports. These filters help you categorize and prioritize findings according to your organization's threat model and security posture. You can use these filters to specify the severity level of a finding, the type of threat detected, or other criteria to tailor the findings to your specific needs.

## Table Usage Guide

The `aws_guardduty_filter` table in Steampipe provides you with information about filters within AWS GuardDuty. This table enables you, as a security analyst, to query filter-specific details, including filter conditions, actions, and associated metadata. You can utilize this table to gather insights on filters, such as filter actions, conditions, and the detector ID to which the filter is associated. The schema outlines for you the various attributes of the GuardDuty filter, including the filter name, detector ID, rank, description, and associated tags.

## Examples

### Basic info
Discover the segments that are being monitored by AWS GuardDuty, including their action priority. This can help prioritize security responses and manage potential threats more effectively.

```sql+postgres
select
  name,
  detector_id,
  action,
  rank
from
  aws_guardduty_filter;
```

```sql+sqlite
select
  name,
  detector_id,
  action,
  rank
from
  aws_guardduty_filter;
```

### List filters that will archive the findings
Discover the segments that will archive findings in AWS GuardDuty. This can be beneficial for understanding which filters are set to archive findings, helping to manage security alerts effectively.

```sql+postgres
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

```sql+sqlite
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
Discover the initial filter that will be applied to your findings in AWS GuardDuty. This is useful for understanding the first layer of scrutiny your data will undergo.

```sql+postgres
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

```sql+sqlite
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
Determine the specifics of a particular filter in AWS GuardDuty to better understand its function and configuration. This is useful for assessing the filter's criteria and optimizing security measures.

```sql+postgres
select
  name,
  jsonb_pretty(finding_criteria) as finding_criteria
from
  aws_guardduty_filter
where
  name = 'filter-1';
```

```sql+sqlite
select
  name,
  finding_criteria
from
  aws_guardduty_filter
where
  name = 'filter-1';
```

### Count the number of filters by region and detector
Assess the distribution of filters across various regions and detectors to better understand the security measures in place. This information can be useful for auditing purposes or for optimizing the distribution of filters.

```sql+postgres
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

```sql+sqlite
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
  count(name) desc;
```