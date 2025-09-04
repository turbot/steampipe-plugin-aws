---
title: "Steampipe Table: aws_sesv2_suppressed_destination - Query AWS SES suppressed destinations using SQL"
description: "Allows users to query the account-level suppression list for AWS Simple Email Service (SES)."
folder: "SES"
---

# Table: aws_sesv2_suppressed_destination - Query AWS SES Suppressed Destinations using SQL

AWS Simple Email Service (SES) is a cloud-based email sending service designed to help digital marketers and application developers send marketing, notification, and transactional emails. The service includes a suppression list to help manage bounces and complaints automatically.

## Table Usage Guide

The `aws_sesv2_suppressed_destination` table in Steampipe provides information about email addresses on the account-level suppression list in AWS SES. This table allows you, as a DevOps engineer or an email administrator, to query suppression details, including the reason for suppression and when the address was last updated. You can use this table to monitor email deliverability health and manage your suppression lists.

## Examples

### Basic info
Explore which email addresses are on the suppression list in a specific region.

```sql+postgres
select
  email_address,
  reason,
  last_update_time,
  region
from
  aws_sesv2_suppressed_destination
where
  region = 'us-east-1';
```

```sql+sqlite
select
  email_address,
  reason,
  last_update_time,
  region
from
  aws_sesv2_suppressed_destination
where
  region = 'us-east-1';
```

### List suppressed destinations ordered by last update time
This query is useful for identifying the most recently added email addresses to the suppression list.

```sql+postgres
select
  email_address,
  reason,
  last_update_time
from
  aws_sesv2_suppressed_destination
order by
  last_update_time desc;
```

```sql+sqlite
select
  email_address,
  reason,
  last_update_time
from
  aws_sesv2_suppressed_destination
order by
  last_update_time desc;
```

### Count suppressed destinations by reason
This query helps in understanding the primary reasons for email delivery failures, such as bounces versus complaints.

```sql+postgres
select
  reason,
  count(*) as total
from
  aws_sesv2_suppressed_destination
group by
  reason
order by
  total desc;
```

```sql+sqlite
select
  reason,
  count(*) as total
from
  aws_sesv2_suppressed_destination
group by
  reason
order by
  total desc;
```

### Find all destinations suppressed due to complaints
This query allows you to get a specific list of all email addresses that were added to the suppression list because the recipient marked an email as spam.

```sql+postgres
select
  email_address,
  last_update_time,
  region
from
  aws_sesv2_suppressed_destination
where
  reason = 'COMPLAINT';
```

```sql+sqlite
select
  email_address,
  last_update_time,
  region
from
  aws_sesv2_suppressed_destination
where
  reason = 'COMPLAINT';
```

### Find suppressed destinations within a date range
This query allows you to filter suppressed destinations that were last updated before a specific date, useful for cleanup or analysis of older suppressions.

```sql+postgres
select
  email_address,
  reason,
  last_update_time,
  region
from
  aws_sesv2_suppressed_destination
where
  last_update_time <= '2024-01-01T00:00:00Z'
order by
  last_update_time desc;
```

```sql+sqlite
select
  email_address,
  reason,
  last_update_time,
  region
from
  aws_sesv2_suppressed_destination
where
  last_update_time <= '2024-01-01 00:00:00'
order by
  last_update_time desc;
```
