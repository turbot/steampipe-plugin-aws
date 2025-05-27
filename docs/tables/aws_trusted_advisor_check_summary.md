---
title: "Steampipe Table: aws_trusted_advisor_check_summary - Query AWS Trusted Advisor Check Summaries using SQL"
description: "A Trusted Advisor check is a specific evaluation or assessment performed by Trusted Advisor in different categories. These checks cover various areas, including cost optimization, security, performance, and fault tolerance. Each check examines a specific aspect of your AWS resources and provides recommendations for improvement."
folder: "Trusted Advisor"
---

# Table: aws_trusted_advisor_check_summary - Query AWS Trusted Advisor Check Summaries using SQL

A Trusted Advisor check is a specific evaluation or assessment performed by Trusted Advisor in different categories. These checks cover various areas, including cost optimization, security, performance, and fault tolerance. Each check examines a specific aspect of your AWS resources and provides recommendations for improvement.

## Table Usage Guide

The `aws_trusted_advisor_check_summary` table in Steampipe allows users to query information about AWS Trusted Advisor Check Summaries. These summaries provide insights into the status and details of Trusted Advisor checks, including the check name, check ID, category, description, status, timestamp, and the number of resources flagged, ignored, processed, and suppressed.

**Important Notes**
- You must specify `language` in a `where` clause in order to use this table.
- Amazon Web Services Support API currently supports the following languages for Trusted Advisor:
  - Chinese, Simplified - zh
  - Chinese, Traditional - zh_TW
  - English - en
  - French - fr
  - German - de
  - Indonesian - id
  - Italian - it
  - Japanese - ja
  - Korean - ko
  - Portuguese, Brazilian - pt_BR
  - Spanish - es

## Examples

### Basic info
Retrieve basic information about AWS Trusted Advisor Check Summaries, including the check name, check ID, category, description, status, timestamp, and the number of resources flagged.

```sql+postgres
select
  name,
  check_id,
  category,
  description,
  status,
  timestamp,
  resources_flagged
from
  aws_trusted_advisor_check_summary
where
  language = 'en';
```

```sql+sqlite
select
  name,
  check_id,
  category,
  description,
  status,
  timestamp,
  resources_flagged
from
  aws_trusted_advisor_check_summary
where
  language = 'en';
```

### Get error check summaries
Retrieve AWS Trusted Advisor Check Summaries with an "error" status. This query helps you identify checks that require attention due to errors.

```sql+postgres
select
  name,
  check_id,
  category,
  status
from
  aws_trusted_advisor_check_summary
where
  language = 'en'
and
  status = 'error';
```

```sql+sqlite
select
  name,
  check_id,
  category,
  status
from
  aws_trusted_advisor_check_summary
where
  language = 'en'
and
  status = 'error';
```

### Get check summaries for the last 5 days
Retrieve AWS Trusted Advisor Check Summaries from the last 5 days. This query allows you to review recent check summaries and their details.

```sql+postgres
select
  name,
  check_id,
  description,
  status,
  timestamp
from
  aws_trusted_advisor_check_summary
where
  language = 'en'
and
  timestamp >= now() - interval '5 day';
```

```sql+sqlite
select
  name,
  check_id,
  description,
  status,
  timestamp
from
  aws_trusted_advisor_check_summary
where
  language = 'en'
and
  timestamp >= datetime('now', '-5 day');
```

### Get resource summaries of each check
Retrieve resource summaries for each AWS Trusted Advisor Check. This includes the number of resources flagged, ignored, processed, and suppressed for each check.

```sql+postgres
select
  name,
  check_id,
  resources_flagged,
  resources_ignored,
  resources_processed,
  resources_suppressed
from
  aws_trusted_advisor_check_summary
where
  language = 'en';
```

```sql+sqlite
select
  name,
  check_id,
  resources_flagged,
  resources_ignored,
  resources_processed,
  resources_suppressed
from
  aws_trusted_advisor_check_summary
where
  language = 'en';
```
